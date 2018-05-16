package main

import (
    "fmt"
    "os"
    "strings"
    "io/ioutil"
)

func readFile(filename string) []string {
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Println("Error: File not found", filename)
        return []string{}
    }
    lines := strings.Split(string(content), "\n")
    return lines
}

func startsWith(str, token string) bool {
    return len(str) >= len(token) && 
        str[:len(token)] == token
}

func exists(filename string) bool {
    _ , err := os.Stat(filename)
    return ! os.IsNotExist(err)
}

// Reads i3 base config and local config: Expects ./config-base
// and (optionally) ./$HOST.config
// Set variables with LET var_name <rest of line = value>
// Use variables in config as @var_name
// Variables set in the local config, $HOST.config, will override
// those set in config-base
func main() {
    hostname, _ := os.Hostname() // Getenv("HOST")
    lines := readFile("./config-base")
    var lines_local = []string{}
    if exists(hostname + ".config") {
        lines_local = readFile(hostname + ".config")
    }
    vars := make(map[string]string)

    // Read and parse variables
    for _, line := range append(lines, lines_local...) {
        if startsWith(line, "LET ") {
            toks := strings.Split(line, " ")
            vars[toks[1]] = strings.Join(toks[2:], " ")
        }
    }

    // Make replacements and send to stdout
    for _, line := range append(lines, lines_local...) {
        if startsWith(line, "LET ") || startsWith(line, "#LET") {
            continue
        }
        if strings.Contains(line, "@") {
            tokens := strings.Split(line, "@")
            for _, t := range tokens[1:] {
                token := strings.Split(t, " ")[0]
                line = strings.Replace(line, "@" + token, vars[token], 1)
            }
        }
        fmt.Println(line)
    }
}

