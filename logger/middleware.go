package logger

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

func RequestMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		status := c.Response().StatusCode()
		attrs := []any{
			"method", c.Method(),
			"path", c.Path(),
			"route", c.Route().Path,
			"status", status,
			"duration_ms", time.Since(start).Milliseconds(),
			"ip", c.IP(),
		}

		if err != nil {
			attrs = append(attrs, "err", err.Error())
		}

		switch {
		case status >= fiber.StatusInternalServerError:
			Error("http request completed", attrs...)
		case status >= fiber.StatusBadRequest:
			Warn("http request completed", attrs...)
		default:
			Info("http request completed", attrs...)
		}

		return err
	}
}
