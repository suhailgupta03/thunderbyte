package core

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/suhailgupta03/thunderbyte/database"
	"log"
	"strconv"
	"strings"
)

type TBAppInterface interface {
	Listen(port int)
}

type TBApp struct {
	Logger         *log.Logger
	DB             *sqlx.DB
	DefaultQueries database.ThunderbyteQueries
}

// Listen It starts the server and listens on the specified address
func (tba *TBApp) Listen(port int) *echo.Echo {
	srv.HideBanner = true
	// Start the server.
	go func() {
		address := ":" + strconv.Itoa(port)
		if err := srv.Start(address); err != nil {
			if strings.Contains(err.Error(), "Server closed") {
				tba.Logger.Println("HTTP server shut down")
			} else {
				tba.Logger.Fatalf("Error starting HTTP server: %v", err)
			}
		} else {
			tba.Logger.Println("HTTP server started on", address)
		}
	}()
	return srv
}
