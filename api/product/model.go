package product

import (
	"app-backend/middleware"
	"app-backend/service"

	"github.com/gofiber/fiber/v2"
)

var (
	dbCtx       = service.DbCtx
	cacheClient = service.CacheClient
	mongoClient = service.MongoDbClient
)

func GetUsername(ctx *fiber.Ctx) (string, error) {
	claim, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return "", ctx.Status(fiber.StatusUnauthorized).JSON(Result{
			Status:    fiber.StatusUnauthorized,
			Message:   "unauthorized",
			MessageTh: "ยืนยันตัวตนไม่สำเร็จ",
			Error:     "unauthorized",
		})
	}
	return claim.Username, nil
}

func GetFirstName(ctx *fiber.Ctx) (string, error) {
	claim, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return "", ctx.Status(fiber.StatusUnauthorized).JSON(Result{
			Status:    fiber.StatusUnauthorized,
			Message:   "unauthorized",
			MessageTh: "ยืนยันตัวตนไม่สำเร็จ",
			Error:     "unauthorized",
		})
	}
	return claim.FirstNameTh, nil
}

func GetLastName(ctx *fiber.Ctx) (string, error) {
	claim, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return "", ctx.Status(fiber.StatusUnauthorized).JSON(Result{
			Status:    fiber.StatusUnauthorized,
			Message:   "unauthorized",
			MessageTh: "ยืนยันตัวตนไม่สำเร็จ",
			Error:     "unauthorized",
		})
	}
	return claim.LastNameTh, nil
}

type Result struct {
	Error     interface{} `json:"error,omitempty"`
	Message   interface{} `json:"message,omitempty"`
	MessageTh interface{} `json:"message_th,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Total     int         `json:"total,omitempty"`
	Count     int         `json:"count,omitempty"`
	Status    int         `json:"status,omitempty"`
}
