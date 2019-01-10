package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/twinj/uuid"
	"gitlab.havana/BDIO/notification-service-v2/db"
	"strconv"
)

// 	обработка вызовов

//	все сообщения
func getCountNoticesByName(w http.ResponseWriter, r *http.Request) {
	// 	проверяем метод
	// 	это GET
	if r.Method == "GET" {
		log.Info("запрос GET")

		// 	получаем параметры
		tempStorage := storage
		// var count string
		var notices []db.Notice
		params := r.URL.Query()

		// 	имя получателя
		if userName, ok := params["user"]; ok {
			// 	добавляем фильтр
			tempStorage = tempStorage.Where("user_to = ? AND read=false", userName[0]).Find(&notices)
			log.Info("notice: ", len(notices))
			// tempStorage.Count(&count)

		}

		w.Write([]byte(fmt.Sprintf("%d", len(notices))))
		// 	// 	возвращаем нотисы
		// json.NewEncoder(w).Encode(count)
	}
}

//	все сообщения
func getNoticeByFilter(w http.ResponseWriter, r *http.Request) {
	// 	проверяем метод
	// 	это GET
	if r.Method == "GET" {
		log.Info("запрос GET")

		// 	задаем переменные
		var request filterRequest
		var notices []db.Notice
		var err error
		tempStorage := storage
		// var err error

		// 	получаем параметры
		params := r.URL.Query()

		// 	имя получателя
		if userTo, ok := params["user_to"]; ok {
			request.UserTo = userTo[0]
			// 	добавляем фильтр
			tempStorage = tempStorage.Where("user_to = ?", request.UserTo)
		}

		// 	имя отправителя
		if userFrom, ok := params["user_from"]; ok {
			request.UserFrom = userFrom[0]
			// 	добавляем фильтр
			tempStorage = tempStorage.Where("user_from = ?", request.UserFrom)
		}

		// 	начальная граница даты
		if dateAfter, ok := params["date_after"]; ok {
			time, err := time.Parse(time.RFC3339, dateAfter[0])
			if err != nil {
				log.Errorf("Ошибка конвертации времени dateAfter: %s", err)
			}
			request.DateAfter = &time

			// 	добавляем фильтр
			tempStorage = tempStorage.Where("date_create > ?", request.DateAfter)
		}

		// 	конечная граница даты
		if dateBefore, ok := params["date_before"]; ok {
			time, err := time.Parse(time.RFC3339, dateBefore[0])
			if err != nil {
				log.Errorf("Ошибка конвертации времени dateAfter: %s", err)
			}
			request.DateBefore = &time

			// 	добавляем фильтр
			tempStorage = tempStorage.Where("date_create < ?", request.DateBefore)

		}

		// 	начальная граница даты
		if messageFilter, ok := params["message_filter"]; ok {
			request.MessageFilter = messageFilter[0]
			// 	добавляем фильтр
			tempStorage = tempStorage.Where("message LIKE ?", "%"+request.MessageFilter+"%")
		}

		// 	конечная граница даты
		if subjectFilter, ok := params["subject_filter"]; ok {
			request.SubjectFilter = subjectFilter[0]

			// 	добавляем фильтр
			tempStorage = tempStorage.Where("subject LIKE ?", "%"+request.SubjectFilter+"%")
		}

		// 	бит прочитанности
		if read, ok := params["read"]; ok {
			if read[0] == "yes" {
				request.Read = true
			} else {
				request.Read = false
			}
			// 	добавляем фильтр
			tempStorage = tempStorage.Where("read = ?", request.Read)
		}

		// 	тип сообщения
		if typeFilter, ok := params["type"]; ok {
			// переводим в число строковый параметр
			request.TypeFilter, err = strconv.Atoi(typeFilter[0])
			// если не произошла ошибка при конвертации
			if err == nil {
				// 	добавляем фильтр
				tempStorage = tempStorage.Where("type = ?", request.TypeFilter)
			}
		}

		log.Warnf("Реквест: %+v", request)

		tempStorage.Order("date_create desc").Find(&notices)

		// 	// 	возвращаем нотисы
		json.NewEncoder(w).Encode(notices)
	}
}

//	вызов получение конкретного сообщения
func getNoticeByID(w http.ResponseWriter, r *http.Request) {
	// 	проверяем метод
	// // 	это GET
	// if r.Method == "GET" {
	// 	log.Info("запрос GET")

	// 	// 	задаем переменные
	// 	var requestID string
	// 	var notices []models.Notification

	// 	// 	получаем параметры
	// 	params := r.URL.Query()

	// 	// 	ID сообщения
	// 	if id, ok := params["note_id"]; ok {
	// 		requestID = id[0]
	// 	}

	// 	// 	транзакция
	// 	tx := storage.Begin()
	// 	// 	обработка ошибки
	// 	var err error
	// 	defer func() {
	// 		if err != nil {
	// 			log.Error("getNoticeByID прерван", err)
	// 			tx.Rollback()
	// 			return
	// 		}
	// 		tx.Commit()
	// 	}()

	// 	// 	получаем данные из селекта
	// 	if err = tx.Selectx(&notices, models.QueryNotificationsByID(requestID)); err == sql.ErrNoRows {
	// 		log.Error("Нет данных, ", err)
	// 	} else if err != nil {
	// 		log.Error(err)
	// 	}

	// 	// 	возвращаем нотисы
	// 	json.NewEncoder(w).Encode(notices)
	// }
}

