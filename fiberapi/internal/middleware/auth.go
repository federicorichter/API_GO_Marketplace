package middleware

import (
    "fiberapi/internal/utils"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v4"
    "log"
    "strings"
)

// JWTMiddleware checks the validity of the JWT token
func JWTMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or malformed JWT"})
        }

        // Remove 'Bearer ' prefix from the token
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
            }
            return utils.JwtKey, nil
        })
        if err != nil || !token.Valid {
            log.Printf("Error parsing token: %v", err)
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired JWT"})
        }

        c.Locals("user", token.Claims.(jwt.MapClaims))
        return c.Next()
    }
}

// AdminMiddleware checks if the user has admin role
func AdminMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        claims := c.Locals("user").(jwt.MapClaims)
        role, ok := claims["role"].(string)
        if !ok || role != "admin" {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
        }
        return c.Next()
    }
}
