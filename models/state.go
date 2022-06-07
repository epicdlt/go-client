package models

import "github.com/shopspring/decimal"

type State struct {
	Balance decimal.Decimal `json:"balance"`
	Nonce   int64           `json:"nonce"`

	// For contracts:
	// Code is the docker image to use, docker images start with `docker:`, webassembly can start with `wasm:`
	Code string `json:"code,omitempty"`
	// Contract state is arbitrary data state for the contract. This is the ID to that state.
	CodeStateID string `json:"codeStateID,omitempty"`

	// ID hash representation
	// ID string `json:"-"`
}

type StateMessage struct {
	State    *State       `json:"state"`
	LastTxID string       `json:"lastTxID"`
	LastTx   *Transaction `json:"lastTx"`
}

func NewState() *State {
	return &State{
		Balance: decimal.Zero,
	}
}

// AddressRecord is a mapping from a public address/key to it's most recent transaction
// To look up address data, get this tx_id value, then get the tx by ID (it's content hash).
type AddressRecord struct {
	TxID string `json:"txID,omitempty"`
	// Contract string `json:"tx_id3"`
}

type ObjectResponse struct {
	Message string      `json:"message"`
	Object  interface{} `json:"object"`
}

type CodeStateResponse struct {
	Message string      `json:"message"`
	State   interface{} `json:"state"`
}
