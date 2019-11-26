package auctions

// BidSerializer is response format
type BidSerializer struct {
	Amount    float64 `json:"amount"`
	CreatedAt int64   `json:"created_at"`
}

// AuctionSerializer is response format
type AuctionSerializer struct {
	Title         string  `json:"title"`
	StartingPrice float64 `json:"starting_price"`
	CreatedAt     int64   `json:"created_at"`
	ExpireAt      int64   `json:"expires_at"`
}

// AuctionSerializerWithBid is response format
type AuctionSerializerWithBid struct {
	Title         string        `json:"title"`
	StartingPrice float64       `json:"starting_price"`
	CreatedAt     int64         `json:"created_at"`
	ExpireAt      int64         `json:"expires_at"`
	LastBid       BidSerializer `json:"last_bid"`
}
