package utils

import (
	"fmt"
	"sort"
)

var (
	separatorArg = "-c"
)

// responsible for validating args and panicking if necessary
func ValidateArgs(args []string) error {
	argsSorted := make([]string, len(args))
	copy(argsSorted, args)
	sort.Strings(argsSorted)
	if i := sort.SearchStrings(argsSorted, separatorArg); i == len(argsSorted) || argsSorted[i] != separatorArg {
		return fmt.Errorf("missing the argument %s", separatorArg)
	}

	return nil
}

// responsible for getting files to monitor
func GetFilesToMonitor(args []string) ([]string, error) {
	pos := getSeparatorPos(args)
	if pos >= 0 {
		return args[0:pos], nil
	}
	return nil, fmt.Errorf("seperator not found %s", separatorArg)
}

func getSeparatorPos(args []string) int {
	for i, item := range args {
		if item == separatorArg {
			return i
		}
	}
	return -1
}

// responsible for getting command to execute
func CommandToExecute(args []string) ([]string, error) {
	pos := getSeparatorPos(args)
	if pos >= 0 {
		return args[pos:len(args)], nil
	}
	return nil, fmt.Errorf("seperator not found %s", separatorArg)
}
