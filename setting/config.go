package setting

import (
	// "github.com/gofiber/fiber/v2"
	"github.com/patcharp/golib/v2/database"
	"github.com/patcharp/golib/v2/server"
	"github.com/patcharp/golib/v2/util"
)

var (
	AppName   string
	BuildTime string

	getEnv = util.GetEnv
)

type Cfg struct {
	// App
	AppName   string
	BuildTime string

	// database
	Db database.Database

	// Server
	Server server.Server

	//line
	LineToken            string
	FeedbackLineToken    string
	AppealBatchLineToken string
}

var theCfg *Cfg

func NewCfg() *Cfg {
	return &Cfg{
		AppName:   AppName,
		BuildTime: BuildTime,
	}
}

func GetCfg() *Cfg {
	if theCfg == nil {
		theCfg = NewCfg()
	}
	return theCfg
}

func (cfg *Cfg) Load() error {
	cfg.Server = server.New(server.Config{
		Host: getEnv("SERVER_HOST", "0.0.0.0"),
		Port: getEnv("SERVER_PORT", "5000"),
	})

	cfg.Db = database.NewWithConfig(
		database.Config{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			Username: getEnv("DB_USERNAME", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", ""),
		},
		database.DriverMySQL,
	)

	//line
	cfg.LineToken = getEnv("LINE_TOKEN", "")

	return nil
}
