package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(input string) []string {
	if input == "" {
		return []string{}
	}

	// Разбиваем строку на слова, учитывая пробельные символы.
	words := strings.Fields(input)

	// Подсчет слов.
	wordCounts := make(map[string]int)
	for _, word := range words {
		wordCounts[word]++
	}

	// Сортировка слов по частоте и лексикографически.
	type wordStat struct {
		word  string
		count int
	}

	stats := make([]wordStat, 0, len(wordCounts))
	for word, count := range wordCounts {
		stats = append(stats, wordStat{word, count})
	}

	sort.Slice(stats, func(i, j int) bool {
		if stats[i].count == stats[j].count {
			return stats[i].word < stats[j].word
		}
		return stats[i].count > stats[j].count
	})

	// Выбираем топ-10 слов.
	result := make([]string, 0, 10)
	for i, stat := range stats {
		if i >= 10 {
			break
		}
		result = append(result, stat.word)
	}

	return result
}
