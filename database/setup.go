package database

import "github.com/Fu-XDU/go_web_template/database/mysql"

func Setup() (err error) {
	err = mysql.Connect()
	return
}
