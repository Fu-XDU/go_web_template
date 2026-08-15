package database

import "go_web_template/database/mysql"

func Setup() (err error) {
	err = mysql.Connect()
	return
}
