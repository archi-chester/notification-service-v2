package tcp

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"

	log "github.com/Sirupsen/logrus"
)

// слушаем порт
func Listener(port int) {
	// начинаем слушать порт
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	// обработка соединения
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(conn net.Conn) {
			// вход в функцию коннекта
			fmt.Println("Point")

			// read the length prefix
			prefix := make([]byte, 4)
			_, err = io.ReadFull(conn, prefix)

			length := binary.BigEndian.Uint32(prefix)
			// verify length if there are restrictions

			message := make([]byte, int(length))
			_, err = io.ReadFull(conn, message)

			fmt.Println(string(message))

			// дешифруем
			plainText, err := decrypt(message)
			if err != nil {
				log.Error("Ошибка шифрования. ", err)
			}

			// 	переводим структуру в байтовый массив
			messageStruct, err := UnpackMessage(plainText)
			if err != nil {
				log.Error("Ошибка конвертирования структуры. ", err)
			}

			fmt.Printf("Получили %s\n", messageStruct.Message)

			err = Exec(conn, messageStruct.Type, messageStruct.Message)
			if err != nil {
				log.Error("Запуска Exec. ", err)
			}
			// // Echo all incoming data.
			// io.Copy(conn, conn)
			// // Читаем данные из порта
			// status, err := bufio.NewReader(conn).ReadString('\n')
			// if err != nil {
			// 	// handle error
			// 	fmt.Println(err)
			// }
			// fmt.Println(status)
			// Shut down the connection.
			conn.Close()
		}(conn)
	}

}

//  посылаем сообщение
func Sender(host string, port int, messageType int, messageBody []byte) error {
	// создаем сендера
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Error("Ошибка создания сендера. ", err)
		return err
	}
	// 	формируем структуру для передачи
	messageStruct, err := CreateMessage(messageType, TCP_PORT, TCP_IP, messageBody)
	if err != nil {
		log.Error("Ошибка формирования структуры сообщения. ", err)
		return err
	}

	// 	переводим структуру в байтовый массив
	plainText, err := PackMessage(&messageStruct)
	if err != nil {
		log.Error("Ошибка конвертирования структуры. ", err)
		return err
	}

	// шифруем
	cipherText, err := encrypt(plainText)
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
	return nil
}

// 	создание сообщения
func CreateMessage(messageType int, sourcePort int, sourceIP string, messageBody []byte) (MessagePackage, error) {
	// 	времянка для возврата
	var messageStruct MessagePackage
	// 	заполняем
	// 	тип сообщения
	if messageType != MESSAGE_TYPE_ERROR {
		messageStruct.Type = messageType
	} else {
		return messageStruct, errors.New("Неверный тип сообщения")
	}
	// 	порт
	if sourcePort >= 0 && sourcePort < 65535 {
		messageStruct.SourcePort = sourcePort
	} else {
		return messageStruct, errors.New("Неверный порт")
	}
	// 	ip
	if len(sourceIP) > 6 && len(sourceIP) < 16 {
		messageStruct.SourceIP = sourceIP
	} else {
		return messageStruct, errors.New("Неверный ip")
	}
	// 	тип сообщения
	if len(messageBody) != 0 && len(messageBody) < 950 {
		messageStruct.Message = messageBody
	} else {
		return messageStruct, errors.New("Неверный размер сообщения")
	}

	// 	возврат
	return messageStruct, nil
}

// 	запаковываем структуру в байтовый массив
func PackMessage(message *MessagePackage) ([]byte, error) {
	// 	времянка для возврата
	var messagePack []byte

	//	Маршалим прочитанное в json
	messagePack, err := json.Marshal(message)
	if err != nil {
		log.Error("Кривой Маршал")
	}
	return messagePack, nil

}

// 	распаковываем байтовый массив в структуру
func UnpackMessage(messagePack []byte) (MessagePackage, error) {
	// 	времянка для возврата
	var messageStruct MessagePackage

	//	Анмаршалим прочитанное в структуру
	err := json.Unmarshal(messagePack, &messageStruct)
	if err != nil {
		log.Error("Кривой анмаршал")
	}
	return messageStruct, nil
}
