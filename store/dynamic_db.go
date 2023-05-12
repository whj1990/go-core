package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/whj1990/go-core/config"
	"github.com/whj1990/go-core/encrypt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DynamicDbItem struct {
	db         *gorm.DB
	expiration int64
}

type DynamicDb struct {
	items             map[string]*DynamicDbItem
	mu                sync.RWMutex
	cacheHours        int
	gcIntervalMinutes int
	gormLogger        logger.Interface
}

func (db *DynamicDb) Get(t, address, username, password, dbName string) (*gorm.DB, error) {
	cachedDb := db.getCachedDbWithRLock(t, address, username, password, dbName)
	if cachedDb != nil {
		return cachedDb, nil
	}
	return db.add(t, address, username, password, dbName)
}

func (db *DynamicDb) add(t, address, username, password, dbName string) (*gorm.DB, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	cachedDb := db.getCachedDb(t, address, username, password, dbName)
	if cachedDb != nil {
		return cachedDb, nil
	}
	var newDb *gorm.DB
	var err error
	switch t {
	case "mysql":
		newDb, err = OpenDB(
			dbName,
			address,
			username,
			password,
			config.GetNacosConfigData().Db.MaxIdleConnects,
			config.GetNacosConfigData().Db.MaxOpenConnects,
			config.GetNacosConfigData().Db.ConnMaxLifetimeHour,
			db.gormLogger,
		)
	case "postgresql":
		newDb, err = OpenPostgresDB(
			dbName,
			address,
			username,
			password,
			config.GetNacosConfigData().Db.MaxIdleConnects,
			config.GetNacosConfigData().Db.MaxOpenConnects,
			config.GetNacosConfigData().Db.ConnMaxLifetimeHour,
			db.gormLogger,
		)
	}
	if err != nil {
		return nil, err
	}
	db.items[getItemKey(t, address, username, password, dbName)] = &DynamicDbItem{
		newDb,
		time.Now().Add(time.Duration(db.cacheHours) * time.Hour).UnixNano(),
	}
	return newDb, nil
}

func (db *DynamicDb) getCachedDb(t, address, username, password, dbName string) *gorm.DB {
	item := db.items[getItemKey(t, address, username, password, dbName)]
	if item != nil {
		item.expiration = time.Now().Add(time.Duration(db.cacheHours) * time.Hour).UnixNano()
		return item.db
	}
	return nil
}

func (db *DynamicDb) getCachedDbWithRLock(t, address, username, password, dbName string) *gorm.DB {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.getCachedDb(t, address, username, password, dbName)
}

func (db *DynamicDb) gcLoop() {
	ticker := time.NewTicker(time.Duration(db.gcIntervalMinutes) * time.Minute)
	for {
		select {
		case <-ticker.C:
			db.closeAndDeleteExpired()
		}
	}
}

func (db *DynamicDb) closeAndDeleteExpired() {
	db.mu.Lock()
	defer db.mu.Unlock()

	now := time.Now().UnixNano()
	zap.L().Warn("开始动态db gc")
	for k, item := range db.items {
		if now > item.expiration {
			sqlDb, err := item.db.DB()
			if err != nil {
				zap.L().Warn("获取动态db失败")
			}
			if err = sqlDb.Close(); err != nil {
				zap.L().Warn("关闭动态db失败")
			}

			delete(db.items, k)
		}
	}
}

func NewDynamicDb(gormLogger logger.Interface) *DynamicDb {
	db := &DynamicDb{
		items:             map[string]*DynamicDbItem{},
		cacheHours:        72,
		gcIntervalMinutes: 5,
		gormLogger:        gormLogger,
	}
	go db.gcLoop()
	return db
}

func getItemKey(t, address, username, password, dbName string) string {
	return encrypt.Md5Encrypt(fmt.Sprintf("%s-%s-%s-%s-%s", t, address, username, password, dbName))
}
