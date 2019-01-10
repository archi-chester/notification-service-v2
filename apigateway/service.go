package apigateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
)

var settings struct {
	gatewayURL            string
	serviceName           string
	serviceKey            string
	description           string
	listenIP              string
	listenPort            int
	repeatRegisterTimeout int
	urls                  []map[string]string
	monitoringURL         string
	versionURL            string
	client                *http.Client
	apiKey                string
}

// ExternalSettings - структура для хранения папок и прочих настроек взаимодействия со сторонними ПК/объектами
type ExternalSettings struct {
	Path             string `json:"path"`
	SystemType       string `json:"system_type"`
	FolderType       string `json:"folder_type"`
	SZIRemoteAddress string `json:"szi_remote_address"`
	SZIRemoteFolder  string `json:"szi_remote_folder"`
	APILink          string `json:"api_link"`
}

// UserInfo - структура для хранения информации о пользователе
type UserInfo struct {
	Login       string   `json:"login"`
	FullName    string   `json:"name"`
	UserRole    string   `json:"user_role"`
	Permissions []string `json:"permissions"`
	// Organization string `json:"organization"`

}

// UserPermissions - дополнительные настройки пользателя
type UserPermissions struct {
	Permissions map[string]interface{}
	Login       string
	Status      string
}

// Init инициализирует модуль работы со шлюзом.
// gatewayURL - ссылка на API шлюза (например, "http://auth.laskovost/apigateway/api/v1"
// serviceName - название сервиса (например, "bdio")
// description - описание сервиса (например, "ПК БД ИО")
// listenIP - IP-адрес, по которому "слушает" процесс
// listenPort - номер порта, по которому "слушает"процесс
// servceKey - ключ процесса, выдается шлюзом API
// repeatRegisterTimeout - таймаут (в минутах) через которое повторять процесс регистрации
// urls - Блок урлов bsURLs представляет собой документ вида
//[
//	{
//		"listen_url": "api/v1/proxy_test/",
//		"target_url": "api_tester/api/v1/proxy_tests1/",
//		"method": "GET",
//		"resource_name": "test_one",
//		"resource_info": "Тестраз"
//	},
//	{
//		"listen_url": "/api/v1/proxy_tests2",
//		"target_url": "/api/v1/proxy_tests2",
//		"method": "GET",
//		"resource_name": "test_two",
//		"resource_info": "Тестдва"
//	},
//	....................
//]
func Init(gatewayURL, serviceName, description, listenIP string, listenPort int, serviceKey string, repeatRegisterTimeout int, urls []map[string]string, monitoringURL string, versionURL string) error {
	settings.gatewayURL = gatewayURL
	settings.serviceName = serviceName
	settings.serviceKey = serviceKey
	settings.description = description
	settings.listenPort = listenPort
	settings.listenIP = listenIP
	settings.repeatRegisterTimeout = repeatRegisterTimeout
	settings.monitoringURL = monitoringURL
	settings.versionURL = versionURL
	settings.urls = urls
	settings.client = new(http.Client)
	// settings.client.Jar, _ = cookiejar.New(nil)

	return nil
}

