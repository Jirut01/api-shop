package product

import (
	"app-backend/model"
	"fmt"
	"time"

	helper "app-backend/helper"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func ImgProductUpload(ctx *fiber.Ctx) error {
	getUid := ctx.Query("uid")
	productCode := ctx.Query("product_code")
	if productCode == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "invalid uuid",
			MessageTh: "รหัสสินค้าไม่ถูกต้อง",
			Error:     "bad request",
		})
	}
	uid, err := uuid.FromString(getUid)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "invalid uuid",
			MessageTh: "รหัสสินค้าไม่ถูกต้อง",
			Error:     "bad request",
		})
	}
	var product model.Product
	if err := dbCtx().Model(&model.Product{}).Where("uid = ? and code = ?", uid, productCode).First(&product).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			logrus.Errorln("system error (get product)->", err)
			// errline := line.SendMsgToLine(ctx.Method(), ctx.Path(), "AppealFileDocUpload", "system error (get appeal)- :"+err.Error())
			// if errline != nil {
			// 	logrus.Errorln(errline)
			// }
			return ctx.Status(fiber.StatusInternalServerError).JSON(Result{
				Status:    fiber.StatusInternalServerError,
				Message:   "internal server error",
				MessageTh: "ระบบมีปัญหา กรุณาติดต่อเจ้าหน้าที่",
				Error:     "internal server error",
			})
		} else {
			return ctx.Status(fiber.StatusNotFound).JSON(Result{
				Status:    fiber.StatusNotFound,
				Message:   "uuid not found",
				MessageTh: "ไม่พบรหัสสินค้า",
				Error:     "not found",
			})
		}
	}

	if product.FileId != "" && product.FileName != "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "duplicate file",
			MessageTh: "มีรูปภาพในระบบแล้ว",
			Error:     "bad request",
		})
	}

	/*----------------------

		check file tupe

	----------------------*/

	file, err := ctx.FormFile("file")
	if err != nil {
		logrus.Errorln("request file err -->", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "invalid request file",
			MessageTh: "รูปภาพ ไม่ถูกต้อง",
			Error:     "bad request",
		})
	}

	obj, err := file.Open()
	if err != nil {
		logrus.Errorln("system error (open file)->", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(Result{
			Status:    fiber.StatusInternalServerError,
			Message:   "internal server error",
			MessageTh: "ระบบมีปัญหา กรุณาติดต่อเจ้าหน้าที่",
			Error:     "internal server error",
		})
	}
	defer obj.Close()
	buffer := make([]byte, file.Size)
	if _, err = obj.Read(buffer); err != nil {
		logrus.Errorln("system error (read file)->", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(Result{
			Status:    fiber.StatusInternalServerError,
			Message:   "internal server error",
			MessageTh: "ระบบมีปัญหา กรุณาติดต่อเจ้าหน้าที่",
			Error:     "internal server error",
		})
	}
	typeFile := mimetype.Detect(buffer).Extension()
	if typeFile != ".jpg" && typeFile != ".png" && typeFile != ".jpeg" {
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "jpg and png only",
			MessageTh: "สามารถใช้ไฟล์ .png และ .jpg ได้ทั้งนั้น",
		})
	}

	fileName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.NewV4().String(), typeFile)

	if err := helper.SaveFile(file, fileName, "assets/img/"); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(Result{
			Status:    fiber.StatusInternalServerError,
			Message:   "internal server error",
			MessageTh: "อัพโหลดรูปภาพไม่สำเร็จ",
			Error:     "internal server error",
		})
	} else {

		//===update product====

		if err := dbCtx().Model(&model.Product{}).Where("uid = ?", uid).Updates(map[string]interface{}{
			"file_id":   fileName,
			"file_name": file.Filename,
		}).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				logrus.Errorln("update product error ->", err)
				return ctx.Status(fiber.StatusInternalServerError).JSON(Result{
					Status:    fiber.StatusInternalServerError,
					Message:   "internal server error",
					MessageTh: "ระบบมีปัญหา กรุณาติดต่อเจ้าหน้าที่",
					Error:     "internal server error",
				})
			}
		}

	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"status":    fiber.StatusOK,
		"message":   "success",
		"messageTh": "อัพโหลดรูปภาพสำเร็จ",
	})
}

func ImgProductRemove(ctx *fiber.Ctx) error {
	getUid := ctx.Query("uid")
	uid, err := uuid.FromString(getUid)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "invalid uuid",
			MessageTh: "รหัสสินค้าไม่ถูกต้อง",
			Error:     "bad request",
		})
	}

	logrus.Info(uid)

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"status":    fiber.StatusOK,
		"message":   "success",
		"messageTh": "ลบรูปสำเร็จ",
	})
}
