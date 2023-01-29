package store

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
)

func OpenPostgresDB(dbName, address, username, password string, maxIdleConns, maxOpenConns, connMaxLifetimeHour int, gormLogger logger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  getPostgresDSN(dbName, address, username, password),
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		CreateBatchSize: 1000,
		Logger:          gormLogger,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	setSqlDBConfig(sqlDB, maxIdleConns, maxOpenConns, connMaxLifetimeHour)
	return db, nil
}

func getPostgresDSN(dbName, address, username, password string) string {
	addrPort := strings.Split(address, ":")
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		addrPort[0], username, password, dbName, addrPort[1])
}
