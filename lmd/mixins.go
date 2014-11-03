package lmd

import (
	"fmt"
	"regexp"
	"strings"
)

// HandleMixins is the primary handling function for the mixin and optional clause parsing exercise
// which is a subset of the overall parse job.
//
// The first thing which the function will do is to extract the parameters which we know are applicable
// to the structured_headers portion of the overall parse job rather than the mixins portion of the
// overall parse job. These parameters are parked into a new map during the remainder of the mixin
// operation.
//
// After the parameters we do not want to handle have been parked, the function then runs through the
// optional clauses parse adding those optional clauses which have been turned on, and taking out of
// the final document those optional clauses which have been turned off.
//
// When the optional clauses have finalized, then the function runs through the simple text mixins
// and replaces all of the keys established in the text of the template with the values established
// in the parameters.
//
// Finally a simple cleanup function is called to compress extraneous white space and then the function
// returns the parsed and corrected contents along with the parked parameters which are relevant to
// the structured_headers phase of the overall parse.
func HandleMixins(contents string, parameters map[string]string) (string, map[string]string) {

	// create a parking_lot variable and park the parameters we know we don't want to mess with
	// during this phase
	var params_parking_lot map[string]string
	parameters, params_parking_lot = prepareParamsParkingLot(parameters)

	// run the optional clauses and mixins
	contents, parameters = runOptionalClauses(contents, parameters)
	contents = runTextMixins(contents, parameters)

	// perform some simple clean up
	contents = cleanUpPostMixins(contents)

	return contents, params_parking_lot
}

// prepareParamsParkingLot pulls out parameters we *know* we don't want to mess with during the mixin
// phase of the overall parse.
func prepareParamsParkingLot(parameters map[string]string) (map[string]string, map[string]string) {

	// define the parameters we want to blacklist in a slice of strings
	parameters_blacklist_strings := []string{"level-[0-9]", "no-reset", "no-indent", "level-style"}

	// compile those strings into a slice of regular expressions.
	parameters_blacklist_regexs := []*regexp.Regexp{}
	for _, to_compile := range parameters_blacklist_strings {
		parameters_blacklist_regexs = append(parameters_blacklist_regexs, regexp.MustCompile(to_compile))
	}

	// prepare a parking lot, loop through each of the parameters, compare against each of the
	// blacklisted parameters and if there's a match add to the parking_lot while deleting from
	// the paramters list
	parking_lot := make(map[string]string)
	for key, val := range parameters {
		for _, black_listed := range parameters_blacklist_regexs {
			if black_listed.MatchString(key) {
				parking_lot[key] = val
				delete(parameters, key)
			}
		}
	}

	return parameters, parking_lot

}

// runOptionalClauses is the guardian of the circuit which needs to be run to ensure that all of the
// optional clauses (nested and single) have been properly handled. it first sends the current parameters
// map to the separateOptionalClauses function to get those which need to be added and those which need
// to be deleted. after that it goes into an infinite loop which runs through the runThisOptionalClause
// pattern for both the to_add and to_delete slices until both have been completely exhausted (which occurs
// after all the "true" and "false" have been expunged from the parameters map).
func runOptionalClauses(contents string, parameters map[string]string) (string, map[string]string) {

	clauses_to_add, clauses_to_rem := separateOptionalClauses(parameters)

	for {

		// first pass
		contents, clauses_to_add = runThisOptionalClause(contents, parameters, clauses_to_add, true)
		contents, clauses_to_rem = runThisOptionalClause(contents, parameters, clauses_to_rem, false)

		// reload the slices against the parameters
		clauses_to_add, clauses_to_rem = separateOptionalClauses(parameters)

		// check if the reloaded slices are of 0 and 0 length, break if yes
		if len(clauses_to_add) == 0 && len(clauses_to_rem) == 0 {
			break
		}

	}

	return contents, parameters
}

