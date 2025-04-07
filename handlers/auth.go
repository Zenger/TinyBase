/*
 * SPDX-License-Identifier: MIT
 * Author: Zenger (https://github.com/Zenger)
 */

package handlers

import (
	"TinyBase/shared"
	"TinyBase/utils"
	"database/sql"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
	"time"
)

type LoginResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AuthHandler(c fiber.Ctx, tbx shared.TinyBaseContext) error {
	var loginRequest LoginRequest
	if err := c.Bind().Body(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if loginRequest.Email == "" || loginRequest.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	query := `SELECT id, password_hash FROM users WHERE email = $1`
	row := tbx.Database.QueryRow(query, loginRequest.Email)
	var id string
	var passwordHash string
	err := row.Scan(&id, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if !utils.CheckPassword(loginRequest.Password, passwordHash, tbx.Settings.App.Salt) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password 2",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": loginRequest.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(tbx.Settings.App.JwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
			"msg":   err.Error(),
		})
	}
	loginResponse := LoginResponse{
		Token: tokenString,
		Email: loginRequest.Email,
	}

	return c.JSON(loginResponse)
}
