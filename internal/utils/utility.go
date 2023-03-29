package utils

import "regexp"

func RegSplit(text string, delim string) []string {
	reg := regexp.MustCompile(delim)
	indexes := reg.FindAllStringIndex(text, -1)
	lastStart := 0
	result := make([]string, len(indexes) + 1)
	for i, element := range indexes {
		result[i] = text[lastStart:element[0]]
		lastStart = element[1]
	}
	result[len(indexes)] = text[lastStart:len(text)]
	return result
}

func TrimSlice(slices []string) []string {
	var results []string
	for _, slice := range slices {
		if len(slice) != 0 {
			results = append(results, slice)
		}
	}
	return results
}