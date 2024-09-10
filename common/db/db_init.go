package db

import (
	"fmt"
	"sync"

	"github.com/n8sPxD/cowIM/common/db/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	client  *gorm.DB
	rwMutex sync.RWMutex
}

func MustNewMySQL(url string) *DB {
	db, err := gorm.Open(mysql.Open(url))
	if err != nil {
		panic("连接至MySQL失败: " + err.Error())
	}
	fmt.Println("连接至MySQL成功")
	return &DB{client: db}
}

func (db *DB) Migrate() error {
	err := db.client.AutoMigrate(
		&models.User{},
		&models.UserConfig{},
	)
	return err
}
