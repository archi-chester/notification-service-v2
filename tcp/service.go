package tcp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	log "github.com/Sirupsen/logrus"
	"gitlab.havana/BDIO/notification-service-v2/serve"
)

// 	выполнение полученной задачи
func Exec(conn net.Conn, messageType int, messageBody []byte) error {
	// 	в зависимости от типа сообщения дергаем конкретную функцию
	switch messageType {
	// 	ошибка
	case MESSAGE_TYPE_ERROR:
		log.Error("Неверный тип сообщения. ")
		return errors.New("Неверный тип сообщения. ")
	// 	заглушка под тест
	case MESSAGE_TYPE_TEST:
		err := sendTestNotice(conn, messageBody)
		if err != nil {
			// handle error
			log.Error("Ошибка обработки testNotice", err)
			return err
		}

	// 	заглушка под тест
	case MESSAGE_TYPE_INSERT:
		err := insertNotice(conn, messageBody)
		if err != nil {
			// handle error
			log.Error("Ошибка обработки insertNotice", err)
			return err
		}
	}
	return nil
}

// 	тестовое сообщение
func sendTestNotice(conn net.Conn, messageBody []byte) error {
	// 	предварительная обработка

	// шифруем
	cipherText, err := encrypt(messageBody)
	if err != nil {
		log.Error("Ошибка шифрования. ", err)
		return err
	}

	// 	отправка обратно
	// 	формируем префикс
	prefix := make([]byte, 4)
	binary.BigEndian.PutUint32(prefix, uint32(len(cipherText)))

	// 	передаем стрим
	_, err = conn.Write(prefix)
	if err != nil {
		// handle error
		log.Error("Ошибка передачи по TCP. Ошибка при передаче префикса. ", err)
		return err
	}
	_, err = conn.Write(cipherText)
	if err != nil {
		// handle error
		log.Error("Ошибка передачи по TCP. Ошибка при передачи сообщения. ", err)
		return err
	}
	// 	возврат ошибки
	return nil
}

// 	добавляем сообщение
func insertNotice(conn net.Conn, messageBody []byte) error {
	// 	предварительная обработка
	var answer, err = serve.InsertNotice(messageBody)
	// 	отправка обратно
	// 	TODO: вызов инсерт
	//
	// plainText = []byte("{'status':'ok', 'message':'" + string(messageBody) + "'}")
	// 	возврат данных клиенту
	// 	шифруем
	cipherText, err := encrypt(answer)
	if err != nil {
		log.Error("Ошибка шифрования. ", err)
		return err
	}

	// 	формируем префикс
	prefix := make([]byte, 4)
	binary.BigEndian.PutUint32(prefix, uint32(len(cipherText)))

	fmt.Printf("Шифровано: %s", string(cipherText))

	// 	передаем стрим
	_, err = conn.Write(prefix)
	if err != nil {
		// handle error
		log.Error("Ошибка передачи по TCP. Ошибка при передаче префикса. ", err)
		return err
	}
	_, err = conn.Write(cipherText)
	if err != nil {
		// handle error
		log.Error("Ошибка передачи по TCP. Ошибка при передачи сообщения. ", err)
		return err
	}
	if err != nil {
		// handle error
		log.Error("Ошибка передачи по TCP. ", err)
		return err
	} else {
		log.Info("Ответ клиенту отправлен.")
	}
	// 	возврат ошибки
	return nil
}
