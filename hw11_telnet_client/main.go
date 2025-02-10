package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "timeout for connection")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatalln("Usage: go-telnet [--timeout=duration] host port")
	}
	host, port := args[0], args[1]
	if host == "" || port == "" {
		log.Fatalln("host and port must not be empty")
	}

	address := net.JoinHostPort(host, port)

	// Создаём контекст с обработкой сигналов SIGINT и SIGTERM.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// Не используем defer cancel() здесь, поскольку при ошибке ниже вызывается log.Fatalf,
	// а отложенный cancel() не выполнится.

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	// При ошибке подключения вызываем cancel() явно, затем завершаем работу.
	if err := client.Connect(); err != nil {
		cancel() // вызываем cancel() вручную, чтобы освободить ресурсы
		log.Fatalf("Failed to connect to %s: %v\n", address, err)
	}

	// При нормальном подключении гарантируем закрытие соединения и отмену контекста.
	defer func() {
		_ = client.Close()
		cancel()
	}()

	log.Printf("...Connected to %s\n", address)

	// Запускаем горутины для отправки и получения данных.
	startRoutine("Send", client.Send, cancel)
	startRoutine("Receive", client.Receive, cancel)

	<-ctx.Done()
	log.Println("...Connection closed. Bye!")
}

func startRoutine(name string, task func() error, cancel func()) {
	go func() {
		if err := task(); err != nil {
			// Вывод ошибки в стандартный поток ошибок.
			fmt.Fprintf(os.Stderr, "%s error: %v\n", name, err)
		}
		// Если задача завершилась (с ошибкой или без), вызываем cancel(),
		// чтобы сигнализировать об окончании работы и остановить клиент.
		cancel()
	}()
}
