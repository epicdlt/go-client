package client

import (
	"context"

	"github.com/epicdlt/go-client/models"
	"github.com/treeder/gotils/v2"
)

func GetState(ctx context.Context, rpcURL, address string) (*models.StateMessage, error) {
	url := rpcURL + "/addr/" + address
	// fmt.Println("url:", url)
	stateM := &models.StateMessage{}
	err := gotils.GetJSON(url, stateM)
	return stateM, err
}

// GetObjectByHash this gets an object by it's content hash. Very similar to IPFS.
func GetObjectBytesByHash(ctx context.Context, rpcURL, hash string) ([]byte, error) {
	url := rpcURL + "/object/" + hash
	// fmt.Println("url:", url)
	return gotils.GetBytes(url)
}

// GetObjectByHash this gets an object by it's content hash. Very similar to IPFS.
func GetObjectByHash(ctx context.Context, rpcURL, hash string, dst interface{}) error {
	url := rpcURL + "/object/" + hash
	return gotils.GetJSON(url, dst)
}
