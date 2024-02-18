package common

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"reflect"
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

type HTTPMethodHandler func(context RequestContext) (interface{}, *HTTPError)
type HTTPMethods map[HTTPMethod]HTTPMethodHandler
type Controller map[RoutePath]HTTPMethods

type controllerDetails struct {
	l *log.Logger
	e *echo.Echo
	c Controller
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
	response := fn.Call([]reflect.Value{args})
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
	for path, methods := range cd.c {
		for method, handler := range methods {
			switch method {
			case GET:
				cd.e.GET(string(path), cd.initIncomingRequestHandler(handler))
				cd.l.Println("GET", string(path))
				break
			case POST:
				cd.e.POST(string(path), cd.initIncomingRequestHandler(handler))
				cd.l.Println("POST", string(path))
				break
			case PUT:
				cd.e.PUT(string(path), cd.initIncomingRequestHandler(handler))
				cd.l.Println("PUT", string(path))
				break
			case DELETE:
				cd.e.DELETE(string(path), cd.initIncomingRequestHandler(handler))
				cd.l.Println("DELETE", string(path))
				break
			case PATCH:
				cd.e.PATCH(string(path), cd.initIncomingRequestHandler(handler))
				cd.l.Println("PATCH", string(path))
				break
			case OPTIONS:
				cd.e.OPTIONS(string(path), cd.initIncomingRequestHandler(handler))
				cd.l.Println("OPTIONS", string(path))
				break
			case HEAD:
				cd.e.HEAD(string(path), cd.initIncomingRequestHandler(handler))
				cd.l.Println("HEAD", string(path))
				break
			case TRACE:
				cd.e.TRACE(string(path), cd.initIncomingRequestHandler(handler))
				cd.l.Println("TRACE", string(path))
				break
			case CONNECT:
				cd.e.CONNECT(string(path), cd.initIncomingRequestHandler(handler))
				cd.l.Println("CONNECT", string(path))
				break
			}
		}
	}
}
