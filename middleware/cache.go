package middleware

import (
	// "encoding/json"
	// "app-backend/package/line"
	"app-backend/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func ClearCheckApiCache() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Next()
		reqbody := c.Request().Body()
		if c.Response().StatusCode() != 200 {
			logrus.Info(string(reqbody))
		}
		switch c.Path() {
		case "/api/product/upload-img-product":
			if err := service.CacheClient().Del(fmt.Sprintf("username:%s", "ball")); err != nil {
				logrus.Errorln("del cache check api error ->", err)
				// errline := line.SendMsgToLine(c.Method(), c.Path(), "ClearCheckApiCache", "del cache check api error:"+err.Error())
				// if errline != nil {
				// 	logrus.Error(errline)
				// }
			}
			// case "/api/v1/opd/service-admissions/dt":
			// 	var payload PayloadvaccineServiceAdmission
			// 	if err := json.Unmarshal(reqbody, &payload); err != nil {
			// 		logrus.Error("unmarshal error ->", err)
			// 		return nil
			// 	}

			// 	for _, v := range payload.Vaccine {
			// 		if err := service.CacheClient().Del(fmt.Sprintf("check_dt:%s:%s", payload.Pid, v.Code)); err != nil {
			// 			logrus.Errorln("del cache check api error ->", err)
			// 			errline := line.SendMsgToLine(c.Method(), c.Path(), "ClearCheckApiCache", "del cache check api error:"+err.Error())
			// 			if errline != nil {
			// 				logrus.Error(errline)
			// 			}
			// 		}
			// 	}
			// case "/api/v1/opd/service-admissions/epi":
			// 	var payload PayloadvaccineServiceAdmission
			// 	if err := json.Unmarshal(reqbody, &payload); err != nil {
			// 		logrus.Error("unmarshal error ->", err)
			// 		return nil
			// 	}

			// 	for _, v := range payload.Vaccine {
			// 		if err := service.CacheClient().Del(fmt.Sprintf("check_epi:%s:%s", payload.Pid, v.Code)); err != nil {
			// 			logrus.Errorln("del cache check api error ->", err)
			// 			errline := line.SendMsgToLine(c.Method(), c.Path(), "ClearCheckApiCache", "del cache check api error:"+err.Error())
			// 			if errline != nil {
			// 				logrus.Error(errline)
			// 			}
			// 		}
			// 	}

			// case "/api/v1/auth/portal/appeal/doc?status=prepared":

			// 	if err := service.CacheClient().Del(fmt.Sprintf("check_appeal:%s", claimJWT.HospitalCode)); err != nil {
			// 		logrus.Errorln("del cache check appeal error ->", err)
			// 		errline := line.SendMsgToLine(c.Method(), c.Path(), "ClearCheckApiCache", "del cache check appeal error:"+err.Error())
			// 		if errline != nil {
			// 			logrus.Error(errline)
			// 		}
			// 	}
		}
		return nil
	}
}
