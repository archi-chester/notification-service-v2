package models

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

// 	создание БД
func CreateDB(storage *DB) {
	// 	создаем схему
	err := storage.CreateSchema(GetAllRegisteredTables())
	if err != nil {
		log.Error("Ошибка при создании новой схемы: ", err)
		os.Exit(1)
	}
}

// 	инициализация БД
func LoadDB(storage *DB, dbServerNameDB string, dbPort int, dbUserName string, dbPassword string, dbName string, dbSslMode string, dbSchemaName string) *DB {

	var err error

	// 	получаем экземпляр БД
	storage, err = InitDB(dbServerNameDB,
		dbPort,
		dbUserName,
		dbPassword,
		dbName,
		dbSslMode,
		dbSchemaName)
	if err != nil {
		log.Error("Ошибка подключения к БД: %s\n", err)
		os.Exit(1)
	} else {
		log.Info("Подключение к БД прошло успешно")
		return storage
	}
	return nil
}
