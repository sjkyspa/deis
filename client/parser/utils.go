package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func safeGetValue(args map[string]interface{}, key string) string {
	if args[key] == nil {
		return ""
	}
	return args[key].(string)
}

func safeGetJSONConfig(args map[string]interface{}, key string) (map[string]interface{}, error) {
	if args[key] == nil {
		return make(map[string]interface{}), nil
	}

	config := args[key].(string)

	configJSON := make(map[string]interface{})
	if strings.HasPrefix(config, "@") {
		configFile := strings.TrimLeft(config, "@")
		content, err := ioutil.ReadFile(configFile)
		if err != nil {
			return make(map[string]interface{}), fmt.Errorf("Please check the file realily exists")
		}

		err = json.Unmarshal(content, &configJSON)
		if err != nil {
			return make(map[string]interface{}), fmt.Errorf("Please ensure the json in config file is valid")
		}
	} else {
		err := json.Unmarshal([]byte(config), &configJSON)

		if err != nil {
			return make(map[string]interface{}), fmt.Errorf("Pleases provide the valid json as the additional parameters")
		}
	}

	return configJSON, nil
}

func responseLimit(limit string) (int, error) {
	if limit == "" {
		return -1, nil
	}

	return strconv.Atoi(limit)
}

// PrintUsage runs if no matching command is found.
func PrintUsage() {
	fmt.Fprintln(os.Stderr, "Found no matching command, try 'deis help'")
	fmt.Fprintln(os.Stderr, "Usage: deis <command> [<args>...]")
}

func printHelp(argv []string, usage string) bool {
	if len(argv) > 1 {
		if argv[1] == "--help" || argv[1] == "-h" {
			fmt.Print(usage)
			return true
		}
	}

	return false
}
