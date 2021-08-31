package helper

import (
	"fmt"
	"regexp"
	"strings"
)

// to find word in string
func RegexFindTerm(source, pattern string) string {
	regex := regexp.MustCompile(pattern)
	return regex.FindString(source)
}

// regex to get word inside string
func RegexFindBetweenTwoPattern(before, after, errStr string) string {

	beforeRegex := regexp.MustCompile(fmt.Sprintf(".*%v", before))
	afterRegex := regexp.MustCompile(fmt.Sprintf("%v.*", after))

	errStr = beforeRegex.ReplaceAllString(errStr, "")
	result := afterRegex.ReplaceAllString(errStr, "")

	return strings.TrimSpace(result)
}
