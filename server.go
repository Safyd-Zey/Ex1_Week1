package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// Структура для хранения истории сообщений
type History struct {
	Messages []string
}

// Добавление сообщения в историю
func (h *History) AddMessage(message string) {
	h.Messages = append(h.Messages, message)
}

// Вывод истории сообщений
func (h *History) ShowHistory(conn net.Conn) {
	conn.Write([]byte("\n== Chat History ==\n"))
	for _, msg := range h.Messages {
		conn.Write([]byte(msg + "\n"))
	}
	conn.Write([]byte("== End of History ==\n\n"))
}

func main() {
	// Инициализация истории сообщений
	history := History{}

	// Создание сервера
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("Server started. Waiting for clients...")

	// Бесконечный цикл для принятия входящих соединений
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}

		// Запуск горутины для обработки клиента
		go handleConnection(conn, &history)
	}
}

// Обработка соединения с клиентом
func handleConnection(conn net.Conn, history *History) {
	defer conn.Close()

	fmt.Println("Client connected:", conn.RemoteAddr())

	// Бесконечный цикл для чтения сообщений от клиента
	for {
		// Чтение сообщения от клиента
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}

		// Преобразование полученных байтов в строку
		message := string(buffer[:n])

		// Обработка команды /history
		if strings.TrimSpace(message) == "/history" {
			history.ShowHistory(conn)
			continue
		}

		// Добавление сообщения в историю
		history.AddMessage(message)

		// Вывод сообщения в консоль сервера
		fmt.Printf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), message)

		// Отправка сообщения клиенту
		conn.Write([]byte("Message received by server: " + message))
	}
}
