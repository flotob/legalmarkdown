package lmd

import (
	"regexp"
	"strconv"
	"strings"
)

// HandleTheHeaders is the primary parser function for parsing a block of structured headers.
//
// The function begins by calling the findTheBlock function with the full contents of the file.
// This function returns a boolean which tells the HandleTheHeaders function whether there is
// a block to parse or not. If there is no block, this function simply returns the full contents
// to the calling function. If there is a block, this is parsed into the pre_block, block, and
// post_block variables by the findTheBlock function.
//
// The block is then passed to the splitTheBlock function which breaks the block into its pieces
// as delimited by whether there is a structured_header at the beginning of the line or not.
// two slices are returned by the splitTheBlock function. The first is the blockBase slice which
// includes only the headers used in the block. This slice is used by the tree parsing to
// determine whether the parser is going up or down the tree or staying at the current location.
// The second slice is the blockAsSlice which contains the full contents of the block broken
// into pieces according to whether there is a structured header or not. This slice is what is
// used by the parser to fill out the contents.
//
// Once these slices are established, they are passed to the runTheHeaders function along with the
// headers map. This function has the primary parsing functionality. It iterates and the headers
// as determined by the trees established within the template file. This function returns a string
// which is the new block reassembled with the structured headers emplaced.
//
// Finally, the contents are reassembled by adding the pre_block variable to the block and post_block
// variables which is a long string that is returned to the calling function.
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

// runTheHeaders is the primary parsing function. It takes a map of the headers which have keys of
// the triggers to look for along with pointers to the Header structs for that particular header.
// It also takes the blockAsSlice slice which is the main contents of the block broken into the block's
// component pieces and the blockBase slice which contains all of the strings comprising the triggers.
// These triggers ease the parser to understand what is happening with regards to going up or down
// the trees.
//
// The function first establishes a map of strings for the cross references which is populated along
// with the structured header parsing.
//
// Then the function establishes two regular expressions. One for the old style ("llll.") headers and
// another for the new style ("l4.") headers. And it also establishes a variable for whether the header
// style is old style or new.
//
// After these preliminaries are established the function loops through the blockAsSlice first calling
// the replaceTheLeader function and then calling the iterateTheLeader function. The result of this
// parsing is then placed back into the blockAsSlice.
//
// Finally the function sends the map of established cross references and the reformulated blockAsSlice
// to the collateTheBlock function before returning to the calling function.
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

// replaceTheLeader has the longest function signature in all of legalmarkdown. It accepts a
// leader to be replaced, the map of headers, a block to replace the header in, the map of
// crossreferences and a boolean as to whether the header is old style or new.
//
// The function establishes a newLeader string which is used as a placeholder that is worked
// on throughout the function. Then it checks if there is a nil pointer within the headers.
// If the header returned from the map of headers is nil it simply returns not performing
// any of its parsing.
//
// Then the function checks whether there is a pre or preval suffix in the beforeVal for the
// header. If that is the case then the assemblePreVal function is called which is a specialized
// function that requires more computation than is necessary for a normal structured headers
// parsing function. If there is no pre or preval then the newLeader variable is simply the
// collation of the current header's beforeVal currtVal and afterVal strings.
//
// Once the newLeader is established either via a simple collation or via the assemblePreVal
// function then the function checks whether there is a cross reference in the block by calling
// the handleCrossReferences function.
//
// Next the function replaces the leader in the block, tightens up the strings and sets the
// indents to the appropriate level. Finally, the assembled block and the cross references map
// is returned to the calling function.
func replaceTheLeader(leader string, headers map[string]*Header, block string, crossref map[string]string, oldStyle bool) (string, map[string]string) {

	var newLeader string
	header := headers[leader]
    if header == nil {
        return block, crossref
    }

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

// assemblePreVal is a complex function which handles the assembly of the newLeader variable
// for the replaceTheLeader function when the header has a pre or preval call. There are two
// main challenges that this function has to overcome. The first is that the level above the
// current level is actually iterated (increased) after it is replaced by the runTheHeaders
// function. So the assemblePreVal function needs to deiterate this level. The second main
// challenge this function has to overcome is that pre and prevals can be nested which requires
// additional looping.
//
// The function starts by establishing some variables, namely the thisCurrtVal reformedLeader and
// prevHeader variables along with the thisHeader variable. Once these variables are established
// then the function checks if there is a recursive call to the function or not. If there is a
// recursive call to be made then it is made here before any further parsing is performed.
//
// Once the nested guard is ran, then the function assembles the leader and cleans up the string.
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
	if !nested {
		if strings.HasSuffix(thisBeforVal, "pre") && thisBeforVal != "pre" {
			reformedLeader = strings.Replace(thisBeforVal, "pre", "", 1)
		} else if strings.HasSuffix(thisBeforVal, "pre (") && thisBeforVal != "pre (" {
			reformedLeader = strings.Replace(thisBeforVal, "pre (", "", 1)
		} else if strings.HasSuffix(thisBeforVal, "preval") && thisBeforVal != "preval" {
			reformedLeader = strings.Replace(thisBeforVal, "preval", "", 1)
		} else {
			reformedLeader = ""
		}
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

// handleCrossReferences is a relatively simply function. It accepts a leader to search for
// a newLeader string, a block to search and the map of cross references. The function searches
// the block to determine whether the appropriate regular expression is matched (depending on
// whether it is old or new style headers) and if a match is found the appropriate newLeader
// that has been assembled by the replaceTheLeader function is placed into the map along with
// the trigger for the cross reference which is pulled out of the regular expression.
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

// iterateTheLeader increases the currtVal of the header by first determining
// whether the next index is going up the tree or not. If the next index is <
// the current index then all the headers at this level and below are reset
// if not then the current level is increased by calling the iterateThisHeader
// function.
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

// iterateThisHeader calls helper functions from the util.go file depending on the
// style of the current header.
func iterateThisHeader(thisHeader *Header) {
    if thisHeader == nil  {
        return
    }
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

// deIterateThisHeader is a helper function which is called by the preval parser and
// simply calls helper functions from the util.go file depending on the style of
// the current header.
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
