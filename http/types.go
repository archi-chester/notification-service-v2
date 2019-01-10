package http

import (
	"net/http"
	"time"
)

// 	структура для фильрованного запроса пользователя
type insertRequest struct {
	Message  string `json:"message"`
	Subject  string `json:"subject"`
	UserFrom string `json:"user_from"`
	UserTo   string `json:"user_to"`
}

// 	структура для фильрованного запроса пользователя
type filterRequest struct {
	// PostId     string   `json:"_id"`
	DateAfter     *time.Time `json:"date_after"`
	DateBefore    *time.Time `json:"date_before"`
	MessageFilter string     `json:"message_filter"`
	SubjectFilter string     `json:"subject_filter"`
	Read          bool       `json:"read"`
	TypeFilter    int        `json:"type"`
	UserFrom      string     `json:"user_from"`
	UserTo        string     `json:"user_to"`
	Links         []Link     `json:"links"`
}

// 	структура для маршрута
type Route struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Method      string           `json:"method"`
	APIVersion  int              `json:"api_version"`
	Pattern     string           `json:"pattern"`
	HandlerFunc http.HandlerFunc `json:"-"`
}

// 	структура для ссылок
type Link struct {
	Type string `json:"type"`
	Text string `json:"text"`
	URL  string `json:"url"`
}
