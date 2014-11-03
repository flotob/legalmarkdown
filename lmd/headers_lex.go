package lmd

import (
	"regexp"
	"strconv"
	"strings"
)

type Header struct {
	trigger  string
	reset    bool
	levelNum int
	indent   int
	style    int
	resetVal string
	beforVal string
	currtVal string
	afterVal string
}

// SetTheHeaders does shit ...
func SetTheHeaders(contents string, parameters map[string]string) map[string]*Header {
	levelStyle := parseLevelStyle(parameters["level-style"])
	delete(parameters, "level-style")
	indentSlice := parseIndents(parameters["no-indent"])
	delete(parameters, "no-indent")
	resetSlice := parseResets(parameters["no-reset"])
	delete(parameters, "no-reset")
	headers := parseHeaders(parameters, levelStyle, indentSlice, resetSlice)
	return headers
}

// return true if "llll." style, else if "l4." style return false
func parseLevelStyle(style string) bool {
	if style == "l1." {
		return false
	}
	return true
}

// split the "no-indent" string into a slice of headers
func parseIndents(indents string) []string {
	indentSlice := strings.Split(indents, ", ")
	return indentSlice
}

// split the "no-reset" string into a slice of headers
func parseResets(resets string) []string {
	resetSlice := strings.Split(resets, ", ")
	return resetSlice
}

// set up structs for the headers and put those into a map for use by the parser
func parseHeaders(parameters map[string]string, levelStyle bool, indentSlice []string, resetSlice []string) map[string]*Header {

	var header *Header
	headers := make(map[string]*Header)

	// set the defaults based on parsing the params.
	for paramKey, paramVal := range parameters {

		header = new(Header)
		header.levelNum, _ = strconv.Atoi(paramKey[len(paramKey)-1:])

		if levelStyle {
			header.trigger = (strings.Repeat("l", header.levelNum) + ".")
		} else {
			header.trigger = ("l" + strconv.Itoa(header.levelNum) + ".")
		}

		header.indent = (2 * header.levelNum)

		header.reset = true

		header.style, header.beforVal, header.currtVal, header.afterVal = defineHeaderStyle(paramVal)
		header.resetVal = header.currtVal

		headers[header.trigger] = header
	}

	// correct resets
	for _, rest := range resetSlice {
		if headers[rest] != nil {
			headers[rest].reset = false
		}
	}

	// set and correct indents
	for _, indnt := range indentSlice {
		if headers[indnt] != nil {
			headers[indnt].indent = 0
		}
	}
	for _, head := range headers {
		baseNum := len(indentSlice)
		if head.indent != 0 && baseNum != 0 {
			head.indent = (head.indent - (2 * baseNum))
		}
		if head.indent < 0 {
			head.indent = 0
		}
	}

	return headers
}

func defineHeaderStyle(h string) (int, string, string, string) {

	// set all the regexps first. roman numerals are parsed first as matches
	// for those need to be extracted before the alphabeticals
	type1 := regexp.MustCompile(`([IVXLCDM]+)\.\z`)   // {{  I. }}
	type2 := regexp.MustCompile(`\(([IVXLCDM]+)\)\z`) // {{ (I) }}
	type3 := regexp.MustCompile(`([ivxlcdm]+)\.\z`)   // {{  i. }}
	type4 := regexp.MustCompile(`\(([ivxlcdm]+)\)\z`) // {{ (i) }}
	type5 := regexp.MustCompile(`([A-Z]+)\.\z`)       // {{  A. }}
	type6 := regexp.MustCompile(`\(([A-Z]+)\)\z`)     // {{ (A) }}
	type7 := regexp.MustCompile(`([a-z]+)\.\z`)       // {{  a. }}
	type8 := regexp.MustCompile(`\(([a-z]+)\)\z`)     // {{ (a) }}
	type9 := regexp.MustCompile(`\(([0-9]+)\)\z`)     // {{ (1) }}
	type0 := regexp.MustCompile(`([0-9]+)\.\z`)       // {{ 1. }} ... also default

	// now run through the sequence
	switch {
	case type1.MatchString(h):
		return 1, type1.ReplaceAllString(h, ""), type1.FindAllStringSubmatch(h, -1)[0][1], ". "
	case type2.MatchString(h):
		return 2, (type2.ReplaceAllString(h, "") + "("), type2.FindAllStringSubmatch(h, -1)[0][1], ") "
	case type3.MatchString(h):
		return 3, type3.ReplaceAllString(h, ""), type3.FindAllStringSubmatch(h, -1)[0][1], ". "
	case type4.MatchString(h):
		return 4, (type4.ReplaceAllString(h, "") + "("), type4.FindAllStringSubmatch(h, -1)[0][1], ") "
	case type5.MatchString(h):
		return 5, type5.ReplaceAllString(h, ""), type5.FindAllStringSubmatch(h, -1)[0][1], ". "
	case type6.MatchString(h):
		return 6, (type6.ReplaceAllString(h, "") + "("), type6.FindAllStringSubmatch(h, -1)[0][1], ") "
	case type7.MatchString(h):
		return 7, type7.ReplaceAllString(h, ""), type7.FindAllStringSubmatch(h, -1)[0][1], ". "
	case type8.MatchString(h):
		return 8, (type8.ReplaceAllString(h, "") + "("), type8.FindAllStringSubmatch(h, -1)[0][1], ") "
	case type9.MatchString(h):
		return 9, (type9.ReplaceAllString(h, "") + "("), type9.FindAllStringSubmatch(h, -1)[0][1], ") "
	case type0.MatchString(h):
		return 0, type0.ReplaceAllString(h, ""), type0.FindAllStringSubmatch(h, -1)[0][1], ". "
	default:
		return 0, "", "1", ". "
	}
}
