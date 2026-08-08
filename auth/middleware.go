package auth

import (
	"urfunavigator/index/logger"
	"urfunavigator/index/models"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/golang-jwt/jwt/v5"
)

const claimsContextKey = "authClaims"

func RequireJWT(cfg Config) fiber.Handler {
	return jwtware.New(jwtware.Config{
		Extractor: extractors.FromAuthHeader("Bearer"),
		SigningKey: jwtware.SigningKey{
			Key:    cfg.Secret,
			JWTAlg: jwt.SigningMethodHS256.Alg(),
		},
		Claims: &Claims{},
		SuccessHandler: func(c fiber.Ctx) error {
			token := jwtware.FromContext(c)
			if token == nil {
				return unauthorized(c)
			}

			claims, ok := token.Claims.(*Claims)
			if !ok {
				return unauthorized(c)
			}

			c.Locals(claimsContextKey, claims)
			logger.Debug("jwt auth succeeded", "login", claims.Login, "user_id", claims.UserID, "role", claims.Role)
			return c.Next()
		},
	})
}

func RequireAdminRole() fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, ok := ClaimsFromContext(c)
		if !ok {
			return unauthorized(c)
		}

		if !CanAccessAdminAPI(claims.Role) {
			logger.Warn("jwt auth failed: insufficient role", "login", claims.Login, "role", claims.Role)
			return c.Status(fiber.StatusForbidden).SendString("Forbidden")
		}

		return c.Next()
	}
}

func ClaimsFromContext(c fiber.Ctx) (*Claims, bool) {
	claims, ok := c.Locals(claimsContextKey).(*Claims)
	return claims, ok && claims != nil
}

func CanAccessAdminAPI(role models.UserRole) bool {
	return role == models.AdmineRole || role == models.WriterRole
}

func unauthorized(c fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
}
