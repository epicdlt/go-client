package models

type ServerTransaction struct {
	FromTxID string // hash of transaction
	FromTx   *Transaction
	ToTxID   string // hash of transaction
	ToTx     *Transaction
}