// RegisterService Функция регистрирует сервис на шлюзе авторизации
func registerService() error {

	serviceInfo := map[string]string{

		"service_name":        settings.serviceName,
		"service_host":        settings.listenIP,
		"service_port":        fmt.Sprintf("%d", settings.listenPort),
		"monitoring":          settings.monitoringURL,
		"version":             settings.versionURL,
		"service_description": settings.description,
	}

	regService := map[string]interface{}{
		"service_info": serviceInfo,
		"service_urls": settings.urls,
	}

	serviceURLs, _ := json.Marshal(regService)

	postfixURL := "/auth_backend/api/v1/register_service"
	fullURL, err := url.Parse(settings.gatewayURL + postfixURL)
	if err != nil {
		return fmt.Errorf("Невозможно разобрать url %s", settings.gatewayURL+postfixURL)
	}

	req, err := http.NewRequest("POST", fullURL.String(), bytes.NewBuffer(serviceURLs))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", settings.serviceKey)

	if err != nil {
		log.Error(err)
		return err
	}

	for i := 0; i < 3; i++ {
		log.Infof("Регистрация на шлюзе авторизации (попытка №%d/3): %s \n", i+1, fullURL.String())
		// Подготавливаем клиент, пихаем  в него куку с api-key

		// regCookie := http.Cookie{
		// 	Name:  "api_key",
		// 	Value: settings.serviceKey,
		// }

		// settings.client.Jar, _ = cookiejar.New(nil)
		// settings.client.Jar.SetCookies(fullURL, []*http.Cookie{&regCookie})
		// log.Debug("Cookies: ", settings.client.Jar.Cookies(fullURL))
		resp, err := settings.client.Do(req)
		if err != nil {
			if i < 3 {
				continue
			} else {
				return fmt.Errorf("Превышен предел допустимых попыток регистрации на шлюзе. Ошибка запроса на шлюз авторизации: %s", err)
			}
		}
		txt, _ := ioutil.ReadAll(resp.Body)
		log.Debug("Ответ от сервера: ", string(txt))
		resp.Body.Close()
		//log.Printf("Вернулся статус: %s", resp.Status)
		if resp.StatusCode != http.StatusOK {
			if i < 3 {
				continue
			}
			return fmt.Errorf("Регистрация сервиса не удалась. Статус: %s", resp.Status)
		}
		return nil

		// Получаем API-ключ и сохраняем его в настройках приложения

		// for _, c := range settings.client.Jar.Cookies(fullURL) {
		// 	if c.Name == "api_key" {
		// 		return c.Value, nil
		// 	}
		// }
	}

	// log.Println(settings.client.Jar.Cookies(fullURL))
	return fmt.Errorf("Регистрация сервиса на шлюзе не удалась. Проверьте, пожалуйста, свои настройки")

}

// RegisterOnAPIGateway Функция выполняет полный цикл регистрации на шлюзе
func RegisterOnAPIGateway() error {

	// В первый запуск регимся сразу же
	err := registerService()
	if err != nil {
		log.Error(err)
	}
	// settings.apiKey = apiKey

	// А далее выполняем ту же процедуру каждую минуту
	for range time.Tick(time.Duration(settings.repeatRegisterTimeout) * time.Second) {

		err := registerService()
		if err != nil {
			log.Error(err)
		}
		// settings.apiKey = apiKey
	}

	return nil
}

// // GetExternalSettings - получение настроек папок с шлюза
// func GetExternalSettings() ([]ExternalSettings, error) {
// 	var fldrs []ExternalSettings
// 	postfixURL := "/monitoring/api/v1/shared_folders"
// 	fullURL, err := url.Parse(settings.gatewayURL + postfixURL)
// 	if err != nil {
// 		return fldrs, fmt.Errorf("Невозможно разобрать url %s", settings.gatewayURL+postfixURL)
// 	}

// 	req, err := http.NewRequest("GET", fullURL.String(), nil)
// 	req.Header.Set("X-API-KEY", settings.serviceKey)
// 	resp, err := settings.client.Do(req)

// 	if err != nil {
// 		return nil, fmt.Errorf("Ошибка запроса настроек папок со шлюза: %s", err)
// 	}
// 	//log.Printf("Вернулся статус: %s", resp.Status)
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("Не удалось получить настройки внешних папок от шлюза. Статус: %s", resp.Status)
// 	}

// 	bs, err := ioutil.ReadAll(resp.Body)
// 	defer resp.Body.Close()
// 	//log.Println("Ответ шлюза: ", string(bs))
// 	err = json.Unmarshal(bs, &fldrs)
// 	if err != nil {
// 		return nil, fmt.Errorf("Невозможно распаковать информацию от шлюза в разрешения")
// 	}

// 	return fldrs, nil

// }

