package store

import (
	"database/sql"
	"fmt"
	"github.com/whj1990/go-core/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"time"
)

func NewReadWriteSeparationDB(gormLogger logger.Interface) (*gorm.DB, error) {
	return openReadWriteSeparationDB(
		config.GetNaCosString("db.name", ""),
		config.GetNaCosString("db.write.address", ""),
		config.GetNaCosString("db.write.username", ""),
		config.GetNaCosString("db.write.password", ""),
		config.GetNaCosString("db.read.address", ""),
		config.GetNaCosString("db.read.username", ""),
		config.GetNaCosString("db.read.password", ""),
		config.GetNaCosInt("db.maxIdleConns", 10),
		config.GetNaCosInt("db.maxOpenConns", 100),
		config.GetNaCosInt("db.connMaxLifetimeHour", 1),
		gormLogger,
	)
}

func NewDB(gormLogger logger.Interface) (*gorm.DB, error) {
	return OpenDB(
		config.GetNaCosString("db.name", ""),
		config.GetNaCosString("db.address", ""),
		config.GetNaCosString("db.username", ""),
		config.GetNaCosString("db.password", ""),
		config.GetNaCosInt("db.maxIdleConns", 10),
		config.GetNaCosInt("db.maxOpenConns", 100),
		config.GetNaCosInt("db.connMaxLifetimeHour", 1),
		gormLogger,
	)
}

func openReadWriteSeparationDB(dbName, writeAddress, writeUsername, writePassword, readAddress, readUsername, readPassword string,
	maxIdleConns, maxOpenConns, connMaxLifetimeHour int, gormLogger logger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(getDSN(dbName, writeAddress, writeUsername, writePassword)), &gorm.Config{
		CreateBatchSize: 1000,
		Logger:          gormLogger,
	})
	db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{mysql.Open(getDSN(dbName, readAddress, readUsername, readPassword))},
	}))
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

func OpenDB(dbName, address, username, password string, maxIdleConns, maxOpenConns, connMaxLifetimeHour int, gormLogger logger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(getDSN(dbName, address, username, password)), &gorm.Config{
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

func getDSN(dbName, address, username, password string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		username, password, address, dbName, true, "Local")
}

func setSqlDBConfig(sqlDB *sql.DB, maxIdleConns, maxOpenConns, connMaxLifetimeHour int) {
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(maxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(maxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetimeHour) * time.Hour)
}
