package types

type Transaction struct {
	Hash      string `json:"hash"`
	From      string `json:"from"`
	To        string `json:"to"`
	Timestamp int64  `json:"timestamp"`
}

func NewTransaction(hash, from, to string, timestamp int64) Transaction {
	return Transaction{
		Hash:      hash,
		From:      from,
		To:        to,
		Timestamp: timestamp,
	}
}
