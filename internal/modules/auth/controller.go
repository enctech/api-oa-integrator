package auth

import (
	"api-oa-integrator/internal/middlewares"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type controller struct {
}

func InitController(e *echo.Echo) {
	g := e.Group("auth")
	c := controller{}
	g.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	g.Use(middlewares.GuardSomePathJWT([]string{"/login", "/user"}))
	g.POST("/login", c.login)
	g.POST("/user", c.register)
	g.DELETE("/user/:id", c.deleteUser, middlewares.AdminOnlyMiddleware())
	g.POST("/refresh", c.register)
}

// login godoc
//
//	@Summary		user login
//	@Description	For user to login into admin
//	@Tags			auth
//	@Accept			application/json
//	@Produce		application/json
//	@Param			request	body	LoginRequest	false	"Request Body"
//	@Router			/auth/login [post]
func (con controller) login(c echo.Context) error {
	req := new(LoginRequest)
	err := c.Bind(req)
	user, err := login(c.Request().Context(), *req)
	if err != nil {
		switch err.Error() {
		case "invalid password":
			return c.String(http.StatusUnauthorized, "")
		default:
			return c.String(http.StatusBadRequest, "")
		}
	}
	return c.JSON(http.StatusCreated, user)
}

// register godoc
//
//	@Summary		create new user
//	@Description	For admin to create new user
//	@Tags			auth
//	@Accept			application/json
//	@Produce		application/json
//	@Param			request	body	CreateUserRequest	false	"Request Body"
//	@Router			/auth/user [post]
func (con controller) register(c echo.Context) error {
	req := new(CreateUserRequest)
	err := c.Bind(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, struct {
		}{})
	}
	res, err := registerUser(c.Request().Context(), *req)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}
	return c.JSON(http.StatusCreated, res)
}

// deleteUser godoc
//
//	@Summary		delete user
//	@Description	For admin to delete user
//	@Tags			auth
//	@Security		Bearer
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id	path	string	true	"Id"
//	@Router			/auth/user/{id} [delete]
func (con controller) deleteUser(c echo.Context) error {
	id := c.Param("id")
	req := new(CreateUserRequest)
	err := c.Bind(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Invalid user id %v", id))
	}
	err = deleteUser(c.Request().Context(), uid)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}
	return c.String(http.StatusOK, "")
}
