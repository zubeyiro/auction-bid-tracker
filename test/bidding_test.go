package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/zubeyiro/auction-bid-tracker/src/server/data"
	"github.com/zubeyiro/auction-bid-tracker/src/server/rest"
)

func TestMultipleBidding(t *testing.T) {
	db := data.Init()
	go rest.StartRestServer(db)

	startingBid := 15000

	for i := 0; i < 10000; i++ {
		startingBid++
		resp, err := http.Post("http://localhost:8080/auctions/2/bids", "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(`{"user_id": %v, "amount": %v}`, 1, startingBid))))

		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", resp.StatusCode, http.StatusOK)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			t.Fatal(err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			t.Fatal(err)
		}
		if !result["status"].(bool) {
			t.Fatal(result)
		}

	}
}
