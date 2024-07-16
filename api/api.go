package api

import (
	"urfunavigator/index/store"

	"github.com/gofiber/fiber/v2"
)

type API struct {
	Port  string
	Store store.Store
}

func NewAPI(port string, store store.Store) *API {
	return &API{
		Port:  port,
		Store: store,
	}
}

func (s *API) Run() error {
	app := fiber.New()

	app.Get("/", s.MainHandler)
	app.Get("/floor", s.FloorHandler)
	app.Get("/institute", s.InstituteHandler)
	app.Get("/institutes", s.InstitutesHandler)
	// app.Get("/graph", s.GraphHandler)
	// app.Get("/stairs", s.StairsHandler)
	app.Get("/points", s.PointsHandler)
	app.Get("/point", s.PointIdHandler)
	app.Get("/path", s.PathHandler)

	return app.Listen(s.Port)
}
