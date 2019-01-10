package main

import (
	"os"
	"os/signal"
	"syscall"

	"gitlab.havana/BDIO/notification-service-v2/apigateway"
	"gitlab.havana/BDIO/notification-service-v2/http"
)

// 	регистрируемся на гейте
func RegServiceOnGateway(listenningIP string, listenningPort int, serviceKey string) {
	// 	приводим порт к числу
	// port, _ := strconv.Atoi(listenningPort)

	//	объявляем переменную для УРЛОВ
	var urls []map[string]string
	urls, healthRoute := http.GetRouteMap()
	log.Warnf("HealthPath: %+v", healthRoute)

	// url := map[string]string{
	// 	"resource_name": "note",
	// 	"target_url":    "/note",
	// 	"listen_url":    "/note",
	// 	"method":        "GET",
	// 	"resource_info": "route.Description",
	// }

	// url2 := map[string]string{
	// 	"resource_name": "insert",
	// 	"target_url":    "/insert_notice",
	// 	"listen_url":    "/insert_notice",
	// 	"method":        "PUT",
	// 	"resource_info": "route.Description",
	// }

	// url3 := map[string]string{
	// 	"resource_name": "init",
	// 	"target_url":    "/init_db",
	// 	"listen_url":    "/init_db",
	// 	"method":        "get",
	// 	"resource_info": "route.Description",
	// }

	// urls = append(urls, url)
	// urls = append(urls, url2)
	// urls = append(urls, url3)

	// 	разрегистрация
	// apigateway.Init("http://laska.havana/laskovost",
	// 	"notification_service",
	// 	"Сервис обмена сообщениями",
	// 	listenningIP,
	// 	int(listenningPort),
	// 	serviceKey,
	// 	30,
	// 	urls,
	// 	healthRoute,
	// 	"версия")
	// apigateway.Unregister()

	log.Info("Инициализация шлюза :", listenningPort, " serviceKey: ", serviceKey)
	// 	инитим
	err := apigateway.Init("http://laska.havana/laskovost",
		"notification_service",
		"Сервис обмена сообщениями",
		listenningIP,
		int(listenningPort),
		serviceKey,
		30,
		urls,
		healthRoute,
		version)

	if err != nil {
		log.Fatal("Ошибка инициализации модуля работы с API-шлюзом")
		return
	}
	// 	регимся
	err = apigateway.RegisterOnAPIGateway()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, os.Kill)
		<-c

		apigateway.Unregister()
		log.Info("Выход из программы")

		os.Exit(0)

	}()
	defer apigateway.Unregister()
}
