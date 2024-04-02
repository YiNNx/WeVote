package models

import "gorm.io/gorm"

type Vote struct {
	gorm.Model
	Username string
	Count    int
}

