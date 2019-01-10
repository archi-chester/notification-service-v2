package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

//	функции инициализации

//	константы инициализации
const (
	//	значение порта по умолчанию
	LISTENING_PORT_HTTP = 10432
	LISTENING_PORT_TCP  = 10433
	//	файл настроек по умолчанию
	SETTINGS_FILE_NAME = "ns.conf"
	//	домашняя папка сервиса
	HOME_DIR = "/opt/notification-service/"
	//	логи
	LOG_FILE_NAME = "notification-service.log"
	LOG_FILE_DIR  = "/opt/notification-service/"
	//
)

//	переменные инициализации

//	Какой порт слушаем
var listeningPortHTTP int
var listeningPortTCP int

//	Объявляем структурку настроек
var mySettings settingsStruct

//	Инициализация данных из командной строки
func initPromptArgs() {
	//	Проверяем сколько аргументов в командной строке

	fmt.Println(len(os.Args) > 1)
	if len(os.Args) > 1 {
		var port, err = strconv.Atoi(os.Args[1])
		if err != nil {
			//	Аргумент один - порт подключения берем по умолчанию
			listeningPortHTTP = LISTENING_PORT_HTTP
			listeningPortTCP = LISTENING_PORT_TCP
		} else {
			//	Аргумента два - первый : порт подкючения
			listeningPortHTTP = int(port)
			listeningPortTCP = int(port) + 1
			log.Warnf("Через командную строку передали порты: %d, %d", listeningPortHTTP, listeningPortTCP)
		}
	}
	// 	Если порта HTTP нет в файле
	if mySettings.ListeningPortHTTP != 0 {
		log.Info("Используем порт HTTP из файла конфигурации: ", mySettings.ListeningPortHTTP)
	} else {
		mySettings.ListeningPortHTTP = listeningPortHTTP
		log.Warn("В файле конфигурации порт не указан, используем: ", listeningPortHTTP)
	}
	// 	Если порта TCP нет в файле
	if mySettings.ListeningPortTCP != 0 {
		log.Info("Используем порт TCP из файла конфигурации: ", mySettings.ListeningPortTCP)
	} else {
		mySettings.ListeningPortTCP = mySettings.ListeningPortHTTP + 1
		log.Warn("В файле конфигурации порт не указан, используем: ", mySettings.ListeningPortHTTP+1)
	}
}

//	программа загрузки настроек
func loadSettings() {

	//	Читаем содержимое файла настроек
	myFile, err := os.Open(HOME_DIR + SETTINGS_FILE_NAME)
	if err != nil {

		log.Warn("Файла конфигурации нет - создаем")
		//	Файла нет - создаем
		myFile, err := os.OpenFile(HOME_DIR+SETTINGS_FILE_NAME, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Error(err)
			return
		}
		defer myFile.Close()

		//	Объявляем для примера один параметр
		mySettings.ListeningPortHTTP = LISTENING_PORT_HTTP
		//	Маршализируем
		buf, err := json.Marshal(mySettings)
		//	Копируем структурку в файлик
		myFile.Write(buf)

		return
	}
	//	Отложенно закрываем
	defer myFile.Close()

	// Получить размер файла
	stat, err := myFile.Stat()
	if err != nil {
		return
	}

	// Чтение файла
	buf := make([]byte, stat.Size())
	_, err = myFile.Read(buf)
	if err != nil {
		return
	}

	//	Маршалим прочитанное в структуру
	err = json.Unmarshal(buf, &mySettings)
	if err != nil {
		log.Error("Кривой анмаршал")
		return
	}

	//	Выводим сообщение
	log.Info("Процесс инициализации завершен")
	log.Warn(mySettings)
}
