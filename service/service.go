package service

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"app-backend/setting"
	"app-backend/model"
)

var (
	dbCtx    *gorm.DB
)

func InitDb() error {
	logrus.Infoln("Init database connection")
	if err := setting.GetCfg().Db.ConnectWithGormConfig(gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}); err != nil {
		return err
	}
	dbCtx = setting.GetCfg().Db.Ctx()
	return nil
}

func MigrateDbSchema() error {
	tables := []interface{}{
		&model.Product{},
	}
	if len(tables) > 0 {
		if err := setting.GetCfg().Db.MigrateDatabase(tables); err != nil {
			return err
		}
	}
	return nil
}

func DbCtx() *gorm.DB {
	return dbCtx
}
