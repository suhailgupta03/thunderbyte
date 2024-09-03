package core

import (
	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type O11Y struct {
	NewRelic *NewRelicConfig
}

type NewRelicConfig struct {
	AppName    string
	Enabled    bool
	LicenseKey string
	UserKey    string
}

func newRelicMiddleware(app *newrelic.Application) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			txnName := c.Request().Method + " " + c.Path()
			txn := app.StartTransaction(txnName)
			defer txn.End()

			// Add transaction to context
			c.Set("newRelicTransaction", txn)

			// This line calls SetWebResponse on the New Relic transaction object (txn),
			// passing in the current HTTP response writer from the Echo context (c.Response().Writer).
			// The SetWebResponse method returns a wrapped version of the response writer that is
			// instrumented to monitor and record details about the HTTP response, such as the
			// status code and headers. The returned writer (w) is capable of capturing this
			// information without altering the behavior of the original response writer.
			w := txn.SetWebResponse(c.Response().Writer)
			c.Response().Writer = w
			// This line replaces the original response writer in the Echo context with the
			// New Relic-instrumented response writer (w). From this point forward, any
			// write operations to the response will be intercepted by New Relic's wrapped writer,
			// allowing it to observe and record the response details. This includes the status code
			// set during the request handling and any headers or body content written to the response

			// Proceed with request
			err := next(c)

			// Record any errors
			if err != nil {
				txn.NoticeError(err)
				c.Error(err)
			}

			return nil
		}
	}
}
