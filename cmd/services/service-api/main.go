package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/jessemolina/lab-go-service/cmd/services/service-api/handlers"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var build = "develop"

func main() {

	// ================================================================
	// LOGGER
	// create a logger to be used across the service

	// create initial logger and defer
	log, err := initLogger("SERVICE-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer log.Sync()

	// ================================================================
	// RUN
	// run the service

	// run service with given logger
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}

}

// run service start up
func run(log *zap.SugaredLogger) error {

	// ================================================================
	// QUOTAS
	// set resource limits

	// set the cpu quota
	if _, err := maxprocs.Set(); err != nil {
		fmt.Errorf("maxprocs: %w", err)
		os.Exit(1)
	}

	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// ================================================================
	// CONFIGURATION
	// set service configurations

	// struct to create dynamic for flags and environment variables
	cfg := struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
		}
	}{
		Version: conf.Version{
			SVN:  build,
			Desc: "copyrights",
		},
	}

	// parse os args to overide cfg
	const prefix = "SERVICE"
	help, err := conf.ParseOSArgs(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// ================================================================
	// STARTUP
	// initiate service

	// log start up and shutdown
	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	// log values used for cfg
	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	// expvar config
	expvar.NewString("build").Set(build)

	// ================================================================
	// DEBUG API
	// enable debug endpoints

	// log debug startup
	log.Infow("startup", "status", "debug router started", cfg.Web.DebugHost)

	// construct mux to serve debug calls
	debugMux := handlers.DebugMux(build, log)

	// start service listening for debug requests
	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debugMux); err != nil {
			log.Errorw("shutdown", "status", "debug router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	// ================================================================
	// SERVICE API
	// enable service endpoints

	// log service startup
	log.Infow("startup", "status", "initializing service API")

	// make a channel with 1 buffer for an os.Signal
	// block on the channel until it receives either SIGINT or SIGTERM
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	//construct the mux for api calls
	// value that implements the http handler
	apiMux := handlers.APIMux(handlers.APIMuxConfig{
		Shutdown: shutdown,
		Log:      log,
	})

	// server struct to be used against the mux
	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	// make a channel for errors coming from the listener
	// use a buffered chanel for the goroutine to exit
	serverErrors := make(chan error, 1)

	// start service listening for api requests
	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// ================================================================
	// SHUTDOWN
	// perform load shedding to shutdown service

	// block main and wait for shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		// log shutdown sequence
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// set deadline for outstanding requests
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// ask listener to shutdown and shed load
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

// build initial zap sugared logger
func initLogger(service string) (*zap.SugaredLogger, error) {

	// create new zap config
	config := zap.NewProductionConfig()

	// overwrite defaults
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": service,
	}

	// build logger
	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	// return sugar formatted logger
	return log.Sugar(), nil
}
