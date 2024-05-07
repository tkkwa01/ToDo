package mysql

import (
	"ToDo/domain"
	"ToDo/driver"
	"fmt"
)

func Init() {
	err := driver.GetRDB().AutoMigrate(
		&domain.User{},
		&domain.Task{},
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
