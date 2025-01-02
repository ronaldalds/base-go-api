package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type InitGORM struct {
	Host     string
	User     string
	Password string
	Database string
	Port     int
	TimeZone string
	Schema   string
	Models   []interface{}
}

type GormStore struct {
	*gorm.DB
}

func newGormStore(init *InitGORM) *GormStore {
	// Cria o Data Source Name (DSN)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s search_path=%s",
		init.Host,
		init.User,
		init.Password,
		init.Database,
		init.Port,
		init.TimeZone,
		init.Schema,
	)

	// Conecta ao banco de dados usando GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}

	// Obtém o banco de dados subjacente (*sql.DB)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get *sql.DB from GORM:", err)
	}

	// Cria um contexto com timeout de 10 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Tenta fazer um ping para garantir que a conexão está ativa
	err = sqlDB.PingContext(ctx)
	if err != nil {
		log.Fatal("failed to ping database:", err)
	}

	// Cria o schema, caso não exista
	_, err = sqlDB.ExecContext(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", init.Schema))
	if err != nil {
		log.Fatal("failed to create schema:", err)
	}

	if len(init.Models) > 0 {
		// Faz a migração do modelo
		err_db := db.AutoMigrate(
			init.Models...,
		)
		if err_db != nil {
			log.Fatal("failed to auto migrate:", err_db)
		}
	}

	fmt.Println("Connected to Gorm and schema is set up.")

	// Cria o GormStore
	gs := &GormStore{db}
	fmt.Printf("%s\n", gs.CreateAdmin().Error())

	// Retorna o GormStore
	return gs
}

func (gs *GormStore) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Access the underlying *sql.DB from GORM and ping it
	sqlDB, err := gs.DB.DB() // Obtém o *sql.DB subjacente do GORM
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db connection error: %v", err)
		log.Fatalf("db connection error: %v", err) // Log the error and terminate the program
		return stats
	}

	err = sqlDB.PingContext(ctx) // Realiza o ping no banco de dados
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db ping failed: %v", err)
		log.Fatalf("db ping failed: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Query for connection pool stats (PostgreSQL example)
	var dbStats struct {
		OpenConnections   int
		InUse             int
		Idle              int
		WaitCount         int64
		WaitDuration      time.Duration
		MaxIdleClosed     int64
		MaxLifetimeClosed int64
	}
	// You can write your own SQL query to fetch database stats
	sqlStats := `
		SELECT 
    		(SELECT count(*) FROM pg_stat_activity WHERE state = 'active') as open_connections,
    		(SELECT count(*) FROM pg_stat_activity WHERE state = 'idle') as idle,
    		(SELECT count(*) FROM pg_stat_activity WHERE wait_event IS NOT NULL) as wait_count
		`
	err = gs.DB.Raw(sqlStats).Scan(&dbStats).Error
	if err != nil {
		log.Printf("Failed to retrieve db stats: %v", err)
	}

	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse) // You can calculate in_use based on your needs
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String() // You can get a duration in some databases
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	return stats
}
