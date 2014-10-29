package lmd

import (
	"fmt"
	"regexp"
	"strings"
)

// HandleMixins .... does shit.
func HandleMixins(contents string, parameters map[string]string) (string, map[string]string) {
	// pull out parameters we *know* we don't want to mess with.
	var params_parking_lot map[string]string
	parameters, params_parking_lot = prepare_params_parking_lot(parameters)
	contents, parameters = run_optional_clauses(contents, parameters)
	contents, parameters = run_text_mixins(contents, parameters)
	contents = clean_up_post_mixins(contents)
	parameters = params_parking_lot
	return contents, parameters
}

// prepare_params_parking_lot pulls out parameters we *know* we don't want to mess with.
func prepare_params_parking_lot(parameters map[string]string) (map[string]string, map[string]string) {

	// define the parameters we want to blacklist in a slice of strings
	parameters_blacklist_strings := []string{"level-[0-9]", "no-reset", "no-indent", "level-style"}

	// compile those strings into a slice of regular expressions.
	parameters_blacklist := []*regexp.Regexp{}
	for _, to_compile := range parameters_blacklist_strings {
		parameters_blacklist = append(parameters_blacklist, regexp.MustCompile(to_compile))
	}

	// prepare a parking lot, loop through each of the parameters, compare against each of the
	//   blacklisted parameters and if there's a match add to the parking_lot while deleting from
	//   the paramters list
	parking_lot := make(map[string]string)
	for key, val := range parameters {
		for _, blackie := range parameters_blacklist {
			if blackie.MatchString(key) {
				parking_lot[key] = val
				delete(parameters, key)
			}
		}
	}

	return parameters, parking_lot

}

// run_optional_clauses ....
func run_optional_clauses(contents string, parameters map[string]string) (string, map[string]string) {
	clauses_added, clauses_deleted := separate_optional_clauses(parameters)
	for {
		contents, clauses_added = run_this_optional_clause(contents, clauses_added, true, parameters)
		contents, clauses_deleted = run_this_optional_clause(contents, clauses_deleted, false, parameters)
		clauses_added, clauses_deleted = separate_optional_clauses(parameters)
		if len(clauses_added) == 0 && len(clauses_deleted) == 0 {
			break
		}
	}
	return contents, parameters
}

func run_text_mixins(contents string, parameters map[string]string) (string, map[string]string) {
	for to_replace, replacer := range parameters {
		mixin_pattern := regexp.MustCompile(fmt.Sprintf(`(\{\{%v\}\})`, to_replace))
		if mixin_pattern.MatchString(contents) {
			contents = mixin_pattern.ReplaceAllString(contents, replacer)
		}
	}
	return contents, parameters
}

func clean_up_post_mixins(contents string) string {
	too_many_lines := regexp.MustCompile(`\n\n+`)
	if too_many_lines.MatchString(contents) {
		contents = too_many_lines.ReplaceAllString(contents, "\n\n")
	}
	too_many_space := regexp.MustCompile(` {2,}`)
	if too_many_space.MatchString(contents) {
		contents = too_many_space.ReplaceAllString(contents, " ")
	}

	return contents
}

func separate_optional_clauses(parameters map[string]string) ([]string, []string) {
	clauses_to_delete := []string{}
	clauses_to_keep := []string{}
	for clause, status := range parameters {
		if status == "true" {
			clauses_to_keep = append(clauses_to_keep, clause)
		} else if status == "false" {
			clauses_to_delete = append(clauses_to_delete, clause)
		}
	}
	return clauses_to_keep, clauses_to_delete
}

func run_this_optional_clause(contents string, clauses []string, add_or_delete bool, parameters map[string]string) (string, []string) {
	for _, clause := range clauses {
		sub_clause, _ := "", ""
		pri_pattern := regexp.MustCompile(fmt.Sprintf(`(?m)(\[\{\{%v\}\}\s*?)(.*?\n*?)(\])`, clause))
		sub_pattern := regexp.MustCompile(`(?m)\[\{\{(\S+?)\}\}\s*?`)
		if pri_pattern.MatchString(contents) {
			da_finding := pri_pattern.FindAllStringSubmatch(contents, -1)
			sub_clause = da_finding[0][2]
		}
		if sub_pattern.MatchString(sub_clause) {
			continue
		}
		if add_or_delete {
			contents = pri_pattern.ReplaceAllString(contents, strings.TrimSpace(sub_clause))
		} else {
			contents = pri_pattern.ReplaceAllString(contents, "")
		}
		if !pri_pattern.MatchString(contents) {
			clauses = removeFromSlice(clauses, clause)
			delete(parameters, clause)
		}
	}
	return contents, clauses
}

// back_out_of_parking_lot adds the parked parameters back into the full parameters list before those
// are sent back to the primary parser.
// func back_out_of_parking_lot(parameters map[string]string, params_parking_lot map[string]string) map[string]string {
//     for key, val := range params_parking_lot {
//         parameters[key] = val
//     }
//     return parameters
// }

func includedInThisSlice(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func removeFromSlice(remove_from_me []string, ele_to_remove string) []string {
	var slice_with_ele_removed []string
	for i, ele := range remove_from_me {
		if ele == ele_to_remove {
			slice_with_ele_removed = append(remove_from_me[:i], remove_from_me[i+1:]...)
		}
	}
	return slice_with_ele_removed
}

func areTheseSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
