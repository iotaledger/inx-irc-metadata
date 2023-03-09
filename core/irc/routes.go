package irc

import (
	"github.com/labstack/echo/v4"
)

const (
	APIRoute = "irc-metadata/v1"

	// ParameterNFTID is used to identify a NFT by its ID.
	ParameterNFTID = "nftID"

	// ParameterNativeTokenID is used to identify a native token by its ID.
	ParameterNativeTokenID = "tokenID"

	RouteIRC27 = "/nfts/:" + ParameterNFTID
	RouteIRC30 = "/tokens/:" + ParameterNativeTokenID
)

func setupRoutes(e *echo.Echo) {
	e.GET(RouteIRC27, func(c echo.Context) error {
		return deps.IRC27Validator.HandleRequest(c)
	})

	e.GET(RouteIRC30, func(c echo.Context) error {
		return deps.IRC30Validator.HandleRequest(c)
	})
}
