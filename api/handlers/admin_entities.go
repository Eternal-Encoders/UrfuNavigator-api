package handlers

import (
	"urfunavigator/index/models"

	"github.com/gofiber/fiber/v3"
)

func adminListHandler[T any](c fiber.Ctx, handler string, list func() ([]T, error)) error {
	items, err := list()
	if err != nil {
		return handleStoreError(c, handler, err)
	}
	return c.JSON(items)
}

// Rooms

// AdminListRoomsHandler lists all rooms.
// @Summary List rooms
// @Tags admin-rooms
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Room
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /admin_api/rooms [get]
func AdminListRoomsHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		return adminListHandler(c, "AdminListRoomsHandler", services.Store.ListRooms)
	}
}

// AdminGetRoomHandler returns a room by ObjectID.
// @Summary Get room
// @Tags admin-rooms
// @Produce json
// @Security BearerAuth
// @Param id path string true "Room ObjectID"
// @Success 200 {object} models.Room
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/rooms/{id} [get]
func AdminGetRoomHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		room, err := services.Store.GetRoom(id)
		if err != nil {
			return handleStoreError(c, "AdminGetRoomHandler", err, "room_id", id.Hex())
		}
		return c.JSON(room)
	}
}

// AdminCreateRoomHandler creates a standalone room document.
// @Summary Create room
// @Tags admin-rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param room body models.Room true "Room payload"
// @Success 201 {object} models.CreatedResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /admin_api/rooms [post]
func AdminCreateRoomHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		var body models.Room
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		userID, err := userIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		stampCreate(&body.BaseDBSchema, userID)
		id, err := services.Store.AddRoom(body)
		if err != nil {
			return handleStoreError(c, "AdminCreateRoomHandler", err)
		}

		return sendCreated(c, *id)
	}
}

// AdminUpdateRoomHandler replaces a room document.
// @Summary Update room
// @Tags admin-rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Room ObjectID"
// @Param room body models.Room true "Room payload"
// @Success 200 {object} models.Room
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/rooms/{id} [put]
func AdminUpdateRoomHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		var body models.Room
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		userID, err := userIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		body.Id = id
		stampUpdate(&body.BaseDBSchema, userID)
		if err := services.Store.UpdateRoom(id, body); err != nil {
			return handleStoreError(c, "AdminUpdateRoomHandler", err, "room_id", id.Hex())
		}

		updated, err := services.Store.GetRoom(id)
		if err != nil {
			return handleStoreError(c, "AdminUpdateRoomHandler", err, "stage", "get_updated")
		}

		return c.JSON(updated)
	}
}

// AdminDeleteRoomHandler deletes a room.
// @Summary Delete room
// @Tags admin-rooms
// @Security BearerAuth
// @Param id path string true "Room ObjectID"
// @Success 204 "No Content"
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/rooms/{id} [delete]
func AdminDeleteRoomHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if _, err := services.Store.RemoveRoom(id); err != nil {
			return handleStoreError(c, "AdminDeleteRoomHandler", err, "room_id", id.Hex())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// Services

// AdminListServicesHandler lists all services.
// @Summary List services
// @Tags admin-services
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Service
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /admin_api/services [get]
func AdminListServicesHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		return adminListHandler(c, "AdminListServicesHandler", services.Store.ListServices)
	}
}

// AdminGetServiceHandler returns a service by ObjectID.
// @Summary Get service
// @Tags admin-services
// @Produce json
// @Security BearerAuth
// @Param id path string true "Service ObjectID"
// @Success 200 {object} models.Service
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/services/{id} [get]
func AdminGetServiceHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		service, err := services.Store.GetService(id)
		if err != nil {
			return handleStoreError(c, "AdminGetServiceHandler", err, "service_id", id.Hex())
		}
		return c.JSON(service)
	}
}

// AdminCreateServiceHandler creates a standalone service document.
// @Summary Create service
// @Tags admin-services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param service body models.Service true "Service payload"
// @Success 201 {object} models.CreatedResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /admin_api/services [post]
func AdminCreateServiceHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		var body models.Service
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		userID, err := userIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		stampCreate(&body.BaseDBSchema, userID)
		id, err := services.Store.AddService(body)
		if err != nil {
			return handleStoreError(c, "AdminCreateServiceHandler", err)
		}

		return sendCreated(c, *id)
	}
}

// AdminUpdateServiceHandler replaces a service document.
// @Summary Update service
// @Tags admin-services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Service ObjectID"
// @Param service body models.Service true "Service payload"
// @Success 200 {object} models.Service
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/services/{id} [put]
func AdminUpdateServiceHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		var body models.Service
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		userID, err := userIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		body.Id = id
		stampUpdate(&body.BaseDBSchema, userID)
		if err := services.Store.UpdateService(id, body); err != nil {
			return handleStoreError(c, "AdminUpdateServiceHandler", err, "service_id", id.Hex())
		}

		updated, err := services.Store.GetService(id)
		if err != nil {
			return handleStoreError(c, "AdminUpdateServiceHandler", err, "stage", "get_updated")
		}

		return c.JSON(updated)
	}
}

// AdminDeleteServiceHandler deletes a service.
// @Summary Delete service
// @Tags admin-services
// @Security BearerAuth
// @Param id path string true "Service ObjectID"
// @Success 204 "No Content"
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/services/{id} [delete]
func AdminDeleteServiceHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if _, err := services.Store.RemoveService(id); err != nil {
			return handleStoreError(c, "AdminDeleteServiceHandler", err, "service_id", id.Hex())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// Graph Points

// AdminListGraphPointsHandler lists all graph points.
// @Summary List graph points
// @Tags admin-graph-points
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.GraphPoint
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /admin_api/graph-points [get]
func AdminListGraphPointsHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		return adminListHandler(c, "AdminListGraphPointsHandler", services.Store.ListGraphPoints)
	}
}

