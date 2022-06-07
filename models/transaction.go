package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Transaction
// From + Nonce make a globally unique ID
type Transaction struct {
	// Inputs:
	Nonce  int64           `json:"nonce"`
	From   string          `json:"from"`
	To     string          `json:"to"`
	Amount decimal.Decimal `json:"amount"`
	// Arguments for contracts
	Args []string `json:"args"`

	// ID is the transaction ID, this also it's CID on IPFS and therefore it's hash.
	// This is based on hashing all the other data EXCEPT this
	ID string `json:"id"`

	// JWS the original transaction the user sent in, signed sealed and delivered
	// JWS string `json:"jws"`

	// Link to the current state document of this account
	StateID string `json:"stateID"`

	// RelatedTx is a link to a transaction when this user is the recipient
	// a tx is put into the recipients ledger with this field filled in which points to the senders transaction
	RelatedTxID string `json:"relatedTxID,omitempty"`

	// Contract stuff
	// Code is a reference to contract code. docker images must start with docker: and webassembly with wasm:
	Code string `json:"code,omitempty"`

	// chain bits
	CreatedAt time.Time `json:"createdAt"`
	// PreviousID is the previous transaction ID in the chain
	PreviousID string `json:"previousID"`
}

type TransactionInput struct {
	JWS string `json:"jws"`
}

type TransactionResponse struct {
	Tx *Transaction `json:"tx"`
}

func NewTransaction() *Transaction {
	return &Transaction{
		Amount: decimal.Zero,
	}
}
