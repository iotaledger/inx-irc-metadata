package irc

import (
	"encoding/json"

	"github.com/labstack/echo/v4"

	"github.com/iotaledger/inx-app/pkg/httpserver"
)

func loadIRC27(c echo.Context) (interface{}, error) {
	nftID, err := httpserver.ParseNFTIDParam(c, ParameterNFTID)
	if err != nil {
		return nil, err
	}

	client := deps.NodeBridge.INXNodeClient()
	ctx := c.Request().Context()

	indexer, err := client.Indexer(ctx)
	if err != nil {
		return nil, err
	}

	_, nftOutput, err := indexer.NFT(ctx, *nftID)
	if err != nil {
		return nil, err
	}

	features, err := nftOutput.ImmutableFeatures.Set()
	if err != nil {
		return nil, err
	}

	metadata := features.MetadataFeature()
	if metadata == nil {
		return nil, httpserver.ErrNotAcceptable
	}

	var irc27 interface{}
	if err := json.Unmarshal(metadata.Data, &irc27); err != nil {
		return nil, httpserver.ErrNotAcceptable
	}

	if err := deps.IRC27Schema.Validate(irc27); err != nil {
		return nil, httpserver.ErrNotAcceptable
	}

	return irc27, nil
}

func loadIRC30(c echo.Context) (interface{}, error) {
	foundryId, err := httpserver.ParseFoundryIDParam(c, ParameterNativeTokenID)
	if err != nil {
		return nil, err
	}

	client := deps.NodeBridge.INXNodeClient()
	ctx := c.Request().Context()

	indexer, err := client.Indexer(ctx)
	if err != nil {
		return nil, err
	}

	_, foundryOutput, err := indexer.Foundry(ctx, *foundryId)
	if err != nil {
		return nil, err
	}

	features, err := foundryOutput.ImmutableFeatures.Set()
	if err != nil {
		return nil, err
	}

	metadata := features.MetadataFeature()
	if metadata == nil {
		return nil, httpserver.ErrNotAcceptable
	}

	var irc30 interface{}
	if err := json.Unmarshal(metadata.Data, &irc30); err != nil {
		return nil, httpserver.ErrNotAcceptable
	}

	if err := deps.IRC30Schema.Validate(irc30); err != nil {
		return nil, httpserver.ErrNotAcceptable
	}

	return irc30, nil
}
