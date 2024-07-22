package api

import (
	"urfunavigator/index/object"
	"urfunavigator/index/store"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type API struct {
	Port           string
	Store          store.Store
	ObjectStore    object.ObjectStore
	AllowedOrigins string
}

func NewAPI(port string, store store.Store, objectStore object.ObjectStore, allowedOrigins string) *API {
	return &API{
		Port:           port,
		Store:          store,
		ObjectStore:    objectStore,
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
	app.Get("/icons/:icon", s.ObjectHandler)

	return app.Listen(s.Port)
}
