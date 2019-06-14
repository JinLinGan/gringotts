package main

import (
	"github.com/jinlingan/gringotts/gringotts-server/cmd"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	cmd.Execute()
}
