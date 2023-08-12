package checker

import "strconv"

func Valid(parameters []string) bool {
	if len(parameters) != 0 && len(parameters) != 1 {
		return false
	}
	if _, err := strconv.Atoi(parameters[0]); err != nil {
		return false
	}
	return true
}
