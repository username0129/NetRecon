package model

import "gorm.io/gorm"

type Initializable interface {
	InsertData(db *gorm.DB) error
}
