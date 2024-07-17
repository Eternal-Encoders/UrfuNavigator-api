package api

import (
	"urfunavigator/index/store"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type API struct {
	Port           string
	Store          store.Store
	AllowedOrigins string
}

func NewAPI(port string, store store.Store, allowedOrigins string) *API {
	return &API{
		Port:           port,
		Store:          store,
		AllowedOrigins: allowedOrigins,
	}
}

func (s *API) Run() error {
	app := fiber.New()

	cors := cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, https://how-to-navigate.ru, https://how-to-navigate.ru:2053/",
	})

	app.Use(cors)

	app.Get("/", s.MainHandler)
	app.Get("/floor", s.FloorHandler)
	app.Get("/institute", s.InstituteHandler)
	app.Get("/institutes", s.InstitutesHandler)
	app.Get("/points", s.PointsHandler)
	app.Get("/point", s.PointIdHandler)
	app.Get("/path", s.PathHandler)

	return app.Listen(s.Port)
}
