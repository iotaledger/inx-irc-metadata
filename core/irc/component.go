package irc

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"go.uber.org/dig"

	"github.com/iotaledger/hive.go/core/app"
	"github.com/iotaledger/inx-app/pkg/httpserver"
	"github.com/iotaledger/inx-app/pkg/nodebridge"
	"github.com/iotaledger/inx-irc-metadata/pkg/daemon"
)

const (
	IRC27SchemaURL = "https://raw.githubusercontent.com/iotaledger/tips/main/tips/TIP-0027/irc27.schema.json"
	IRC30SchemaURL = "https://raw.githubusercontent.com/iotaledger/tips/main/tips/TIP-0030/irc30.schema.json"
)

func init() {
	CoreComponent = &app.CoreComponent{
		Component: &app.Component{
			Name:     "IRC",
			Params:   params,
			DepsFunc: func(cDeps dependencies) { deps = cDeps },
			Provide:  provide,
			Run:      run,
		},
	}
}

var (
	CoreComponent *app.CoreComponent
	deps          dependencies
)

type dependencies struct {
	dig.In
	NodeBridge  *nodebridge.NodeBridge
	IRC27Schema *jsonschema.Schema `name:"irc27Schema"`
	IRC30Schema *jsonschema.Schema `name:"irc30Schema"`
}

func provide(c *dig.Container) error {
	type outDeps struct {
		dig.Out
		IRC27Schema *jsonschema.Schema `name:"irc27Schema"`
		IRC30Schema *jsonschema.Schema `name:"irc30Schema"`
	}

	return c.Provide(func() outDeps {
		irc27, err := jsonschema.Compile(IRC27SchemaURL)
		if err != nil {
			panic(err)
		}

		irc30, err := jsonschema.Compile(IRC30SchemaURL)
		if err != nil {
			panic(err)
		}

		return outDeps{
			IRC27Schema: irc27,
			IRC30Schema: irc30,
		}
	})

}

func run() error {
	// create a background worker that handles the API
	if err := CoreComponent.Daemon().BackgroundWorker("API", func(ctx context.Context) {
		CoreComponent.LogInfo("Starting API ... done")

		e := httpserver.NewEcho(CoreComponent.Logger(), nil, ParamsRestAPI.DebugRequestLoggerEnabled)

		CoreComponent.LogInfo("Starting API server ...")

		setupRoutes(e)
		go func() {
			CoreComponent.LogInfof("You can now access the API using: http://%s", ParamsRestAPI.BindAddress)
			if err := e.Start(ParamsRestAPI.BindAddress); err != nil && !errors.Is(err, http.ErrServerClosed) {
				CoreComponent.LogErrorfAndExit("Stopped REST-API server due to an error (%s)", err)
			}
		}()

		ctxRegister, cancelRegister := context.WithTimeout(ctx, 5*time.Second)

		advertisedAddress := ParamsRestAPI.BindAddress
		if ParamsRestAPI.AdvertiseAddress != "" {
			advertisedAddress = ParamsRestAPI.AdvertiseAddress
		}

		if err := deps.NodeBridge.RegisterAPIRoute(ctxRegister, APIRoute, advertisedAddress); err != nil {
			CoreComponent.LogErrorfAndExit("Registering INX api route failed: %s", err)
		}
		cancelRegister()

		CoreComponent.LogInfo("Starting API server ... done")
		<-ctx.Done()
		CoreComponent.LogInfo("Stopping API ...")

		ctxUnregister, cancelUnregister := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelUnregister()

		//nolint:contextcheck // false positive
		if err := deps.NodeBridge.UnregisterAPIRoute(ctxUnregister, APIRoute); err != nil {
			CoreComponent.LogWarnf("Unregistering INX api route failed: %s", err)
		}

		shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCtxCancel()

		//nolint:contextcheck // false positive
		if err := e.Shutdown(shutdownCtx); err != nil {
			CoreComponent.LogWarn(err)
		}

		CoreComponent.LogInfo("Stopping API ... done")
	}, daemon.PriorityStopRestAPI); err != nil {
		CoreComponent.LogPanicf("failed to start worker: %s", err)
	}

	return nil
}
