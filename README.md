### Start Anvil node in terminal 1
```bash
anvil
```
### Create contracts
```bash
# Install OpenZeppelin for ERC721 token
forge install OpenZeppelin/openzeppelin-contracts

# Design Overview:
# User
#  ↓
# Marketplace (entry point / orchestrator)
#  └─ Delegates auctions → AuctionHouse
#                           ├─ createAuction
#                           ├─ placeBid
#                           └─ endAuction
```
```
src
├── NFT.sol
├── NFTAuctionHouse.sol
├── NFTMarketplace.sol
```

### Deploy contracts to the running node in terminal 2
```bash
export PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 # Anvil’s default first account

forge script script/Deploy.s.sol:Deploy --rpc-url http://127.0.0.1:8545 --broadcast
```
---
### Set up Go backend
- Initialize a Go module
```bash
go mod init nft-auction-back
```

- Create main (entry) file
```bash
touch main.go
```
```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, NFT Auction Go Backend!")
}
```
- Run the go backend
```bash
go run main.go
```

### Install Gin
```bash
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/cors
go mod tidy # cleanup dependencies if needed
```

### API routing
```
backend
├──	api/
	├── server.go        # setup + router init
	├── routes.go        # route definitions
	├── handlers/
	│   ├── listing.go   # handle listing requests
	│   ├── auction.go   # handle auction requests
	│   └── bid.go       # handle bid requests
```
### Install gorm
```bash
cd backend
go get gorm.io/gorm
go get gorm.io/driver/mysql #mysql driver
go get gorm.io/driver/postgres # postgres driver
```

### Define models and repositories
```
backend
|__api
|__models
|	├── auction.go
|	├── bid.go 
|	├── listing.go 
|__repo
	├── auction_repo.go
	├── bid_repo.go 
	├── listing_repo.go 
```

### Set up MySQL
```bash
sudo apt install mysql-server
sudo systemctl start mysql
# open MySQL CLI
sudo mysql
# delete database if wanna start from scratch
mysql > DROP DATABASE IF EXISTS nft_marketplace;
# create database, user/password, and privileges
mysql > CREATE DATABASE nft_marketplace;
mysql > CREATE USER 'user'@'%' IDENTIFIED BY 'password';
mysql > GRANT ALL PRIVILEGES ON nft_marketplace.* TO 'user'@'%';
mysql > FLUSH PRIVILEGES;
```

### Connect to DB
```go
dsn := "user:password@tcp(127.0.0.1:3306)/nft_marketplace?charset=utf8mb4&parseTime=True&loc=Local"

db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
if err != nil {
   log.Fatal(err)
}

err = db.AutoMigrate(&models.Listing{}, &models.Auction{}, &models.Bid{})
if err != nil {
   log.Fatal(err)
}

return db
```

### Ethereum Events and Logs
An Ethereum block is basically a container with:

- Block metadata (header)
- A list of transactions
- Execution results (receipts, logs, gas usage, state changes)
```
Block
├── Header
│   ├── blockNumber
│   ├── parentHash
│   ├── timestamp
│   ├── miner / validator
│   ├── gasLimit
│   ├── gasUsed
│   ├── stateRoot
│   ├── transactionsRoot
│   └── receiptsRoot
│
├── Transactions[]
│   ├── tx1
│   ├── tx2
│   └── tx3
│
└── (Derived during execution)
    └── Receipts[]
        ├── receipt for tx1
        │      ├── status
        │      ├── gasUsed
        │      └── logs[]
        │          ├── Log #0
        │          ├── address        // contract address
        │          ├── topics[]       // indexed event fields + event signature
        │          ├── data           // non-indexed fields
        │          ├── blockNumber
        │          ├── txHash
        │          ├── logIndex
        │          └── Log #1
        ├── receipt for tx2
        └── receipt for tx3
```
```solidity
event Transfer(address indexed from, address indexed to, uint256 value);
```
corresponding `vLog`
```go
types.Log{
    Address: common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
    // Topics is an array of 32-byte hashes([]common.hash). Each topic is a `bytes32` value
	Topics: []common.Hash{
        // Topic[0]: keccak256("Transfer(address,address,uint256)")
        common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55aebc4a0b6c"),
        
        // Topic[1]: indexed from (address)
        common.HexToHash("0x000000000000000000000000aabbccddeeff0011223344556677889900aabbcc"),
        
        // Topic[2]: indexed to (address)
        common.HexToHash("0x000000000000000000000000ffeeddccbbaa99887766554433221100ffeeddcc"),
    },
    Data: common.FromHex(
        // non-indexed params (e.g., uint256 value)
        "0x00000000000000000000000000000000000000000000000000000000000003e8",
    ),
    BlockNumber: 12345678,
    TxHash:      common.HexToHash("0xabc123..."),
    TxIndex:     0,
    BlockHash:   common.HexToHash("0xdef456..."),
    Index:       0,
    Removed:     false,
}
```
---
### Full Stack Workflow
```bash
# Start local ethereum dev node
anvil

# Deploy contracts to the local network
export PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 # Anvil’s default first account
forge script script/Deploy.s.sol:Deploy --rpc-url http://127.0.0.1:8545 --broadcast

# Update contract addresses in backend and frontend's configuration
frontend
|__config.js

backend
|__config
     |__env_vars.go

# Run Go backend
cd backend
go run main.go

# Launch the frontend (serves static files on a local server)
cd frontend
npx serve .

# Open MetaMask and import accounts that created by anvil via Add wallet -> Import an account and enter private key

# Access http://localhost:3000 and connect to an added account

# Mint a NFT by enter ipfs://fake-uri and click Mint button then approve in MetaMask

# Approve NFT market and List the minted NFT, tokenId starts from 1.

```
---
### Reliable Ethereum Event Synchronization (Using Polling and Safe Block Confirmations)
Real-time subscriptions alone (SubscribeFilterLogs) are unreliable because:
- the backend may go down and miss events
- websocket connections can disconnect silently
- Ethereum blocks near the chain tip may be reorganized (reorgs)

