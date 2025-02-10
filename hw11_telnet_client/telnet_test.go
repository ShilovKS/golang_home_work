package main

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestTelnetClient проверяет базовую работу клиента, используя локальный TCP-сервер.
func TestTelnetClient(t *testing.T) {
	// Запускаем локальный сервер.
	listener, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer listener.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	// Клиентская часть.
	go func() {
		defer wg.Done()

		// В качестве источника для Send используем буфер.
		inBuf := bytes.NewBufferString("hello\n")
		// Буфер для получения данных.
		outBuf := &bytes.Buffer{}

		// Создаём клиента с таймаутом 10s.
		client := NewTelnetClient(listener.Addr().String(), 10*time.Second, io.NopCloser(inBuf), outBuf)
		require.NoError(t, client.Connect())
		defer client.Close()

		// Отправляем данные.
		require.NoError(t, client.Send())
		// Получаем данные.
		require.NoError(t, client.Receive())
		// Проверяем, что получили ожидаемую строку.
		require.Equal(t, "world\n", outBuf.String())
	}()

	// Серверная часть.
	go func() {
		defer wg.Done()

		conn, err := listener.Accept()
		require.NoError(t, err)
		defer conn.Close()

		// Читаем данные, ожидаем "hello\n".
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		require.NoError(t, err)
		require.Equal(t, "hello\n", string(buf[:n]))

		// Отправляем ответ "world\n".
		_, err = conn.Write([]byte("world\n"))
		require.NoError(t, err)
	}()

	wg.Wait()
}
