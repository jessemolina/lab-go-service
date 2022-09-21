package main

import (
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
	"github.com/jessemolina/ultimate-service/app/services/service-api/handlers"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var build = "develop"

func main() {

	// ================================================================
	// LOGGER

	// create initial logger and defer
	log, err := initLogger("SERVICE-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer log.Sync()

	// ================================================================
	// SHUTDOWN

	// run service with given logger
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}


}

// run service start up
func run(log *zap.SugaredLogger) error {

	// ================================================================
	// GOMAXPROCS

	// set the cpu quota
	if _, err := maxprocs.Set(); err != nil {
		fmt.Errorf("maxprocs: %w", err)
		os.Exit(1)
	}

	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// ================================================================
	// CONFIGURATION

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
	// DEBUG
	//

	// log debug startup
	log.Infow("startup", "status", "debug router started", cfg.Web.DebugHost)

	// TODO construct mux to serve debug calls
	debugMux := handlers.DebugStandardLibraryMux()

	// TODO start service listening for debug requests
	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debugMux); err != nil {
			log.Errorw("shutdown", "status", "debug router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	// ================================================================
	// SHUTDOWN

	// make a channel with 1 buffer for an os.Signal
	// block on the channel until it receives either SIGINT or SIGTERM
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

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
