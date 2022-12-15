package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/jessemolina/lab-go-service/pkg/web"
	"go.uber.org/zap"
)

// ================================================================
// FUNCTIONS

// A web Middlerware is a function that accepts and returns a handler.
func Logger(log *zap.SugaredLogger) web.Middleware {

	// Create an anonymous middleware function;
	// a function that accepts and returns a web.Handler.
	m := func(handler web.Handler) web.Handler {

		// Create an anonymous function for web.Handler;
		// closure enables the use of parameters that exist outside of the scope of this function.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			traceID := "00000000000000000000"
			statuscode := http.StatusOK
			now := time.Now()

			log.Infow("request started", "traceid", traceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Infow("request completed", "traceid", traceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr, "statuscode", statuscode, "since", time.Since(now))

			return err
		}

		return h
	}

	return m
}