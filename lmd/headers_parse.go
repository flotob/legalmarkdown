package lmd

import (
	"regexp"
	"strconv"
	"strings"
)

// HandleTheHeaders is the primary parser function for parsing a block of structured headers.
func HandleTheHeaders(contents string, headers map[string]*Header) string {

	to_run, pre_block, block, post_block := findTheBlock(contents)
	if !to_run {
		return contents
	}

	blockAsSlice, blockBase := splitTheBlock(block)

	block = runTheHeaders(headers, blockAsSlice, blockBase)
	contents = pre_block + "\n" + block + "\n\n" + post_block

	return contents

}

// findTheBlock regexes the whole contents against the blockpattern "```". if a match is found,
// it returns the portion of the contents prior to the block, the contents of the block with the backticks
// removed, the portion of the contents following the block along with (the first return value) true.
// if no match is found false, two empty strings and the original contents are returned.
func findTheBlock(contents string) (bool, string, string, string) {
	blockPattern := regexp.MustCompile(`(?sm)(^` + "```" + `+\s*\n?)(.*?\n?)(^` + "```" + `+\s*\n?|\z)`)

	if blockPattern.MatchString(contents) {
		matches := blockPattern.FindAllStringSubmatch(contents, -1)
		non_matches := strings.Split(contents, matches[0][0])
		return true, non_matches[0], matches[0][2], non_matches[1]
	}

	return false, "", "", ""
}

// splitTheBlock takes a block of text, splits it into a slice based on the new lines. then it refactors
// that block, taking out empty lines from the text and adding lines which do not begin with the
// header pattern to the last element in the assembled slice. it returns a copy of the assembled slice as
// well as a copy of the slice which only has the leaders (which is needed during the tree parsing phase.
func splitTheBlock(block string) ([]string, []string) {

	blockAsSlice := []string{}
	blockBase := []string{}

	headerPatternOld := regexp.MustCompile(`\Al+\.`)
	headerPatternNew := regexp.MustCompile(`\Al[0-9]+\.`)
	blankPattern := regexp.MustCompile(`\A\s*\z`)

	for _, line := range strings.Split(block, "\n") {
		if headerPatternOld.MatchString(line) {
			blockAsSlice = append(blockAsSlice, line)
			leader := strings.TrimSpace(headerPatternOld.FindAllString(line, 1)[0])
			blockBase = append(blockBase, leader)
		} else if headerPatternNew.MatchString(line) {
			blockAsSlice = append(blockAsSlice, line)
			leader := strings.TrimSpace(headerPatternNew.FindAllString(line, 1)[0])
			blockBase = append(blockBase, leader)
		} else if blankPattern.MatchString(line) {
			continue
		} else {
			blockAsSlice[len(blockAsSlice)-1] = (blockAsSlice[len(blockAsSlice)-1] + "\n\n" + line)
		}
	}

	return blockAsSlice, blockBase
}

// runTheHeaders ...
func runTheHeaders(headers map[string]*Header, blockAsSlice []string, blockBase []string) string {

	crossref := make(map[string]string)
	headerPatternOld := regexp.MustCompile(`\Al+.`)
	headerPatternNew := regexp.MustCompile(`\Al[0-9]+.`)
	oldStyle := true

	if headers["l1."] != nil {
		oldStyle = false
	}

	for i, block := range blockAsSlice {
		if oldStyle {
			leader := headerPatternOld.FindAllString(block, 1)[0]
			block, crossref = replaceTheLeader(leader, headers, block, crossref, oldStyle)
		} else {
			leader := headerPatternNew.FindAllString(block, 1)[0]
			block, crossref = replaceTheLeader(leader, headers, block, crossref, oldStyle)
		}
		iterateTheLeader(headers, blockBase, i)
		blockAsSlice[i] = block
	}

	return collateTheBlock(blockAsSlice, crossref)
}

// replaceTheLeader ...
func replaceTheLeader(leader string, headers map[string]*Header, block string, crossref map[string]string, oldStyle bool) (string, map[string]string) {

	var newLeader string
	header := headers[leader]

	thisBeforVal := strings.TrimSpace(header.beforVal)

	if strings.HasSuffix(thisBeforVal, "pre") || strings.HasSuffix(thisBeforVal, "pre (") || strings.HasSuffix(thisBeforVal, "preval") {
		newLeader = assemblePreVal(leader, headers, false, oldStyle)
	} else {
		newLeader = header.beforVal + header.currtVal + header.afterVal
	}

	leader, crossref = handleCrossReferences(leader, newLeader, block, crossref, oldStyle)

	block = strings.Replace(block, leader, newLeader, 1)
	block = strings.Replace(block, "  ", " ", -1)

	indents := strings.Repeat(" ", header.indent)
	hasMultipleLines := regexp.MustCompile(`\n\n`)
	if hasMultipleLines.MatchString(block) {
		block = hasMultipleLines.ReplaceAllString(block, ("\n\n" + indents))
	}
	block = indents + block

	return block, crossref
}

