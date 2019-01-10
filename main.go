package main

import (
	"github.com/jinzhu/gorm"
	"gitlab.havana/BDIO/notification-service-v2/db"
	"gitlab.havana/BDIO/notification-service-v2/http"
	"gitlab.havana/BDIO/notification-service-v2/serve"
	"gitlab.havana/BDIO/notification-service-v2/tcp"
)

// 	переменная для ошибки
// var err error
var version string

// 	Storage - экземпляр БД
// var Storage *models.DB
var data *gorm.DB

func main() {

	//	создаем/загружаем файл настроек
	loadSettings()

	// 	подключаемся к БД
	data, err := db.InitDB(data, mySettings.ServerDB, mySettings.PortDB, mySettings.UserNameDB, mySettings.NameDB, mySettings.PasswordDB, mySettings.SslModeDB, mySettings.SchemaNameDB)
	if err != nil {
		log.Error("Er", err)
		return
	} else {
		log.Info("Подключено к БД.")
	}
	defer db.CloseDB(data)

	// 	грузим базу
	// Storage = models.LoadDB(Storage,
	// 	mySettings.ServerDB,
	// 	mySettings.PortDB,
	// 	mySettings.UserNameDB,
	// 	mySettings.PasswordDB,
	// 	mySettings.NameDB,
	// 	mySettings.SslModeDB,
	// 	mySettings.SchemaNameDB)

	//	инитим аргументы из командной строки
	initPromptArgs()

	// Включаем роутер
	http.CreateWebServer(mySettings.ListeningPortHTTP, data)

	// 	Передаем экземпляр БД
	serve.InitDB(data)

	// 	Включаем прослушку TCP
	go tcp.Listener(mySettings.ListeningPortTCP)

	// Регистрируемся на сервере
	RegServiceOnGateway(mySettings.GatewayIP, mySettings.ListeningPortHTTP, mySettings.XAPIKey)

	logFile.Close()

}
