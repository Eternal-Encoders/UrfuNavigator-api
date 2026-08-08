package api

import (
	"urfunavigator/index/api/handlers"
	"urfunavigator/index/auth"
	"urfunavigator/index/logger"
	"urfunavigator/index/models"

	"github.com/gofiber/contrib/v3/swaggerui"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

type API struct {
	DefaultPath    string
	Port           string
	AllowedOrigins []string
}

func NewAPI(
	defaultPath string,
	port string,
	allowedOrigins []string,
) *API {
	return &API{
		DefaultPath:    defaultPath,
		Port:           port,
		AllowedOrigins: allowedOrigins,
	}
}

// MainHandler health check.
// @Summary Health check
// @Produce plain
// @Success 200 {string} string "OK"
// @Router / [get]
func MainHandler(c fiber.Ctx) error {
	return c.SendString("OK")
}

func (s *API) Run(dataService models.DataService, jwtCfg auth.Config) error {
	app := fiber.New()

	cors := cors.New(cors.Config{
		AllowOrigins: s.AllowedOrigins,
	})

	swaggerCfg := swaggerui.Config{
		FilePath: "./docs/swagger.json",
		Title:    "UrfuNavigator Api Doc",
		CacheAge: 0,
	}

	app.Use(cors)
	app.Use(logger.RequestMiddleware())
	app.Use(swaggerui.New(swaggerCfg))

	app.Get(s.DefaultPath+"/", MainHandler)

	// Api For Frontend
	commonGroup := app.Group("/api")

	commonGroup.Get("/floor", handlers.FloorHandler(dataService))
	commonGroup.Get("/building", handlers.BuildingHandler(dataService))
	commonGroup.Get("/buildings", handlers.BuildingsHandler(dataService))
	commonGroup.Get("/points", handlers.PointsHandler(dataService))
	commonGroup.Get("/point", handlers.PointIdHandler(dataService))
	commonGroup.Get("/search", handlers.SearchHandler(dataService))
	commonGroup.Get("/path", handlers.PathHandler(dataService))
	commonGroup.Get("/icons/:icon", handlers.ObjectHandler(dataService))
	commonGroup.Post("/login", handlers.LoginHandler(dataService, jwtCfg))

	// Admin API
	admin := app.Group("/admin_api", auth.RequireJWT(jwtCfg), auth.RequireAdminRole())

	// Icons
	admin.Delete("/delete_icon", handlers.DeleteObjectHandler(dataService))
	admin.Post("/upload_icon", handlers.UploadObjectHandler(dataService))

	// Users
	admin.Get("/users", handlers.AdminListUsersHandler(dataService))
	admin.Get("/users/:id", handlers.AdminGetUserHandler(dataService))
	admin.Post("/users", handlers.AdminCreateUserHandler(dataService))
	admin.Put("/users/:id", handlers.AdminUpdateUserHandler(dataService))
	admin.Delete("/users/:id", handlers.AdminDeleteUserHandler(dataService))

	// Buildings
	admin.Get("/buildings", handlers.AdminListBuildingsHandler(dataService))
	admin.Get("/buildings/:id", handlers.AdminGetBuildingHandler(dataService))
	admin.Post("/buildings", handlers.AdminCreateBuildingHandler(dataService))
	admin.Put("/buildings/:id", handlers.AdminUpdateBuildingHandler(dataService))
	admin.Delete("/buildings/:id", handlers.AdminDeleteBuildingHandler(dataService))
	admin.Post("/buildings/:id/floors", handlers.AdminAddBuildingFloorHandler(dataService))
	admin.Post("/buildings/:id/color-schemes", handlers.AdminAddBuildingColorSchemaHandler(dataService))
	admin.Post("/buildings/:id/gps", handlers.AdminAddBuildingGpsHandler(dataService))
	admin.Delete("/buildings/:id/floors/:floorId", handlers.AdminRemoveBuildingFloorHandler(dataService))
	admin.Delete("/buildings/:id/color-schemes/:schemaId", handlers.AdminRemoveBuildingColorSchemaHandler(dataService))
	admin.Delete("/buildings/:id/gps/:floorId", handlers.AdminRemoveBuildingGpsHandler(dataService))

	// Floors
	admin.Get("/floors", handlers.AdminListFloorsHandler(dataService))
	admin.Get("/floors/:id", handlers.AdminGetFloorHandler(dataService))
	admin.Post("/floors", handlers.AdminCreateFloorHandler(dataService))
	admin.Put("/floors/:id", handlers.AdminUpdateFloorHandler(dataService))
	admin.Delete("/floors/:id", handlers.AdminDeleteFloorHandler(dataService))
	admin.Post("/floors/:id/rooms", handlers.AdminAddFloorRoomHandler(dataService))
	admin.Post("/floors/:id/services", handlers.AdminAddFloorServiceHandler(dataService))
	admin.Post("/floors/:id/graph-points", handlers.AdminAddFloorGraphPointHandler(dataService))
	admin.Delete("/floors/:id/rooms/:roomId", handlers.AdminRemoveFloorRoomHandler(dataService))
	admin.Delete("/floors/:id/services/:serviceId", handlers.AdminRemoveFloorServiceHandler(dataService))
	admin.Delete("/floors/:id/graph-points/:pointId", handlers.AdminRemoveFloorGraphPointHandler(dataService))

	// Rooms
	admin.Get("/rooms", handlers.AdminListRoomsHandler(dataService))
	admin.Get("/rooms/:id", handlers.AdminGetRoomHandler(dataService))
	admin.Post("/rooms", handlers.AdminCreateRoomHandler(dataService))
	admin.Put("/rooms/:id", handlers.AdminUpdateRoomHandler(dataService))
	admin.Delete("/rooms/:id", handlers.AdminDeleteRoomHandler(dataService))

	// Services
	admin.Get("/services", handlers.AdminListServicesHandler(dataService))
	admin.Get("/services/:id", handlers.AdminGetServiceHandler(dataService))
	admin.Post("/services", handlers.AdminCreateServiceHandler(dataService))
	admin.Put("/services/:id", handlers.AdminUpdateServiceHandler(dataService))
	admin.Delete("/services/:id", handlers.AdminDeleteServiceHandler(dataService))

	// Graph Points
	admin.Get("/graph-points", handlers.AdminListGraphPointsHandler(dataService))
	admin.Get("/graph-points/:id", handlers.AdminGetGraphPointHandler(dataService))
	admin.Post("/graph-points", handlers.AdminCreateGraphPointHandler(dataService))
	admin.Put("/graph-points/:id", handlers.AdminUpdateGraphPointHandler(dataService))
	admin.Delete("/graph-points/:id", handlers.AdminDeleteGraphPointHandler(dataService))

	// Color Schemes
	admin.Get("/color-schemes", handlers.AdminListColorSchemesHandler(dataService))
	admin.Get("/color-schemes/:id", handlers.AdminGetColorSchemeHandler(dataService))
	admin.Post("/color-schemes", handlers.AdminCreateColorSchemeHandler(dataService))
	admin.Put("/color-schemes/:id", handlers.AdminUpdateColorSchemeHandler(dataService))
	admin.Delete("/color-schemes/:id", handlers.AdminDeleteColorSchemeHandler(dataService))

	return app.Listen(s.Port)
}
