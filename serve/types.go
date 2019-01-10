package serve

// 	структура для фильрованного запроса пользователя
type insertRequest struct {
	Body     string `json:"body"`
	Subject  string `json:"subject"`
	Type     int    `json:"type"`
	UserFrom string `json:"user_from"`
	UserTo   string `json:"user_to"`
}