// GetLoginFromSession получает логин пользователя из сессии
func GetLoginFromSession(userSession string) (string, error) {
	postfixURL := "/user_session_info/api/v1/check_session"
	fullURL, err := url.Parse(settings.gatewayURL + postfixURL)
	if err != nil {
		return "", fmt.Errorf("Невозможно разобрать url %s", settings.gatewayURL+postfixURL)
	}

	req, err := http.NewRequest("GET", fullURL.String(), nil)
	userCookie := &http.Cookie{
		Name:  "session",
		Value: userSession,
	}
	req.AddCookie(userCookie)
	req.Header.Set("X-API-KEY", settings.serviceKey)

	//req.Header.Set("Content-Type", "application/json")

	resp, err := settings.client.Do(req)

	if err != nil {
		return "", fmt.Errorf("Ошибка запроса проверки валидности сессии на шлюз авторизации: %s", err)
	}
	//log.Printf("Вернулся статус: %s", resp.Status)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Проверка валидности сессии пользователя на шлюзе авторизации не удалась. Статус: %s", resp.Status)
	}

	bs, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	//log.Println(string(bs))
	var login struct {
		Login string `json:"login"`
	}

	err = json.Unmarshal(bs, &login)
	if err != nil {
		return "", fmt.Errorf("Невозможно распаковать ответ сервера. %s", err)
	}

	return login.Login, nil
}

// GetUserInfo получает информацию о пользователе для аудита
func GetUserInfo(userSession string) (UserInfo, error) {
	var ui UserInfo
	log.Debug("Получаем с сервера информацию о пользователе для аудита")
	postfixURL := "/user_session_info/api/v1/user_info"
	fullURL, err := url.Parse(settings.gatewayURL + postfixURL)
	log.Debugf("Путь получения инфы: %s", fullURL.String())
	if err != nil {
		return ui, fmt.Errorf("Невозможно разобрать url %s", settings.gatewayURL+postfixURL)
	}

	req, err := http.NewRequest("GET", fullURL.String(), nil)
	userCookie := &http.Cookie{
		Name:  "session",
		Value: userSession,
	}
	req.AddCookie(userCookie)
	req.Header.Set("X-API-KEY", settings.serviceKey)

	//req.Header.Set("Content-Type", "application/json")

	resp, err := settings.client.Do(req)

	if err != nil {
		return ui, fmt.Errorf("Ошибка запроса проверки валидности сессии на шлюз авторизации: %s", err)
	}
	//log.Printf("Вернулся статус: %s", resp.Status)
	if resp.StatusCode != http.StatusOK {
		return ui, fmt.Errorf("Проверка валидности сессии пользователя на шлюзе авторизации не удалась. Статус: %s", resp.Status)
	}

	bs, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	//log.Println(string(bs))
	err = json.Unmarshal(bs, &ui)
	if err != nil {
		return ui, fmt.Errorf("Невозможно распаковать ответ сервера. %s", err)
	}

	log.Debug("Получена следующая информация: ", ui)

	return ui, nil
}

// GetGroupsBySession получает с шлюза информацию о группах пользователях и другие настройки
func GetGroupsBySession(userSession string) (UserPermissions, error) {
	var up UserPermissions
	postfixURL := "/user_session_info/api/v1/get_trbd_groups"
	fullURL, err := url.Parse(settings.gatewayURL + postfixURL)
	if err != nil {
		return up, fmt.Errorf("Невозможно разобрать url %s", settings.gatewayURL+postfixURL)
	}

	req, err := http.NewRequest("GET", fullURL.String(), nil)

	userCookie := &http.Cookie{
		Name:  "session",
		Value: userSession,
	}
	req.AddCookie(userCookie)
	req.Header.Set("X-API-KEY", settings.serviceKey)

	//req.Header.Set("Content-Type", "application/json")
	log.Debugf("ломимся на урл %s с куками: %v", fullURL, req.Cookies())
	resp, err := settings.client.Do(req)

	if err != nil {
		log.Errorf("Ошибка запроса групп пользователя на шлюз: %s", err)
		return up, fmt.Errorf("Ошибка запроса групп пользователя на шлюз: %s", err)
	}
	//log.Printf("Вернулся статус: %s", resp.Status)

	bs, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Не удалось получить группы пользователя с шлюза. Статус: %s\nТело запроса: %s", resp.Status, string(bs))
		return up, fmt.Errorf("Не удалось получить группы пользователя с шлюза. Статус: %s\nТело запроса: %s", resp.Status, string(bs))
	}

	//log.Println("Ответ шлюза: ", string(bs))
	err = json.Unmarshal(bs, &up)
	if err != nil {
		return up, fmt.Errorf("Невозможно распаковать информацию от шлюза в разрешения")
	}

	return up, nil
}

