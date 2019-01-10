Структура файла конфигурации:

    В папке /opt/notification-service файл ns.conf - json вида:
        
        {
        "ServerDB":"10.46.2.54",
        "PortDB":5432,
        "NameDB":"laskovost",
        "UserNameDB":"postgres",
        "PasswordDB":"postgres",
        "SslModeDB":"disable",
        "SchemaNameDB":"notification_service",
        "ListeningIP":"",
        "ListeningPort":10444
        }

    Если файла нет - сервис создаст пустой каркас.


Существующие пути:

	Получаем число непрочитанных:
	
	    Метод: GET,
		URL: /count?user=USER, где USER - для кого сообщения
		Возврат: Число,
		Функция: getCountNoticesByName

	Получаем список сообщений c учетом фильтра:
	
		Метод: GET,
		URL: /note?ПАРАМЕТРЫ, где ПАРАМЕТРЫ: subject_filter, message_filter,
		    date_before, date_after, user_to, user_from, read 
		    Если параметр опущен, он не применяется
		Возврат: json с полями - _id, date_create, date_read, opened, message, 
		    read, subject, user_from, user_to
		Функция: getNoticeByFilter,
		
	Проверка статуса сервиса сообщений:
	
		Метод: GET,
		URL: /health,
		Возврат: json вида - {'status': 'ok'}
		Функция: healthFunc,
		
	Добавляем сообщение:
	
		Метод: PUT,
		URL: /insert_notice, Параметры в теле json c полями: message, subject, 
		    user_from, user_to
		Возврат: json с полями - _id, date_create, date_read, opened, message, 
		    read, subject, user_from, user_to
		Функция: insertNotice,
		
	Изменение статуса сообщения - свернуто/развернуто:
	
		Метод: PATCH,
		URL: /switch_open_state/{note_id}, где {note_id} - ИД сообщения
		Функция: switchOpenState,
		
	Простановка у сообщения статуса прочитано и даты:
	
		Метод: PATCH,
		URL: /set_read/{note_id}, где {note_id} - ИД сообщения
		Функция: setRead,