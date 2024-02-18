package core

import (
	"github.com/labstack/echo/v4"
	"github.com/suhailgupta03/thunderbyte/common"
	"io"
	"log"
	"os"
)

type TBFactory struct {
}

type FactoryCreate struct {
	ControllerConfig []*common.ControllerConfig
	Providers        []interface{}
	Imports          []*common.Module
}

type TBFactoryInterface interface {
	Create(fc *FactoryCreate) *TBApp
}

var (
	srv    = echo.New()
	logger = log.New(io.MultiWriter(os.Stdout), "", log.Ldate|log.Ltime|log.Lshortfile)
)

// Create It returns a pointer to a new TBApp
func (tbf *TBFactory) Create(fc *FactoryCreate) *TBApp {
	if len(fc.ControllerConfig) == 0 {
		logger.Fatalf("ControllerConfig is required")
	}

	for _, cc := range fc.ControllerConfig {
		if cc.Controllers != nil {
			module := common.Module{
				E:                srv,
				L:                logger,
				ControllerConfig: cc,
				Providers:        fc.Providers,
			}
			common.InitModule([]*common.Module{&module}, srv, logger, nil)
		}
	}
	common.InitModule(fc.Imports, srv, logger, nil)
	return &TBApp{
		Logger: logger,
	}
}
