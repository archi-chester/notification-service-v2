package serve

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"github.com/twinj/uuid"
	"gitlab.havana/BDIO/notification-service-v2/db"
)

// 	Экземпляр БД
var storage *gorm.DB

func InitDB(db *gorm.DB) {
	// 	делаем экземпляр стораджа БД для пакета
	storage = db
}

// 	добавляем сообщение
func InsertNotice(data []byte) ([]byte, error) {
	// 	переменные
	var request insertRequest
	var newNotice db.Notice
	var notices []db.Notice
	tempStorage := storage

	//	анмаршалим запрос в структурку
	err := json.Unmarshal(data, &request)
	if err != nil {
		log.Errorf("Неверный формат запроса: %+v", err)
		return nil, err
	}
	// 	собираем запись для добавления
	// 	тема сообщения
	if len(request.Subject) > 0 {
		newNotice.Subject = request.Subject
	} else {
		log.Errorf("Пустая тема сообщения: %+v", err)
		return nil, err
	}

	// 	тело сообщения
	if len(request.Body) > 0 {
		newNotice.Message = request.Body
	} else {
		log.Errorf("Пустое тело сообщения: %+v", err)
		return nil, err
	}

	// 	имя получателя
	if len(request.UserTo) > 0 {
		newNotice.UserTo = request.UserTo
	} else {
		log.Errorf("Пустой получатель: %+v", err)
		return nil, err
	}

	// 	имя отправителя
	if len(request.UserFrom) > 0 {
		newNotice.UserFrom = request.UserFrom
	} else {
		log.Errorf("Пустой отправитель: %+v", err)
		return nil, err
	}

	// 	тип сообщения
	if request.Type > 0 {
		newNotice.Type = request.Type
	} else {
		log.Errorf("Пустой тип сообщения: %+v", err)
		return nil, err
	}

	// newNotice.Type = request.Type

	// UUID
	newNotice.ID = getNewUUID()

	log.Infof("%+v", newNotice)
	//	добавление новой записи
	// 	добавляем запись
	tempStorage.Create(&newNotice)

	// 	возвращаем нотисы
	tempStorage.Find(&notices)

	// 	помещаем ответ в массив символов для возврата
	data, err = json.Marshal(notices)

	return data, err
}

// GetNewUUID - функция отдает новый UUID
func getNewUUID() string {
	return uuid.NewV4().String()
}
