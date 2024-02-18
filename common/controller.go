package common

import (
	"github.com/labstack/echo/v4"
	"log"
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
type RequestContext struct {
	Body        interface{}
	QueryParams QueryParams
	Path        string
	Headers     Headers
}
type ServiceName string
type InjectedServicesMap map[ServiceName]interface{}
type HTTPMethodHandler func(context RequestContext, injectedServicesMap *InjectedServicesMap) (interface{}, *HTTPError)
type HTTPMethods map[HTTPMethod]HTTPMethodHandler
type Controller map[RoutePath]HTTPMethods

type ControllerConfig struct {
	ModulePath  RoutePath
	Controllers Controller
}

type controllerDetails struct {
	l                   *log.Logger
	e                   *echo.Echo
	c                   *ControllerConfig
	injectedServicesMap *InjectedServicesMap
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
	args := reflect.ValueOf(extractRequestContext(c))
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
		cd.l.Printf("Request for %s and %s failed with status code %d", c.Path(), c.Request().Method, statusCode)
		cd.l.Printf("Request took %d ms", (time.Now().UnixNano()-requestStart)/1000000)
		return c.JSON(statusCode, errorResp{
			Error:      controllerError.Message,
			Code:       statusCode,
			StatusText: statusText,
		})
	}
	cd.l.Printf("Request for %s and %s succeeded", c.Path(), c.Request().Method)
	cd.l.Printf("Request took %d ms", (time.Now().UnixNano()-requestStart)/1000000)
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
			cd.l.Printf("Controller with path '%s' has no module path. It is recommended to always have a module path.\n", pathToRegister)
		}
		for method, handler := range methods {
			switch method {
			case GET:
				cd.e.GET(pathToRegister, cd.initIncomingRequestHandler(handler))
				cd.l.Println("GET", pathToRegister)
				break
			case POST:
				cd.e.POST(pathToRegister, cd.initIncomingRequestHandler(handler))
				cd.l.Println("POST", pathToRegister)
				break
			case PUT:
				cd.e.PUT(pathToRegister, cd.initIncomingRequestHandler(handler))
				cd.l.Println("PUT", pathToRegister)
				break
			case DELETE:
				cd.e.DELETE(pathToRegister, cd.initIncomingRequestHandler(handler))
				cd.l.Println("DELETE", pathToRegister)
				break
			case PATCH:
				cd.e.PATCH(pathToRegister, cd.initIncomingRequestHandler(handler))
				cd.l.Println("PATCH", pathToRegister)
				break
			case OPTIONS:
				cd.e.OPTIONS(pathToRegister, cd.initIncomingRequestHandler(handler))
				cd.l.Println("OPTIONS", pathToRegister)
				break
			case HEAD:
				cd.e.HEAD(pathToRegister, cd.initIncomingRequestHandler(handler))
				cd.l.Println("HEAD", pathToRegister)
				break
			case TRACE:
				cd.e.TRACE(pathToRegister, cd.initIncomingRequestHandler(handler))
				cd.l.Println("TRACE", pathToRegister)
				break
			case CONNECT:
				cd.e.CONNECT(pathToRegister, cd.initIncomingRequestHandler(handler))
				cd.l.Println("CONNECT", pathToRegister)
				break
			}
		}
	}
}
