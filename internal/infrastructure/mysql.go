package infrastructure

import (
	"fmt"
	"github.com/IkezawaYuki/popple/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetMysqlConnection() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Env.DatabaseUser,
		config.Env.DatabasePass,
		config.Env.DatabaseHost,
		config.Env.DatabaseName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

type Transaction interface {
	Commit() error
	Rollback()
	GetTx() *gorm.DB
}

type Tx struct {
	tx *gorm.DB
}

func (t *Tx) Commit() error {
	return t.tx.Commit().Error
}

func (t *Tx) Rollback() {
	t.tx.Rollback()
}

func (t *Tx) GetTx() *gorm.DB {
	return t.tx
}
