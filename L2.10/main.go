package main

import (
	"flag"
	"fmt"
)

type Config struct {
	ColumnSort     int
	StringsIsNums  bool
	ReverseSort    bool
	UniqueStrings  bool
	MonthSort      bool
	BlankTrailing  bool
	CheckAfterSort bool
}

func main() {
	config := parseArgs()
	fmt.Println(config)
}

func parseArgs() *Config {
	config := &Config{}

	flag.Parse()
	return config
}
