package acb

import (
	"fmt"
)

const (
	OpBuy         = "Buy"
	OpSell        = "Sell"
	OpCancel      = "Cancellation"
	OpSplit       = "Stock Split"
	OpExpiration  = "Expiration"
	OpAssignation = "Assignation"
	OpTransfer    = "Transfer"
)

type ACBResult struct {
	Holdings    []*Holding
	Gains       []*GainLoss
	GainsByYear []*GainSummary
}

// ComputeACB computes ACBs for all sell transactions within the given list of JournalEntry's
func ComputeACB(journal []*JournalEntry) (*ACBResult, error) {
	holdings := make(map[string]*Holding, 300)
	var gains []*GainLoss
	var date DateTime
	for i := len(journal) - 1; i >= 0; i-- {
		tx := journal[i]
		if d := tx.TradeDate; d.After(date.Time) {
			// new trade date
			date = d
		}
		if gain, err := updateHoldings(holdings, tx); err != nil {
			return nil, err
		} else if gain != nil {
			gains = append(gains, gain)
		}
	}
	results := &ACBResult{
		Gains: gains,
	}
	if len(gains) > 0 {
		results.GainsByYear = getGainsByYear(gains)
	}
	for _, holding := range holdings {
		results.Holdings = append(results.Holdings, holding)
	}
	return results, nil
}

func getGainsByYear(gains []*GainLoss) []*GainSummary {
	year := gains[0].TradeDate.Year()
	symbolGains := make(map[string]*GainSummary, len(gains))
	var summary []*GainSummary
	for _, gain := range gains {
		if year < gain.TradeDate.Year() {
			summary = tallyAnnualGains(summary, symbolGains)
			symbolGains = make(map[string]*GainSummary, len(gains))
			year = gain.TradeDate.Year()
		}
		if g, ok := symbolGains[gain.Symbol]; !ok {
			symbolGains[gain.Symbol] = &GainSummary{
				Year:          gain.TradeDate.Year(),
				AccountNumber: gain.AccountNumber,
				Market:        gain.Market,
				Symbol:        gain.Symbol,
				Gain:          gain.Gain,
			}
		} else {
			g.Gain += gain.Gain
		}
	}
	return tallyAnnualGains(summary, symbolGains)
}

func tallyAnnualGains(summary []*GainSummary, symbolGains map[string]*GainSummary) []*GainSummary {
	var year int
	var yearTotal DollarAmount
	for _, g := range symbolGains {
		year = g.Year
		summary = append(summary, g)
		yearTotal += g.Gain
	}
	if year > 0 {
		summary = append(summary, &GainSummary{Year: year, AnnualTotal: yearTotal})
	}
	return summary
}

func updateHoldings(holdings map[string]*Holding, tx *JournalEntry) (*GainLoss, error) {
	if tx.Symbol == "" {
		return nil, nil
	}
	var holding *Holding
	if h, ok := holdings[tx.Symbol]; !ok {
		holding = &Holding{
			AccountNumber: tx.AccountNumber,
			Market:        tx.Market,
			Symbol:        tx.Symbol,
		}
		holdings[tx.Symbol] = holding
	} else {
		holding = h
	}
	holding.TradeDate = tx.TradeDate
	var err error
	var gainLoss *GainLoss
	switch tx.Operation {
	case OpBuy:
		err = handleBuy(holding, tx)
	case OpSell:
		gainLoss, err = handleSell(holding, tx)
	case OpSplit:
		handleSplit(holding, tx)
	case OpCancel:
		handleCancel(holding, tx)
	case OpExpiration:
		handleExpiredOrAssigned(holding, tx)
	case OpAssignation:
		handleExpiredOrAssigned(holding, tx)
	case OpTransfer:
		handleTransfer(holding, tx)
	}

	if err != nil {
		return nil, err
	}
	if holding.Quantity == 0 {
		delete(holdings, tx.Symbol)
	}
	return gainLoss, nil
}

// handleTransfer handles transferred shares
// TODO: nbdb does not provide cost in tranferred units in transaction history ...
func handleTransfer(holding *Holding, tx *JournalEntry) {
	holding.Quantity += tx.Quantity
	if tx.Quantity > 0 {
		// transfer in
		holding.TransferInQuantity += tx.Quantity
	}
	// ignore transfer out for now
}

func handleBuy(holding *Holding, tx *JournalEntry) error {
	if err := assertSymbolNotEmpty(tx); err != nil {
		return err
	}
	holding.BookValue += -1 * tx.NetAmount
	holding.Quantity += tx.Quantity
	holding.Price = tx.Price
	holding.MarketValue = holding.Quantity * holding.Price
	holding.ACB = DollarAmount(holding.BookValue / holding.Quantity)
	return nil
}

func handleSell(holding *Holding, tx *JournalEntry) (*GainLoss, error) {
	if err := assertSymbolNotEmpty(tx); err != nil {
		return nil, err
	}
	gain := &GainLoss{
		AccountNumber: tx.AccountNumber,
		TradeDate:     tx.TradeDate,
		Market:        tx.Market,
		Symbol:        tx.Symbol,
		Quantity:      tx.Quantity,
		Price:         tx.Price,
		Cost:          float64(holding.ACB) * tx.Quantity,
		Proceed:       tx.NetAmount,
	}
	gain.Gain = DollarAmount(gain.Proceed - gain.Cost)

	// holding
	holding.BookValue += -1 * tx.NetAmount
	holding.Quantity -= tx.Quantity
	holding.Price = tx.Price
	holding.MarketValue = holding.Quantity * holding.Price

	// handle the special case where the shares were obtained from transfer-ins and the price were unknown
	if holding.TransferInQuantity >= tx.Quantity {
		holding.TransferInQuantity -= tx.Quantity
		gain.Cost = gain.Proceed
		gain.Gain = DollarAmount(0)
	}
	return gain, nil
}

func handleSplit(holding *Holding, tx *JournalEntry) {
	holding.Quantity += tx.Quantity
	// don't know price for current trade date
	holding.ACB = DollarAmount(holding.BookValue / holding.Quantity)
}

func handleExpiredOrAssigned(holding *Holding, tx *JournalEntry) {
	holding.Quantity += tx.Quantity
	// TODO if holding.Quantity > 0 ...  calculate capital loss
}

func handleCancel(holding *Holding, tx *JournalEntry) {
	if tx.NetAmount == 0 {
		return
	}
	if tx.NetAmount > 0 {
		// cancel buy
		holding.Quantity -= tx.Quantity
		holding.BookValue -= tx.NetAmount
	} else {
		// cancel sell
		holding.Quantity += tx.Quantity
		holding.BookValue += tx.NetAmount
	}
	holding.Price = tx.Price
	holding.MarketValue = holding.Quantity * holding.Price
	holding.ACB = DollarAmount(holding.BookValue / holding.Quantity)
}

func assertSymbolNotEmpty(tx *JournalEntry) error {
	if tx.Symbol == "" {
		return fmt.Errorf("no symbol specified in a journal entry: tradeDate=%s, operation=%s", tx.TradeDate.Format("2006-01-02"), tx.Operation)
	}
	return nil
}
