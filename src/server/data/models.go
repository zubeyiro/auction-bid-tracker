package data

import (
	"sync"
)

// User struct
type User struct {
	CreatedAt int64
}

// Auction struct
type Auction struct {
	UserID        uint64
	Title         string
	StartingPrice float64
	Closed        bool
	Sold          bool // auction can be closed without being sold
	CreatedAt     int64
	ExpireAt      int64 // every auction has an expiration time, countdown
	LastBidID     uint64
	WinnerID      uint64
	Mtx           *sync.RWMutex
}

// Bid struct
type Bid struct {
	UserID    uint64
	AuctionID uint64
	Amount    float64
	CreatedAt int64
}

// Data is overall data struct
type Data struct {
	Users          map[uint64]User
	Auctions       map[uint64]Auction
	ClosedAuctions map[uint64]Auction
	Bids           map[uint64]Bid
	LastUserID     uint64
	LastAuctionID  uint64
	LastBidID      uint64
	Mtx            *sync.RWMutex
}
