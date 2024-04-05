package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Создание объекта для чтения с консоли
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Connecting to server...")

	// Подключение к серверу
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server.")

	// Флаг, показывающий, присоединен ли клиент к чату
	joined := false

	// Бесконечный цикл для чтения ввода пользователя
	for {
		fmt.Print("You: ")

		// Считывание ввода пользователя
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err.Error())
			return
		}

		// Обрезаем пробелы в начале и в конце строки
		message = strings.TrimSpace(message)

		if !joined {
			// Если пользователь не присоединен к чату, отправляем команду /join
			if message == "/join" {
				fmt.Println("You have joined the chat.")
				joined = true
				continue
			} else {
				fmt.Println("Please join the chat first by typing '/join'.")
				continue
			}
		}

		// Обработка команды для просмотра истории
		if message == "/history" {
			// Отправка команды на сервер
			_, err := conn.Write([]byte(message + "\n"))
			if err != nil {
				fmt.Println("Error writing:", err.Error())
				return
			}

			// Чтение истории от сервера
			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
				return
			}
			fmt.Println("Chat History:")
			fmt.Println(string(buffer[:n]))
			continue
		}

		// Отправка сообщения серверу
		_, err = conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Error writing:", err.Error())
			return
		}

		// Чтение ответа от сервера
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}
		fmt.Println("Server:", string(buffer[:n]))
	}
}
