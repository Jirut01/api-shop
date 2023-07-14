package setting

import (
	// "github.com/gofiber/fiber/v2"
	"crypto/rsa"

	"github.com/gofiber/fiber/v2"
	"github.com/patcharp/golib/v2/cache"
	"github.com/patcharp/golib/v2/crypto"
	"github.com/patcharp/golib/v2/database"
	"github.com/patcharp/golib/v2/server"
	"github.com/patcharp/golib/v2/util"
	"github.com/patcharp/golib/v2/util/httputil"
	"github.com/sirupsen/logrus"
)

const (
	Dev  = "dev"
	Prod = "prod"
	Test = "test"
)

var (
	AppName   string
	BuildTime string
	atoi      = util.AtoI
	getEnv = util.GetEnv
	genSecret = crypto.GenSecretString
)

type Cfg struct {
	// App
	AppName   string
	BuildTime string
	Debug       bool
	Production  bool
	PodName     string

	// database
	Db database.Database
	Cache cache.Redis
	MongoDbUser      string
	MongoDbPasssword string

	MongoDbName             string
	MongoDbHost01           string
	MongoDbHost02           string
	MongoDbHost03           string
	MongoDbAuthSource       string
	MongoDbPort             string
	MongoDbDirectConnection string
	MongoDbTimeoutMS        string
	MongoDbAppName          string

	// Server
	Server server.Server
	PrivateKey *rsa.PrivateKey
	Secret     string

	//line
	LineToken            string

	//admin
	AdminUsername string
	AdminPassword string
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
	// App
	cfg.Production = getEnv("ENV", Prod) == Prod
	cfg.Debug = getEnv("DEBUG", "false") == "true"
	cfg.PodName = getEnv("MY_POD_NAME", "local")
	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Server
	fiberCfg := fiber.Config{
		Prefork:               getEnv("HTTP_PRE_FORK", "true") == "true",
		ServerHeader:          getEnv("HTTP_SERVER_HEADER", ""),
		ProxyHeader:           getEnv("HTTP_PROXY_HEADER", httputil.HeaderXForwardedFor),
		ReduceMemoryUsage:     getEnv("HTTP_REDUCE_MEMORY_USAGE", "true") == "true",
		DisableStartupMessage: getEnv("HTTP_DISABLE_STARTUP_MESSAGE", "true") == "true",
		BodyLimit:             25 * 1024 * 1024,
		ReadBufferSize:        16000,
	}


	cfg.Server = server.New(server.Config{
		Host: getEnv("SERVER_HOST", "0.0.0.0"),
		Port: getEnv("SERVER_PORT", "5000"),
		Config: &fiberCfg,
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

	cfg.Cache = cache.NewWithCfg(cache.Config{
		Host:     getEnv("CACHE_HOST", "127.0.0.1"),
		Port:     getEnv("CACHE_PORT", "6379"),
		Password: getEnv("CACHE_PASSWORD", ""),
		Db:       atoi(getEnv("CACHE_DB", ""), 0),
	})

	//line
	cfg.LineToken = getEnv("LINE_TOKEN", "")


	//Admin
	cfg.AdminUsername = getEnv("ADMIN_USERNAME", "")
	cfg.AdminPassword = getEnv("ADMIN_PASSWORD", "")


	cfg.Secret = getEnv("SERVER_SECRET", genSecret(32))

	//mongo
	cfg.MongoDbHost01 = getEnv("MONGO_DB_HOST_01", "")
	cfg.MongoDbHost02 = getEnv("MONGO_DB_HOST_02", "")
	cfg.MongoDbHost03 = getEnv("MONGO_DB_HOST_03", "")
	cfg.MongoDbHost03 = getEnv("MONGO_DB_HOST_03", "")
	cfg.MongoDbUser = getEnv("MONGO_DB_USER", "")
	cfg.MongoDbPasssword = getEnv("MONGO_DB_PASSWORD", "")
	cfg.MongoDbAuthSource = getEnv("MONGO_DB_AUTH_SOURCE", "")
	cfg.MongoDbName = getEnv("MONGO_DB_NAME", "")

	cfg.MongoDbPort = getEnv("MONGO_DB_PORT", "")
	cfg.MongoDbDirectConnection = getEnv("MONGO_DB_DIRECT_CONECTION", "")
	cfg.MongoDbTimeoutMS = getEnv("MONGO_DB_TIMEOUT_MS", "")
	cfg.MongoDbAppName = getEnv("MONGO_DB_APP_NAME", "")
	
	return nil
}
