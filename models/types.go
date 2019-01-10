package models

import (
	"time"
)

// 	типы данных для БД
// 	Notification - структура для хранения и работы с сообщением
type Notification struct {
	ID         string     `json:"_id" db:"id"`
	DateCreate string     `json:"date_create" db:"date_create"`
	DateRead   *time.Time `json:"date_read" db:"date_read"`
	Opened     bool       `json:"opened" db:"opened"`
	Message    string     `json:"message" db:"message"`
	Read       bool       `json:"read" db:"read"`
	Subject    string     `json:"subject" db:"subject"`
	UserFrom   string     `json:"user_from" db:"user_from"`
	UserTo     string     `json:"user_to" db:"user_to"`
}

// 	NoticeRequest - структура для обработки запроса
type NoticeRequest struct {
	ID       string `json:"_id" db:"id"`
	Message  string `json:"message" db:"message"`
	Subject  string `json:"subject" db:"subject"`
	UserFrom string `json:"user_from" db:"user_from"`
	UserTo   string `json:"user_to" db:"user_to"`
}

// CountNotification - структура для хранения и работы с локальными группами
type CountNotification struct {
	Count string `json:"count" db:"count"`
}
