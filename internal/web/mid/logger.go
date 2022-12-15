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

// Create a zap.SugaredLogger middlware function;
// a function that accepts and returns a handler.
// This function enables tracing and logging before/after
// top layered web.Handler.
func Logger(log *zap.SugaredLogger) web.Middleware {

	// Create an anonymous middleware function;
	// a function that accepts and returns a web.Handler.
	m := func(handler web.Handler) web.Handler {

		// Create an anonymous function for web.Handler;
		// closure enables the use of parameters that exist outside of the scope of this function.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context is missing this value, request the service
			// to be shutdown gracefully
			v, err := web.GetValues(ctx)
			if err != nil {
				return err
			}

			log.Infow("request started", "traceid", v.TraceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr)

			err = handler(ctx, w, r)

			log.Infow("request completed", "traceid", v.TraceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr, "statuscode", v.StatusCode, "since", time.Since(v.Now))

			return err
		}

		return h
	}

	return m
}