// assemblePreVal ...
func assemblePreVal(leader string, headers map[string]*Header, nested bool, oldStyle bool) string {
	var thisCurrtVal string
	var reformedLeader string
	var prevHeader *Header
	thisHeader := headers[leader]
	if oldStyle {
		prevHeader = headers[leader[1:]]
	} else {
		tmp, _ := strconv.Atoi(leader[1:])
		prevHeader = headers[strconv.Itoa(tmp-1)]
	}
	prevBeforVal := strings.TrimSpace(prevHeader.beforVal)
	thisBeforVal := strings.TrimSpace(thisHeader.beforVal)

	// if val is "Section preval" then the result should be "Section 10x"
	if strings.HasSuffix(thisBeforVal, "pre") && thisBeforVal != "pre" {
		reformedLeader = strings.Replace(thisBeforVal, "pre", "", 1)
	} else if strings.HasSuffix(thisBeforVal, "pre (") && thisBeforVal != "pre (" {
		reformedLeader = strings.Replace(thisBeforVal, "pre (", "", 1)
	} else if strings.HasSuffix(thisBeforVal, "preval") && thisBeforVal != "preval" {
		reformedLeader = strings.Replace(thisBeforVal, "preval", "", 1)
	} else {
		reformedLeader = ""
	}

	// nesting pre and preval is possible, this section controls that.
	if strings.HasSuffix(prevBeforVal, "pre") || strings.HasSuffix(prevBeforVal, "pre (") || strings.HasSuffix(prevBeforVal, "preval") {
		reformedLeader = reformedLeader + assemblePreVal(leader[1:], headers, true, oldStyle)
	} else {
		reformedLeader = reformedLeader + deIterateThisHeader(prevHeader.currtVal, prevHeader.style)
	}
	if !nested {
		thisCurrtVal = thisHeader.currtVal
	} else {
		thisCurrtVal = deIterateThisHeader(thisHeader.currtVal, thisHeader.style)
	}

	// start to build the string.
	reformedLeader = reformedLeader + prevHeader.afterVal
	reformedLeader = strings.TrimSpace(reformedLeader)

	// most of the building happens here.
	if strings.HasSuffix(thisBeforVal, "pre (") {
		reformedLeader = reformedLeader + "(" + thisCurrtVal + thisHeader.afterVal
	} else if strings.HasSuffix(thisBeforVal, "preval") {
		reformedLeader = strings.Replace(reformedLeader, ".", "", -1)
		reformedLeader = reformedLeader + "0" + thisCurrtVal + thisHeader.afterVal
	} else {
		reformedLeader = reformedLeader + thisCurrtVal + thisHeader.afterVal
	}

	// cleanup
	reformedLeader = strings.Replace(reformedLeader, ".(", "(", -1)
	reformedLeader = strings.Replace(reformedLeader, ". (", "(", -1)
	reformedLeader = strings.Replace(reformedLeader, ") )(", ")(", -1)

	return reformedLeader

}

// handleCrossReferences ...
func handleCrossReferences(leader string, newLeader string, block string, crossref map[string]string, oldStyle bool) (string, map[string]string) {

	var cross string
	var crossTmp string
	hasCrossRef := regexp.MustCompile(`\Al+. \|(.+?)\|`)
	if !oldStyle {
		hasCrossRef = regexp.MustCompile(`\Al[0-9]+. \|(.+?)\|`)
	}

	if hasCrossRef.MatchString(block) {
		leader = hasCrossRef.FindAllString(block, 1)[0]
		cross = hasCrossRef.FindAllStringSubmatch(block, 1)[0][1]
		crossTmp = strings.TrimSpace(newLeader)
		if strings.HasSuffix(crossTmp, ".") {
			crossref[cross] = crossTmp[:len(crossTmp)-1]
		} else {
			crossref[cross] = crossTmp
		}
	}

	return leader, crossref
}

// iterateTheLeader ...
func iterateTheLeader(headers map[string]*Header, blockBase []string, index int) {

	var nextLeader string
	thisLeader := blockBase[index]
	thisHeader := headers[blockBase[index]]

	if !(index >= len(blockBase)-1) {
		nextLeader = blockBase[index+1]
	} else {
		nextLeader = thisLeader
	}

	if isGoingUp(thisLeader, nextLeader) {
		resetThisAndJuniors(headers, nextLeader, thisLeader)
	} else {
		iterateThisHeader(thisHeader)
	}

}

// isGoingDown analyzes the current leader and the next leader to determine
// whether the tree is going up (returns true), down (returns false), or staying
// at the same level (returns false. up the tree indicates that thisLeader (e.g.,
// "ll.") is underneath nextLeader (e.g., "l.")
func isGoingUp(thisLeader string, nextLeader string) bool {
	if thisLeader > nextLeader {
		return true
	}
	return false
}

// resetThisAndJuniors ...
func resetThisAndJuniors(headers map[string]*Header, nextLeader string, thisLeader string) {
	for lead, header := range headers {
		if lead > nextLeader {
			if header.reset {
				header.currtVal = header.resetVal
			} else {
				if header.trigger == thisLeader {
					iterateThisHeader(header)
				}
			}
		}
	}
}

// iterateThisHeader ...
func iterateThisHeader(thisHeader *Header) {
	switch thisHeader.style {
	case 1, 2:
		thisHeader.currtVal = next_roman_upper(thisHeader.currtVal)
	case 3, 4:
		thisHeader.currtVal = next_roman_lower(thisHeader.currtVal)
	case 5, 6, 7, 8:
		thisHeader.currtVal = next_lettering(thisHeader.currtVal)
	case 9, 0:
		thisHeader.currtVal = next_numbering(thisHeader.currtVal)
	}
}

// deIterateThisHeader ...
func deIterateThisHeader(thisHeader string, style int) string {
	switch style {
	case 1, 2:
		return prev_roman_upper(thisHeader)
	case 3, 4:
		return prev_roman_lower(thisHeader)
	case 5, 6, 7, 8:
		return prev_lettering(thisHeader)
	default:
		return prev_numbering(thisHeader)
	}
}

// collateTheBlock replaces the cross references and joins up the block into a
// single string
func collateTheBlock(blockAsSlice []string, crossref map[string]string) string {
	collatedBlock := strings.Join(blockAsSlice, "\n\n")
	for pointer, replacer := range crossref {
		pointer = "|" + pointer + "|"
		collatedBlock = strings.Replace(collatedBlock, pointer, replacer, -1)
	}
	return collatedBlock
}
