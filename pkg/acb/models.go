package acb

import (
	"fmt"
	"strconv"
	"time"
)

// JournalEntry a record of security transaction
type JournalEntry struct {
	AccountNumber           string   `csv:"Account number"`
	AccountDescription      string   `csv:"Account description"`
	TradeDate               DateTime `csv:"Trade date"`
	SettlementDate          DateTime `csv:"Settlement date"`
	ProcessingDate          DateTime `csv:"Processing date"`
	Market                  string   `csv:"Market"`
	Symbol                  string   `csv:"Symbol"`
	Description             string   `csv:"Description"`
	Operation               string   `csv:"Operation"`
	Quantity                float64  `csv:"Quantity"`
	Price                   float64  `csv:"Price"`
	Commission              float64  `csv:"Commission"`
	NetAmount               float64  `csv:"Net amount"`
	BalanceAtSettlementDate float64  `csv:"Balance as at settlement date"`
	CurrentBalance          float64  `csv:"Current balance"`
}

// Holding a security holding
type Holding struct {
	AccountNumber string       `csv:"Account number"`
	TradeDate     DateTime     `csv:"Trade date"`
	Market        string       `csv:"Market"`
	Symbol        string       `csv:"Symbol"`
	Quantity      float64      `csv:"Quantity"`
	Price         float64      `csv:"Price"`
	ACB           DollarAmount `csv:"ACB"`
	MarketValue   float64      `csv:"Market value"`
	BookValue     float64      `csv:"Book value"`
}

// GainLoss a sell transaction resulting in a gain or loss
type GainLoss struct {
	AccountNumber string       `csv:"Account number"`
	TradeDate     DateTime     `csv:"Trade date"`
	Market        string       `csv:"Market"`
	Symbol        string       `csv:"Symbol"`
	Quantity      float64      `csv:"Quantity"`
	Price         float64      `csv:"Price"`
	Cost          float64      `csv:"Cost"`
	Proceed       float64      `csv:"Proceed"`
	Gain          DollarAmount `csv:"Gain"`
}

// GainSummary a tally of gains (or losses) for a ticker symbol for a given year
type GainSummary struct {
	Year          int          `csv:"Year"`
	AccountNumber string       `csv:"Account number"`
	Market        string       `csv:"Market"`
	Symbol        string       `csv:"Symbol"`
	Gain          DollarAmount `csv:"Gain"`
	AnnualTotal   DollarAmount `csv:"Annual total"`
}

// DateTime a time.Time with custom marshaling format to/from csv field: dd/mm/yyyy
type DateTime struct {
	time.Time
}

// MarshalCSV Convert the internal date as CSV string
func (date *DateTime) MarshalCSV() (string, error) {
	return date.Time.Format("02/01/2006"), nil
}

// String string representation of this date
func (date *DateTime) String() string {
	return date.String() // Redundant, just for example
}

// UnmarshalCSV Convert the CSV string as internal date
func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("02/01/2006", csv)
	return err
}

// DollarAmount a float64 that represents a dollar amount (with 2 decimal places)
type DollarAmount float64

func (d *DollarAmount) MarshalCSV() (string, error) {
	return d.String(), nil
}
func (d *DollarAmount) String() string {
	return fmt.Sprintf("%.2f", *d)
}

// UnmarshalCSV Convert the CSV string as internal float
func (d *DollarAmount) UnmarshalCSV(csv string) (err error) {
	if parsed, err := strconv.ParseFloat(csv, 64); err != nil {
		return err
	} else {
		*d = DollarAmount(parsed)
	}
	return nil
}