//	вызов получения сообщений
func insertNotice(w http.ResponseWriter, r *http.Request) {

	log.Info("insertNotice")
	// 	проверяем метод
	// 	это PUT
	if r.Method == "PUT" {

		// 	задаем переменные
		var request insertRequest
		var newNotice db.Notice
		var notices []db.Notice
		tempStorage := storage

		// 	получаем параметры
		// params := r.URL.Query()
		log.Infof("%+v", r)
		//	//////////////////////////////////
		//	PARSE FORM
		// err := r.ParseForm()
		// if err != nil {
		// 	log.Error("Err")
		// }
		// log.Infof("PostForm: %+v", r.Form.Get("user_to"))
		// for key, values := range r.Form { // range over map
		// 	for _, value := range values { // range over []string
		// 		log.Infof("%+v %+v", key, value)
		// 	}
		// }
		///////////////////////////////////////////

		// 	читаем тело запроса в буфер
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("Невозможно прочитать тело запроса: %+v", err)
			http.Error(w, fmt.Sprintf("{\"status\":\"error\",\"error\":\"невозможно прочитать тело запроса: %s\"}", err), http.StatusUnprocessableEntity)
			return
		}

		//	анмаршалим буфер в структурку
		err = json.Unmarshal(bs, &request)
		if err != nil {
			log.Errorf("Неверный формат запроса: %+v", err)
			http.Error(w, fmt.Sprintf("{\"status\":\"error\",\"error\":\"неверный формат запроса: %s\"}", err), http.StatusUnprocessableEntity)
			return
		}

		// 	собираем запись для добавления
		// 	тема сообщения
		if len(request.Subject) > 0 {
			newNotice.Subject = request.Subject
		} else {
			log.Errorf("Пустая тема сообщения: %+v", err)
			http.Error(w, fmt.Sprintf("{\"status\":\"error\",\"error\":\"пустая тема сообщения: %s\"}", err), http.StatusUnprocessableEntity)
		}

		// 	тело сообщения
		if len(request.Message) > 0 {
			newNotice.Message = request.Message
		} else {
			log.Errorf("Пустое тело сообщения: %+v", err)
			http.Error(w, fmt.Sprintf("{\"status\":\"error\",\"error\":\"пустое тело сообщения: %s\"}", err), http.StatusUnprocessableEntity)
		}

		// 	имя получателя
		if len(request.UserTo) > 0 {
			newNotice.UserTo = request.UserTo
		} else {
			log.Errorf("Пустой получатель: %+v", err)
			http.Error(w, fmt.Sprintf("{\"status\":\"error\",\"error\":\"пустой получатель: %s\"}", err), http.StatusUnprocessableEntity)
		}

		// 	имя отправителя
		if len(request.UserFrom) > 0 {
			newNotice.UserFrom = request.UserFrom
		} else {
			log.Errorf("Пустой отправитель: %+v", err)
			http.Error(w, fmt.Sprintf("{\"status\":\"error\",\"error\":\"пустой отправитель: %s\"}", err), http.StatusUnprocessableEntity)
		}

		// UUID
		newNotice.ID = getNewUUID()

		//	добавление новой записи
		// 	проверяем есть ли запись, тру - запись новая
		log.Warnf("\n СТРОКА: %+v", newNotice)
		// 	добавляем запись
		tempStorage.Create(&newNotice)

		// 	возвращаем нотисы
		tempStorage.Find(&notices)
		json.NewEncoder(w).Encode(notices)
	}
}

// GetNewUUID - функция отдает новый UUID
func getNewUUID() string {
	return uuid.NewV4().String()
}

func healthFunc(w http.ResponseWriter, r *http.Request) {
	// log.Info("healthFunc")
	// 	вернули ок
	w.Write([]byte("{\"status\": \"ok\"}"))

}

// 	ответ при опросе
func initDBFunc(w http.ResponseWriter, r *http.Request) {
	// 	дергаем инициализацию таблицы
	// models.CreateDB(storage)

}

// 	изменение статуса сообщения - свернуто/развернуто
func switchOpenState(w http.ResponseWriter, r *http.Request) {
	// 	это PATCH
	if r.Method == "PATCH" {
		// 	переменные
		tempStorage := storage
		var notice db.Notice

		// 	получаем параметры
		vars := mux.Vars(r)

		// 	имя получателя
		if noteID, ok := vars["note_id"]; ok {
			tempStorage.Where("id = ?", noteID).First(&notice)
			// opened := !notice.Opened
			tempStorage.Model(&notice).Where("id = ?", noteID).Update("opened", !notice.Opened)
		}
	}
}

// 	простановка у сообщения статуса прочитано и даты
func setRead(w http.ResponseWriter, r *http.Request) {
	// 	это PATCH
	if r.Method == "PATCH" {
		// 	переменные
		tempStorage := storage
		var notice db.Notice

		// 	получаем параметры
		vars := mux.Vars(r)
		// var noteID string = vars["note_id"]

		// 	имя получателя
		if noteID, ok := vars["note_id"]; ok {
			// 	получаем текущее
			time := time.Now()
			// tempStorage = tempStorage.Where("id = ?", noteID)
			tempStorage.Model(&notice).Where("id = ?", noteID).Updates(map[string]interface{}{"read": true, "date_read": &time})
		}

		log.Info("UPD:%+v", notice)
	}
}
