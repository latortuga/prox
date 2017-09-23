package prox

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ParseEnvFile reads environment variables that should be set on all processes
// from the ".env" file.
//
// The format of the ".env" file is expected to be a newline separated list of
// key=value pairs which represent the environment variables that should be used
// by all started processes. Trimmed lines which are empty or start with a "#"
// are ignored and can be used to add comments.
//
// All values are expanded using the Environment. If a value refers to a
// variable that is not set in e then it is replaced with the empty string.
func (e Environment) ParseEnvFile(r io.Reader) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		line = e.Expand(line)
		e.Set(line)
	}

	return s.Err()
}

// Environment is a set of key value pairs that are used to set environment
// variables for processes.
type Environment map[string]string

// SystemEnv creates a new Environment from the operating systems environment
// variables.
func SystemEnv() Environment {
	return NewEnv(os.Environ())
}

// NewEnv creates a new Environment and immediately sets all given key=value
// pairs.
func NewEnv(values []string) Environment {
	env := Environment{}
	env.SetAll(values)
	return env
}

// Set splits the input string at the first "=" character (if any) and sets the
// resulting key and value on e.
func (e Environment) Set(s string) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) == 1 {
		parts[1] = ""
	}

	e[parts[0]] = strings.TrimSpace(parts[1])
}

// SetAll assigns a list of key=value pairs on e.
func (e Environment) SetAll(vars []string) {
	for _, v := range vars {
		e.Set(v)
	}
}

// List returns all variables of e as a list of key=value pairs.
func (e Environment) List() []string {
	vars := make([]string, 0, len(e))
	for key, value := range e {
		vars = append(vars, fmt.Sprintf("%s=%s", key, value))
	}
	return vars
}

// Expand replaces ${var} or $var in the input string with the corresponding
// values of e. If the variable is not found in e then an empty string is
// returned.
func (e Environment) Expand(input string) string {
	return os.Expand(input, func(key string) string {
		return e[key]
	})
}
