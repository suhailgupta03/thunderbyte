package common

import (
	"github.com/labstack/echo/v4"
	"log"
	"path"
)

type Module struct {
	E                *echo.Echo
	L                *log.Logger
	ControllerConfig *ControllerConfig
	Providers        []interface{}
	Imports          []*Module
}

// InitModule It initializes the module by registering routes
func InitModule(modules []*Module, srv *echo.Echo, logger *log.Logger, basePath *string) {
	for _, module := range modules {
		if module != nil {
			if module.ControllerConfig.ModulePath == "" {
				logger.Fatalf("ModulePath is missing for one of the controller configs in imports")
			}
			if module.E == nil {
				module.E = srv
			}
			if module.L == nil {
				module.L = logger
			}
			if module.ControllerConfig == nil {
				logger.Fatalf("ControllerConfig is required")
			}
			if basePath != nil {
				module.ControllerConfig.ModulePath = RoutePath(path.Join(*basePath, string(module.ControllerConfig.ModulePath)))
			}
			cd := controllerDetails{
				l: module.L,
				e: module.E,
				c: module.ControllerConfig,
			}
			logger.Printf("Initializing module %s", module.ControllerConfig.ModulePath)
			cd.registerRoutes()
			if len(module.Imports) > 0 {
				// Recursively initialize the imports
				// Does the nesting of routes
				newBasePath := string(module.ControllerConfig.ModulePath)
				InitModule(module.Imports, srv, logger, &newBasePath)
			}
		}
	}
}
