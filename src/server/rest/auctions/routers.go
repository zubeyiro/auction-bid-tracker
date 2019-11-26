package auctions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zubeyiro/auction-bid-tracker/src/server/data"
	"github.com/zubeyiro/auction-bid-tracker/src/utils"
)

var database *data.Data

// RegisterRouters creates all routers for /auctions path
func RegisterRouters(router *gin.RouterGroup, db *data.Data) {
	database = db
	router.GET("/:AuctionID/bids", getBidsOfAuction)       // get bids for auction
	router.POST("/:AuctionID/bids", createBid)             // new bid
	router.GET("/:AuctionID", getAuction)                  // get auction overview
	router.GET("/:AuctionID/bids/:UserID", getBiddedItems) // get bidded items of user
}

func createBid(c *gin.Context) {
	var param AuctionIDValidator
	var body CreateBidValidator

	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewHTTPResponse(false, err.Error()))

		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewHTTPResponse(false, err.Error()))

		return
	}

	if bidID, err := database.AddBid(body.UserID, param.AuctionID, body.Amount); err != nil {
		c.JSON(http.StatusOK, utils.NewHTTPResponse(false, err.Error()))
	} else {
		c.JSON(http.StatusOK, utils.NewHTTPResponse(true, bidID))
	}

	return
}

func getBidsOfAuction(c *gin.Context) {
	var param AuctionIDValidator

	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewHTTPResponse(false, err.Error()))

		return
	}

	if list, err := database.GetBidsOfAuction(param.AuctionID); err != nil {
		c.JSON(http.StatusOK, utils.NewHTTPResponse(false, err.Error()))
	} else {
		var bids []BidSerializer

		for _, value := range list {
			bids = append(bids, BidSerializer{
				Amount:    value.Amount,
				CreatedAt: value.CreatedAt,
			})
		}

		c.JSON(http.StatusOK, utils.NewHTTPResponse(true, bids))
	}

	return
}

func getAuction(c *gin.Context) {
	var param AuctionIDValidator

	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewHTTPResponse(false, err.Error()))

		return
	}

	if auction, bid, err := database.GetAuctionOverview(param.AuctionID); err != nil {
		c.JSON(http.StatusOK, utils.NewHTTPResponse(false, err.Error()))
	} else {
		ret := AuctionSerializerWithBid{
			Title:         auction.Title,
			StartingPrice: auction.StartingPrice,
			CreatedAt:     auction.CreatedAt,
			ExpireAt:      auction.ExpireAt,
		}

		if bid != nil {
			ret.LastBid = BidSerializer{
				Amount:    bid.Amount,
				CreatedAt: bid.CreatedAt,
			}
		}

		c.JSON(http.StatusOK, utils.NewHTTPResponse(true, ret))
	}

	return
}

func getBiddedItems(c *gin.Context) {
	var param UserIDValidator

	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewHTTPResponse(false, err.Error()))

		return
	}

	if list, err := database.GetBiddedItemsOfUser(param.UserID); err != nil {
		c.JSON(http.StatusOK, utils.NewHTTPResponse(false, err.Error()))
	} else {
		var ret []AuctionSerializer

		for _, val := range list {
			ret = append(ret, AuctionSerializer{
				Title:         val.Title,
				StartingPrice: val.StartingPrice,
				CreatedAt:     val.CreatedAt,
				ExpireAt:      val.ExpireAt,
			})
		}

		c.JSON(http.StatusOK, utils.NewHTTPResponse(true, ret))
	}

	return
}
