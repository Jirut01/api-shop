package model

import (
	"github.com/patcharp/golib/v2/database"
)

type Users struct {
	database.Model
	Username     string `json:"username"  gorm:"index"`
	Password     uint64 `json:"password"  gorm:"index"`
	FirstNameTh  string `json:"first_name_th"`
	LastNameTh   string `json:"last_name_th"`
	FirstNameEng string `json:"first_name_eng"`
	LastNameEng  string `json:"last_name_eng"`
	LoginFail    int    `json:"login_fail"`
	Actived      bool   `json:"actived"`
}