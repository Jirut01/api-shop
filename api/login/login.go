package login

import (
	"app-backend/middleware"
	"app-backend/model"
	validatormsg "app-backend/package/validator_msg"
	"app-backend/service"
	"fmt"
	"github.com/patcharp/golib/v2/helper"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/fasthash/fnv1a"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RegistorPayload struct {
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
	FirstNameTh  string `json:"first_name_th" validate:"required"`
	LastNameTh   string `json:"last_name_th" validate:"required"`
	FirstNameEng string `json:"first_name_eng"`
	LastNameEng  string `json:"last_name_eng"`
}

func RegistorUser(ctx *fiber.Ctx) error {
	payload := RegistorPayload{}
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "invalid request body",
			MessageTh: "เนื้อหาคำขอไม่ถูกต้อง",
			Error:     "bad request",
		})
	}

	var validate = validator.New()
	if err := validate.Struct(payload); err != nil {
		errorMsg := validatormsg.ValidatorMsg(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "missing required parameter" + errorMsg,
			MessageTh: "เนื้อหาคำขอไม่ครบถ้วน โปรดระบุ" + errorMsg,
			Error:     "bad request",
		})
	}
	var user model.Users
	if err := dbCtx().Model(&model.Users{}).Where("username=?", payload.Username).First(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusInternalServerError).JSON(Result{
				Status:    fiber.StatusInternalServerError,
				Message:   "internal server error",
				MessageTh: "ระบบมีปัญหา กรุณาติดต่อเจ้าหน้าที่",
				Error:     "internal server error",
			})
		} else {
			hashPassword := fnv1a.HashString64(payload.Password)
			user = model.Users{
				Username:     payload.Username,
				Password:     hashPassword,
				FirstNameTh:  payload.FirstNameTh,
				LastNameTh:   payload.LastNameTh,
				FirstNameEng: payload.FirstNameEng,
				LastNameEng:  payload.LastNameEng,
				Actived:      true,
			}
			if err := dbCtx().Model(&model.Users{}).Create(&user).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusInternalServerError).JSON(Result{
						Status:    fiber.StatusInternalServerError,
						Message:   "internal server error",
						MessageTh: "ระบบมีปัญหา กรุณาติดต่อเจ้าหน้าที่",
						Error:     "internal server error",
					})
				}
			}
		}
	} else {
		return ctx.Status(fiber.StatusConflict).JSON(Result{
			Status:    fiber.StatusConflict,
			Message:   "dupplicate username",
			MessageTh: "มีชื่อผู้ใช้งานในระบบแล้ว",
			Error:     "conflict",
		})
	}

	resp := map[string]interface{}{
		"username": user.Username,
	}
	return ctx.Status(fiber.StatusOK).JSON(Result{
		Status:    fiber.StatusOK,
		Message:   "success",
		MessageTh: "ลงทะเบียนสำเร็จ",
		Data:      resp,
	})
}

type LoginPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Login(ctx *fiber.Ctx) error {
	payload := LoginPayload{}
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "invalid request body",
			MessageTh: "เนื้อหาคำขอไม่ถูกต้อง",
			Error:     "bad request",
		})
	}

	if payload.Username == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "missing required parameter, please enter your username",
			MessageTh: "เนื้อหาคำขอไม่ครบถ้วน โปรดใส่ชื่อเข้าสู่ระบบ",
			Error:     "bad request",
		})
	}

	if payload.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "missing required parameter, please enter your password",
			MessageTh: "เนื้อหาคำขอไม่ครบถ้วน โปรดใส่รหัสผ่าน",
			Error:     "bad request",
		})
	}

	var countWrongPass int

	hashPassword := fnv1a.HashString64(payload.Password)
	var user model.Users
	if err := dbCtx().Model(&model.Users{}).Where("username=?", payload.Username).First(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusInternalServerError).JSON(Result{
				Status:    fiber.StatusInternalServerError,
				Message:   "internal server error",
				MessageTh: "ระบบมีปัญหา กรุณาติดต่อเจ้าหน้าที่",
				Error:     "internal server error",
			})
		} else {
			if err := service.CacheClient().Get(fmt.Sprintf("username:%s", payload.Username), &countWrongPass); err != nil {
				if !service.CacheClient().IsKeyNotFound(err) {
					logrus.Errorln("set cache check api error ->", err)
					return ctx.Status(fiber.StatusInternalServerError).JSON(Result{
						Status:    fiber.StatusInternalServerError,
						Message:   "internal server error",
						MessageTh: "ระบบมีปัญหา กรุณาติดต่อเจ้าหน้าที่",
						Error:     "internal server error",
					})
				} else {
					if err := service.CacheClient().Set(fmt.Sprintf("username:%s", payload.Username), 1, 5*time.Minute); err != nil {
						logrus.Errorln("set cache check api error ->", err)
					}
				}
			}
			if countWrongPass >= 5 {
				return ctx.Status(fiber.StatusForbidden).JSON(Result{
					Status:    fiber.StatusForbidden,
					Message:   "you have tried to login too many times.please try again in 5 minutes",
					MessageTh: "คุณป้อนรหัสผ่านผิดเกินจำนวนครั้งที่กำหนด กรุณาลองใหม่อีกครั้งในอีก 5 นาที",
					Error:     "internal server error",
				})
			} else {
				if countWrongPass < 4 {
					if err := service.CacheClient().Set(fmt.Sprintf("username:%s", payload.Username), countWrongPass+1, 5*time.Minute); err != nil {
						logrus.Errorln("set cache check api error ->", err)
					}
				}
				return ctx.Status(fiber.StatusBadRequest).JSON(Result{
					Status:    fiber.StatusBadRequest,
					Message:   "invalid username or password",
					MessageTh: "ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง กรุณาลองใหม่อีกครั้ง",
					Error:     "not found",
				})
			}
		}
	}
	if err := service.CacheClient().Get(fmt.Sprintf("username:%s", payload.Username), &countWrongPass); err != nil {
		if !service.CacheClient().IsKeyNotFound(err) {
			logrus.Errorln("get cache check api error ->", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(Result{
				Status:    fiber.StatusInternalServerError,
				Message:   "internal server error",
				MessageTh: "ระบบมีปัญหา กรุณาติดต่อเจ้าหน้าที่",
				Error:     "internal server error",
			})
		} else {
			if err := service.CacheClient().Set(fmt.Sprintf("username:%s", payload.Username), 1, 5*time.Minute); err != nil {
				logrus.Errorln("set cache check api error ->", err)
				// errline := line.SendMsgToLine(ctx.Method(), ctx.Path(), "NhsoRegistorEP", "set cache wrong pass error :"+err.Error())
				// if errline != nil {
				// 	logrus.Error(errline)
				// }
			}
		}
	}
	if !user.Actived {
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "invalid username or password",
			MessageTh: "ชื่อผู้ใช้ถูกระงับการใช้งานโปรดติดต่อเจ้าหน้าที่",
			Error:     "not found",
		})
	}
	if countWrongPass >= 5 {
		if user.LoginFail < 5 {
			if err := dbCtx().Model(&model.Users{}).Where("username=?", payload.Username).Updates(map[string]interface{}{
				"login_fail": user.LoginFail + 1,
			}).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusInternalServerError).JSON(Result{
						Status:    fiber.StatusInternalServerError,
						Message:   "internal server error",
						MessageTh: "ระบบมีปัญหา กรุณาติดต่อเจ้าหน้าที่",
						Error:     "internal server error",
					})
				}
			}
			return ctx.Status(fiber.StatusForbidden).JSON(Result{
				Status:    fiber.StatusForbidden,
				Message:   "you have tried to login too many times.please try again in 5 minutes",
				MessageTh: "คุณป้อนรหัสผ่านผิดเกินจำนวนครั้งที่กำหนด กรุณาลองใหม่อีกครั้งในอีก 5 นาที",
				Error:     "forbidden",
			})
		} else {
			if err := dbCtx().Model(&model.Users{}).Where("username=?", payload.Username).Updates(map[string]interface{}{
				"actived": false,
			}).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusInternalServerError).JSON(Result{
						Status:    fiber.StatusInternalServerError,
						Message:   "internal server error",
						MessageTh: "ระบบมีปัญหา กรุณาติดต่อเจ้าหน้าที่",
						Error:     "internal server error",
					})
				}
			}
			return ctx.Status(fiber.StatusBadRequest).JSON(Result{
				Status:    fiber.StatusBadRequest,
				Message:   "invalid username or password",
				MessageTh: "ชื่อผู้ใช้ถูกระงับการใช้งานโปรดติดต่อเจ้าหน้าที่",
				Error:     "not found",
			})
		}
	}

	if hashPassword != user.Password {
		if countWrongPass < 5 {
			if err := service.CacheClient().Set(fmt.Sprintf("username:%s", payload.Username), countWrongPass+1, 5*time.Minute); err != nil {
				logrus.Errorln("set cache check api error ->", err)
				// errline := line.SendMsgToLine(ctx.Method(), ctx.Path(), "NhsoRegistorEP", "set cache wrong pass error :"+err.Error())
				// if errline != nil {
				// 	logrus.Error(errline)
				// }
			}
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "invalid username or password",
			MessageTh: "ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง กรุณาลองใหม่อีกครั้ง",
			Error:     "not found",
		})
	} else {
		if err := service.CacheClient().Set(fmt.Sprintf("username:%s", payload.Username), 0, 5*time.Minute); err != nil {
			logrus.Errorln("set cache check api error ->", err)
			// errline := line.SendMsgToLine(ctx.Method(), ctx.Path(), "NhsoRegistorEP", "set cache wrong pass error :"+err.Error())
			// if errline != nil {
			// 	logrus.Error(errline)
			// }
		}
	}

	sessionId := uuid.NewV4().String()
	expired := time.Now().Add(time.Hour * 24)
	requestGentoken := middleware.RequestGenerateToken{
		FirstNameTh: user.FirstNameTh,
		LastNameTh:  user.LastNameTh,
		Username:    user.Username,
	}

	newToken, err := middleware.GenerateToken(sessionId, "user_token", requestGentoken, &expired)
	if err != nil {
		logrus.Errorln(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(Result{
			Status:    fiber.StatusInternalServerError,
			Message:   "internal server error",
			MessageTh: "ระบบมีปัญหา กรุณาติดต่อเจ้าหน้าที่",
			Error:     "internal server error",
		})
	}
	ctx.Cookie(&fiber.Cookie{
		Name:     "user_token",
		Value:    newToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Expires:  expired,
	})

	return ctx.Status(fiber.StatusOK).JSON(Result{
		Status:    fiber.StatusOK,
		Message:   "success",
		MessageTh: "ลงชื่อเข้าใช้สำเร็จ",
		Data: map[string]interface{}{
			"first_name_th": user.FirstNameTh,
			"last_name_th":  user.LastNameTh,
			"username":      user.Username,
		},
	})

}

func Logout(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name:     "user_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Expires:  time.Now().Add(1 * time.Microsecond),
	})
	return helper.HttpOk(ctx, "success")
}

func GetUser(ctx *fiber.Ctx) error {

	claim, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(Result{
			Status:    fiber.StatusUnauthorized,
			Message:   "unauthorized",
			MessageTh: "ยืนยันตัวตนไม่สำเร็จ",
			Error:     "unauthorized",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(Result{
		Status:    fiber.StatusOK,
		Message:   "success",
		Data: map[string]interface{}{
			"first_name_th": claim.FirstNameTh,
			"last_name_th":  claim.LastNameTh,
			"username":      claim.Username,
		},
	})
}
