package db

import (
	"time"
)

var schemaName string = "notification_service"

type Notice struct {
	ID        string     `json:"_id" db:"id" gorm:"primary_key"`
	CreatedAt *time.Time `json:"date_create" db:"date_create" gorm:"column:date_create"`
	DateRead  *time.Time `json:"date_read" db:"date_read" gorm:"default:null"`
	Opened    bool       `json:"opened" db:"opened" gorm:"default:false"`
	Message   string     `json:"message" db:"message"`
	Read      bool       `json:"read" db:"read" gorm:"default:false"`
	Subject   string     `json:"subject" db:"subject"`
	Type      int        `json:"type" db:"type" gorm:"default:0"`
	UserFrom  string     `json:"user_from" db:"user_from"`
	UserTo    string     `json:"user_to" db:"user_to"`
}

// 	настройки БД
type SettingsDB struct {
	HostName string
	Port     int
	UserName string
	Password string
	DBName   string
	SslMode  string
	Schema   string
}

// set User's table name to be `profiles`
func (Notice) TableName() string {
	return schemaName + ".notifications"
}
