package lmd

import (
	"strconv"
)

// next_lettering will iterate a given string to give the next logical lettering point
// please note this only works with english language letters "a" to "z" and "A" to "Z"
// this is not meant to be a general purpose letter iterator, but rather is purpose
// built for legal documents which use letter based denotations.
//
// if the final code point (converted to integer type) is 90 or 122, that means
// that the final letter is "Z" or "z" in which case we subtract 25 from the code
// point and then double the letter which will return the final bit as "AA" or "aa"
// in keeping with fairly typical legal syntax.
func next_lettering(previous string) string {

	codepoint := []byte(previous)
	finlpoint := codepoint[len(codepoint)-1]

	if int(finlpoint) == 90 || int(finlpoint) == 122 {
		codepoint[len(codepoint)-1] = finlpoint - 25
		codepoint = append(codepoint, codepoint[len(codepoint)-1])
	} else {
		codepoint[len(codepoint)-1] = finlpoint + 1
	}

	return string(codepoint)
}

// prev_lettering is the opposite of next_lettering.
func prev_lettering(previous string) string {

	codepoint := []byte(previous)
	finlpoint := codepoint[len(codepoint)-1]
	var penlpoint byte

	if len(codepoint) >= 2 {
		penlpoint = codepoint[len(codepoint)-2]
	} else {
		penlpoint = byte(0)
	}

	if finlpoint == penlpoint {
		if int(finlpoint) == 65 || int(finlpoint) == 97 {
			codepoint[len(codepoint)-2] = finlpoint + 25
			codepoint = codepoint[:len(codepoint)-1]
		} else {
			codepoint[len(codepoint)-1] = finlpoint - 1
		}
	} else {
		if !(int(finlpoint) == 65 || int(finlpoint) == 97) {
			codepoint[len(codepoint)-1] = finlpoint - 1
		}
	}

	return string(codepoint)

}

// next_roman_upper increases a roman number by one position for upper cased roman
// numerals by using the from_roman_to_arabic function which returns an integer
// that is subsequently increased by one and then changed back to a roman number
// via the from_arabic_to_roman function.
func next_roman_upper(previous string) string {
	prev_as_digit := from_roman_to_arabic_upper(previous)
	next_as_digit := prev_as_digit + 1
	return from_arabic_to_roman_upper(next_as_digit)
}

// prev_roman_upper is the opposite of next_roman_upper
func prev_roman_upper(previous string) string {
	prev_as_digit := from_roman_to_arabic_upper(previous)
	next_as_digit := prev_as_digit - 1
	return from_arabic_to_roman_upper(next_as_digit)
}

// next_roman_lower increases a roman number by one position for lower cased roman
// numerals by using the from_roman_to_arabic function which returns an integer
// that is subsequently increased by one and then changed back to a roman number
// via the from_arabic_to_roman function.
func next_roman_lower(previous string) string {
	prev_as_digit := from_roman_to_arabic_lower(previous)
	next_as_digit := prev_as_digit + 1
	return from_arabic_to_roman_lower(next_as_digit)
}

// prev_roman_lower isthe opposite of next_roman_lower
func prev_roman_lower(previous string) string {
	prev_as_digit := from_roman_to_arabic_lower(previous)
	next_as_digit := prev_as_digit - 1
	return from_arabic_to_roman_lower(next_as_digit)
}

// next_numbering increases numbers which are presented as strings for the package
// the function uses the strconv package first to convert the string to an integer
// then it increases the integer by one and then converts it back to a string.
func next_numbering(previous string) string {
	prev_as_digit, _ := strconv.Atoi(previous)
	next_as_digit := prev_as_digit + 1
	return strconv.Itoa(next_as_digit)
}

// prev_numbering is the opposite of next_numbering
func prev_numbering(previous string) string {
	prev_as_digit, _ := strconv.Atoi(previous)
	next_as_digit := prev_as_digit - 1
	return strconv.Itoa(next_as_digit)
}

