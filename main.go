package main

import (
	"app-backend/service"
	"app-backend/setting"
	"app-backend/api"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	// "github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	a := kingpin.New(filepath.Base(os.Args[0]), fmt.Sprintf("%s %s", "ProgramName", "Version"))
	a.HelpFlag.Short('h')

	// Start
	startCmd := a.Command("start", "Start server command.")
	// migrateHospitalToken := a.Command("migrate_hopatal_token", "Start migrate")

	switch kingpin.MustParse(a.Parse(os.Args[1:])) {
	case startCmd.FullCommand():
		if err := cmdStartServer(); err != nil {
			logrus.Errorln("Start server error ->", err)
		}
		_ = cleanUpResources()
		logrus.Infoln("Server terminated")

	// case migrateHospitalToken.FullCommand():
	// 	if err := setting.GetCfg().Load(); err != nil {
	// 		logrus.Error("load configuration", err.Error())
	// 	}
	// 	// Initialize service
	// 	if err := service.InitDb(); err != nil {
	// 		logrus.Error("database connection", err.Error())
	// 	}
	// 	if err := service.MigrateDbSchema(); err != nil {
	// 		logrus.Error("migrate database schema", err.Error())
	// 	}
	// 	dbCtx := service.DbCtx
	// 	migrateToken(dbCtx)
	}


	// app := fiber.New()
	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello, World!")
	// })

	// app.Listen(":3000")
}


func cleanUpResources() error {
	logrus.Infoln("Closing database connection")
	if err := setting.GetCfg().Db.Close(); err != nil {
		return errors.New(fmt.Sprint("close database connection", err.Error()))
	}
	logrus.Infoln("Closing redis cache connection")
	if err := service.CacheClient().Client.Close(); err != nil {
		return errors.New(fmt.Sprint("close redis cache connection", err.Error()))
	}

	return nil
}

func cmdStartServer() error {
	// Load configuration

	if err := setting.GetCfg().Load(); err != nil {
		return errors.New(fmt.Sprint("load configuration", err.Error()))
	}
	
	// Initialize service
	if err := service.InitDb(); err != nil {
		return errors.New(fmt.Sprint("database connection", err.Error()))
	}
	if err := service.MigrateDbSchema(); err != nil {
		return errors.New(fmt.Sprint("migrate database schema", err.Error()))
	}
	if err := service.InitCache(); err != nil {
		return errors.New(fmt.Sprint("cache connection", err.Error()))
	}
	if err := service.InitMongo(); err != nil {
		return errors.New(fmt.Sprint("mongodb connection ->", err.Error()))
	}
	if err := service.InitRSA(); err != nil {
		return errors.New(fmt.Sprint("Init RSA key", err.Error()))
	}
	// Api router register
	if err := api.Register(); err != nil {
		return errors.New(fmt.Sprint("api route register", err.Error()))
	}
	
	logrus.Infoln("Server start", time.Now().Format(time.ANSIC))

	return setting.GetCfg().Server.Run()
}
