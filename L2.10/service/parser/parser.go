package parser

import (
	"fmt"
	"os"
	"sort/external/model"

	"github.com/akamensky/argparse"
)

func ParseArgs() (model.Config, error) {
	parser := argparse.NewParser("sort", "UNIX sort strings")
	config := model.Config{}

	HelpStringsIsNums := "sort by numerical value (strings are interpreted as numbers)."
	StringsIsNums := parser.Flag("n", "numeric", &argparse.Options{Required: false, Help: HelpStringsIsNums})

	HelpReverseSort := "sort in reverse order"
	ReverseSort := parser.Flag("r", "revers", &argparse.Options{Required: false, Help: HelpReverseSort})

	HelpUniqueStrings := "do not display duplicate lines (only unique ones)"
	UniqueStrings := parser.Flag("u", "unique", &argparse.Options{Required: false, Help: HelpUniqueStrings})

	HelpMonthSort := "sort by month name"
	MonthSort := parser.Flag("M", "month", &argparse.Options{Required: false, Help: HelpMonthSort})

	HelpColumnSort := "sort by column (column) No. N (default separator is tab)."
	ColumnSort := parser.Int("k", "column", &argparse.Options{Required: false, Help: HelpColumnSort})

	HelpFilePath := "set filepath to sort"
	IncludePath := parser.String("i", "input", &argparse.Options{Required: true, Help: HelpFilePath})

	HelpOutputPath := "set output path"
	OutputPath := parser.String("o", "output", &argparse.Options{Default: "/cache/output.txt", Required: false, Help: HelpOutputPath})

	err := parser.Parse(os.Args)

	config.Numeric = *StringsIsNums
	config.ReverseSort = *ReverseSort
	config.UniqueStrings = *UniqueStrings
	config.MonthSort = *MonthSort
	config.ColumnSort = *ColumnSort
	config.IncludePath = *IncludePath
	config.OutPutPath = *OutputPath

	if err != nil {
		return model.Config{}, fmt.Errorf("Error ParseArgs: %v", err)
	}

	return config, nil
}
