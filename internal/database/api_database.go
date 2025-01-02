package database

import (
	"github.com/ronaldalds/base-go-api/internal/models"
	"github.com/ronaldalds/base-go-api/internal/settings"
)

type Database struct {
	GormStore  *GormStore
	RedisStore *RedisStore
}

var DB Database

func DbLoad() {
	dbSql := &InitGORM{
		Host:     settings.Env.SqlHost,
		User:     settings.Env.SqlUsername,
		Password: settings.Env.SqlPassword,
		Database: settings.Env.SqlDatabase,
		Port:     settings.Env.SqlPort,
		TimeZone: settings.Env.TimeZone,
		Schema:   settings.Env.SqlSchema,
		Models: []interface{}{
			&models.User{},
			&models.Role{},
			&models.Permission{},
		},
	}
	dbRedis := &InitRedis{
		Host:     settings.Env.RedisHost,
		Port:     settings.Env.RedisPort,
		Password: settings.Env.RedisPassword,
		DB:       settings.Env.RedisDb,
	}
	DB = Database{
		GormStore:  newGormStore(dbSql),
		RedisStore: newRedisStore(dbRedis),
	}
}
