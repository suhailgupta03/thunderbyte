package core

import (
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/suhailgupta03/thunderbyte/database"
	"github.com/zerodha/logf"
	"strconv"
	"strings"
)

type TBAppInterface interface {
	Listen(port int)
}

type TBApp struct {
	Logger         *logf.Logger
	DB             *sqlx.DB
	DefaultQueries database.ThunderbyteQueries
}

// Listen It starts the server and listens on the specified address
func (tba *TBApp) Listen(port int) *echo.Echo {
	srv.HideBanner = true
	// Initialize the request validator
	srv.Validator = &RequestValidator{validator: validator.New()}
	// Start the server.
	go func() {
		address := ":" + strconv.Itoa(port)
		if err := srv.Start(address); err != nil {
			if strings.Contains(err.Error(), "Server closed") {
				tba.Logger.Info("HTTP server shut down")
			} else {
				tba.Logger.Fatal("Error starting HTTP server", "error", err)
			}
		} else {
			tba.Logger.Info("HTTP server started", "port", address)
		}
	}()
	return srv
}
