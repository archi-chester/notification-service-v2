package db

import (
	"errors"
	"fmt"
	// 	log
	log "github.com/Sirupsen/logrus"

	// 	postgres
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// 	глобальные переменные
var DB *gorm.DB
var err error

// 	инициализация БД
func InitDB(db *gorm.DB, hostName string, port int, userName string, dbName string, password string, sslMode string, schema string) (*gorm.DB, error) {
	// 	присваиваем хост подключения
	if len(hostName) == 0 {
		log.Error("Пустой хост")
		return nil, errors.New("При формировании структуры настроек данных передали пустой хост")
	}

	// 	присваиваем порт поключения
	if port < 1 && port > 65535 {
		log.Error("Пустой порт")
		return nil, errors.New("При формировании структуры настроек данных передали пустой порт")
	}

	// 	присваиваем имя пользователя
	if len(userName) == 0 {
		log.Error("Пустое имя пользователя")
		return nil, errors.New("При формировании структуры настроек данных передали пустое имя пользователя")
	}

	// 	присваиваем название базы
	if len(dbName) == 0 {
		log.Error("Пустое название базы")
		return nil, errors.New("При формировании структуры настроек данных передали пустое название базы")
	}

	// 	присваиваем название базы
	if len(sslMode) == 0 {
		log.Error("Пустое название базы")
		return nil, errors.New("При формировании структуры настроек данных передали пустое название схемы")
	}

	// 	присваиваем название базы
	if len(schema) == 0 {
		log.Error("Пустое название базы")
		return nil, errors.New("При формировании структуры настроек данных передали пустое название схемы")
	}

	// 	Открываем БД
	db, err = gorm.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", hostName, port, userName, dbName, password, sslMode))
	if err != nil {
		log.Error("Не удалось подключиться к базе: ", err)
		return nil, err
	}

	// 	Создаем схему
	initSchema(db, schema)

	initTable(db)

	log.Info("Подключение к БД успешно")

	return db, err
}

// 	инициализация схемы
func initSchema(db *gorm.DB, schema string) error {
	var count int
	//	проверяю есть ли схема
	db.Table("pg_catalog.pg_namespace").Where("nspname = ?", schema).Count(&count)
	if count == 0 {
		// 	схемы нет - создаем
		db.Exec("CREATE SCHEMA " + schema)
		// db.Exec("SET search_path TO " + schema)
		log.Info("создали схему")
		// 	зациклили
		initSchema(db, schema)
	}

	// gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	// 	return schema + "." + defaultTableName
	// }

	return err
}

// 	инициализация схемы
func initTable(db *gorm.DB) error {
	// notice := Notice{
	// 	UID:        "123",
	// 	DateCreate: "21313",
	// 	DateRead:   nil,
	// 	Opened:     false,
	// 	Message:    "Message",
	// 	Read:       false,
	// 	Subject:    "Subject",
	// 	UserFrom:   "UserFrom",
	// 	UserTo:     "UserTo",
	// }
	var notices []Notice
	// var notice Notice

	// db.Find(&notices)
	// log.Warnf("Notices: %+v", notices)
	//	проверяю есть ли схема
	// fmt.Println("Hey")
	// db.CreateTable(&notice)
	// if count == 0 {
	// 	// 	схемы нет - создаем
	// 	db.Exec("CREATE SCHEMA " + schema)
	// 	db.Exec("SET search_path TO " + schema)
	// 	log.Info("создали схему")
	// 	// 	зациклили
	// 	initSchema(db, schema)
	// }

	db.Find(&notices)
	log.Warnf("%+v", notices)

	return err
}

// 	закрываем базу
func CloseDB(db *gorm.DB) error {
	//	вызов закрытия
	err = db.Close()
	// 	если ошибка - пишем в лог
	if err != nil {
		log.Error(err)
	}

	return err
}
