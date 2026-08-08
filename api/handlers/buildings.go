package handlers

import (
	"fmt"
	"urfunavigator/index/models"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func loadBuildingIcon(objectStore models.ObjectStore, iconName string) (models.BuildingIconResponse, error) {
	iconURL, err := objectStore.GetBuildingIconURL(iconName)
	if err != nil {
		return models.BuildingIconResponse{}, err
	}

	return models.BuildingIconResponse{
		Name:      iconName,
		URL:       iconURL.URL,
		ExpiresAt: iconURL.ExpiresAt,
	}, nil
}

func loadBuildingIcons(objectStore models.ObjectStore, iconNames []string) (map[string]models.BuildingIconResponse, error) {
	unique := make(map[string]struct{}, len(iconNames))
	for _, name := range iconNames {
		if name == "" {
			continue
		}
		unique[name] = struct{}{}
	}

	if len(unique) == 0 {
		return map[string]models.BuildingIconResponse{}, nil
	}

	results := make(map[string]models.BuildingIconResponse, len(unique))
	for iconName := range unique {
		icon, err := loadBuildingIcon(objectStore, iconName)
		if err != nil {
			return nil, fmt.Errorf("icon %s: %w", iconName, err)
		}
		results[iconName] = icon
	}

	return results, nil
}

func toBuildingResponse(building models.Building, icon models.BuildingIconResponse, floors map[bson.ObjectID]*models.Floor) models.BuildingResponse {
	return models.ToBuildingResponse(building, icon, floors)
}

// BuildingHandler returns a single building with resolved icon.
// @Summary Get building
// @Tags buildings
// @Produce json
// @Param id query string true "Building ObjectID"
// @Success 200 {object} models.BuildingResponse
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/building [get]
func BuildingHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		var query models.BuildingQuery
		if err := c.Bind().Query(&query); err != nil {
			logHandlerError(c, "BuildingHandler", err, "stage", "bind_query")
			return c.Status(fiber.StatusBadRequest).SendString("Request must contain id query parameter")
		}

		buildingID, err := models.ParseObjectID(query.Id)
		if err != nil {
			logHandlerWarn(c, "BuildingHandler", "invalid building id", "building_id", query.Id)
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		building, err := services.Store.GetBuilding(buildingID)
		if err != nil {
			logHandlerError(c, "BuildingHandler", err, "building_id", buildingID.Hex(), "stage", "get_building")
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetBuilding")
		}

		icons, iconErr := loadBuildingIcons(services.ObjectStore, []string{building.Icon})
		if iconErr != nil {
			logHandlerWarn(c, "BuildingHandler", "building icon not found",
				"building_id", buildingID.Hex(),
				"icon", building.Icon,
				"err", iconErr.Error(),
			)
			return c.Status(fiber.StatusNotFound).SendString("Cannot find building icon")
		}

		icon, ok := icons[building.Icon]
		if !ok {
			logHandlerWarn(c, "BuildingHandler", "building icon missing after load",
				"building_id", buildingID.Hex(),
				"icon", building.Icon,
			)
			return c.Status(fiber.StatusNotFound).SendString("Cannot find building icon")
		}

		floors, err := services.Store.GetFloors(building.Floors)
		if err != nil {
			logHandlerError(c, "BuildingHandler", err, "building_id", buildingID.Hex(), "stage", "get_floors")
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetFloors")
		}

		logHandlerDebug(c, "BuildingHandler", "building loaded",
			"building_id", buildingID.Hex(),
			"icon", building.Icon,
		)

		return c.JSON(toBuildingResponse(*building, icon, *floors))
	}
}

// BuildingsHandler lists all buildings with resolved icons.
// @Summary List buildings
// @Tags buildings
// @Produce json
// @Success 200 {array} models.BuildingResponse
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/buildings [get]
func BuildingsHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		buildings, err := services.Store.ListBuildings()
		if err != nil {
			logHandlerError(c, "BuildingsHandler", err, "stage", "list_buildings")
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in ListBuildings")
		}

		iconNames := make([]string, 0, len(buildings))
		for _, building := range buildings {
			iconNames = append(iconNames, building.Icon)
		}

		icons, iconErr := loadBuildingIcons(services.ObjectStore, iconNames)
		if iconErr != nil {
			logHandlerWarn(c, "BuildingsHandler", "failed to load building icons", "err", iconErr.Error())
			return c.Status(fiber.StatusNotFound).SendString("For some buildings icons were not found")
		}

		floorIDs := make([]bson.ObjectID, 0)
		for _, building := range buildings {
			floorIDs = append(floorIDs, building.Floors...)
		}

		floors, err := services.Store.GetFloors(floorIDs)
		if err != nil {
			logHandlerError(c, "BuildingsHandler", err, "stage", "get_floors")
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetFloors")
		}

		response := make([]models.BuildingResponse, 0, len(buildings))
		for _, building := range buildings {
			icon, ok := icons[building.Icon]
			if !ok {
				logHandlerWarn(c, "BuildingsHandler", "building icon missing after load",
					"building_id", building.Id.Hex(),
					"icon", building.Icon,
				)
				return c.Status(fiber.StatusNotFound).SendString("For some buildings icons were not found")
			}

			response = append(response, toBuildingResponse(building, icon, *floors))
		}

		logHandlerDebug(c, "BuildingsHandler", "buildings listed",
			"count", len(response),
			"unique_icons", len(icons),
		)
		return c.JSON(response)
	}
}

func mapValuesToSlice[T any](items map[bson.ObjectID]*T) []T {
	result := make([]T, 0, len(items))
	for _, item := range items {
		if item != nil {
			result = append(result, *item)
		}
	}
	return result
}