// from_roman_to_arabic_upper is a convenience function which converts upper cased
// roman numerals to integers.
//
// via: http://rosettacode.org/wiki/Roman_numerals/Decode#Go
func from_roman_to_arabic_upper(romans string) int {
	var arabic int

	// set a map from the roman runes to arabic integers
	var roman_to_arabic_map = map[rune]int{
		'I': 1,
		'V': 5,
		'X': 10,
		'L': 50,
		'C': 100,
		'D': 500,
		'M': 1000,
	}

	// loop through the roman string which is passed by rune
	last_digit := 1000
	for _, roman := range romans {
		digit := roman_to_arabic_map[roman]
		if last_digit < digit {
			arabic -= 2 * last_digit
		}
		last_digit = digit
		arabic += digit
	}

	return arabic
}

// from_arabic_to_roman_upper is a convenience function which converts upper cased
// integers to roman numerals.
//
// via: http://rosettacode.org/wiki/Roman_numerals/Encode#Go
func from_arabic_to_roman_upper(arabic int) string {

	var (
		m0 = []string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"}
		m1 = []string{"", "X", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"}
		m2 = []string{"", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"}
		m3 = []string{"", "M", "MM", "MMM", "IV", "V", "VI", "VII", "VIII", "IX"}
		m4 = []string{"", "X", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"}
		m5 = []string{"", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"}
		m6 = []string{"", "M", "MM", "MMM"}
	)

	if arabic < 1 || arabic >= 4e6 {
		return ""
	}

	return m6[arabic/1e6] + m5[arabic%1e6/1e5] + m4[arabic%1e5/1e4] + m3[arabic%1e4/1e3] + m2[arabic%1e3/1e2] + m1[arabic%100/10] + m0[arabic%10]
}

// from_roman_to_arabic_lower is a convenience function which converts lower cased
// roman numerals to integers.
//
// via: http://rosettacode.org/wiki/Roman_numerals/Decode#Go
func from_roman_to_arabic_lower(romans string) int {
	var arabic int

	// set a map from the roman runes to arabic integers
	var roman_to_arabic_map = map[rune]int{
		'i': 1,
		'v': 5,
		'x': 10,
		'l': 50,
		'c': 100,
		'd': 500,
		'm': 1000,
	}

	// loop through the roman string which is passed by rune
	last_digit := 1000
	for _, roman := range romans {
		digit := roman_to_arabic_map[roman]
		if last_digit < digit {
			arabic -= 2 * last_digit
		}
		last_digit = digit
		arabic += digit
	}

	return arabic
}

// from_arabic_to_roman_upper is a convenience function which converts lower cased
// integers to roman numerals.
//
// via: http://rosettacode.org/wiki/Roman_numerals/Encode#Go
func from_arabic_to_roman_lower(arabic int) string {

	var (
		m0 = []string{"", "i", "ii", "iii", "iv", "v", "vi", "vii", "viii", "ix"}
		m1 = []string{"", "x", "xx", "xxx", "xl", "l", "lx", "lxx", "lxxx", "xc"}
		m2 = []string{"", "c", "cc", "ccc", "cd", "d", "dc", "dcc", "dccc", "cm"}
		m3 = []string{"", "m", "mm", "mmm", "iv", "v", "vi", "vii", "viii", "xi"}
		m4 = []string{"", "x", "xx", "xxx", "xl", "l", "lx", "lxx", "lxxx", "xc"}
		m5 = []string{"", "c", "cc", "ccc", "cd", "d", "dc", "dcc", "dccc", "cm"}
		m6 = []string{"", "m", "mm", "mmm"}
	)

	if arabic < 1 || arabic >= 4e6 {
		return ""
	}

	return m6[arabic/1e6] + m5[arabic%1e6/1e5] + m4[arabic%1e5/1e4] + m3[arabic%1e4/1e3] + m2[arabic%1e3/1e2] + m1[arabic%100/10] + m0[arabic%10]
}
