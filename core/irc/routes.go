package irc

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/iotaledger/inx-app/pkg/httpserver"
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
		resp, err := loadIRC27(c)
		if err != nil {
			return err
		}

		return httpserver.JSONResponse(c, http.StatusOK, resp)
	})

	e.POST(RouteIRC30, func(c echo.Context) error {
		resp, err := loadIRC30(c)
		if err != nil {
			return err
		}

		return httpserver.JSONResponse(c, http.StatusOK, resp)
	})
}
