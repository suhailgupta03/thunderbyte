package common

import (
	"fmt"
	"github.com/knadh/koanf/v2"
	"github.com/labstack/echo/v4"
	"github.com/suhailgupta03/smtppool"
	"github.com/suhailgupta03/thunderbyte/database"
	"github.com/suhailgupta03/thunderbyte/otp/store/redis"
	"github.com/zerodha/logf"
	"path"
	"reflect"
)

type Module struct {
	E                *echo.Echo
	L                *logf.Logger
	ControllerConfig *ControllerConfig
	Providers        []interface{}
	Imports          []*Module
}

type InitModuleParams struct {
	Srv      *echo.Echo
	DBConfig *database.DBConfig
	Redis    *redis.Redis
	K        *koanf.Koanf
	SMTPPool *smtppool.Pool
	Logger   *logf.Logger
}

// InitModule It initializes the module by registering routes
func InitModule(modules []*Module, moduleParams *InitModuleParams, basePath *string) {
	logger := moduleParams.Logger
	srv := moduleParams.Srv

	for _, module := range modules {
		if module != nil {
			fmt.Println(module.ControllerConfig.ModulePath)
			if module.ControllerConfig.ModulePath == "" {
				logger.Fatal("ModulePath is missing for one of the controller configs in imports")
			}
			if module.E == nil {
				module.E = srv
			}
			if module.L == nil {
				module.L = logger
			}
			if module.ControllerConfig == nil {
				logger.Fatal("ControllerConfig is required")
			}
			if basePath != nil {
				module.ControllerConfig.ModulePath = RoutePath(path.Join(*basePath, string(module.ControllerConfig.ModulePath)))
			}
			// Create a map of services to be injected
			// into the controller
			serviceMap := make(InjectedServicesMap)
			for _, p := range module.Providers {
				ssType := reflect.TypeOf(p)
				serviceMap[ServiceName(ssType.Name())] = p
			}
			cd := controllerDetails{
				l:                   module.L,
				e:                   module.E,
				c:                   module.ControllerConfig,
				injectedServicesMap: &serviceMap,
				dbConfig:            moduleParams.DBConfig,
				redis:               moduleParams.Redis,
				smtpPool:            moduleParams.SMTPPool,
				k:                   moduleParams.K,
			}
			logger.Info("Initializing module", "path", module.ControllerConfig.ModulePath)
			cd.registerRoutes()
			if len(module.Imports) > 0 {
				// Recursively initialize the imports
				// Does the nesting of routes
				newBasePath := string(module.ControllerConfig.ModulePath)
				InitModule(module.Imports, moduleParams, &newBasePath)
			}
		}
	}
}
