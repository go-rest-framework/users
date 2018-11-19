package users

import "github.com/jinzhu/gorm"

type Profile struct {
	gorm.Model
	Firstname  string
	Middlename string
	Lastname   string
	Phone      string
}
