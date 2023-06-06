package model

import (
	"github.com/patcharp/golib/v2/database"
	// uuid "github.com/satori/go.uuid"
)

type Product struct {
	database.Model
	Code        string `json:"code" gorm:"index:,unique"`
	Name        string `json:"name"`
	Price       float64 `json:"price"`
	Description string `json:"description,omitempty"`
	FileName      string `json:"file_name"`
	FileId      string `json:"file_id"`
}
