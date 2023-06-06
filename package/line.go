package line

import (
	"fmt"
	"app-backend/setting"

	"github.com/juunini/simple-go-line-notify/notify"
	"github.com/sirupsen/logrus"
)

func SendMsgToLine(medthod, path, fnc, msg string) error {
	if err := notify.SendText(setting.GetCfg().LineToken, fmt.Sprintf("Method : %s\nPath : %s\nFunction : %s\nMessage :%s", medthod, path, fnc, msg)); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// func SendFeedbackToLine(subject, detail, rating, name, hcode, hname, telephone_number, email string) error {
// 	if err := notify.SendText(setting.GetCfg().FeedbackLineToken, fmt.Sprintf("\nSubject: %s\nDetail: %s\nRating: %s\nName: %s\nHcode: %s\nHname: %s\nTelephone Number: %s\nEmail: %s", subject, detail, rating, name, hcode, hname, telephone_number, email)); err != nil {
// 		logrus.Error(err)
// 		return err
// 	}
// 	return nil
// }

// func SendAppealBatchToLine(date string, total int, msg interface{}) error {
// 	if err := notify.SendText(setting.GetCfg().AppealBatchLineToken, fmt.Sprintf("\nDate : %s\nTotal : %d\nDetail :%v", date, total, msg)); err != nil {
// 		logrus.Error(err)
// 		return err
// 	}
// 	return nil
// }
