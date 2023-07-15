package login

import (
	"app-backend/model"
	line "app-backend/package"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	// uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func GetDocUser(ctx *fiber.Ctx) error {

	userList := []model.Users{}
	if err := dbCtx().Model(&model.Users{}).Find(&userList).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			logrus.Errorln("find user err ->", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(Result{
				Status:    fiber.StatusInternalServerError,
				Message:   "internal server error",
				MessageTh: "ระบบมีปัญหา กรุณาติดต่อเจ้าหน้าที่",
				Error:     "internal server error",
			})
		}
	}

	//=============== gen excel =====================================

	fExcel := excelize.NewFile()
	sheetName := "sheet1"

	// Set Header
	columnHeader := 'A'
	header := []string{"ลำดับ", "ชื่อ-นามสกุล (ภาษาไทย)", "ชื่อ-นามสกุล (ภาษาอังกฤษ)", "สถานะ", "วันทีเข้าร่วม"}
	for i := 0; i < len(header); i++ {
		fExcel.SetCellValue(sheetName, fmt.Sprintf(`%s1`, string(columnHeader)), header[i])
		columnHeader++
	}

	// set column width
	fExcel.SetColWidth(sheetName, "A", string(columnHeader), 30)

	// style headder
	style, _ := fExcel.NewStyle(&excelize.Style{Alignment: &excelize.Alignment{Vertical: "center", Horizontal: "center"}, Font: &excelize.Font{Size: 14}})
	fExcel.SetCellStyle(sheetName, "A1", fmt.Sprintf(`%s1`, string(columnHeader)), style)

	bodyData := 'A'
	for i := 0; i < len(userList); i++ {
		bodyData = 'A'
		fExcel.SetCellValue(sheetName, fmt.Sprintf(`%s%d`, string(bodyData), i+2), i+1)
		bodyData++
		fExcel.SetCellValue(sheetName, fmt.Sprintf(`%s%d`, string(bodyData), i+2), fmt.Sprintf(`%s %s`, userList[i].FirstNameTh, userList[i].LastNameTh))
		bodyData++
		fExcel.SetCellValue(sheetName, fmt.Sprintf(`%s%d`, string(bodyData), i+2), fmt.Sprintf(`%s %s`, userList[i].FirstNameEng, userList[i].LastNameEng))
		bodyData++
		fExcel.SetCellValue(sheetName, fmt.Sprintf(`%s%d`, string(bodyData), i+2), userList[i].Actived)
		bodyData++
		fExcel.SetCellValue(sheetName, fmt.Sprintf(`%s%d`, string(bodyData), i+2), userList[i].CreatedAt.Local().Format("02/01/2006 15:04 น."))
	}

	// style body
	style, _ = fExcel.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 13}})
	fExcel.SetCellStyle(sheetName, "A2", fmt.Sprintf(`%s%d`, string(bodyData), len(userList)+1), style)
	buffer, err := fExcel.WriteToBuffer()
	if err != nil {
		errline := line.SendMsgToLine(ctx.Method(), ctx.Path(), "GetDocUser", "write file to buffer error :"+err.Error())
		if errline != nil {
			logrus.Errorln(errline)
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(Result{
			Status:    fiber.StatusInternalServerError,
			Message:   "system error (open file)",
			MessageTh: "ไม่สามารถเปิดไฟล์ได้",
			Error:     "internal server error",
		})
	}
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), "user_list")
	ctx.Set(fiber.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s.xlsx", fileName))
	return ctx.SendStream(buffer)
}
