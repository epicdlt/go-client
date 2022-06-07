package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/epicdlt/go-client/models"
	"github.com/shopspring/decimal"
	"github.com/treeder/gotils/v2"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	// The latency in milliseconds
	MLatencyMs = stats.Float64("tx/submit/latency", "The latency in milliseconds per submit", "ms")

	// Counts/groups the lengths of lines read in.
	MLineLengths = stats.Int64("tx/submit/abc", "The distribution of line lengths", "By")
)
var (
	KeyMethod, _ = tag.NewKey("method")
	KeyStatus, _ = tag.NewKey("status")
	KeyError, _  = tag.NewKey("error")
)
var (
	LatencyView = &view.View{
		Name:        "demo/latency",
		Measure:     MLatencyMs,
		Description: "The distribution of the latencies",

		// Latency in buckets:
		// [>=0ms, >=25ms, >=50ms, >=75ms, >=100ms, >=200ms, >=400ms, >=600ms, >=800ms, >=1s, >=2s, >=4s, >=6s]
		Aggregation: view.Distribution(0, 25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000, 4000, 6000),
		TagKeys:     []tag.Key{KeyMethod}}
)

func init() {
	// Register the views
	if err := view.Register(LatencyView); err != nil {
		log.Fatalf("Failed to register views: %v", err)
	}
}

// SubmitTx ...
// codeRef is a docker image or wasm reference
func SubmitTx(ctx context.Context, rpcURL, privateKey string, amount decimal.Decimal, to string, codeRef string, args []string) (*models.Transaction, error) {
	ctx, err := tag.New(ctx, tag.Insert(KeyMethod, "SubmitTx"), tag.Insert(KeyStatus, "OK"))
	if err != nil {
		return nil, err
	}
	startTime := time.Now()
	client, err := NewClient(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to %q: %v", rpcURL, err)
	}
	defer client.Close()

	if codeRef != "" {
		// then we're deploying, so to isn't required
	} else {
		if to == "" {
			return nil, errors.New("The recipient address cannot be empty")
		}
	}
	defer func() {
		stats.Record(ctx, MLatencyMs.M(sinceInMilliseconds(startTime)))
	}()
	// todo: verify valid TO address
	// start := time.Now()
	pubHashed, err := PublicKeyFromPrivate(ctx, privateKey)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("BENCH PublicKeyFromPrivate time: %v\n", time.Since(start))

	// start = time.Now()
	nonceURL := rpcURL + "/addr/" + pubHashed
	state := &models.StateMessage{}
	err = gotils.GetJSON(nonceURL, state)
	if err != nil {
		return nil, (err)
	}
	// fmt.Printf("BENCH GetJSON time: %v\n", time.Since(start))
	// fmt.Printf("state: %+v amount: %v nonce: %v\n", state.State, state.State.Balance, state.State.Nonce)
	t := &models.Transaction{
		From:   pubHashed,
		To:     to,
		Amount: amount,
		Nonce:  state.State.Nonce + 1,
		Code:   codeRef,
		Args:   args, // for contracts
	}
	// fmt.Println(t)

	serialized, err := Sign(ctx, privateKey, t)
	if err != nil {
		return nil, (err)
	}
	// fmt.Println("jws", serialized)

	jsonValue, err := json.Marshal(map[string]interface{}{
		"jws": serialized, // todo: maybe call this jws-full and also accept compact under jws
	})
	if err != nil {
		return nil, (err)
	}

	// start = time.Now()
	url := rpcURL + "/tx/" + pubHashed
	// fmt.Println(url)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, (err)
	}
	defer resp.Body.Close()
	// fmt.Printf("BENCH Post time: %v\n", time.Since(start))
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, (err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Error response %v: %v", resp.StatusCode, string(bodyBytes))
	}

	err = json.Unmarshal(bodyBytes, t)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response, but it was successful...: %v", err)
	}
	return t, err
}

func sinceInMilliseconds(startTime time.Time) float64 {
	return float64(time.Since(startTime).Nanoseconds()) / 1e6
}
