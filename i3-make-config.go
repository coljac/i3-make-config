package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var default_directory string = "/home/coljac/build/conf/i3"
var default_filename string = "config-base"
var i3_config_location string = "/.config/i3/config"

func readFile(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error: File not found", filename)
		return []string{}
	}
	lines := strings.Split(string(content), "\n")
	return lines
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// Reads i3 base config and local config: Expects ./config-base
// and (optionally) ./$HOST.config
// Set variables with LET var_name <rest of line = value>
// Use variables in config as @var_name
// Variables set in the local config, $HOST.config, will override
// those set in config-base
func main() {
	var save_config = flag.Bool("w", false,
		"Save config to i3 default config file location.")
	flag.Parse()
	if len(flag.Args()) > 0 {
		default_directory = flag.Args()[0]
	}
	// fmt.Println(*save_config, default_directory)

	output_writer := os.Stdout
	if *save_config {
		file, err := os.Create(os.Getenv("HOME") + "/" +
			i3_config_location)
		if err != nil {
			fmt.Println("error:", err)
			return
		} else {
			output_writer = file
		}
		defer file.Close()
	}

	hostname, _ := os.Hostname() // Getenv("HOST")
	lines := readFile(default_directory + "/" + default_filename)

	var lines_local = []string{}
	if exists(default_directory + "/" + hostname + ".config") {
		lines_local = readFile(default_directory + "/" + hostname + ".config")
	}
	vars := make(map[string]string)

	// Read and parse variables
	for _, line := range append(lines, lines_local...) {
		if strings.HasPrefix(line, "LET ") {
			toks := strings.Split(line, " ")
			vars[toks[1]] = strings.Join(toks[2:], " ")
		}
	}

	// Make replacements and send to stdout
	for _, line := range append(lines, lines_local...) {
		if strings.HasPrefix(line, "LET ") ||
			strings.HasPrefix(line, "#LET") {
			continue
		}
		if strings.Contains(line, "@@") {
            line = strings.Replace(line, "@@", "@", 1)
        } else if strings.Contains(line, "@") {
			tokens := strings.Split(line, "@")
			for _, t := range tokens[1:] {
				token := ""
				for _, c := range t {
					isletter, _ := regexp.MatchString("[A-z0-9_]",
						string(c))
					if isletter {
						token += string(c)
					} else {
						break
					}
				}
				// token := strings.Split(t, " ")[0]
				line = strings.Replace(line, "@"+token, vars[token], 1)
			}
		}
		output_writer.WriteString(line + "\n")
	}
}
