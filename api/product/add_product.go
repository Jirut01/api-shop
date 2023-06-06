package product

import (
	"gorm.io/gorm"

	"app-backend/model"

	"github.com/gofiber/fiber/v2"
	"github.com/patcharp/golib/v2/helper"
	"github.com/sirupsen/logrus"
)

func AddProduct(ctx *fiber.Ctx) error {
	payload := model.Product{}

	if err := ctx.BodyParser(&payload); err != nil {
		return helper.HttpErrBadRequest(ctx)
	}

	if payload.Code == "" || payload.Name == "" || payload.Price == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "Invalid Request",
			MessageTh: "กรอกข้อมูลให้ครบถ้วน",
			Error:     "data requied",
		})
	}

	// c, _ := json.MarshalIndent(payload, "", "  ")
	// fmt.Println("payload:",string(c))

	if err := dbCtx().Model(&model.Product{}).Create(&payload).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			logrus.Errorln("create oproduct err ->", err)
			// errline := line.SendMsgToLine(ctx.Method(), ctx.Path(), "AppealFileDocUpload", "system error (create openland_storage_error err)- :"+err.Error())
			// if errline != nil {
			// 	logrus.Errorln(errline)
			// }
			return ctx.Status(fiber.StatusBadRequest).JSON(Result{
				Status:    fiber.StatusBadRequest,
				Message:   "Code Duplicate",
				MessageTh: "รหัสสินค้าซ้ำ",
				Error:     err,
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"status":    fiber.StatusOK,
		"message":   "success",
		"messageTh": "สำเร็จ",
	})
}
