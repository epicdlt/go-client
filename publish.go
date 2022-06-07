package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/epicdlt/go-client/models"
)

func Publish(ctx context.Context, rpcURL string, msg *models.BroadcastMsg) error {
	client, err := NewClient(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to %q: %v", rpcURL, err)
	}
	defer client.Close()

	jsonValue, err := json.Marshal(msg)
	if err != nil {
		return (err)
	}
	url := rpcURL + "/verify"
	// fmt.Println(url)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return (err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return (err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("error response %v: %v", resp.StatusCode, string(bodyBytes))
	}

	// fmt.Println("response:", string(bodyBytes))
	return nil
}
