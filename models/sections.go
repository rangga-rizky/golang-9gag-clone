package models

import (
	"github.com/jinzhu/gorm"
)

type Section struct {
	gorm.Model
	Name string `json:"name"`
}

func GetSections() []*Section {

	var sections []*Section
	GetDB().Table("sections").Find(&sections)
	return sections
}
