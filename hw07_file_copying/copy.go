package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileInfo              = errors.New("failed to get file info")
	ErrLimit                 = errors.New("limit might not be lower than 0")
	ErrFileOverlap           = errors.New("source and destination files must be different")
)

// Copy копирует данные из файла fromPath в файл toPath, начиная с позиции offset.
// Если limit равен 0, копируется весь файл (от offset до EOF), иначе копируются не более limit байт.
// Если offset больше размера файла – возвращается ошибка.
func Copy(fromPath, toPath string, offset, limit int64) error {
	// Получаем информацию об исходном файле.
	srcInfo, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", ErrFileInfo)
	}

	// Если целевой файл уже существует, проверяем, что это не тот же файл.
	if dstInfo, err := os.Stat(toPath); err == nil {
		if os.SameFile(srcInfo, dstInfo) {
			return ErrFileOverlap
		}
	}

	// Открываем исходный файл только для чтения.
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer fromFile.Close()

	// Проверяем, что исходный файл является обычным файлом.
	if !srcInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > srcInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	if limit < 0 {
		return ErrLimit
	}

	// Создаём (или перезаписываем) целевой файл.
	toFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer toFile.Close()

	// Смещаем указатель чтения исходного файла на offset.
	if _, err := fromFile.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek in source file: %w", err)
	}

	// Определяем число байт для копирования: если limit == 0, копируем до EOF,
	// иначе берем минимальное значение между limit и оставшимся размером файла.
	bytesToCopy := srcInfo.Size() - offset
	if limit > 0 && limit < bytesToCopy {
		bytesToCopy = limit
	}

	// Создаём прогресс-бар.
	bar := pb.Full.Start64(bytesToCopy)
	// Оборачиваем исходный файл в прокси-ридер, который обновляет прогресс-бар.
	barReader := bar.NewProxyReader(fromFile)
	defer bar.Finish()

	// Копируем данные.
	_, err = io.CopyN(toFile, barReader, bytesToCopy)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	return nil
}
