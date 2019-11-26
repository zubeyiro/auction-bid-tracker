package auctions

// CreateBidValidator is request body of /items (POST)
type CreateBidValidator struct {
	UserID uint64  `json:"user_id" binding:"required,numeric,min=0"`
	Amount float64 `json:"amount" binding:"required,numeric,min=0"`
}

// AuctionIDValidator is only for AuctionID validation
type AuctionIDValidator struct {
	AuctionID uint64 `json:"AuctionID" binding:"required,min=1"`
}

// UserIDValidator is for UserID validation
type UserIDValidator struct {
	UserID uint64 `json:"UserID" binding:"required,min=1"`
}
