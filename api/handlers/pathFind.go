package handlers

import (
	"urfunavigator/index/models"

	"github.com/gofiber/fiber/v3"
)

// PathHandler returns a path between two point ObjectIDs.
// @Summary Find path between two points
// @Tags routing
// @Produce json
// @Param from query string true "Start point ObjectID"
// @Param to query string true "End point ObjectID"
// @Success 200 {object} models.PathJSONResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /api/path [get]
func PathHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		var query models.PathQuery
		if err := c.Bind().Query(&query); err != nil {
			logHandlerError(c, "PathHandler", err, "stage", "bind_query")
			return c.Status(fiber.StatusBadRequest).SendString("Request must contain from and to query parameters")
		}

		startID, err := models.ParseObjectID(query.From)
		if err != nil {
			logHandlerWarn(c, "PathHandler", "invalid start point id", "from", query.From)
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		endID, err := models.ParseObjectID(query.To)
		if err != nil {
			logHandlerWarn(c, "PathHandler", "invalid end point id", "to", query.To)
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		start, startErr := services.Store.GetGraphPoint(startID)
		end, endErr := services.Store.GetGraphPoint(endID)
		if startErr != nil {
			logHandlerError(c, "PathHandler", startErr, "from", startID.Hex(), "stage", "get_start_point")
			return c.Status(fiber.StatusInternalServerError).SendString(startErr.Error())
		}
		if endErr != nil {
			logHandlerError(c, "PathHandler", endErr, "to", endID.Hex(), "stage", "get_end_point")
			return c.Status(fiber.StatusInternalServerError).SendString(endErr.Error())
		}

		path, pathErr := services.GeoService.FindPath(
			*start,
			*end,
			services.Store.GetGraphByFloorID,
			services.Store.GetFloorElevation,
			services.Store.GetExitPointsByBuilding,
			services.Store.GetGraphPoint,
		)
		if pathErr != nil {
			logHandlerError(c, "PathHandler", pathErr,
				"from", startID.Hex(),
				"to", endID.Hex(),
				"start_building", start.BuildingId.Hex(),
				"end_building", end.BuildingId.Hex(),
				"stage", "find_path",
			)
			return c.Status(fiber.StatusInternalServerError).SendString(pathErr.Error())
		}

		logHandlerDebug(c, "PathHandler", "path calculated",
			"from", startID.Hex(),
			"to", endID.Hex(),
			"start_building", start.BuildingId.Hex(),
			"end_building", end.BuildingId.Hex(),
			"segments", len(path),
		)

		return c.JSON(models.PathJSONResponse{
			Result: models.ToPathResultResponse(path),
		})
	}
}
