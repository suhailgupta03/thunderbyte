package core

import (
	"github.com/labstack/echo/v4"
	"log"
	"strings"
)

type TBAppInterface interface {
	Listen(port int)
}

type TBApp struct {
	logger log.Logger
}

// Listen It starts the server and listens on the specified address
func (tba *TBApp) Listen(port int) *echo.Echo {
	var srv = echo.New()
	srv.HideBanner = true
	// Start the server.
	go func() {
		address := ":" + string(port)
		if err := srv.Start(address); err != nil {
			if strings.Contains(err.Error(), "Server closed") {
				tba.logger.Println("HTTP server shut down")
			} else {
				tba.logger.Fatalf("error starting HTTP server: %v", err)
			}
		} else {
			tba.logger.Println("HTTP server started on", address)
		}
	}()
	return srv
}
