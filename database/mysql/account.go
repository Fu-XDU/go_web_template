package mysql

import (
	"go_web_template/model"
)

func InsertUser(user *model.User) (ID uint, err error) {
	err = db.Create(user).Error
	if err != nil {
		return
	}
	ID = user.ID
	return
}
