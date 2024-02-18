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
	Controllers common.Controller
	Providers   []interface{}
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
	module := common.Module{
		E:           srv,
		L:           logger,
		Controllers: fc.Controllers,
		Providers:   fc.Providers,
	}
	module.InitModule()
	return &TBApp{
		Logger: logger,
	}
}
