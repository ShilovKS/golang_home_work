package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	runes := []rune(s)
	var result strings.Builder

	var (
		escaped  bool
		prevRune rune
	)

	for i := 0; i < len(runes); i++ {
		currentRune := runes[i]

		switch {
		case escaped:
			if err := handleEscapedRune(&result, currentRune); err != nil {
				return "", err
			}
			prevRune = currentRune
			escaped = false

		case currentRune == '\\':
			escaped = true

		case unicode.IsDigit(currentRune):
			newIndex, err := handleDigit(&result, runes, prevRune, i)
			if err != nil {
				return "", err
			}
			i = newIndex

		default:
			result.WriteRune(currentRune)
			prevRune = currentRune
		}
	}

	if escaped {
		return "", ErrInvalidString
	}

	return result.String(), nil
}

// handleEscapedRune обрабатывает символ после экранирования.
func handleEscapedRune(result *strings.Builder, r rune) error {
	if r != '\\' && !unicode.IsDigit(r) {
		return ErrInvalidString
	}
	result.WriteRune(r)
	return nil
}

// handleDigit обрабатывает цифру и возвращает новый индекс.
func handleDigit(result *strings.Builder, runes []rune, prevRune rune, index int) (int, error) {
	if index == 0 {
		return index, ErrInvalidString
	}

	if unicode.IsDigit(prevRune) && !isEscaped(runes, index-1) {
		return index, ErrInvalidString
	}

	numberStr, newIndex, err := collectNumber(runes, index)
	if err != nil {
		return index, err
	}

	count, err := strconv.Atoi(numberStr)
	if err != nil {
		return index, ErrInvalidString
	}

	if count == 0 {
		if err := removeLastRune(result); err != nil {
			return index, err
		}
	} else {
		result.WriteString(strings.Repeat(string(prevRune), count-1))
	}

	return newIndex, nil
}

// collectNumber собирает число из последовательности цифр.
func collectNumber(runes []rune, index int) (string, int, error) {
	start := index
	for index+1 < len(runes) && unicode.IsDigit(runes[index+1]) && !isEscaped(runes, index+1) {
		index++
	}
	// Проверяем, что число не состоит из более чем одной цифры без экранирования
	if index > start && !isEscaped(runes, start) {
		return "", index, ErrInvalidString
	}
	return string(runes[start : index+1]), index, nil
}

// isEscaped проверяет, экранирован ли символ по индексу.
func isEscaped(runes []rune, index int) bool {
	count := 0
	for i := index - 1; i >= 0; i-- {
		if runes[i] == '\\' {
			count++
		} else {
			break
		}
	}
	return count%2 == 1
}

// removeLastRune удаляет последний символ из результата.
func removeLastRune(result *strings.Builder) error {
	output := result.String()
	if len(output) == 0 {
		return ErrInvalidString
	}
	runes := []rune(output)
	result.Reset()
	result.WriteString(string(runes[:len(runes)-1]))
	return nil
}
