package handlers

import (
	"urfunavigator/index/models"

	"github.com/gofiber/fiber/v3"
)

// FloorHandler returns floor layout with resolved rooms, services, and graph points.
// @Summary Get floor layout
// @Tags floors
// @Produce json
// @Param id query string true "Floor ObjectID"
// @Success 200 {object} models.FloorDetailResponse
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Server error"
// @Router /api/floor [get]
func FloorHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		var query models.FloorQuery
		if err := c.Bind().Query(&query); err != nil {
			logHandlerError(c, "FloorHandler", err, "stage", "bind_query")
			return c.Status(fiber.StatusBadRequest).SendString("Request must contain id query parameter")
		}

		floorID, err := models.ParseObjectID(query.Id)
		if err != nil {
			logHandlerWarn(c, "FloorHandler", "invalid floor id", "floor_id", query.Id)
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		floor, err := services.Store.GetFloor(floorID)
		if err != nil {
			logHandlerError(c, "FloorHandler", err, "floor_id", floorID.Hex(), "stage", "get_floor")
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetFloor")
		}

		roomsMap, err := services.Store.GetRooms(floor.Rooms)
		if err != nil {
			logHandlerError(c, "FloorHandler", err, "floor_id", floorID.Hex(), "stage", "get_rooms")
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetRooms")
		}

		servicesMap, err := services.Store.GetServices(floor.Services)
		if err != nil {
			logHandlerError(c, "FloorHandler", err, "floor_id", floorID.Hex(), "stage", "get_services")
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetServices")
		}

		graphMap, err := services.Store.GetGraphPoints(floor.Graph)
		if err != nil {
			logHandlerError(c, "FloorHandler", err, "floor_id", floorID.Hex(), "stage", "get_graph_points")
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetGraphPoints")
		}

		response := models.ToFloorDetailResponse(
			*floor,
			mapValuesToSlice(*roomsMap),
			mapValuesToSlice(*servicesMap),
			mapValuesToSlice(*graphMap),
		)

		logHandlerDebug(c, "FloorHandler", "floor loaded",
			"floor_id", floorID.Hex(),
			"rooms", len(response.Rooms),
			"services", len(response.Services),
			"graph_points", len(response.Graph),
		)

		return c.JSON(response)
	}
}
