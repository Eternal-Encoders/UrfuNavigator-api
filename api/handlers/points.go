package handlers

import (
	"urfunavigator/index/models"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// PointsHandler returns map points with optional filters.
// @Summary List points with filters
// @Tags points
// @Produce json
// @Param buildingId query string false "Filter by building ObjectID"
// @Param floorId query string false "Filter by floor ObjectID"
// @Param type query string false "Filter by point type"
// @Param name query string false "Filter by name (case-insensitive substring)"
// @Param length query int false "Max results (default 40)"
// @Success 200 {array} models.GraphPointResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /api/points [get]
func PointsHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		var query models.PointsQuery
		if err := c.Bind().Query(&query); err != nil {
			logHandlerError(c, "PointsHandler", err, "stage", "bind_query")
			return c.Status(fiber.StatusBadRequest).SendString("Something went wrong in QueryParser")
		}

		filter := bson.M{}

		if query.BuildingId != nil {
			buildingID, err := models.ParseObjectID(*query.BuildingId)
			if err != nil {
				logHandlerWarn(c, "PointsHandler", "invalid building id", "building_id", *query.BuildingId)
				return c.Status(fiber.StatusBadRequest).SendString(err.Error())
			}
			filter["buildingId"] = buildingID
		}

		if query.FloorId != nil {
			floorID, err := models.ParseObjectID(*query.FloorId)
			if err != nil {
				logHandlerWarn(c, "PointsHandler", "invalid floor id", "floor_id", *query.FloorId)
				return c.Status(fiber.StatusBadRequest).SendString(err.Error())
			}
			filter["floorId"] = floorID
		}

		if query.Type != nil {
			filter["types"] = bson.M{"$in": []string{*query.Type}}
		}

		if query.Name != nil {
			filter["names.name"] = bson.M{
				"$regex":   *query.Name,
				"$options": "i",
			}
		}

		limit := 40
		if query.Length != nil {
			limit = *query.Length
		}

		points, err := services.Store.FindGraphPoints(filter, limit)
		if err != nil {
			logHandlerError(c, "PointsHandler", err, "stage", "find_graph_points", "limit", limit)
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in FindGraphPoints")
		}

		logHandlerDebug(c, "PointsHandler", "points listed", "count", len(points), "limit", limit)
		return c.JSON(models.ToGraphPointResponses(points))
	}
}

// PointIdHandler returns one point by ObjectID.
// @Summary Get point by id
// @Tags points
// @Produce json
// @Param id query string true "Point ObjectID"
// @Success 200 {object} models.GraphPointResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /api/point [get]
func PointIdHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		var query models.PointIdQuery
		if err := c.Bind().Query(&query); err != nil {
			logHandlerError(c, "PointIdHandler", err, "stage", "bind_query")
			return c.Status(fiber.StatusBadRequest).SendString("Request must contain id query parameter")
		}

		pointID, err := models.ParseObjectID(query.Id)
		if err != nil {
			logHandlerWarn(c, "PointIdHandler", "invalid point id", "point_id", query.Id)
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		point, err := services.Store.GetGraphPoint(pointID)
		if err != nil {
			logHandlerError(c, "PointIdHandler", err, "point_id", pointID.Hex(), "stage", "get_graph_point")
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetGraphPoint")
		}

		logHandlerDebug(c, "PointIdHandler", "point loaded", "point_id", pointID.Hex())
		return c.JSON(models.ToGraphPointResponse(*point))
	}
}

// SearchHandler runs full-text search on points by name.
// @Summary Search points by name
// @Tags points
// @Produce json
// @Param name query string true "Search string"
// @Param length query int false "Max results (1–40, default 40)"
// @Success 200 {array} models.GraphPointResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /api/search [get]
func SearchHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		var query models.SearchQuery
		if err := c.Bind().Query(&query); err != nil {
			logHandlerError(c, "SearchHandler", err, "stage", "bind_query")
			return c.Status(fiber.StatusBadRequest).SendString("Request must contain name query parameter")
		}

		length := 40
		if query.Length != nil {
			rawLength := *query.Length
			if rawLength > 40 {
				rawLength = 40
			}
			if rawLength < 1 {
				rawLength = 1
			}
			length = rawLength
		}

		points, err := services.Store.GetBySearchEngine(query.Name, length)
		if err != nil {
			logHandlerError(c, "SearchHandler", err, "stage", "search", "query", query.Name, "limit", length)
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in Search")
		}

		logHandlerDebug(c, "SearchHandler", "search completed",
			"query", query.Name,
			"limit", length,
			"count", len(points),
		)

		if len(points) == 0 {
			return c.JSON([]models.GraphPointResponse{})
		}

		return c.JSON(models.ToGraphPointResponses(points))
	}
}
