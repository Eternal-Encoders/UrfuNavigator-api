package handlers

import (
	"errors"
	"time"
	"urfunavigator/index/auth"
	"urfunavigator/index/models"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func parseParamID(c fiber.Ctx, param string) (bson.ObjectID, error) {
	return models.ParseObjectID(c.Params(param))
}

func userIDFromContext(c fiber.Ctx) (bson.ObjectID, error) {
	claims, ok := auth.ClaimsFromContext(c)
	if !ok {
		return bson.ObjectID{}, errors.New("missing auth claims")
	}
	return models.ParseObjectID(claims.UserID)
}

func stampCreate(base *models.BaseDBSchema, userID bson.ObjectID) {
	now := time.Now()
	if base.Id.IsZero() {
		base.Id = bson.NewObjectID()
	}
	base.CreatedAt = now
	base.UpdateAt = now
	base.LastUpdatedBy = userID
}

func stampUpdate(base *models.BaseDBSchema, userID bson.ObjectID) {
	base.UpdateAt = time.Now()
	base.LastUpdatedBy = userID
}

func handleStoreError(c fiber.Ctx, handler string, err error, attrs ...any) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusNotFound).SendString("Not found")
	}
	logHandlerError(c, handler, err, attrs...)
	return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
}

func sendCreated(c fiber.Ctx, id bson.ObjectID) error {
	return c.Status(fiber.StatusCreated).JSON(models.CreatedResponse{Id: id})
}
