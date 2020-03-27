package CLIParser

import (
	"os"
	"strings"
)

type CLIFlag struct {
	Name  string
	Value string
}

func IsPresentedCLIFlagValueByName(CLIFlag []CLIFlag, name string) bool {
	for _, curFlag := range CLIFlag {
		if curFlag.Name == name {
			return true
		}
	}
	return false
}

func GetCLIFlagValueByName(CLIFlags []CLIFlag, name string) string {
	for _, curFlag := range CLIFlags {
		if curFlag.Name == name {
			return curFlag.Value
		}
	}
	return "" // return empty string if flag not found
}

//-------------------------------------------------------------------------------
func GetAction() string {
	if len(os.Args) > 1 {
		action := ""
		for curArgIndex, curArg := range os.Args {
			if curArgIndex > 0 {
				if curArg[0] != '-' {
					action = curArg
					break
				}
			}
		}
		return action
	} else {
		return "" // return empty string if command not found
	}
}

func GetParams() map[string]string {
	var result map[string]string
	result = make(map[string]string)
	return result
}

// parse parameters like -y=someText and return it all.
// name is begin from first letter of argument. Ends on "="
// Value begin from "=" and end with space
func ParseAllCLIFlags() map[string]string {
	var result map[string]string
	result = make(map[string]string)
	for _, curCLIArgument := range os.Args {
		if curCLIArgument[0] == '-' {
			equalPos := strings.Index(curCLIArgument, "=")
			name := curCLIArgument[1:equalPos]
			value := curCLIArgument[equalPos+1:]
			result[name] = value
		}
	}
	return result
}

func GetAllVariantsOfFlagSeparatedBy(flag string, separator string) []string {
	return strings.Split(flag, separator)
}
