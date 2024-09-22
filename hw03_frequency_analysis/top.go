package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(source string) []string {
	if isEmptyString(&source) {
		return nil
	}

	cleanedText := prepareString(source)

	allWords := make(map[string]int)
	split := strings.Fields(cleanedText)

	for _, word := range split {
		_, ok := allWords[word]
		if ok {
			allWords[word]++
		} else {
			allWords[word] = 1
		}
	}

	keys := make([]string, 0, len(allWords))

	for key := range allWords {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		if allWords[keys[i]] > allWords[keys[j]] {
			return true
		} else if allWords[keys[i]] == allWords[keys[j]] {
			return keys[i] < keys[j]
		}

		return false
	})

	return keys[:10]
}

func prepareString(source string) string {
	source = strings.ToLower(source)

	source = regexp.MustCompile(`\s-\s`).ReplaceAllString(source, " ") // удаляем одиночные дефисы с пробелами
	source = regexp.MustCompile(`-`).ReplaceAllString(source, "")      // удаляем оставшиеся одиночные дефисы

	// Регулярное выражение для удаления знаков препинания после слова
	re := regexp.MustCompile(`([a-zA-Zа-яА-ЯёЁ0-9\-]{2,})([.,!?;:]+)(\s)`)
	cleanedText := re.ReplaceAllString(source, "$1$3") // сохраняем слово и пробел

	// Удаляем любые оставшиеся знаки препинания
	return regexp.MustCompile(`[.,!?;:]+`).ReplaceAllString(cleanedText, "")
}

func isEmptyString(s *string) bool {
	if s == nil || *s == "" {
		return true
	}
	return false
}
