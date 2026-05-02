### Start Anvil node in terminal 1
```bash
anvil
```
### Deploy contracts to the running node in terminal 2
```bash
export PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 # Anvil’s default first account

forge script script/Deploy.s.sol:Deploy --rpc-url http://127.0.0.1:8545 --broadcast
```

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
# create table, user and password
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

### Foundry Test functions
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


