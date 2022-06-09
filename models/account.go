package models

// Account for details about a destination account.
// An account can have a name like @treeder and the account details will have the list of signers and any other metadata.
type Account struct {

	// Name eg: @username, must be globally unique
	Name string `json:"name"`

	// Signers list of public keys for an account
	Signers []string `json:"signers"`
}
