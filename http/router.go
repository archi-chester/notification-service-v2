package http

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
)

// 	Экземпляр БД
var storage *gorm.DB

// создание роутера
func CreateRouter(r []Route) *mux.Router {
	// создаем экземпляр роутера
	router := mux.NewRouter().StrictSlash(true)

	log.Info("createRouter")
	var handler http.Handler
	//	создаем роутер

	for _, route := range r {
		// 	ссылка на функцию
		handler = route.HandlerFunc

		handler = AddHeaders(handler)
		router.
			Methods([]string{route.Method}...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	// log.Warnf("Router: %+v", router)
	return router
}

// Routes2Map Конвертирует пути в слайс мапов (для работы со шлюзом) + возвращает отдельный урл для хелсчека
func GetRouteMap() ([]map[string]string, string) {

	// Данные для регистрации урлов
	var resources []map[string]string
	healthRoute := ""
	for _, route := range Routes {
		res := map[string]string{
			"resource_name": route.Name,
			"target_url":    strings.Split(route.Pattern, "/{")[0],
			"listen_url":    strings.Split(route.Pattern, "/{")[0],
			"method":        route.Method,
			"resource_info": route.Description,
		}
		if route.Name == "health" {
			healthRoute = res["listen_url"]
		}

		resources = append(resources, res)
	}

	return resources, healthRoute
}

// AddHeaders adds all needed headers
func AddHeaders(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Here we adding headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
		// http.SetCookie(w, &http.Cookie{Name: "api_key", Value: app.GetAPIKey()})
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		inner.ServeHTTP(w, r)
		//log.Println(w.Header().Get("api_key"))
	})
}

// Serve запуск веб-сервера
func CreateWebServer(listeningPort int, db *gorm.DB) {
	// 	делаем экземпляр стораджа БД для пакета
	storage = db

	// log.Warnf("Storage: %+v", Storage)

	// 	создаем экземпляр роутера
	router := CreateRouter(Routes)

	// headersOk := handlers.AllowedHeaders([]string{"Content-Type", "X-Forwarded-For"})
	// originsOk := handlers.AllowedOrigins([]string{"*"})
	// methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
	log.Infof("Запуск сервера %s:%d", "", listeningPort)
	// err := http.ListenAndServe(fmt.Sprintf("%s:%d", "", listeningPort), handlers.CORS(headersOk, originsOk, methodsOk)(router))

	// 	стартуем сервер
	go http.ListenAndServe(fmt.Sprintf(":%d", listeningPort), router)
}
