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
	pos := getFirstSeparatorPos(args)
	if pos >= 0 {
		return args[0:pos], nil
	}
	return nil, fmt.Errorf("seperator not found %s", separatorArg)
}

func getFirstSeparatorPos(args []string) int {
	for i, item := range args {
		if item == separatorArg {
			return i
		}
	}
	return -1
}

func getAllSeparatorPos(args []string) []int {
	separatorPositions := make([]int, 0)
	for i, item := range args {
		if item == separatorArg {
			separatorPositions = append(separatorPositions, i)
		}
	}
	return separatorPositions
}

// responsible for getting command to execute
func CommandToExecute(args []string) ([][]string, error) {
	positions := getAllSeparatorPos(args)
	commands := make([][]string, 0)
	for i := 0; i < len(positions); i++ {
		if i == len(positions) -1{
			commands = append(commands, args[positions[i]:len(args)])
		}else{
			commands = append(commands, args[positions[i]:positions[i+1]])
		}
	}
	if len(commands) == 0 {
		return nil, fmt.Errorf("seperator not found %s", separatorArg)
	}
	return commands, nil
}
