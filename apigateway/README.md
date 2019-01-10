package apigateway

ФУНКЦИИ

```func Init(gatewayURL, serviceName, description, listenIP string, listenPort int, serviceKey string, repeatRegisterTimeout int, urls []map[string]string, monitoringURL string) error```

Init инициализирует модуль работы со шлюзом. gatewayURL - ссылка на API
    шлюза (например, "http://auth.laskovost/apigateway/api/v1" serviceName -
    название сервиса (например, "bdio") description - описание сервиса
    (например, "ПК БД ИО") listenIP - IP-адрес, по которому "слушает"
    процесс listenPort - номер порта, по которому "слушает"процесс servceKey
    - ключ процесса, выдается шлюзом API repeatRegisterTimeout - таймаут (в
    минутах) через которое повторять процесс регистрации urls - Блок урлов
    bsURLs представляет собой документ вида 
    ```[

        {
                "listen_url": "api/v1/proxy_test/",
                "target_url": "api_tester/api/v1/proxy_tests1/",
                "method": "GET",
                "resource_name": "test_one",
                "resource_info": "Тестраз"
        },
        {
                "listen_url": "/api/v1/proxy_tests2",
                "target_url": "/api/v1/proxy_tests2",
                "method": "GET",
                "resource_name": "test_two",
                "resource_info": "Тестдва"
        },
        ....................

    ]```     
monitoringURL - урл для мониторинга процесса со стороны шлюза (должен возвращаться 200 на GET-запрос шлюза)

```func RegisterOnAPIGateway() error```

RegisterOnAPIGateway Функция выполняет полный цикл регистрации на шлюзе

```func GetGroupsBySession(userSession string) (map[string]interface{}, error)```

GetGroupsBySession получает с шлюза информацию о группах пользователях и другие настройки

```func GetLoginFromSession(userSession string) (string, error)```

GetLoginFromSession получает логин пользователя из сессии

```func GetUserInfo(userSession string) (map[string]string, error)```

GetUserInfo получает информацию о пользователе для аудита
