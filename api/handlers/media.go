package handlers

import (
	"urfunavigator/index/models"
	"urfunavigator/index/object"

	"github.com/gofiber/fiber/v3"
)

// ObjectHandler returns a presigned URL for downloading a building icon SVG from object storage.
// @Summary Get building icon presigned URL
// @Tags media
// @Produce json
// @Param icon path string true "Icon filename, must end with .svg"
// @Success 200 {object} models.BuildingIconResponse
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /api/icons/{icon} [get]
func ObjectHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		iconName := c.Params("icon")
		if err := object.ValidateIconFileName(iconName); err != nil {
			logHandlerWarn(c, "ObjectHandler", "invalid icon filename", "icon", iconName, "err", err.Error())
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		iconURL, err := services.ObjectStore.GetBuildingIconURL(iconName)
		if err != nil {
			logHandlerError(c, "ObjectHandler", err, "icon", iconName, "stage", "get_building_icon_url")
			return c.Status(fiber.StatusNotFound).SendString("Cannot get icon from object storage")
		}

		logHandlerDebug(c, "ObjectHandler", "icon presigned url issued", "icon", iconName)
		return c.JSON(models.BuildingIconResponse{
			Name:      iconName,
			URL:       iconURL.URL,
			ExpiresAt: iconURL.ExpiresAt,
		})
	}
}

// DeleteObjectHandler removes a building icon from object storage.
// @Summary Delete building icon SVG
// @Tags admin-media
// @Accept json
// @Produce plain
// @Security BearerAuth
// @Param icon body models.MediaDeleteBody true "Icon filename, must end with .svg"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Router /api/delete_icon [delete]
func DeleteObjectHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		var body models.MediaDeleteBody
		if err := c.Bind().Body(&body); err != nil {
			logHandlerError(c, "DeleteObjectHandler", err, "stage", "bind_body")
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if err := object.ValidateIconFileName(body.Filename); err != nil {
			logHandlerWarn(c, "DeleteObjectHandler", "invalid icon filename", "icon", body.Filename, "err", err.Error())
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if err := services.ObjectStore.RemoveBuildingIcon(body.Filename); err != nil {
			logHandlerError(c, "DeleteObjectHandler", err, "icon", body.Filename, "stage", "remove_building_icon")
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		logHandlerInfo(c, "DeleteObjectHandler", "icon deleted", "icon", body.Filename)
		return c.SendString("ok")
	}
}

// UploadObjectHandler uploads a building icon to object storage.
// @Summary Upload building icon SVG
// @Tags admin-media
// @Accept json
// @Produce plain
// @Security BearerAuth
// @Param icon body models.MediaUploadBody true "Icon filename and byte array data"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Router /api/upload_icon [post]
func UploadObjectHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		var body models.MediaUploadBody
		if err := c.Bind().Body(&body); err != nil {
			logHandlerError(c, "UploadObjectHandler", err, "stage", "bind_body")
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if err := object.ValidateIconFileName(body.Filename); err != nil {
			logHandlerWarn(c, "UploadObjectHandler", "invalid icon filename", "icon", body.Filename, "err", err.Error())
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if err := services.ObjectStore.PutBuildingIcon(body.Filename, body.Image); err != nil {
			logHandlerError(c, "UploadObjectHandler", err, "icon", body.Filename, "stage", "put_building_icon")
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		logHandlerInfo(c, "UploadObjectHandler", "icon uploaded", "icon", body.Filename, "bytes", len(body.Image))
		return c.SendString("ok")
	}
}
