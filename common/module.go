package common

import (
	"github.com/labstack/echo/v4"
	"log"
)

type Module struct {
	E           *echo.Echo
	L           *log.Logger
	Controllers Controller
	Providers   []interface{}
}

// InitModule It initializes the module by registering routes
func (m *Module) InitModule() {
	cd := controllerDetails{
		l: m.L,
		e: m.E,
		c: m.Controllers,
	}
	cd.registerRoutes()
}
