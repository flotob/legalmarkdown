package lmd

import (
	"regexp"
	"strconv"
	"strings"
)

// Header is the primary type which establishes how the structured headers will be parsed.
// Trigger is the string at the beginning of the line in the parsing block which triggers
// this header. Reset is a boolean which lets the parser know whether this header is to be
// reset or not. Level Number is an integer of which level of the header this particular
// header is. Indent is the number of spaces which the header should be indented. Style
// is pulled from the case switch in the defineHeaderStyle function. ResetVal is the value
// which the header is reset to -- this is always the value which the header is initiated
// to. BeforeVal is the string which is placed before the currrent value. CurrtVal is the
// current value of the header which is iterated as the tree parser performs its work.
// AfterVal is the string which goes after the currtVal -- typically it is only one char.
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

// SetTheHeaders is the primary parser of structured headers. It parses the styles and header
// structure and establishes a map of triggers and Header structs delimiting the structured
// headers.
//
// It begins by parsing the "level-style" parameter to determine if the original style "llll."
// headers are used by the template or if the new style headers "l4." are used.
//
// Second the function parses the "no-indent" parameter to determine which of the headers are
// not to be indented. These are put into a slice.
//
// Third it parses the "no-reset" parameter to determine which of the headers are not to be
// reset. These are also put into a slice.
//
// Finally the function calls the main parsing function "parseHeaders" which returns
// the map of structs and triggers for each of the relevant headers by the parser. This
// map is what is returned to the calling function.
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

// defineHeaderStyle sets the style integer according to the 10 relevant types which are
// allowed in legalmarkdown. type1 sets the *Header.style to 1 and is for uppercase
// roman numerals followed by a period. type2 sets the *Header.style to 2 and is for
// uppercase roman numerals encased in parentheses. type3 sets the *Header.style to
// 3 and is for lowercase roman numerals followed by a period. type4 sets the *Header
// .style to 4 and is for lowercase roman numerals encased in parentheses. type5 sets
// the *Header.style to 5 and is for uppercase letters followed by a period. type6
// sets the *Header.style to 6 and is for uppercase letters encased in parentheses.
// type7 sets the *Headers.style to 7 and is for lowercase letters followed by a
// period. type8 sets the *Header.style to 8 and is for lowercase letters encased in
// parentheses. type9 sets the *Header.style to 9 and is for numbers encased in
// parentheses. type0 sets the *Header.style to 0 and is for numbers followed by a
// period. This is also the dot.
//
// in addition to returning the style integer, the following strings are parsed from
// the parameter which is passed to the function: the *Header.beforVal, the *Header.
// currtVal and the *Header.afterVal all of which are simply pulled from the string
// parse which is also used to determine the style type.
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
