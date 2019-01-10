package http

var Routes = []Route{
	Route{
		"getCountNoticesByName",
		"Получаем число непрочитанных",
		"GET",
		1,
		"/count",
		getCountNoticesByName,
	},
	Route{
		"getNoticeByFilter",
		"Получаем список сообщений c учетом фильтра",
		"GET",
		1,
		"/note",
		getNoticeByFilter,
	},
	Route{
		"health",
		"Проверка статуса сервиса сообщений",
		"GET",
		1,
		"/health",
		healthFunc,
	},
	Route{
		"insertNotice",
		"Добавляем сообщение",
		"PUT",
		1,
		"/insert_notice",
		insertNotice,
	},
	Route{
		"switchOpenState",
		"Изменение статуса сообщения - свернуто/развернуто",
		"PATCH",
		1,
		"/switch_open_state/{note_id}",
		switchOpenState,
	},
	Route{
		"setRead",
		"Простановка у сообщения статуса прочитано и даты",
		"PATCH",
		1,
		"/set_read/{note_id}",
		setRead,
	},
}
