package apigateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SetObjectlistResponsible - функция обрабатывает событие2 "Назначение ответственного"
func SetObjectlistResponsible(bsEventRequest []byte, userSession string) error {
	var eR = make(map[string]interface{})

	err := json.Unmarshal(bsEventRequest, &eR)
	if err != nil {
		return fmt.Errorf("ошибка распаковки документа: %s", err)
	}
	// Проверяем необходимые поля в запросе
	_, ok := eR["owner"]
	if !ok {
		return errors.New("отсутствует параметр owner (владелец перечня)")
	}
	_, ok = eR["responsible_users"]
	if !ok {
		return errors.New("отсутствует или неверный формат параметр responsible_users (список назначенных пользователей)")
	}
	_, ok = eR["object_list"]
	if !ok {
		return errors.New("отсутствует параметр object_list (uuid перечня)")
	}
	// Генерим урл
	requestURL := settings.gatewayURL + "/events/new_task"
	// Отправляем запрос
	return sendEventRequest(bsEventRequest, requestURL, userSession)
}

func sendEventRequest(bodyRequest []byte, eventURL string, userSession string) error {
	req, err := http.NewRequest(http.MethodPost, eventURL, bytes.NewReader(bodyRequest))
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: userSession,
	})
	req.AddCookie(&http.Cookie{
		Name:  "X-API-KEY",
		Value: settings.apiKey,
	})

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса на шлюзе: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("невозможно получить ответ от шлюза: %s", err)
		}
		return fmt.Errorf("ответ шлюза: %s", string(bs))
	}

	return nil
}
