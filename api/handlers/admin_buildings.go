package handlers

import (
	"urfunavigator/index/models"

	"github.com/gofiber/fiber/v3"
)

// AdminListBuildingsHandler lists all buildings.
// @Summary List buildings
// @Tags admin-buildings
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Building
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /admin_api/buildings [get]
func AdminListBuildingsHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		buildings, err := services.Store.ListBuildings()
		if err != nil {
			return handleStoreError(c, "AdminListBuildingsHandler", err)
		}
		return c.JSON(buildings)
	}
}

// AdminGetBuildingHandler returns a building by ObjectID.
// @Summary Get building
// @Tags admin-buildings
// @Produce json
// @Security BearerAuth
// @Param id path string true "Building ObjectID"
// @Success 200 {object} models.Building
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/buildings/{id} [get]
func AdminGetBuildingHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		building, err := services.Store.GetBuilding(id)
		if err != nil {
			return handleStoreError(c, "AdminGetBuildingHandler", err, "building_id", id.Hex())
		}

		return c.JSON(building)
	}
}

// AdminCreateBuildingHandler creates a new building.
// @Summary Create building
// @Tags admin-buildings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param building body models.Building true "Building payload"
// @Success 201 {object} models.CreatedResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /admin_api/buildings [post]
func AdminCreateBuildingHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		var body models.Building
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		userID, err := userIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		stampCreate(&body.BaseDBSchema, userID)
		id, err := services.Store.AddBuilding(body)
		if err != nil {
			return handleStoreError(c, "AdminCreateBuildingHandler", err)
		}

		return sendCreated(c, *id)
	}
}

// AdminUpdateBuildingHandler replaces a building document.
// @Summary Update building
// @Tags admin-buildings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Building ObjectID"
// @Param building body models.Building true "Building payload"
// @Success 200 {object} models.Building
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/buildings/{id} [put]
func AdminUpdateBuildingHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		var body models.Building
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		userID, err := userIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		body.Id = id
		stampUpdate(&body.BaseDBSchema, userID)
		if err := services.Store.UpdateBuilding(id, body); err != nil {
			return handleStoreError(c, "AdminUpdateBuildingHandler", err, "building_id", id.Hex())
		}

		updated, err := services.Store.GetBuilding(id)
		if err != nil {
			return handleStoreError(c, "AdminUpdateBuildingHandler", err, "stage", "get_updated")
		}

		return c.JSON(updated)
	}
}

// AdminDeleteBuildingHandler deletes a building and cascades child removal.
// @Summary Delete building
// @Tags admin-buildings
// @Security BearerAuth
// @Param id path string true "Building ObjectID"
// @Success 204 "No Content"
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/buildings/{id} [delete]
func AdminDeleteBuildingHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if _, err := services.Store.RemoveBuilding(id); err != nil {
			return handleStoreError(c, "AdminDeleteBuildingHandler", err, "building_id", id.Hex())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// AdminAddBuildingFloorHandler creates a floor and links it to a building.
// @Summary Add floor to building
// @Tags admin-buildings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Building ObjectID"
// @Param floor body models.Floor true "Floor payload"
// @Success 201 {object} models.CreatedResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/buildings/{id}/floors [post]
func AdminAddBuildingFloorHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		buildingID, err := parseParamID(c, "id")
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

		stampCreate(&body.BaseDBSchema, userID)
		id, err := services.Store.AddBuildingFloor(buildingID, body)
		if err != nil {
			return handleStoreError(c, "AdminAddBuildingFloorHandler", err, "building_id", buildingID.Hex())
		}

		return sendCreated(c, *id)
	}
}

// AdminAddBuildingColorSchemaHandler creates a color schema and links it to a building.
// @Summary Add color schema to building
// @Tags admin-buildings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Building ObjectID"
// @Param schema body models.BuildingColorSchema true "Color schema payload"
// @Success 201 {object} models.CreatedResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/buildings/{id}/color-schemes [post]
func AdminAddBuildingColorSchemaHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		buildingID, err := parseParamID(c, "id")
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

		stampCreate(&body.BaseDBSchema, userID)
		id, err := services.Store.AddBuildingColorSchemaLink(buildingID, body)
		if err != nil {
			return handleStoreError(c, "AdminAddBuildingColorSchemaHandler", err, "building_id", buildingID.Hex())
		}

		return sendCreated(c, *id)
	}
}

// AdminAddBuildingGpsHandler appends a GPS mapping entry to a building.
// @Summary Add GPS entry to building
// @Tags admin-buildings
// @Accept json
// @Produce plain
// @Security BearerAuth
// @Param id path string true "Building ObjectID"
// @Param gps body models.BuildingGps true "GPS payload"
// @Success 204 "No Content"
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/buildings/{id}/gps [post]
func AdminAddBuildingGpsHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		buildingID, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		var body models.BuildingGps
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if err := services.Store.AddBuildingGps(buildingID, body); err != nil {
			return handleStoreError(c, "AdminAddBuildingGpsHandler", err, "building_id", buildingID.Hex())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// AdminRemoveBuildingFloorHandler unlinks and deletes a floor from a building.
// @Summary Remove floor from building
// @Tags admin-buildings
// @Security BearerAuth
// @Param id path string true "Building ObjectID"
// @Param floorId path string true "Floor ObjectID"
// @Success 204 "No Content"
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/buildings/{id}/floors/{floorId} [delete]
func AdminRemoveBuildingFloorHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		buildingID, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		floorID, err := parseParamID(c, "floorId")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if err := services.Store.RemoveBuildingFloorLink(buildingID, floorID); err != nil {
			return handleStoreError(c, "AdminRemoveBuildingFloorHandler", err,
				"building_id", buildingID.Hex(), "floor_id", floorID.Hex())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// AdminRemoveBuildingColorSchemaHandler unlinks and deletes a color schema from a building.
// @Summary Remove color schema from building
// @Tags admin-buildings
// @Security BearerAuth
// @Param id path string true "Building ObjectID"
// @Param schemaId path string true "Color schema ObjectID"
// @Success 204 "No Content"
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/buildings/{id}/color-schemes/{schemaId} [delete]
func AdminRemoveBuildingColorSchemaHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		buildingID, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		schemaID, err := parseParamID(c, "schemaId")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if err := services.Store.RemoveBuildingColorSchemaLink(buildingID, schemaID); err != nil {
			return handleStoreError(c, "AdminRemoveBuildingColorSchemaHandler", err,
				"building_id", buildingID.Hex(), "schema_id", schemaID.Hex())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// AdminRemoveBuildingGpsHandler removes a GPS entry from a building by floorId.
// @Summary Remove GPS entry from building
// @Tags admin-buildings
// @Security BearerAuth
// @Param id path string true "Building ObjectID"
// @Param floorId path string true "Floor ObjectID referenced in GPS entry"
// @Success 204 "No Content"
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/buildings/{id}/gps/{floorId} [delete]
func AdminRemoveBuildingGpsHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		buildingID, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		floorID, err := parseParamID(c, "floorId")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if err := services.Store.RemoveBuildingGps(buildingID, floorID); err != nil {
			return handleStoreError(c, "AdminRemoveBuildingGpsHandler", err,
				"building_id", buildingID.Hex(), "floor_id", floorID.Hex())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
