package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type responseBodyWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func WithRequestResponseLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			req := c.Request()
			res := c.Response()

			var reqBody []byte
			if req.Body != nil {
				buf := new(bytes.Buffer)
				_, err := buf.ReadFrom(req.Body)
				if err == nil {
					reqBody = buf.Bytes()
					req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
				}
			}

			log.Info().
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Str("remote_ip", c.RealIP()).
				Str("user_agent", req.UserAgent()).
				Str("request_body", string(reqBody)).
				Msg("Incoming HTTP request")

			rbw := &responseBodyWriter{
				ResponseWriter: res.Writer,
				body:           new(bytes.Buffer),
			}

			res.Writer = rbw

			err := next(c)

			log.Info().
				Int("status", res.Status).
				Str("latency", time.Since(start).String()).
				Str("response_body", rbw.body.String()).
				Msg("Outgoing HTTP response")

			return err
		}
	}
}
