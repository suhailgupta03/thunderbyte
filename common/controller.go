package common

import (
	"github.com/knadh/koanf/v2"
	"github.com/labstack/echo/v4"
	"github.com/suhailgupta03/thunderbyte/database"
	"github.com/suhailgupta03/thunderbyte/otp/store/redis"
	"github.com/zerodha/logf"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type RoutePath string
type HTTPMethod int

const (
	GET HTTPMethod = iota
	POST
	PUT
	DELETE
	PATCH
	OPTIONS
	HEAD
	TRACE
	CONNECT
)

type QueryParams map[string][]string
type Headers map[string][]string
type AppContext struct {
	RequestContext    RequestContext
	HTTPServerContext echo.Context
	DBConfig          *database.DBConfig
	Redis             *redis.Redis
	Logger            *logf.Logger
	K                 *koanf.Koanf
	Q                 interface{}
}
type RequestContext struct {
	Body        interface{}
	QueryParams QueryParams
	Path        string
	Headers     Headers
}
type ServiceName string
type InjectedServicesMap map[ServiceName]interface{}
type HTTPMethodHandler func(context AppContext, injectedServicesMap *InjectedServicesMap) (interface{}, *HTTPError)
type HandlerWithMiddlewares struct {
	Handler    HTTPMethodHandler
	Middleware []echo.MiddlewareFunc
}
type HTTPMethods map[HTTPMethod]HandlerWithMiddlewares
type Controller map[RoutePath]HTTPMethods

type ControllerConfig struct {
	ModulePath  RoutePath
	Controllers Controller
}

type controllerDetails struct {
	l                   *logf.Logger
	e                   *echo.Echo
	c                   *ControllerConfig
	injectedServicesMap *InjectedServicesMap
	dbConfig            *database.DBConfig
	redis               *redis.Redis
	k                   *koanf.Koanf
}

// okResp It is a response struct for successful requests
type okResp struct {
	Data interface{} `json:"data"`
}

// errorResp It is a response struct for failed requests
type errorResp struct {
	Error      string `json:"error"`
	Code       int    `json:"code"`
	StatusText string `json:"statusText"`
}

type HTTPError struct {
	// Code It is the HTTP status code used by standard http package in Go
	Code int
	// Message It is the error message set by the caller
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}

// extractRequestContext It extracts the request context from the echo context
// and returns it as a RequestContext
func extractRequestContext(c echo.Context) RequestContext {
	paramsMap := make(QueryParams)
	for key, value := range c.QueryParams() {
		paramsMap[key] = value
	}
	headersMap := make(Headers)
	for key, value := range c.Request().Header {
		headersMap[key] = value
	}
	return RequestContext{
		Body:        c.Request().Body,
		QueryParams: paramsMap,
		Path:        c.Path(),
		Headers:     headersMap,
	}
}

func (cd *controllerDetails) handleIncomingRequest(c echo.Context, handler HTTPMethodHandler) error {
	requestStart := time.Now().UnixNano()
	fn := reflect.ValueOf(handler)
	appContext := AppContext{
		RequestContext:    extractRequestContext(c),
		HTTPServerContext: c,
		DBConfig:          cd.dbConfig,
		Redis:             cd.redis,
		K:                 cd.k,
		Logger:            cd.l,
	}
	args := reflect.ValueOf(appContext)
	serviceMap := reflect.ValueOf(cd.injectedServicesMap)
	response := fn.Call([]reflect.Value{args, serviceMap})
	data, e := response[0].Interface(), response[1].Interface()
	if controllerError, ok := e.(*HTTPError); ok && controllerError != nil {
		statusText := http.StatusText(controllerError.Code)
		statusCode := controllerError.Code
		if statusText == "" {
			statusText = http.StatusText(http.StatusInternalServerError)
			statusCode = http.StatusInternalServerError
		}
		cd.l.Error("Request failed", "path", c.Path(), "method", c.Request().Method, "code", statusCode)
		cd.l.Info("Request time", "ms", (time.Now().UnixNano()-requestStart)/1000000)
		return c.JSON(statusCode, errorResp{
			Error:      controllerError.Message,
			Code:       statusCode,
			StatusText: statusText,
		})
	}
	cd.l.Info("Request success", "path", c.Path(), "method", c.Request().Method)
	cd.l.Info("Request time", "ms", (time.Now().UnixNano()-requestStart)/1000000)
	return c.JSON(200, okResp{Data: data})
}

func (cd *controllerDetails) initIncomingRequestHandler(handler HTTPMethodHandler) func(echo.Context) error {
	return func(c echo.Context) error {
		return cd.handleIncomingRequest(c, handler)
	}
}

func (cd *controllerDetails) registerRoutes() {
	for path, methods := range cd.c.Controllers {
		modulePath := string(cd.c.ModulePath)
		pathToRegister := string(path)
		if !strings.HasPrefix(pathToRegister, "/") {
			pathToRegister = "/" + pathToRegister
		}
		if strings.TrimSpace(modulePath) != "" {
			modulePath = strings.TrimSuffix(modulePath, "/")
			pathToRegister = strings.TrimPrefix(pathToRegister, "/")
			pathToRegister = modulePath + "/" + pathToRegister
		} else {
			cd.l.Warn("A controller has no module path. It is recommended to always have a module path", "path", pathToRegister)
		}
		for method, handler := range methods {
			switch method {
			case GET:
				cd.e.GET(pathToRegister, cd.initIncomingRequestHandler(handler.Handler), handler.Middleware...)
				cd.l.Info("Registered", "GET", pathToRegister)
				break
			case POST:
				cd.e.POST(pathToRegister, cd.initIncomingRequestHandler(handler.Handler), handler.Middleware...)
				cd.l.Info("Registered", "POST", pathToRegister)
				break
			case PUT:
				cd.e.PUT(pathToRegister, cd.initIncomingRequestHandler(handler.Handler), handler.Middleware...)
				cd.l.Info("Registered", "PUT", pathToRegister)
				break
			case DELETE:
				cd.e.DELETE(pathToRegister, cd.initIncomingRequestHandler(handler.Handler), handler.Middleware...)
				cd.l.Info("Registered", "DELETE", pathToRegister)
				break
			case PATCH:
				cd.e.PATCH(pathToRegister, cd.initIncomingRequestHandler(handler.Handler), handler.Middleware...)
				cd.l.Info("Registered", "PATCH", pathToRegister)
				break
			case OPTIONS:
				cd.e.OPTIONS(pathToRegister, cd.initIncomingRequestHandler(handler.Handler), handler.Middleware...)
				cd.l.Info("Registered", "OPTIONS", pathToRegister)
				break
			case HEAD:
				cd.e.HEAD(pathToRegister, cd.initIncomingRequestHandler(handler.Handler), handler.Middleware...)
				cd.l.Info("Registered", "HEAD", pathToRegister)
				break
			case TRACE:
				cd.e.TRACE(pathToRegister, cd.initIncomingRequestHandler(handler.Handler), handler.Middleware...)
				cd.l.Info("Registered", "TRACE", pathToRegister)
				break
			case CONNECT:
				cd.e.CONNECT(pathToRegister, cd.initIncomingRequestHandler(handler.Handler), handler.Middleware...)
				cd.l.Info("Registered", "CONNECT", pathToRegister)
				break
			}
		}
	}
}
