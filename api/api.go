package api

import (
	"urfunavigator/index/geo"
	"urfunavigator/index/logs"
	"urfunavigator/index/object"
	"urfunavigator/index/store"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type API struct {
	DefaultPath    string
	Port           string
	Store          store.Store
	ObjectStore    object.ObjectStore
	GEoService     geo.Geo
	AllowedOrigins string
	Logger         logs.Log
}

func NewAPI(
	defaultPath string,
	port string,
	store store.Store,
	objectStore object.ObjectStore,
	geo geo.Geo,
	allowedOrigins string,
	logger logs.Log,
) *API {
	return &API{
		DefaultPath:    defaultPath,
		Port:           port,
		Store:          store,
		ObjectStore:    objectStore,
		GEoService:     geo,
		AllowedOrigins: allowedOrigins,
		Logger:         logger,
	}
}

func (s *API) Run() error {
	app := fiber.New()

	cors := cors.New(cors.Config{
		AllowOrigins: s.AllowedOrigins,
	})

	app.Use(cors)

	app.Get(s.DefaultPath+"/", s.MainHandler)
	app.Get(s.DefaultPath+"/floor", s.FloorHandler)
	app.Get(s.DefaultPath+"/institute", s.InstituteHandler)
	app.Get(s.DefaultPath+"/institutes", s.InstitutesHandler)
	app.Get(s.DefaultPath+"/points", s.PointsHandler)
	app.Get(s.DefaultPath+"/point", s.PointIdHandler)
	app.Get(s.DefaultPath+"/search", s.SearchHandler)
	app.Get(s.DefaultPath+"/path", s.PathHandler)
	app.Get(s.DefaultPath+"/icons/:icon", s.ObjectHandler)

	return app.Listen(s.Port)
}
