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

	txn, _ := database.D().Begin()
	user, err := database.New(database.D()).WithTx(txn).CreateUser(ctx, database.CreateUserParams{
		Username:   sql.NullString{String: in.Username, Valid: true},
		Password:   sql.NullString{String: hp, Valid: true},
		Permission: sql.NullString{String: in.Permission, Valid: true},
	})
	if err != nil {
		zap.L().Sugar().Errorf("Error create user %v", err)
		return LoginResponse{}, err
	}
	err = txn.Commit()
	if err != nil {
		zap.L().Sugar().Errorf("Error commit user txn %v", err)
		return LoginResponse{}, err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username.String
	claims["permission"] = user.Permission.String
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
	txn, _ := database.D().Begin()
	user, err := database.New(database.D()).WithTx(txn).GetUser(ctx, sql.NullString{String: in.Username, Valid: true})

	if !checkPassword(in.Password, user.Password.String) {
		return LoginResponse{}, errors.New("invalid password")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username.String
	claims["permission"] = user.Permission.String
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
		Permission:   user.Permission.String,
	}, nil
}

func deleteUser(ctx context.Context, in uuid.UUID) error {
	txn, _ := database.D().Begin()
	_, err := database.New(database.D()).WithTx(txn).DeleteUser(ctx, in)
	if err != nil {
		return err
	}
	err = txn.Commit()

	return err
}
