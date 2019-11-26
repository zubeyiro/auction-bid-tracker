# auction-bid-tracker
Concurrent bidding system

## Endpoints

### Create Bid
`/auctions/:AuctionID/bids (POST)`

Request body
```
{
	"user_id": 2,
	"amount": 5000
}
```

Response format
```
{
    "data": [BID_ID],
    "status": true // if status is false, that means bid could not be saved and data contains error message
}
```

### List bids for auction
`/auctions/:AuctionID/bids (GET)`

Response format
```
{
    "data": [
        {
            "amount": [BID_AMOUNT],
            "created_at": [UNIX_TIMESTAMP]
        },
        ...
    ],
    "status": true
}
```

### Get auction with last (highest) bid
`/auctions/1 (GET)`

Response format
```
{
    "data": {
        "title": [BID_TITLE],
        "starting_price": [BID_STARTING_PRICE],
        "created_at": [UNIX_TIMESTAMP],
        "expires_at": [UNIX_TIMESTAMP],
        "last_bid": {
            "amount": [BID_AMOUNT],
            "created_at": [UNIX_TIMESTAMP]
        }
    },
    "status": true
}
```

### Get users bidded auctions
`/auctions/users/bids/:UserID`

Response format
```
{
    "data": [
        {
            "title": [BID_TITLE],
            "starting_price": [BID_STARTING_PRICE],
            "created_at": [UNIX_TIMESTAMP],
            "expires_at": [UNIX_TIMESTAMP]
        },
        ...
    ],
    "status": true
}
```

## Approach

I have used Gin for RESTful API and created simple shared data object which contains Users, Auctions and Bids under Data struct `/src/server/data/models.go`. Data object and Auction object has to be prevented from data races so I have used mutexes for those structs. I have also used atomic operations for increasing ID's. All data objects are being processed under data package `/src/server/data/data.go`.

Data object is being initialized before REST, then its being passed as pointer to REST listener. It can also be shared with i.e. socket server for real-time purposes as well.

REST service has also validators and serializers, and it also has basic response structure for standardization `/src/utils/response.go`.

## Test & Run

For testing, run;

`go test -v test/bidding_test.go`

For running RESTful API, run;

`go run src/main.go`