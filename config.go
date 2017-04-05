package gapp

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config interface {
	Int(key string) int
	String(key string) string
	Bool(key string) bool
	Duration(key string) time.Duration
}

type Map []MapEntry

type MapEntry struct {
	Key     string
	Default interface{}
}

type config struct {
	ints      map[string]int
	strings   map[string]string
	bools     map[string]bool
	durations map[string]time.Duration
}

func New(prefix string, m Map) Config {
	result := &config{
		ints:      make(map[string]int),
		strings:   make(map[string]string),
		bools:     make(map[string]bool),
		durations: make(map[string]time.Duration),
	}
	for _, entry := range m {
		switch d := entry.Default.(type) {
		case string:
			result.strings[entry.Key] = loadString(prefix+entry.Key, d)
		case int:
			result.ints[entry.Key] = loadInt(prefix+entry.Key, d)
		case bool:
			result.bools[entry.Key] = loadBool(prefix+entry.Key, d)
		case time.Duration:
			result.durations[entry.Key] = loadDuration(prefix+entry.Key, d)
		default:
			panic(fmt.Errorf("invalid default type for key %s", entry.Key))
		}
	}

	return result
}

func (c *config) Int(key string) int {
	return c.ints[key]
}

func (c *config) String(key string) string {
	return c.strings[key]
}

func (c *config) Bool(key string) bool {
	return c.bools[key]
}

func (c *config) Duration(key string) time.Duration {
	return c.durations[key]
}

func loadString(varName string, dflt string) string {
	val := os.Getenv(varName)
	if val != "" {
		return val
	}
	return dflt
}

func loadInt(varName string, dflt int) int {
	val := os.Getenv(varName)
	if val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			panic(fmt.Sprintf("error parsing int %s: %s", varName, val))
		}
		return i
	}
	return dflt
}

func loadDuration(varName string, dflt time.Duration) time.Duration {
	val := os.Getenv(varName)
	if val != "" {
		dur, err := time.ParseDuration(val)
		if err != nil {
			panic(fmt.Sprintf("error parsing duration %s: %s", varName, val))
		}
		return dur
	}
	return dflt
}

func loadBool(varName string, dflt bool) bool {
	val := os.Getenv(varName)
	if val != "" {
		logToFile, err := strconv.ParseBool(val)
		if err != nil {
			panic(fmt.Sprintf("error parsing bool %s: %s", varName, val))
		}
		return logToFile
	}

	return dflt
}
