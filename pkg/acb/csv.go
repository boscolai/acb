package acb

import (
	"encoding/csv"
	"github.com/gocarina/gocsv"
	"io"
	"os"
)

// GenerateAcbCsv generates 3 corresponding csv files related ACB calculation for the input csv file (infile):
//
//   - {infile}_gains.csv: a list of SELL transactions with Gains (or loss)
//   - {infile}_gains_summary: a list of Gains (Loss) by ticker symbols
//   - {infile}_holdings.csv: a list of current holdings
func GenerateAcbCsv(journalFile string, delimiter rune) error {
	// set the pipe as the delimiter for writing
	gocsv.TagSeparator = string(delimiter)
	// set ";" as delimiter for reading
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = delimiter
		return r // Allows use pipe as delimiter
	})

	inFile, err := os.Open(journalFile)
	if err != nil {
		return err
	}
	defer inFile.Close()

	var journal []*JournalEntry
	if err = gocsv.UnmarshalFile(inFile, &journal); err != nil {
		return err
	}
	results, err := ComputeACB(journal)
	if err != nil {
		return err
	}
	if err = writeCsv(journalFile+"_holdings.csv", results.Holdings); err != nil {
		return err
	}
	if err = writeCsv(journalFile+"_gains.csv", results.Gains); err != nil {
		return err
	}
	if err = writeCsv(journalFile+"_gains_summary.csv", results.GainsByYear); err != nil {
		return err
	}
	return nil
}

func writeCsv(filename string, records interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = gocsv.MarshalFile(records, f); err != nil {
		return err
	}
	return nil
}
