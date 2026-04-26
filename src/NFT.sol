// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;

import "openzeppelin-contracts/contracts/token/ERC721/ERC721.sol"; // forge install OpenZeppelin/openzeppelin-contracts
import "openzeppelin-contracts/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "openzeppelin-contracts/contracts/access/Ownable.sol";

contract NFT is ERC721, ERC721URIStorage, Ownable {
    uint256 private _tokenCount;

    uint256 public constant MAX_SUPPLY = 10000;
    uint256 public mintPrice = 0.01 ether;

    event Minted(address indexed minter, uint256 indexed tokenId, string uri);

    constructor() ERC721("GopherNFT", "GONFT") Ownable(msg.sender) {}

    function setMintPrice(uint256 price) public onlyOwner {
        require(price > 0, "Price must be greater than zero");
        mintPrice = price;
    }

    function mint(string memory uri) public payable returns (uint256) {
        require(_tokenCount < MAX_SUPPLY, "Max supply reached");
        require(msg.value >= mintPrice, "Insufficient payment");

        uint256 tokenId = _tokenCount++;
        _safeMint(msg.sender, tokenId);
        _setTokenURI(tokenId, uri);

        emit Minted(msg.sender, tokenId, uri);

        return tokenId;
    }

    function withdraw() public onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");

        (bool success, ) = owner().call{value: balance}("");
        require(success, "Withdrawal failed");
    }

    /* The following are required overridden functions */
    
    // Both ERC721 and ERC721URIStorage implement tokenURI. We need to specify which one to use. 
    function tokenURI(uint256 tokenId) 
        public view override(ERC721, ERC721URIStorage) returns (string memory)
    {
        return super.tokenURI(tokenId);
    }

    // Both ERC721 and ERC721URIStorage implement supportsInterface. We need to specify which one to use.
    function supportsInterface(bytes4 interfaceId)
        public view override(ERC721, ERC721URIStorage) returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}