Instead of trusting live events immediately:
- Periodically poll logs using FilterLogs
- Only process blocks older than a confirmation threshold
- Persist the last processed block in the database
- Replay from that checkpoint after restart
```
          ┌──────────────────────────┐
          │        Start Poller      │
          └────────────┬─────────────┘
                       │
                       ▼
          ┌──────────────────────────┐
          │ Load lastProcessedBlock  │
          │      (from DB)           │
          └────────────┬─────────────┘
                       │
                       ▼
        ┌───────────────────────────────┐
        │        Poll Loop (every N s)  │
        └────────────┬──────────────────┘
                     │
                     ▼
        ┌───────────────────────────────┐
        │   Get latest block from RPC   │
        └────────────┬──────────────────┘
                     │
                     ▼
        ┌───────────────────────────────┐
        │ safeBlock = latest - N_conf   │
        └────────────┬──────────────────┘
                     │
                     ▼
        ┌───────────────────────────────┐
        │ safeBlock <= lastProcessed?   │
        └───────┬───────────────┬───────┘
                │ YES           │ NO
                ▼               ▼
        ┌──────────────┐  ┌──────────────────────────┐
        │   Sleep      │  │  Process [last+1 → safe] │
        └──────────────┘  └────────────┬─────────────┘
                                       │
                                       ▼
                        ┌────────────────────────────┐
                        │  Split into batches (step) │
                        └────────────┬───────────────┘
                                     │
                                     ▼
                    ┌────────────────────────────────┐
                    │ FilterLogs(from, to) per batch │
                    └────────────┬───────────────────┘
                                 │
                                 ▼
                    ┌────────────────────────────────┐
                    │      For each log:             │
                    │  - dedup (txHash+logIndex)     │
                    │  - decode via ABI              │
                    │  - apply business logic        │
                    └────────────┬───────────────────┘
                                 │
                                 ▼
                    ┌────────────────────────────────┐
                    │ Update lastProcessedBlock = to │
                    └────────────┬───────────────────┘
                                 │
                                 ▼
                           (loop again)
```
---
### Chainlink Price Feed From Go (ETH to USD for display)
Read directly from the Chainlink ETH/USD feed contract
On Ethereum mainnet, the ETH/USD feed is:
```
0x5f4ec3df9cbd43714fe2740f5e3616155c5b8419
```
This contract exposes:
```solidity
latestRoundData()
```
- Generate Go Binding
  - Install abigen 
  ```bash
  cd backend 
  go install github.com/ethereum/go-ethereum/cmd/abigen@latest
  ```
  - Create `backend\abi\AggregatorV3Interface.json`
  ```json
  [
    {
        "inputs": [],
        "name": "latestRoundData",
        "outputs": [
        { "internalType": "uint80", "name": "roundId", "type": "uint80" },
        { "internalType": "int256", "name": "answer", "type": "int256" },
        { "internalType": "uint256", "name": "startedAt", "type": "uint256" },
        { "internalType": "uint256", "name": "updatedAt", "type": "uint256" },
        { "internalType": "uint80", "name": "answeredInRound", "type": "uint80" }
        ],
        "stateMutability": "view",
        "type": "function"
    }
  ]
  ```
  - Generate binding
  ```bash
  cd backend
  
  abigen \
  --abi ./abi/AggregatorV3Interface.json \
  --pkg chainlink \
  --out ./service/chainlink.go

  go mod tidy # clean up and install any missing dependencies for the generated chainlink.go
  ```
- Create `price_service.go` and use `chainlink.go`
  ```go
  // wei -> usd
  feed, err := NewChainlink(config.ChainlinkFeedAddress, client)
  data, err := feed.LatestRoundData(nil)

  eth2usdExchangeRate := new(big.Float).Quo(
      new(big.Float).SetInt(data.Answer),
      big.NewFloat(1e8), // Chainlink ETH/USD decimals
  )

  eth2usdExchangeRatef64, _ := eth2usdExchangeRate.Float64()

  ethVal := new(big.Float).Quo(
	  new(big.Float).SetInt(wei),
	  big.NewFloat(1e18), // ETH decimals
  )

  ethFloat, _ := ethVal.Float64()

  usd := ethFloat * eth2usdExchangeRatef64
  ```
---
### Foundry Unit Tests
- Must start with `test`
```solidity
function testMint() public {}
function testWithdraw() public {}
```
- Fuzz tests, Foundry generates random inputs and run many times
```solidity
function testFuzz_Mint(uint256 amount) public {}
```

### Key Foundry Test concepts
- vm.prank(addr) → simulate caller, applies to *ONE* single call
- vm.startPrank(addr)/vm.stopPrank(), applies to *ALL* calls until stopped
- vm.deal(addr, amount) → give ETH
- vm.expectRevert() → expect failure
- assertEq() → assertions

### Foundry built-in hook
- setUp(), it runs before every single test
```solidity
function setUp() public {}
```


