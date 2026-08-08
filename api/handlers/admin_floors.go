package handlers

import (
	"urfunavigator/index/models"

	"github.com/gofiber/fiber/v3"
)

// AdminListFloorsHandler lists all floors.
// @Summary List floors
// @Tags admin-floors
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Floor
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /admin_api/floors [get]
func AdminListFloorsHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		floors, err := services.Store.ListFloors()
		if err != nil {
			return handleStoreError(c, "AdminListFloorsHandler", err)
		}
		return c.JSON(floors)
	}
}

// AdminGetFloorHandler returns a floor by ObjectID.
// @Summary Get floor
// @Tags admin-floors
// @Produce json
// @Security BearerAuth
// @Param id path string true "Floor ObjectID"
// @Success 200 {object} models.Floor
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/floors/{id} [get]
func AdminGetFloorHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		floor, err := services.Store.GetFloor(id)
		if err != nil {
			return handleStoreError(c, "AdminGetFloorHandler", err, "floor_id", id.Hex())
		}

		return c.JSON(floor)
	}
}

// AdminCreateFloorHandler creates a standalone floor document.
// @Summary Create floor
// @Tags admin-floors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param floor body models.Floor true "Floor payload"
// @Success 201 {object} models.CreatedResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /admin_api/floors [post]
func AdminCreateFloorHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		var body models.Floor
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		userID, err := userIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		stampCreate(&body.BaseDBSchema, userID)
		id, err := services.Store.AddFloor(body)
		if err != nil {
			return handleStoreError(c, "AdminCreateFloorHandler", err)
		}

		return sendCreated(c, *id)
	}
}

// AdminUpdateFloorHandler replaces a floor document.
// @Summary Update floor
// @Tags admin-floors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Floor ObjectID"
// @Param floor body models.Floor true "Floor payload"
// @Success 200 {object} models.Floor
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/floors/{id} [put]
func AdminUpdateFloorHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		var body models.Floor
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		userID, err := userIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		body.Id = id
		stampUpdate(&body.BaseDBSchema, userID)
		if err := services.Store.UpdateFloor(id, body); err != nil {
			return handleStoreError(c, "AdminUpdateFloorHandler", err, "floor_id", id.Hex())
		}

		updated, err := services.Store.GetFloor(id)
		if err != nil {
			return handleStoreError(c, "AdminUpdateFloorHandler", err, "stage", "get_updated")
		}

		return c.JSON(updated)
	}
}

// AdminDeleteFloorHandler deletes a floor and cascades child removal.
// @Summary Delete floor
// @Tags admin-floors
// @Security BearerAuth
// @Param id path string true "Floor ObjectID"
// @Success 204 "No Content"
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/floors/{id} [delete]
func AdminDeleteFloorHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if _, err := services.Store.RemoveFloor(id); err != nil {
			return handleStoreError(c, "AdminDeleteFloorHandler", err, "floor_id", id.Hex())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// AdminAddFloorRoomHandler creates a room and links it to a floor.
// @Summary Add room to floor
// @Tags admin-floors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Floor ObjectID"
// @Param room body models.Room true "Room payload"
// @Success 201 {object} models.CreatedResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/floors/{id}/rooms [post]
func AdminAddFloorRoomHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		floorID, err := parseParamID(c, "id")
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

		stampCreate(&body.BaseDBSchema, userID)
		id, err := services.Store.AddFloorRoom(floorID, body)
		if err != nil {
			return handleStoreError(c, "AdminAddFloorRoomHandler", err, "floor_id", floorID.Hex())
		}

		return sendCreated(c, *id)
	}
}

// AdminAddFloorServiceHandler creates a service and links it to a floor.
// @Summary Add service to floor
// @Tags admin-floors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Floor ObjectID"
// @Param service body models.Service true "Service payload"
// @Success 201 {object} models.CreatedResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/floors/{id}/services [post]
func AdminAddFloorServiceHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		floorID, err := parseParamID(c, "id")
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

		stampCreate(&body.BaseDBSchema, userID)
		id, err := services.Store.AddFloorService(floorID, body)
		if err != nil {
			return handleStoreError(c, "AdminAddFloorServiceHandler", err, "floor_id", floorID.Hex())
		}

		return sendCreated(c, *id)
	}
}

// AdminAddFloorGraphPointHandler creates a graph point and links it to a floor.
// @Summary Add graph point to floor
// @Tags admin-floors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Floor ObjectID"
// @Param point body models.GraphPoint true "Graph point payload"
// @Success 201 {object} models.CreatedResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/floors/{id}/graph-points [post]
func AdminAddFloorGraphPointHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		floorID, err := parseParamID(c, "id")
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

		stampCreate(&body.BaseDBSchema, userID)
		id, err := services.Store.AddFloorGraphPoint(floorID, body)
		if err != nil {
			return handleStoreError(c, "AdminAddFloorGraphPointHandler", err, "floor_id", floorID.Hex())
		}

		return sendCreated(c, *id)
	}
}

// AdminRemoveFloorRoomHandler unlinks and deletes a room from a floor.
// @Summary Remove room from floor
// @Tags admin-floors
// @Security BearerAuth
// @Param id path string true "Floor ObjectID"
// @Param roomId path string true "Room ObjectID"
// @Success 204 "No Content"
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/floors/{id}/rooms/{roomId} [delete]
func AdminRemoveFloorRoomHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		floorID, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		roomID, err := parseParamID(c, "roomId")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if err := services.Store.RemoveFloorRoomLink(floorID, roomID); err != nil {
			return handleStoreError(c, "AdminRemoveFloorRoomHandler", err,
				"floor_id", floorID.Hex(), "room_id", roomID.Hex())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// AdminRemoveFloorServiceHandler unlinks and deletes a service from a floor.
// @Summary Remove service from floor
// @Tags admin-floors
// @Security BearerAuth
// @Param id path string true "Floor ObjectID"
// @Param serviceId path string true "Service ObjectID"
// @Success 204 "No Content"
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/floors/{id}/services/{serviceId} [delete]
func AdminRemoveFloorServiceHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		floorID, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		serviceID, err := parseParamID(c, "serviceId")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if err := services.Store.RemoveFloorServiceLink(floorID, serviceID); err != nil {
			return handleStoreError(c, "AdminRemoveFloorServiceHandler", err,
				"floor_id", floorID.Hex(), "service_id", serviceID.Hex())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// AdminRemoveFloorGraphPointHandler unlinks and deletes a graph point from a floor.
// @Summary Remove graph point from floor
// @Tags admin-floors
// @Security BearerAuth
// @Param id path string true "Floor ObjectID"
// @Param pointId path string true "Graph point ObjectID"
// @Success 204 "No Content"
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/floors/{id}/graph-points/{pointId} [delete]
func AdminRemoveFloorGraphPointHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		floorID, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		pointID, err := parseParamID(c, "pointId")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if err := services.Store.RemoveFloorGraphPointLink(floorID, pointID); err != nil {
			return handleStoreError(c, "AdminRemoveFloorGraphPointHandler", err,
				"floor_id", floorID.Hex(), "point_id", pointID.Hex())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
