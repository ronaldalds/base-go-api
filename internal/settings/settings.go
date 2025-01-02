package settings

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ronaldalds/base-go-api/internal/i18n"
)

type EnvSettings struct {
	// SQL
	SqlUsername string
	SqlPassword string
	SqlHost     string
	SqlPort     int
	SqlDatabase string
	SqlSchema   string
	// Redis
	RedisDb       int
	RedisHost     string
	RedisPort     int
	RedisPassword string
	// JWT
	JwtSecret      string
	JwtExpireAcess time.Duration
	// APP
	TimeUCT  time.Location
	TimeZone string
	Port     int
	// SUPER USER
	SuperName     string
	SuperUsername string
	SuperPass     string
	SuperEmail    string
	SuperPhone    string
}

var Env EnvSettings

// Load reads and validates environment variables
func Load() {
	Env = EnvSettings{
		// SQL
		SqlUsername: getEnv("SQL_USERNAME", true),
		SqlPassword: getEnv("SQL_PASSWORD", true),
		SqlHost:     getEnv("SQL_HOST", true),
		SqlPort:     getEnvAsInt("SQL_PORT", true),
		SqlDatabase: getEnv("SQL_DATABASE", true),
		SqlSchema:   getEnv("SQL_SCHEMA", true),
		// Redis
		RedisDb:       getEnvAsInt("REDIS_DB", true),
		RedisHost:     getEnv("REDIS_HOST", true),
		RedisPort:     getEnvAsInt("REDIS_PORT", true),
		RedisPassword: getEnv("REDIS_PASSWORD", true),
		// JWT
		JwtSecret:      getEnv("JWT_SECRET", true),
		JwtExpireAcess: getEnvAsTime("JWT_EXPIRE_ACCESS", false, 10080),
		// APP
		TimeUCT:  getUCT("TIMEZONE", false, "America/Fortaleza"),
		TimeZone: getEnv("TIMEZONE", false, "America/Fortaleza"),
		Port:     getEnvAsInt("PORT", false, 3000),
		// SUPER USER
		SuperName:     getEnv("SUPER_NAME", true, "Admin"),
		SuperUsername: getEnv("SUPER_USERNAME", true, "admin"),
		SuperPass:     getEnv("SUPER_PASS", true, "admin"),
		SuperEmail:    getEnv("SUPER_EMAIL", true, "ronald.ralds@gmail"),
		SuperPhone:    getEnv("SUPER_PHONE", true, "+558892200365"),
	}
}

func getUCT(key string, required bool, defaultValue ...string) time.Location {
	value := os.Getenv(key)

	if value == "" {
		if required {
			panic(fmt.Sprintf("variable %s is required", key))
		}
		if len(defaultValue) > 0 {
			location, err := time.LoadLocation(value)
			if err != nil {
				panic(fmt.Sprintf("invalid timezone: %s", err.Error()))
			}
			return *location
		}
	}
	location, err := time.LoadLocation(value)
	if err != nil {
		panic(fmt.Sprintf("invalid timezone: %s", err.Error()))
	}
	return *location
}

func getEnv(key string, required bool, defaultValue ...string) string {
	value := os.Getenv(key)

	if value == "" {
		if required {
			panic(fmt.Sprintf(i18n.ERR_VARIABLE_IS_REQUIRED, key))
		}
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	return value
}

func getEnvAsInt(key string, required bool, defaultValue ...int) int {
	valueStr := os.Getenv(key)

	if valueStr == "" {
		if required {
			panic(fmt.Sprintf(i18n.ERR_VARIABLE_IS_REQUIRED, key))
		}
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic(fmt.Sprintf(i18n.ERR_CONVERT_TO_INTEGER, key, err))
	}
	return value
}

func getEnvAsTime(key string, required bool, defaultValue ...int) time.Duration {
	valueStr := os.Getenv(key)

	if valueStr == "" {
		if required {
			panic(fmt.Sprintf(i18n.ERR_VARIABLE_IS_REQUIRED, key))
		}
		if len(defaultValue) > 0 {
			return time.Duration(defaultValue[0]) * time.Minute
		}
		return 0
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic(fmt.Sprintf(i18n.ERR_CONVERT_TO_INTEGER, key, err))
	}
	return time.Duration(value) * time.Minute
}
