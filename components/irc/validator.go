//nolint:structcheck
package irc

import (
	"context"
	"encoding/json"
	"net/http"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/labstack/echo/v4"
	"github.com/santhosh-tekuri/jsonschema/v5"

	//nolint:golint // this is on purpose
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"

	"github.com/iotaledger/hive.go/ierrors"
	"github.com/iotaledger/inx-app/pkg/httpserver"
)

var (
	ErrLoadMetadataNotFound = ierrors.New("metadata not found")
	ErrLoadMetadataInvalid  = ierrors.New("invalid metadata")
)

type MetadataValidator[K comparable] struct {
	schema *jsonschema.Schema
	cache  *lru.Cache[K, *cachedMetadata]

	parseKeyFunc  func(c echo.Context) (K, error)
	loadValueFunc func(ctx context.Context, key K) ([]byte, error)
}

func NewMetadataValidator[K comparable](schemaURL string, cacheSize int, parseKeyFunc func(c echo.Context) (K, error), loadValueFunc func(ctx context.Context, key K) ([]byte, error)) (*MetadataValidator[K], error) {
	schema, err := jsonschema.Compile(schemaURL)
	if err != nil {
		return nil, err
	}

	cache, err := lru.New[K, *cachedMetadata](cacheSize)
	if err != nil {
		return nil, err
	}

	return &MetadataValidator[K]{
		schema:        schema,
		cache:         cache,
		parseKeyFunc:  parseKeyFunc,
		loadValueFunc: loadValueFunc,
	}, nil
}

func (v *MetadataValidator[K]) get(ctx context.Context, key K) (*cachedMetadata, error) {
	cached, has := v.cache.Get(key)
	if has {
		return cached, nil
	}

	bytes, err := v.loadValueFunc(ctx, key)
	if err != nil {
		if ierrors.Is(err, ErrLoadMetadataNotFound) {
			return v.storeNotFound(key), nil
		}
		if ierrors.Is(err, ErrLoadMetadataInvalid) {
			return v.storeInvalid(key), nil
		}

		return nil, err
	}

	return v.validate(key, bytes), nil
}

func (v *MetadataValidator[K]) validate(key K, bytes []byte) *cachedMetadata {
	var metadata interface{}
	if err := json.Unmarshal(bytes, &metadata); err != nil {
		return v.storeInvalid(key)
	}

	if err := v.schema.Validate(metadata); err != nil {
		return v.storeInvalid(key)
	}

	cached := newCachedMetadata(bytes)
	v.cache.Add(key, cached)

	return cached
}

func (v *MetadataValidator[K]) storeNotFound(key K) *cachedMetadata {
	cached := newCachedError(echo.ErrNotFound)
	v.cache.Add(key, cached)

	return cached
}

func (v *MetadataValidator[K]) storeInvalid(key K) *cachedMetadata {
	cached := newCachedError(httpserver.ErrNotAcceptable)
	v.cache.Add(key, cached)

	return cached
}

func (v *MetadataValidator[K]) HandleRequest(c echo.Context) error {
	key, err := v.parseKeyFunc(c)
	if err != nil {
		return err
	}

	cached, err := v.get(c.Request().Context(), key)
	if err != nil {
		return err
	}

	if cached.Error != nil {
		return cached.Error
	}

	return c.Blob(http.StatusOK, echo.MIMEApplicationJSONCharsetUTF8, cached.Bytes)
}

type cachedMetadata struct {
	Bytes []byte
	Error error
}

func newCachedError(err error) *cachedMetadata {
	return &cachedMetadata{
		Bytes: nil,
		Error: err,
	}
}

func newCachedMetadata(bytes []byte) *cachedMetadata {
	return &cachedMetadata{
		Bytes: bytes,
		Error: nil,
	}
}
