package service

import (
	"app-backend/model"
	"app-backend/setting"
	"context"
	"fmt"
	"github.com/patcharp/golib/v2/crypto"
	"time"

	"github.com/patcharp/golib/v2/cache"
	"github.com/patcharp/golib/v2/util"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

var (
	dbCtx    *gorm.DB
	cacheCtx cache.Redis
	mongoCtx *mongo.Database
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
		&model.Users{},
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

func InitCache() error {
	logrus.Infoln("Init cache connection")
	if err := setting.GetCfg().Cache.Ping(); err != nil {
		return err
	}
	cacheCtx = setting.GetCfg().Cache
	return nil
}

func CacheClient() *cache.Redis {
	return &cacheCtx
}

func InitMongo() error {
	// mongodb://user:password@moph-claim2-mongodb-01, moph-claim2-mongodb-02, moph-claim2-mongodb-03/?authSource=admin
	var dsn string
	if util.GetEnv("ENV", "") == "local" {
		dsn = fmt.Sprintf("mongodb://%s:%s//?directConnection=%s&serverSelectionTimeoutMS=%s&appName=%s",
			setting.GetCfg().MongoDbHost01,
			setting.GetCfg().MongoDbPort,
			setting.GetCfg().MongoDbDirectConnection,
			setting.GetCfg().MongoDbTimeoutMS,
			setting.GetCfg().MongoDbAppName,
		)
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(dsn))
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	if err = client.Connect(ctx); err != nil {
		return err
	}
	mongoCtx = client.Database(setting.GetCfg().MongoDbName)
	logrus.Info("Init mongodb connection")
	return nil
}

func MongoDbClient() *mongo.Database {
	return mongoCtx
}

func InitRSA() error {
	if err := cacheCtx.Get("config:rsa", &setting.GetCfg().PrivateKey); err != nil {
		logrus.Infoln("No configured rsa found, generating new one")
		setting.GetCfg().PrivateKey, err = crypto.InitRSAKey(2048)
		if err != nil {
			return err
		}
		if err := cacheCtx.Set("config:rsa", setting.GetCfg().PrivateKey, 0); err != nil {
			return err
		}
	}
	return nil
}
