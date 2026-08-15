package mysql

import (
	// "go_web_template/model"
	mingfumysql "github.com/Fu-XDU/mingfu_go_common/database/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() (err error) {
	db, err = mingfumysql.Connect(mingfumysql.NewConnOptionsFromFlags(), nil, initMysql)

	if err != nil {
		return
	}
	return
}

func initMysql(db *gorm.DB) (err error) {
	// Template: Initialize tables here
	// err = db.AutoMigrate(&model.User{})
	return
}
