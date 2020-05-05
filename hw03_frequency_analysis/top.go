package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"regexp"
	"sort"
	"strings"
)

func Top10(rawData string) []string {

	var repeatableWords = extractMapRepeatableWords(rawData)

	return extractTopWords(repeatableWords, 10)
}

// prepare
// clear \n\r\t -> \s
// \s{2,} -> \s
func prepareString(rawData string) string {
	replacer := regexp.MustCompile("[\n|\t|\r]+")
	spaceReplacer := regexp.MustCompile(`[\s]{2,}`)
	preparedString := spaceReplacer.ReplaceAllString(replacer.ReplaceAllString(rawData, " "), " ")
	preparedString = strings.Trim(preparedString, "\n\r\t ")

	return preparedString
}

func extractMapRepeatableWords(rawData string) map[string]int {
	// prepare
	var preparedString = prepareString(rawData)
	var words = strings.Split(preparedString, " ")
	var capWords = 0

	if len(words) > 1 {
		capWords = len(words)
	}

	repeatableWords := make(map[string]int, capWords)
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

	for word, repeats := range repeatableWords {
		_, ok := mapPositions[repeats]
		if !ok {
			mapPositions[repeats] = []string{word}
		} else {
			mapPositions[repeats] = append(mapPositions[repeats], word)
		}
	}

	topPositions := extractPositions(mapPositions, topSize)

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
