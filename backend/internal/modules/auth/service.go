package auth

import (
	"api-oa-integrator/internal/database"
	"context"
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

func registerUser(ctx context.Context, in CreateUserRequest) (LoginResponse, error) {
	hp, err := hashPassword(in.Password)
	if err != nil {
		return LoginResponse{}, err
	}

	user, err := database.New(database.D()).CreateUser(ctx, database.CreateUserParams{
		Name:        sql.NullString{String: in.Name, Valid: true},
		Username:    sql.NullString{String: in.Username, Valid: true},
		Password:    sql.NullString{String: hp, Valid: true},
		Permissions: in.Permissions,
	})
	if err != nil {
		zap.L().Sugar().Errorf("Error create user %v", err)
		return LoginResponse{}, err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username.String
	claims["permissions"] = user.Permissions
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

	userToken, err := token.SignedString([]byte(viper.GetString("app.secret")))
	claims["isForRefresh"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24 * 60).Unix()
	refreshToken, err := token.SignedString([]byte(viper.GetString("app.secret")))
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		UserId:       user.ID.String(),
		Username:     user.Username.String,
		Token:        userToken,
		RefreshToken: refreshToken,
	}, nil
}

func login(ctx context.Context, in LoginRequest) (LoginResponse, error) {
	user, err := database.New(database.D()).GetUser(ctx, sql.NullString{String: in.Username, Valid: true})
	if !checkPassword(in.Password, user.Password.String) {
		return LoginResponse{}, errors.New("invalid password")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username.String
	claims["permissions"] = user.Permissions
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

	userToken, err := token.SignedString([]byte(viper.GetString("app.secret")))
	claims["isForRefresh"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24 * 60).Unix()
	refreshToken, err := token.SignedString([]byte(viper.GetString("app.secret")))
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		UserId:       user.ID.String(),
		Username:     user.Username.String,
		Token:        userToken,
		RefreshToken: refreshToken,
		Permissions:  user.Permissions,
	}, nil
}

func deleteUser(ctx context.Context, in uuid.UUID) error {
	_, err := database.New(database.D()).DeleteUser(ctx, in)
	return err
}

func getUsers(ctx context.Context) ([]UsersResponse, error) {
	users, err := database.New(database.D()).GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	var res []UsersResponse
	for _, user := range users {
		res = append(res, UsersResponse{
			Name:        user.Name.String,
			UserId:      user.ID.String(),
			Username:    user.Username.String,
			Permissions: user.Permissions,
		})
	}

	return res, nil
}
