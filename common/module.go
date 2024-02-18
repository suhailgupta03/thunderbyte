package common

import (
	"github.com/labstack/echo/v4"
	"log"
)

type Module struct {
	E                *echo.Echo
	L                *log.Logger
	ControllerConfig *ControllerConfig
	Providers        []interface{}
}

// InitModule It initializes the module by registering routes
func (m *Module) InitModule() {
	cd := controllerDetails{
		l: m.L,
		e: m.E,
		c: m.ControllerConfig,
	}
	cd.registerRoutes()
}
