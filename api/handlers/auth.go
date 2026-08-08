package handlers

import (
	"urfunavigator/index/auth"
	"urfunavigator/index/models"

	"github.com/gofiber/fiber/v3"
)

// LoginHandler authenticates an admin user and returns a JWT bearer token.
// @Summary Admin login
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Router /api/login [post]
func LoginHandler(services models.DataService, jwtCfg auth.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		var body models.LoginRequest
		if err := c.Bind().Body(&body); err != nil {
			logHandlerError(c, "LoginHandler", err, "stage", "bind_body")
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if body.Login == "" || body.Password == "" {
			logHandlerWarn(c, "LoginHandler", "missing login or password")
			return c.Status(fiber.StatusBadRequest).SendString("login and password are required")
		}

		user, err := services.Store.GetUserByLogin(body.Login)
		if err != nil {
			logHandlerWarn(c, "LoginHandler", "login failed: user lookup error", "login", body.Login, "err", err.Error())
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		valid, err := auth.VerifyPassword(body.Password, user.PasswordHash)
		if err != nil || !valid {
			logHandlerWarn(c, "LoginHandler", "login failed: invalid credentials", "login", body.Login, "user_id", user.Id.Hex())
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		if !auth.CanAccessAdminAPI(user.Role) {
			logHandlerWarn(c, "LoginHandler", "login failed: insufficient role", "login", body.Login, "role", user.Role)
			return c.Status(fiber.StatusForbidden).SendString("Forbidden")
		}

		token, expiresAt, err := auth.GenerateToken(user, jwtCfg)
		if err != nil {
			logHandlerError(c, "LoginHandler", err, "login", body.Login, "stage", "generate_token")
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate token")
		}

		logHandlerInfo(c, "LoginHandler", "login succeeded", "login", body.Login, "user_id", user.Id.Hex(), "role", user.Role)
		return c.JSON(models.LoginResponse{
			Token:     token,
			ExpiresAt: expiresAt,
			Login:     user.Login,
			Role:      user.Role,
		})
	}
}
