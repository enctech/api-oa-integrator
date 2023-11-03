package middlewares

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

func AdminOnlyMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")

			token, err := jwt.Parse(strings.Replace(tokenString, "bearer ", "", 1), func(token *jwt.Token) (interface{}, error) {
				return []byte(viper.GetString("app.secret")), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, "Unauthorized")
			}

			claims := token.Claims.(jwt.MapClaims)

			if claims["permission"] != "admin" {
				return c.JSON(http.StatusUnauthorized, "Unauthorized")
			}
			c.Set("username", claims["username"])
			c.Set("permission", claims["permission"])

			return next(c)
		}
	}
}
