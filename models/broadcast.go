package models

type Broadcast struct {
	FromTxID  string       `json:"from_tx_id"`
	FromTx    *Transaction `json:"from_tx"`
	ToTxID    string       `json:"to_tx_id"`
	ToTx      *Transaction `json:"to_tx"`
	FromState *State       `json:"from_state"`
	ToState   *State       `json:"to_state"`
	CodeState []byte       `json:"code_state"`
}

type BroadcastMsg struct {
	PubKey string `json:"pub_key"`
	JWS    string `json:"jws"`
}