// AdminGetGraphPointHandler returns a graph point by ObjectID.
// @Summary Get graph point
// @Tags admin-graph-points
// @Produce json
// @Security BearerAuth
// @Param id path string true "Graph point ObjectID"
// @Success 200 {object} models.GraphPoint
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/graph-points/{id} [get]
func AdminGetGraphPointHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		point, err := services.Store.GetGraphPoint(id)
		if err != nil {
			return handleStoreError(c, "AdminGetGraphPointHandler", err, "point_id", id.Hex())
		}
		return c.JSON(point)
	}
}

// AdminCreateGraphPointHandler creates a standalone graph point document.
// @Summary Create graph point
// @Tags admin-graph-points
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param point body models.GraphPoint true "Graph point payload"
// @Success 201 {object} models.CreatedResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /admin_api/graph-points [post]
func AdminCreateGraphPointHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		var body models.GraphPoint
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		userID, err := userIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		stampCreate(&body.BaseDBSchema, userID)
		id, err := services.Store.AddGraphPoint(body)
		if err != nil {
			return handleStoreError(c, "AdminCreateGraphPointHandler", err)
		}

		return sendCreated(c, *id)
	}
}

// AdminUpdateGraphPointHandler replaces a graph point document.
// @Summary Update graph point
// @Tags admin-graph-points
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Graph point ObjectID"
// @Param point body models.GraphPoint true "Graph point payload"
// @Success 200 {object} models.GraphPoint
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/graph-points/{id} [put]
func AdminUpdateGraphPointHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		var body models.GraphPoint
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		userID, err := userIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		body.Id = id
		stampUpdate(&body.BaseDBSchema, userID)
		if err := services.Store.UpdateGraphPoint(id, body); err != nil {
			return handleStoreError(c, "AdminUpdateGraphPointHandler", err, "point_id", id.Hex())
		}

		updated, err := services.Store.GetGraphPoint(id)
		if err != nil {
			return handleStoreError(c, "AdminUpdateGraphPointHandler", err, "stage", "get_updated")
		}

		return c.JSON(updated)
	}
}

// AdminDeleteGraphPointHandler deletes a graph point.
// @Summary Delete graph point
// @Tags admin-graph-points
// @Security BearerAuth
// @Param id path string true "Graph point ObjectID"
// @Success 204 "No Content"
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/graph-points/{id} [delete]
func AdminDeleteGraphPointHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if _, err := services.Store.RemoveGraphPoint(id); err != nil {
			return handleStoreError(c, "AdminDeleteGraphPointHandler", err, "point_id", id.Hex())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// Color Schemes

// AdminListColorSchemesHandler lists all building color schemes.
// @Summary List color schemes
// @Tags admin-color-schemes
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.BuildingColorSchema
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /admin_api/color-schemes [get]
func AdminListColorSchemesHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		return adminListHandler(c, "AdminListColorSchemesHandler", services.Store.ListBuildingColorSchemes)
	}
}

// AdminGetColorSchemeHandler returns a color scheme by ObjectID.
// @Summary Get color scheme
// @Tags admin-color-schemes
// @Produce json
// @Security BearerAuth
// @Param id path string true "Color scheme ObjectID"
// @Success 200 {object} models.BuildingColorSchema
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/color-schemes/{id} [get]
func AdminGetColorSchemeHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		schema, err := services.Store.GetBuildingColorSchema(id)
		if err != nil {
			return handleStoreError(c, "AdminGetColorSchemeHandler", err, "schema_id", id.Hex())
		}
		return c.JSON(schema)
	}
}

// AdminCreateColorSchemeHandler creates a standalone color scheme document.
// @Summary Create color scheme
// @Tags admin-color-schemes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param schema body models.BuildingColorSchema true "Color schema payload"
// @Success 201 {object} models.CreatedResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /admin_api/color-schemes [post]
func AdminCreateColorSchemeHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		var body models.BuildingColorSchema
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		userID, err := userIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		stampCreate(&body.BaseDBSchema, userID)
		id, err := services.Store.AddBuildingColorSchema(body)
		if err != nil {
			return handleStoreError(c, "AdminCreateColorSchemeHandler", err)
		}

		return sendCreated(c, *id)
	}
}

// AdminUpdateColorSchemeHandler replaces a color scheme document.
// @Summary Update color scheme
// @Tags admin-color-schemes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Color scheme ObjectID"
// @Param schema body models.BuildingColorSchema true "Color schema payload"
// @Success 200 {object} models.BuildingColorSchema
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/color-schemes/{id} [put]
func AdminUpdateColorSchemeHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		var body models.BuildingColorSchema
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		userID, err := userIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		body.Id = id
		stampUpdate(&body.BaseDBSchema, userID)
		if err := services.Store.UpdateBuildingColorSchema(id, body); err != nil {
			return handleStoreError(c, "AdminUpdateColorSchemeHandler", err, "schema_id", id.Hex())
		}

		updated, err := services.Store.GetBuildingColorSchema(id)
		if err != nil {
			return handleStoreError(c, "AdminUpdateColorSchemeHandler", err, "stage", "get_updated")
		}

		return c.JSON(updated)
	}
}

// AdminDeleteColorSchemeHandler deletes a color scheme.
// @Summary Delete color scheme
// @Tags admin-color-schemes
// @Security BearerAuth
// @Param id path string true "Color scheme ObjectID"
// @Success 204 "No Content"
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/color-schemes/{id} [delete]
func AdminDeleteColorSchemeHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if _, err := services.Store.RemoveBuildingColorSchema(id); err != nil {
			return handleStoreError(c, "AdminDeleteColorSchemeHandler", err, "schema_id", id.Hex())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
