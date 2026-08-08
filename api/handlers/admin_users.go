package handlers

import (
	"urfunavigator/index/auth"
	"urfunavigator/index/models"

	"github.com/gofiber/fiber/v3"
)

func ToUserAdminResponse(user models.User) models.UserAdminResponse {
	return models.UserAdminResponse{
		AdminEntityResponse: models.ToAdminEntityResponse(user.BaseDBSchema),
		Login:               user.Login,
		Role:                user.Role,
	}
}

// AdminListUsersHandler lists all users.
// @Summary List users
// @Tags admin-users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.UserAdminResponse
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /admin_api/users [get]
func AdminListUsersHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		users, err := services.Store.ListUsers()
		if err != nil {
			return handleStoreError(c, "AdminListUsersHandler", err)
		}

		response := make([]models.UserAdminResponse, len(users))
		for i, user := range users {
			response[i] = ToUserAdminResponse(user)
		}
		return c.JSON(response)
	}
}

// AdminGetUserHandler returns a user by ObjectID.
// @Summary Get user
// @Tags admin-users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ObjectID"
// @Success 200 {object} models.UserAdminResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/users/{id} [get]
func AdminGetUserHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		user, err := services.Store.GetUser(id)
		if err != nil {
			return handleStoreError(c, "AdminGetUserHandler", err, "user_id", id.Hex())
		}

		return c.JSON(ToUserAdminResponse(*user))
	}
}

// AdminCreateUserHandler creates a new user.
// @Summary Create user
// @Tags admin-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body models.CreateUserRequest true "User payload"
// @Success 201 {object} models.CreatedResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /admin_api/users [post]
func AdminCreateUserHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		var body models.CreateUserRequest
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if body.Login == "" || body.Password == "" {
			return c.Status(fiber.StatusBadRequest).SendString("login and password are required")
		}

		userID, err := userIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		hash, err := auth.HashPassword(body.Password)
		if err != nil {
			return handleStoreError(c, "AdminCreateUserHandler", err, "stage", "hash_password")
		}

		user := models.User{
			BaseDBSchema: models.BaseDBSchema{
				DisplayableName: body.DisplayableName,
			},
			Login:        body.Login,
			PasswordHash: hash,
			Role:         body.Role,
		}
		stampCreate(&user.BaseDBSchema, userID)

		id, err := services.Store.AddUser(user)
		if err != nil {
			return handleStoreError(c, "AdminCreateUserHandler", err, "stage", "add_user")
		}

		return sendCreated(c, *id)
	}
}

// AdminUpdateUserHandler updates an existing user.
// @Summary Update user
// @Tags admin-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ObjectID"
// @Param user body models.UpdateUserRequest true "User update payload"
// @Success 200 {object} models.UserAdminResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/users/{id} [put]
func AdminUpdateUserHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		var body models.UpdateUserRequest
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		actorID, err := userIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		user, err := services.Store.GetUser(id)
		if err != nil {
			return handleStoreError(c, "AdminUpdateUserHandler", err, "user_id", id.Hex())
		}

		if body.DisplayableName != nil {
			user.DisplayableName = *body.DisplayableName
		}
		if body.Login != nil {
			user.Login = *body.Login
		}
		if body.Role != nil {
			user.Role = *body.Role
		}
		if body.Password != nil && *body.Password != "" {
			hash, hashErr := auth.HashPassword(*body.Password)
			if hashErr != nil {
				return handleStoreError(c, "AdminUpdateUserHandler", hashErr, "stage", "hash_password")
			}
			user.PasswordHash = hash
		}

		stampUpdate(&user.BaseDBSchema, actorID)
		if err := services.Store.UpdateUser(id, *user); err != nil {
			return handleStoreError(c, "AdminUpdateUserHandler", err, "user_id", id.Hex())
		}

		updated, err := services.Store.GetUser(id)
		if err != nil {
			return handleStoreError(c, "AdminUpdateUserHandler", err, "stage", "get_updated")
		}

		return c.JSON(ToUserAdminResponse(*updated))
	}
}

// AdminDeleteUserHandler deletes a user.
// @Summary Delete user
// @Tags admin-users
// @Security BearerAuth
// @Param id path string true "User ObjectID"
// @Success 204 "No Content"
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /admin_api/users/{id} [delete]
func AdminDeleteUserHandler(services models.DataService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseParamID(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if _, err := services.Store.RemoveUser(id); err != nil {
			return handleStoreError(c, "AdminDeleteUserHandler", err, "user_id", id.Hex())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
