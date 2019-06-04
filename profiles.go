package users

import "github.com/jinzhu/gorm"

type Profile struct {
	gorm.Model
	Firstname  string `json:"firstname"`
	Middlename string `json:"middlename"`
	Lastname   string `json:"lastname"`
	Phone      string `json:"phone"`
}
