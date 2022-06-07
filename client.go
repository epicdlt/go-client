package client

import (
	"context"

	"github.com/epicdlt/go-client/models"
)

// Client is an interface for the web3 RPC API.
type Client interface {

	// SendTransaction ...
	SendTransaction(ctx context.Context, t *models.Transaction) error
	// SendRawTransaction sends the signed raw transaction bytes.
	SendRawTransaction(ctx context.Context, tx []byte) error
	// // Call executes a call without submitting a transaction.
	// Call(ctx context.Context, msg CallMsg) ([]byte, error)
	Close()
}

// NewClient creates a new client, backed by url (supported schemes "http", "https", "ws" and "wss").
func NewClient(url string) (Client, error) {
	return &client_{url: url}, nil
}

type client_ struct {
	url string
}

func (c *client_) Close() {
}

func (c *client_) SendTransaction(ctx context.Context, t *models.Transaction) error {
	// Sign it into JWS, send it
	return nil //c.SendRawTransaction(ctx, tx)
}

func (c *client_) SendRawTransaction(ctx context.Context, tx []byte) error {
	return nil
}
