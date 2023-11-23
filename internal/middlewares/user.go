package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

func GuardWithJWT() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(viper.GetString("app.secret")),
	})
}

func GuardSomePathJWT(paths []string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(viper.GetString("app.secret")),
		Skipper: func(c echo.Context) bool {
			for _, path := range paths {
				if strings.HasSuffix(c.Path(), path) {
					return true
				}
			}
			return false
		},
	})
}

func AdminOnlyMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")

			token, err := jwt.Parse(strings.Replace(tokenString, "Bearer ", "", 1), func(token *jwt.Token) (interface{}, error) {
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
