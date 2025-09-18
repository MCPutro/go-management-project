package handler

import (
	"context"
	"errors"
	"github.com/MCPutro/go-management-project/internal/config/constant"
	"github.com/MCPutro/go-management-project/internal/model"
	"github.com/MCPutro/go-management-project/internal/usecase"
	"github.com/MCPutro/go-management-project/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

type UserHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}

type userHandler struct {
	userUsecase usecase.UserUsecase
	//logger      *zap.Logger
}

func NewUserHandler(userUsecase usecase.UserUsecase) UserHandler {
	return &userHandler{userUsecase: userUsecase}
}

func (h *userHandler) Register(c *fiber.Ctx) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		//h.logger.Warn("Invalid request body for register", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input, cannot parse JSON",
		})
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, email, and password are required",
		})
	}

	// Simulate password hashing
	user.Password = user.Email + "hashed_" + user.Password

	// Untuk testing: inject user_id=1 ke context
	ctx := context.WithValue(c.Context(), constant.UserIDKey, int64(1))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := h.userUsecase.CreateUser(ctx, &user); err != nil {
		//h.logger.Error("Failed to register user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Generate dummy JWT token
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDAwMDAwMDB9.fake-signature"

	//h.logger.Info("User registered successfully", zap.Int64("user_id", user.ID), zap.String("email", user.Email))

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
		"token": token,
	})
}

func (h *userHandler) Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Dummy logic — asumsi user ditemukan
	// Di real app: cek ke DB, verify password, generate JWT

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDAwMDAwMDB9.fake-signature"

	//h.logger.Info("User logged in", zap.String("email", req.Email))

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func (h *userHandler) CreateUser(c *fiber.Ctx) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		//h.logger.Warn("Invalid request body for create user", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input, cannot parse JSON",
		})
	}

	if user.Name == "" || user.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name and email are required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	if err := h.userUsecase.CreateUser(ctx, &user); err != nil {
		//h.logger.Error("Failed to create user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Hapus password sebelum kirim response
	user.Password = ""

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *userHandler) GetUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		//h.logger.Warn("Invalid user ID format", zap.String("id", idStr))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	user, err := h.userUsecase.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		//h.logger.Error("Failed to get user", zap.Int64("id", id), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	user.Password = ""
	return c.JSON(user)
}

func (h *userHandler) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	user.ID = id // set ID from URL

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	if err := h.userUsecase.UpdateUser(ctx, &user); err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		//h.logger.Error("Failed to update user", zap.Int64("id", id), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	user.Password = ""
	return c.JSON(user)
}

func (h *userHandler) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	if err := h.userUsecase.DeleteUser(ctx, id, 1); err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		//h.logger.Error("Failed to delete user", zap.Int64("id", id), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
