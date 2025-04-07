package commands

import (
	"fmt"
	"github.com/boscolai/acb/pkg/acb"
	"github.com/spf13/cobra"
)

const longDescText = `
For each given csv file containing security transactions in reversed chronicle order, 
this tool computes ACB values for each sell and outputs 3 corresponding csv files:

  - {infile}_gains.csv: a list of SELL transactions with Gains (or loss)
  - {infile}_gains_summary: a list of Gains (Loss) by ticker symbols
  - {infile}_holdings.csv: a list of current holdings

Usage example:
    acb -d \; account1.csv account2.csv
`

var (
	delimiter string
	rootCmd   = &cobra.Command{
		Use:   "acb csv_file1 [csv_file2 [csv_file2 ...]]",
		Short: "Compute average cost base (ACB) of securities",
		Long:  longDescText,
		RunE:  exec,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&delimiter, "delimiter", "d", ",", "field delimiter in csv file")
}

func exec(command *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no input csv file given")
	}
	for _, file := range args {
		if err := acb.GenerateAcbCsv(file, rune(delimiter[0])); err != nil {
			return fmt.Errorf("acb failed for file %s: %s", file, err.Error())
		}
	}
	return nil
}
