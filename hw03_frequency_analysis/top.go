package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"regexp"
	"sort"
	"strings"
)

func Top10(rawData string) []string {
	// prepare
	preparedData := strings.Split(prepareString(rawData), " ")
	repeatableWords := make(map[string]int, len(preparedData))

	for _, word := range preparedData {
		repeatableWords[word]++
	}

	t := extractListTopWords(repeatableWords, 10)

	return t
}

// prepare
// clear \n\r\t -> \s
// \s{2,} -> \s
func prepareString(rawData string) string {
	replacer := regexp.MustCompile("[\n|\t|\r]+")
	spaceReplacer := regexp.MustCompile(`[\s]{2,}`)
	preparedString := spaceReplacer.ReplaceAllString(replacer.ReplaceAllString(rawData, " "), " ")

	return preparedString
}

func extractListTopWords(repeatableWords map[string]int, topSize int) []string {
	var mapPositions = make(map[int][]string, topSize)

	for word, repeats := range repeatableWords {
		_, ok := mapPositions[repeats]
		if !ok {
			mapPositions[repeats] = []string{word}
		} else {
			mapPositions[repeats] = append(mapPositions[repeats], word)
		}
	}

	var topPositions = extractTopPositions(mapPositions, topSize)
	var topWords = extractTopWords(mapPositions, topPositions, topSize)

	return topWords
}

func extractTopPositions(mapPositions map[int][]string, topSize int) []int {
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

	return wordsPositions[:size-1]
}

func extractTopWords(mapPositions map[int][]string, topPositions []int, topSize int) []string {

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
