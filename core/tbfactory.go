package core

import (
	"github.com/knadh/koanf/v2"
	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/suhailgupta03/smtppool"
	"github.com/suhailgupta03/thunderbyte/common"
	"github.com/suhailgupta03/thunderbyte/database"
	"github.com/suhailgupta03/thunderbyte/otp/store/redis"
	"github.com/zerodha/logf"
	"time"
)

type TBFactory struct {
}

type FactoryCreate struct {
	DBConfig         *database.DBConfig
	K                *koanf.Koanf
	SMTPPool         *smtppool.Pool
	Redis            *redis.Redis
	ControllerConfig []*common.ControllerConfig
	Providers        []interface{}
	Imports          []*common.Module
	O11Y             *O11Y
}

type TBFactoryInterface interface {
	Create(fc *FactoryCreate) *TBApp
}

var (
	srv    = echo.New()
	logger = logf.New(logf.Opts{
		EnableColor:          true,
		Level:                logf.DebugLevel,
		CallerSkipFrameCount: 3,
		EnableCaller:         true,
		TimestampFormat:      time.RFC3339Nano,
		DefaultFields:        []any{"scope", "example"},
	})
)

// Create It returns a pointer to a new TBApp
func (tbf *TBFactory) Create(fc *FactoryCreate) *TBApp {
	if len(fc.ControllerConfig) == 0 {
		logger.Fatal("ControllerConfig is required")
	}

	if fc.DBConfig != nil {
		database.ForRoot(fc.DBConfig, &logger)
	}

	if fc.O11Y != nil {
		if fc.O11Y.NewRelic != nil && fc.O11Y.NewRelic.Enabled {
			app, err := newrelic.NewApplication(
				newrelic.ConfigAppName(fc.O11Y.NewRelic.AppName),
				newrelic.ConfigLicense(fc.O11Y.NewRelic.LicenseKey),
				newrelic.ConfigAppLogForwardingEnabled(true),
			)
			if err != nil {
				logger.Error("Failed to initialize NewRelic", "error", err)
			}
			srv.Use(newRelicMiddleware(app))
		}
	}
	for _, cc := range fc.ControllerConfig {
		if cc.Controllers != nil {
			module := common.Module{
				E:                srv,
				L:                &logger,
				ControllerConfig: cc,
				Providers:        fc.Providers,
			}
			common.InitModule([]*common.Module{&module}, &common.InitModuleParams{
				Logger:   &logger,
				Srv:      srv,
				DBConfig: fc.DBConfig,
				Redis:    fc.Redis,
				SMTPPool: fc.SMTPPool,
				K:        fc.K,
			}, nil)
		}
	}

	common.InitModule(fc.Imports, &common.InitModuleParams{
		Logger:   &logger,
		Srv:      srv,
		DBConfig: fc.DBConfig,
		Redis:    fc.Redis,
		SMTPPool: fc.SMTPPool,
		K:        fc.K,
	}, nil)

	srv.Validator = common.NewRequestValidator()

	return &TBApp{
		Logger:         &logger,
		DB:             fc.DBConfig.GetDB(),
		DefaultQueries: fc.DBConfig.GetDefaultQueries(),
		DBConfig:       fc.DBConfig,
	}
}
