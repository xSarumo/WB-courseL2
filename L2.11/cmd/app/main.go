package main

import (
	"fmt"
	"slices"
)

func anagramSet(input []string) map[string][]string {
	resultMap := make(map[string][]string, len(input)/2)
	for _, word := range input {
		rWord := []rune(word)

		slices.Sort(rWord)

		sWord := string(rWord)
		if _, ok := resultMap[sWord]; ok {
			resultMap[sWord] = append(resultMap[sWord], word)
			continue
		}

		resultMap[sWord] = make([]string, 0, len(input)/2)
		resultMap[sWord] = append(resultMap[sWord], word)
	}
	return resultMap
}

func PrintResultMap(resultMap map[string][]string) {
	for _, value := range resultMap {
		if len(value) > 1 {
			fmt.Println(value[0], value)
		}
	}
}

func main() {
	input := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}

	PrintResultMap(anagramSet(input))
}