// runTextMixins is a very simple function. it loops through the remaining parameters (it is called
// after the optional clauses have run) and replaces the keys with the values from the parameters
// map which remains.
func runTextMixins(contents string, parameters map[string]string) string {

	for to_replace, replacer := range parameters {
		mixin_pattern := regexp.MustCompile(fmt.Sprintf(`(\{\{%v\}\})`, to_replace))
		if mixin_pattern.MatchString(contents) {
			contents = mixin_pattern.ReplaceAllString(contents, replacer)
		}
	}

	return contents
}

// cleanUpPostMixins is a simple function which compresses excessive whitespace.
func cleanUpPostMixins(contents string) string {

	// when there are more than two new_lines, squeeze those into two.
	too_many_lines := regexp.MustCompile(`\n\n+`)
	if too_many_lines.MatchString(contents) {
		contents = too_many_lines.ReplaceAllString(contents, "\n\n")
	}

	// when there are more than two spaces, squeeze those into one.
	too_many_space := regexp.MustCompile(` {2,}`)
	if too_many_space.MatchString(contents) {
		contents = too_many_space.ReplaceAllString(contents, " ")
	}

	return contents
}

// separateOptionalClauses checks the parameters map and pulls out those parameters with a
// key:value of _:true and adds those values to the clauses_to_keep slice. those parameters
// with a key:value of _:false are added to the clauses_to_delete slice. both slices are then
// returned to the calling function.
func separateOptionalClauses(parameters map[string]string) ([]string, []string) {

	clauses_to_dele := []string{}
	clauses_to_keep := []string{}

	for clause, status := range parameters {
		if status == "true" {
			clauses_to_keep = append(clauses_to_keep, clause)
		} else if status == "false" {
			clauses_to_dele = append(clauses_to_dele, clause)
		}
	}

	return clauses_to_keep, clauses_to_dele
}

// runThisOptionalClause is the primary optional clause parsing function. it parses the contents against
// a set of optional clauses, running through the slice of strings it is given in the clauses slice one
// time. in general this function will be called a minimum of twice if there are optional clauses in the
// template file -- once for those optional clauses turned on and once for optional clauses turned off.
func runThisOptionalClause(contents string, parameters map[string]string, clauses []string, add_or_delete bool) (string, []string) {

	for _, clause := range clauses {

		// there are two relevant regex's one for primary optional clauses and one to check if there are
		// nested optional clauses
		pri_pattern := regexp.MustCompile(fmt.Sprintf(`(?sm)\[\{\{%v\}\}\s*?(.*?\n*?)\]`, clause))
		sub_pattern := regexp.MustCompile(`(?sm)\[\{\{(\S+?)\}\}\s*?`)

		// first we check if there is a match within the overall content. are there any optional clauses
		// at all? if there are then we dump the found group from the regex into the sub_clause string.
		var sub_clause string
		if pri_pattern.MatchString(contents) {
			sub_clause = pri_pattern.FindAllStringSubmatch(contents, -1)[0][1]
		}

		// if the sub_clause (what is between the square brackets) has another nested optional clause
		// then the pattern will break because we do not want to do anything until we have found the
		// inner-most nested optional clause which will make our non-greedy regex work.
		if sub_pattern.MatchString(sub_clause) {
			continue
		}

		// the add_or_delete variable is a boolean if it is an add then the subclause will replace the
		// overall found pattern, if it is false then the whole thing will be replaced with an empty
		// string.
		if add_or_delete {
			contents = pri_pattern.ReplaceAllString(contents, strings.TrimSpace(sub_clause))
		} else {
			contents = pri_pattern.ReplaceAllString(contents, "")
		}

		// the final step is to delete the matched clause from the parameters if it is not found in the
		// contents any longer. this is because we will refresh the loop overall in the runOptionalClauses
		// function so we don't want to mess with this key any longer.
		if !pri_pattern.MatchString(contents) {
			delete(parameters, clause)
		}
	}

	return contents, clauses
}