// GetAPIKey - выдает полученный API-ключ приложения
func GetAPIKey() string {
	return settings.apiKey
}

// Unregister - дерегистрация сервиса из шлюза
func Unregister() error {
	log.Info("Отзываем регистрацию сервиса")
	postfixURL := "/monitoring/api/v1/dead_services"
	fullURL, err := url.Parse(settings.gatewayURL + postfixURL)
	if err != nil {
		return fmt.Errorf("Невозможно разобрать url %s", settings.gatewayURL+postfixURL)
	}

	var unregisterReq = []struct {
		ServiceName string `json:"service_name"`
		Instance    string `json:"instance"`
	}{
		{
			ServiceName: settings.serviceName,
			Instance:    fmt.Sprintf("%s:%d", settings.listenIP, settings.listenPort),
		},
	}

	request, _ := json.Marshal(unregisterReq)

	req, err := http.NewRequest("DELETE", fullURL.String(), bytes.NewBuffer(request))

	if err != nil {
		log.Error(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", settings.serviceKey)
	// regCookie := http.Cookie{
	// 	Name:  "api_key",
	// 	Value: settings.apiKey,
	// }

	// settings.client.Jar, _ = cookiejar.New(nil)
	// settings.client.Jar.SetCookies(fullURL, []*http.Cookie{&regCookie})
	// log.Debug("Cookies: ", settings.client.Jar.Cookies(fullURL))
	_, err = settings.client.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil

}

// MakeExternalCall - вызов внешнего урла через шлюз
func MakeExternalCall(method string, urla string, body []byte) (*http.Response, error) {
	var b *bytes.Reader
	if body == nil {
		b = nil
	} else {
		b = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, urla, b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-KEY", settings.serviceKey)

	// settings.client.Jar, _ = cookiejar.New(nil)
	return settings.client.Do(req)
}

// SendTRBDRequest - вызов внешнего урла через шлюз
// func SendTRBDRequest(method string, tableURL string, body io.Reader, userSession string) (int, []byte, http.Header, error) {
// 	gatewayTRBDURL := "/trbd/api/v1"

// 	// Обработка урла ТРБД
// 	fullURL, err := url.Parse(settings.gatewayURL + gatewayTRBDURL + tableURL)
// 	log.Debug("FullURL: " + fullURL.String())

// 	if err != nil {
// 		return 0, nil, nil, err
// 	}
// 	req, err := http.NewRequest(method, fullURL.String(), body)
// 	if err != nil {
// 		log.Debugf("Ошибка: %s", err)
// 	}
// 	userCookie := &http.Cookie{
// 		Name:  "session",
// 		Value: userSession,
// 	}
// 	req.AddCookie(userCookie)
// 	req.Header.Set("X-API-KEY", settings.serviceKey)

// 	//req.Header.Set("Content-Type", "application/json")

// 	resp, err := settings.client.Do(req)
// 	if err != nil {
// 		return 0, nil, nil, err
// 	}
// 	bs, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return 0, nil, nil, err
// 	}
// 	defer resp.Body.Close()

// 	return resp.StatusCode, bs, resp.Header, nil
// }
