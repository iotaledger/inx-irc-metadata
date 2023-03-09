package irc

import (
	"encoding/json"

	"github.com/labstack/echo/v4"

	"github.com/iotaledger/inx-app/pkg/httpserver"
)

func loadIRC27(c echo.Context) (*IRC27, error) {
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

	irc27 := new(IRC27)
	if err := json.Unmarshal(metadata.Data, irc27); err != nil {
		return nil, httpserver.ErrNotAcceptable
	}

	return irc27, httpserver.ErrInvalidParameter
}

func loadIRC30(c echo.Context) (*IRC30, error) {
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

	irc30 := new(IRC30)
	if err := json.Unmarshal(metadata.Data, irc30); err != nil {
		return nil, httpserver.ErrNotAcceptable
	}

	return irc30, httpserver.ErrInvalidParameter
}
