package httpserver

import (
	"github.com/labstack/echo/v4"
)

type HttpServer struct {
	httpServerPort string
	Echo           *echo.Echo
}

func NewHttpServer(httpServerPort string) *HttpServer {
	echoInstanse := echo.New()
	return &HttpServer{
		httpServerPort: httpServerPort,
		Echo:           echoInstanse,
	}
}

func (s *HttpServer) Run() error {
	return s.Echo.Start(s.httpServerPort)
}

func (s *HttpServer) Shutdown() error {
	return s.Echo.Close()
}
