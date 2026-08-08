package handlers

import (
	"urfunavigator/index/logger"

	"github.com/gofiber/fiber/v3"
)

func requestContext(c fiber.Ctx) []any {
	return []any{
		"method", c.Method(),
		"path", c.Path(),
		"route", c.Route().Path,
		"ip", c.IP(),
	}
}

func logHandlerError(c fiber.Ctx, handler string, err error, attrs ...any) {
	args := append([]any{"handler", handler, "err", err.Error()}, requestContext(c)...)
	args = append(args, attrs...)
	logger.Error("handler failed", args...)
}

func logHandlerWarn(c fiber.Ctx, handler string, msg string, attrs ...any) {
	args := append([]any{"handler", handler}, requestContext(c)...)
	args = append(args, attrs...)
	logger.Warn(msg, args...)
}

func logHandlerDebug(c fiber.Ctx, handler string, msg string, attrs ...any) {
	args := append([]any{"handler", handler}, requestContext(c)...)
	args = append(args, attrs...)
	logger.Debug(msg, args...)
}

func logHandlerInfo(c fiber.Ctx, handler string, msg string, attrs ...any) {
	args := append([]any{"handler", handler}, requestContext(c)...)
	args = append(args, attrs...)
	logger.Info(msg, args...)
}
