package data

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Init should run when server starts
func Init() *Data {
	data := Data{
		Users:          make(map[uint64]User),
		Auctions:       make(map[uint64]Auction),
		ClosedAuctions: make(map[uint64]Auction),
		Bids:           make(map[uint64]Bid),
		LastUserID:     0,
		LastAuctionID:  0,
		LastBidID:      0,
		Mtx:            &sync.RWMutex{},
	}

	// Create initial data
	userID1 := data.AddUser()
	userID2 := data.AddUser()
	userID3 := data.AddUser()
	userID4 := data.AddUser()

	itemID1, err := data.AddAuction(userID1, "Car for sale", 500)

	if err != nil {
		fmt.Println("Error while creating initial data")
	} else {
		data.AddBid(userID2, itemID1, 500.5)
		data.AddBid(userID3, itemID1, 510)
		data.AddBid(userID4, itemID1, 650)
	}

	itemID2, err := data.AddAuction(userID2, "House for sale", 10000)

	if err != nil {
		fmt.Println("Error while creating initial data")
	} else {
		data.AddBid(userID1, itemID2, 12500.21)
		data.AddBid(userID3, itemID2, 13000)
		data.AddBid(userID4, itemID2, 13500)
	}

	return &data
}

// AddUser creates new user and returns id
func (d *Data) AddUser() uint64 {
	d.Mtx.Lock()
	defer d.Mtx.Unlock()

	userID := atomic.AddUint64(&d.LastUserID, 1)
	d.Users[userID] = User{
		CreatedAt: time.Now().Unix(),
	}

	return userID
}

// AddAuction creates new auction and returns id
func (d *Data) AddAuction(userID uint64, title string, startingPrice float64) (uint64, error) {
	if _, userExists := d.Users[userID]; userExists {
		d.Mtx.Lock()
		defer d.Mtx.Unlock()

		auctionID := atomic.AddUint64(&d.LastAuctionID, 1)
		d.Auctions[auctionID] = Auction{
			UserID:        userID,
			Title:         title,
			StartingPrice: startingPrice,
			Closed:        false,
			Sold:          false,
			LastBidID:     0,
			WinnerID:      0,
			CreatedAt:     time.Now().Unix(),
			ExpireAt:      time.Now().Add(1 * time.Hour).Unix(), // 1 hour expiration
			Mtx:           &sync.RWMutex{},
		}

		return auctionID, nil
	}

	return 0, errors.New("User does not exist")
}

// AddBid creates new bid and returns id
func (d *Data) AddBid(userID, auctionID uint64, amount float64) (uint64, error) {
	if auction, auctionExists := d.Auctions[auctionID]; auctionExists {
		if _, userExists := d.Users[userID]; userExists {
			d.Mtx.Lock()
			auction.Mtx.Lock()
			defer d.Mtx.Unlock()
			defer auction.Mtx.Unlock()

			if auction.UserID == userID {
				return 0, errors.New("You cannot bid on your own auction")
			}

			if auction.LastBidID > 0 {
				if d.Bids[auction.LastBidID].Amount >= amount {
					return 0, errors.New("Bid cannot be less then last bid")
				}
			}

			bidID := atomic.AddUint64(&d.LastBidID, 1)
			d.Bids[bidID] = Bid{
				UserID:    userID,
				AuctionID: auctionID,
				Amount:    amount,
				CreatedAt: time.Now().Unix(),
			}
			auction.LastBidID = bidID
			d.Auctions[auctionID] = auction

			return bidID, nil
		}

		return 0, errors.New("User does not exist")
	}

	return 0, errors.New("Item does not exist")
}

// CloseAuction closes auction, this can be called by goroutine which runs every minute
func (d *Data) CloseAuction(auctionID uint64) bool {
	if auction, auctionExists := d.Auctions[auctionID]; auctionExists {
		auction.Mtx.Lock()
		d.Mtx.Lock()
		defer auction.Mtx.Unlock()
		defer d.Mtx.Unlock()

		auction.Closed = true
		if auction.LastBidID > 0 {
			auction.WinnerID = auction.LastBidID
			auction.Sold = true
		}

		d.ClosedAuctions[auctionID] = auction
		delete(d.Auctions, auctionID)

		return true
	}

	return false
}

// GetBidsOfAuction lists bids of auction
func (d *Data) GetBidsOfAuction(auctionID uint64) ([]Bid, error) {
	var bids []Bid

	if auction, auctionExists := d.Auctions[auctionID]; auctionExists {
		auction.Mtx.Lock()
		defer auction.Mtx.Unlock()

		for _, value := range d.Bids {
			if value.AuctionID == auctionID {
				bids = append(bids, Bid{
					Amount:    value.Amount,
					CreatedAt: value.CreatedAt,
				})
			}
		}

		return bids, nil
	}

	return bids, errors.New("Auction does not exist")
}

// GetAuctionOverview returns overview of single auction
func (d *Data) GetAuctionOverview(auctionID uint64) (*Auction, *Bid, error) {
	if auction, auctionExists := d.Auctions[auctionID]; auctionExists {
		auction.Mtx.Lock()
		defer auction.Mtx.Unlock()

		var bid Bid

		if auction.LastBidID > 0 {
			bid = d.Bids[auction.LastBidID]
		}

		return &auction, &bid, nil
	}

	return nil, nil, errors.New("Auction does not exist")
}

//GetBiddedItemsOfUser returns bidded items of user
func (d *Data) GetBiddedItemsOfUser(userID uint64) (map[uint64]Auction, error) {
	var list = make(map[uint64]Auction)

	if _, userExists := d.Users[userID]; userExists {
		for _, value := range d.Bids {
			if value.UserID == userID {
				list[value.AuctionID] = d.Auctions[value.AuctionID]
			}
		}

		return list, nil
	}

	return list, errors.New("User does not exist")
}
