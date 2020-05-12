package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"regexp"
	"sort"
	"strings"
)

func Top10(rawData string) []string {
	var repeatableWords = extractRepeatableWords(rawData)

	return extractTopWords(repeatableWords, 10)
}

func prepareString(rawData string) string {
	var preparers = []func(value string) string{
		func(value string) string {
			return regexp.MustCompile("[\n|\t|\r]+").ReplaceAllString(value, " ")
		},
		//// punctuation
		//func(value string) string {
		//	//\p{P}
		//	return regexp.MustCompile(`([\p{Pc}\p{Pe}\p{Pf}\p{Pi}\p{Po}])`).ReplaceAllString(value, " ")
		//},
		func(value string) string {
			return regexp.MustCompile(`\s{2,}`).ReplaceAllString(value, " ")
		},
		strings.TrimSpace,
		//strings.ToLower,
	}

	var preparedString = rawData
	for _, strategy := range preparers {
		preparedString = strategy(preparedString)
	}

	return preparedString
}

func extractRepeatableWords(rawData string) map[string]int {
	// prepare
	var preparedString = prepareString(rawData)
	var words = strings.Split(preparedString, " ")

	repeatableWords := make(map[string]int, len(words))
	for _, word := range words {
		cleanedWord := strings.TrimSpace(word)
		if cleanedWord != "" {
			repeatableWords[cleanedWord]++
		}
	}

	return repeatableWords
}

func extractTopWords(repeatableWords map[string]int, topSize int) []string {
	var mapPositions = make(map[int][]string)

	// разворачиваем и получаем позиции и слова на них
	for word, repeats := range repeatableWords {
		_, ok := mapPositions[repeats]
		if !ok {
			mapPositions[repeats] = []string{word}
		} else {
			mapPositions[repeats] = append(mapPositions[repeats], word)
		}
	}

	// получаем срез топ позиций
	topPositions := extractPositions(mapPositions, topSize)

	// получаем топ слова на топ позициях
	return extractWords(mapPositions, topPositions, topSize)
}

func extractPositions(mapPositions map[int][]string, topSize int) []int {
	var wordsPositions sort.IntSlice = make([]int, len(mapPositions))
	i := 0
	for position := range mapPositions {
		wordsPositions[i] = position
		i++
	}
	sort.Sort(sort.Reverse(wordsPositions))

	var size = topSize
	if topSize > len(wordsPositions) {
		size = len(wordsPositions)
	}

	return wordsPositions[:size]
}

func extractWords(mapPositions map[int][]string, topPositions []int, topSize int) []string {
	var listWords = make([]string, 0)
	var i = 0
	for _, position := range topPositions {
		words, ok := mapPositions[position]
		if ok {
			sort.Strings(words)

			for _, word := range words {
				if i >= topSize {
					return listWords
				} else {
					listWords = append(listWords, word)
				}
				i++
			}
		}
	}

	return listWords
}
