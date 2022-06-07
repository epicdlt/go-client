package models

type Broadcast struct {
	FromTxID  string       `json:"fromTxID"`
	FromTx    *Transaction `json:"fromTx"`
	ToTxID    string       `json:"toTxID"`
	ToTx      *Transaction `json:"toTx"`
	FromState *State       `json:"fromState"`
	ToState   *State       `json:"toState"`
	CodeState []byte       `json:"codeState"`
}

type BroadcastMsg struct {
	PubKey string `json:"pubKey"`
	JWS    string `json:"jws"`
}
