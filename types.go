package main

//	подключение к БД
const (
//	Файл настроек
// USERS_FILE_NAME = "users.conf"
//	роутинг
)

//	структура настроек
type settingsStruct struct {
	ServerDB          string
	PortDB            int
	NameDB            string
	UserNameDB        string
	PasswordDB        string
	SslModeDB         string
	SchemaNameDB      string
	ListeningIP       string
	ListeningPortHTTP int
	ListeningPortTCP  int
	GatewayIP         string
	XAPIKey           string
}

// 	структура сообщения
type notificationStruct struct {
	DateCreate string `json:"date_create"`
	DateRead   string `json:"date_read"`
	IsOpened   bool   `json:"is_opened"`
	Message    string `json:"message"`
	Read       bool   `json:"read"`
	Subject    string `json:"subject"`
	UserFrom   string `json:"user_from"`
	UserTo     string `json:"user_to"`
	Id         string `json:"_id"`
}
