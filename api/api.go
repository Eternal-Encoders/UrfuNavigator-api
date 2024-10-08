package api

import (
	"urfunavigator/index/geo"
	"urfunavigator/index/object"
	"urfunavigator/index/store"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type API struct {
	Port           string
	Store          store.Store
	ObjectStore    object.ObjectStore
	GEoService     geo.Geo
	AllowedOrigins string
}

func NewAPI(port string, store store.Store, objectStore object.ObjectStore, geo geo.Geo, allowedOrigins string) *API {
	return &API{
		Port:           port,
		Store:          store,
		ObjectStore:    objectStore,
		GEoService:     geo,
		AllowedOrigins: allowedOrigins,
	}
}

func (s *API) Run() error {
	app := fiber.New()

	cors := cors.New(cors.Config{
		AllowOrigins: s.AllowedOrigins,
	})

	app.Use(cors)

	app.Get("/", s.MainHandler)
	app.Get("/floor", s.FloorHandler)
	app.Get("/institute", s.InstituteHandler)
	app.Get("/institutes", s.InstitutesHandler)
	app.Get("/points", s.PointsHandler)
	app.Get("/point", s.PointIdHandler)
	app.Get("/search", s.SearchHandler)
	app.Get("/path", s.PathHandler)
	app.Get("/icons/:icon", s.ObjectHandler)

	return app.Listen(s.Port)
}
