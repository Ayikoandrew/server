# Liora Server - Backend API for Payment Application

The robust Go backend powering **Liora**, a modern financial application. Provides secure authentication, user management, and payment infrastructure with planned smart contract integration using **abigen** for Go bindings and zkSync Layer 2 scaling.

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
![Ethereum](https://img.shields.io/badge/ethereum-%23627EEA.svg?style=for-the-badge&logo=ethereum&logoColor=white)
![Solidity](https://img.shields.io/badge/solidity-%23363636.svg?style=for-the-badge&logo=solidity&logoColor=white)

## 🏦 About Liora

Liora is a comprehensive financial ecosystem consisting of:
- **💱 Flutter Mobile App** ([Ayikoandrew/liora](https://github.com/Ayikoandrew/liora)) - Cross-platform payment interface
- **🚀 Go Backend Server** (this repository) - Secure API with smart contract integration
- **🌐 Smart Contracts** (planned) - Blockchain-based payment infrastructure
- **⚡ zkSync Integration** (planned) - Scalable L2 payment processing
- **🔗 Abigen Bindings** (planned) - Type-safe smart contract interactions

## 🔮 Planned Features - Blockchain Integration

### 🔗 Smart Contract Integration with Abigen
- **📜 Contract Bindings** - Auto-generated Go types from Solidity contracts
- **🔒 Type Safety** - Compile-time contract interaction validation
- **⚡ Efficient Calls** - Direct contract method invocation from Go
- **📡 Event Listening** - Real-time blockchain event monitoring
- **🔄 Transaction Management** - Automated gas estimation and transaction handling

### 💰 Blockchain Payment Features
- **💳 On-Chain Wallets** - Smart contract-managed user wallets
- **💸 Instant Transfers** - Direct blockchain-based payments
- **🎯 Multi-Token Support** - ERC-20 token handling via abigen
- **📊 Transaction Tracking** - Real-time blockchain transaction monitoring
- **🔄 Automated Settlements** - Smart contract-based payment processing

### ⚡ zkSync Layer 2 Features
- **🚄 Fast Transactions** - Near-instant payment processing via L2 contracts
- **💸 Low-Cost Operations** - Minimal fees through zkSync Era integration
- **🔄 Batch Processing** - Efficient bulk payment handling with abigen
- **🌉 L1/L2 Bridge** - Seamless asset transfers between layers

## 🏗️ Enhanced Architecture with Blockchain

```
Liora Ecosystem with Blockchain Integration
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (Flutter)                       │
│  🏠 Dashboard  💰 Payments  🔄 Transfers  📈 DeFi          │
└─────────────────────┬───────────────────────────────────────┘
                      │ REST API + WebSocket
┌─────────────────────▼───────────────────────────────────────┐
│                 Backend Server (Go)                        │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐   │
│  │    Auth     │ │  Payments   │ │    Blockchain       │   │
│  │   Service   │ │   Service   │ │     Service         │   │
│  └─────────────┘ └─────────────┘ └─────────┬───────────┘   │
│                                           │               │
│  ┌─────────────────────────────────────────▼─────────────┐ │
│  │              Abigen Bindings                         │ │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────────┐ │ │
│  │  │   Payments  │ │    Users    │ │     zkSync      │ │ │
│  │  │  Contract   │ │  Registry   │ │   Contracts     │ │ │
│  │  └─────────────┘ └─────────────┘ └─────────────────┘ │ │
│  └─────────────────────────────────────────┬─────────────┘ │
└──────────────────────────────────────────────┼───────────────┘
                                             │
┌──────────────────────────────────────────────▼───────────────┐
│                    Blockchain Layer                          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐   │
│  │  Ethereum   │ │   zkSync    │ │      IPFS           │   │
│  │  Mainnet    │ │    Era      │ │   (Metadata)        │   │
│  └─────────────┘ └─────────────┘ └─────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### Enhanced Server Directory Structure
```
├── main.go                 # Application entry point
├── api/                    # HTTP handlers and routing
│   └── server.go          # Main server setup and endpoints
├── blockchain/             # Blockchain integration layer
│   ├── client.go          # Ethereum/zkSync client setup
│   ├── contracts/         # Generated abigen bindings
│   │   ├── payments.go    # Payment contract bindings
│   │   ├── registry.go    # User registry bindings
│   │   └── zksync.go      # zkSync contract bindings
│   ├── events.go          # Blockchain event listeners
│   ├── transactions.go    # Transaction management
│   └── wallet.go          # Wallet operations
├── contracts/             # Solidity smart contracts
│   ├── Payments.sol       # Main payment contract
│   ├── UserRegistry.sol   # User verification contract
│   └── zkSync/            # zkSync-specific contracts
│       ├── L2Payments.sol # Layer 2 payment processing
│       └── Bridge.sol     # L1/L2 bridge contract
├── scripts/               # Development and deployment scripts
│   ├── generate-bindings.sh # Abigen binding generation
│   ├── deploy-contracts.sh  # Contract deployment
│   └── migrate.sh          # Database migrations
├── database/              # Data persistence layer
├── middleware/            # HTTP middleware stack
├── security/              # Security utilities
├── functions/             # Business logic layer
├── types/                 # Data structures and models
└── config/               # Configuration management
```

## 🛠️ Enhanced Tech Stack

### Blockchain Integration Stack
- **Smart Contract Language**: Solidity with OpenZeppelin
- **Go Bindings**: Abigen for type-safe contract interactions
- **Ethereum Client**: go-ethereum (geth) for blockchain connectivity
- **Layer 2**: zkSync Era for scalable payments
- **Development Tools**: Foundry for contract development
- **Token Standards**: ERC-20 (payments), ERC-721 (achievements)

### Abigen Workflow
```bash
# Generate Go bindings from compiled contracts
abigen --abi contracts/build/Payments.abi \
       --pkg contracts \
       --type Payments \
       --out blockchain/contracts/payments.go

abigen --abi contracts/build/UserRegistry.abi \
       --pkg contracts \
       --type UserRegistry \
       --out blockchain/contracts/registry.go
```

## 🚀 Enhanced Setup with Blockchain

### Additional Prerequisites
- **Ethereum Node** or **Infura/Alchemy** API key
- **Abigen** (part of go-ethereum)

### Blockchain Development Setup

1. **Install blockchain development tools**
   ```bash
   # Install go-ethereum (includes abigen)
   go install github.com/ethereum/go-ethereum/cmd/abigen@latest
   
   # Install Node.js dependencies for smart contracts
   npm install --save-dev hardhat @openzeppelin/contracts
   ```

2. **Environment configuration with blockchain**
   ```bash
   # Add to your .env file
   
   # Ethereum Configuration
   ETH_RPC_URL=https://mainnet.infura.io/v3/your-project-id
   ETH_PRIVATE_KEY=your-wallet-private-key-for-deployment
   ETH_CHAIN_ID=1
   
   # zkSync Configuration
   ZKSYNC_RPC_URL=https://mainnet.era.zksync.io
   ZKSYNC_CHAIN_ID=324
   
   # Contract Addresses (after deployment)
   PAYMENTS_CONTRACT=0x742d35Cc6634C0532925a3b8D1C11C0532925a3b
   REGISTRY_CONTRACT=0x742d35Cc6634C0532925a3b8D1C11C0532925a3b
   ZKSYNC_PAYMENTS_CONTRACT=0x742d35Cc6634C0532925a3b8D1C11C0532925a3b
   
   # API Keys
   INFURA_PROJECT_ID=your-infura-project-id
   ALCHEMY_API_KEY=your-alchemy-api-key
   ```

3. **Generate contract bindings**
   ```bash
   # Make the script executable
   chmod +x scripts/generate-bindings.sh
   
   # Generate Go bindings from smart contracts
   ./scripts/generate-bindings.sh
   ```

4. **Run with blockchain integration**
   ```bash
   # Development mode with blockchain
   go run main.go --enable-blockchain
   
   # Production build
   go build -tags blockchain -o liora-server
   ./liora-server
   ```

## 📋 Smart Contract Integration Examples

### Using Abigen Bindings in Go

```go
package blockchain

import (
    "context"
    "crypto/ecdsa"
    "math/big"
    
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
    
    "github.com/Ayikoandrew/server/blockchain/contracts"
)

type BlockchainService struct {
    client          *ethclient.Client
    paymentsContract *contracts.Payments
    registryContract *contracts.UserRegistry
    privateKey      *ecdsa.PrivateKey
}

func NewBlockchainService(rpcURL, privateKeyHex, paymentsAddr, registryAddr string) (*BlockchainService, error) {
    // Connect to Ethereum client
    client, err := ethclient.Dial(rpcURL)
    if err != nil {
        return nil, err
    }
    
    // Load private key for transactions
    privateKey, err := crypto.HexToECDSA(privateKeyHex)
    if err != nil {
        return nil, err
    }
    
    // Initialize contract instances using abigen bindings
    paymentsContract, err := contracts.NewPayments(
        common.HexToAddress(paymentsAddr), 
        client,
    )
    if err != nil {
        return nil, err
    }
    
    registryContract, err := contracts.NewUserRegistry(
        common.HexToAddress(registryAddr), 
        client,
    )
    if err != nil {
        return nil, err
    }
    
    return &BlockchainService{
        client:          client,
        paymentsContract: paymentsContract,
        registryContract: registryContract,
        privateKey:      privateKey,
    }, nil
}

// Process payment using smart contract
func (bs *BlockchainService) ProcessPayment(fromAddr, toAddr string, amount *big.Int) error {
    auth, err := bind.NewKeyedTransactorWithChainID(bs.privateKey, big.NewInt(1))
    if err != nil {
        return err
    }
    
    // Call smart contract method through abigen binding
    tx, err := bs.paymentsContract.Transfer(
        auth,
        common.HexToAddress(fromAddr),
        common.HexToAddress(toAddr),
        amount,
    )
    if err != nil {
        return err
    }
    
    log.Printf("Payment transaction sent: %s", tx.Hash().Hex())
    return nil
}

// Register user on blockchain
func (bs *BlockchainService) RegisterUser(userID, walletAddr string) error {
    auth, err := bind.NewKeyedTransactorWithChainID(bs.privateKey, big.NewInt(1))
    if err != nil {
        return err
    }
    
    // Convert userID to bytes32 for contract
    userIDBytes := [32]byte{}
    copy(userIDBytes[:], userID)
    
    tx, err := bs.registryContract.RegisterUser(
        auth,
        userIDBytes,
        common.HexToAddress(walletAddr),
    )
    if err != nil {
        return err
    }
    
    log.Printf("User registration transaction: %s", tx.Hash().Hex())
    return nil
}

// Listen for payment events
func (bs *BlockchainService) ListenForPayments(ctx context.Context) error {
    // Create event filter using abigen
    iterator, err := bs.paymentsContract.FilterTransfer(
        &bind.FilterOpts{Context: ctx},
        []common.Address{}, // from addresses (empty = all)
        []common.Address{}, // to addresses (empty = all)
    )
    if err != nil {
        return err
    }
    defer iterator.Close()
    
    for iterator.Next() {
        event := iterator.Event
        log.Printf("Payment detected: %s -> %s, Amount: %s", 
            event.From.Hex(), 
            event.To.Hex(), 
            event.Amount.String(),
        )
        
        // Process the payment event in your business logic
        // e.g., update database, send notifications, etc.
    }
    
    return iterator.Error()
}
```

### Smart Contract Example (Solidity)

```solidity
// contracts/Payments.sol
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract Payments is ERC20, Ownable {
    mapping(address => bool) public verifiedUsers;
    mapping(bytes32 => address) public userRegistry;
    
    event Transfer(address indexed from, address indexed to, uint256 amount);
    event UserRegistered(bytes32 indexed userId, address indexed wallet);
    
    constructor() ERC20("Liora Token", "LIORA") {
        _mint(msg.sender, 1000000 * 10**decimals());
    }
    
    function registerUser(bytes32 userId, address wallet) external onlyOwner {
        userRegistry[userId] = wallet;
        verifiedUsers[wallet] = true;
        emit UserRegistered(userId, wallet);
    }
    
    function transfer(address to, uint256 amount) public override returns (bool) {
        require(verifiedUsers[msg.sender], "User not verified");
        require(verifiedUsers[to], "Recipient not verified");
        
        return super.transfer(to, amount);
    }
}
```

## 📜 Abigen Integration Scripts

### Generate Bindings Script
```bash
#!/bin/bash
# scripts/generate-bindings.sh

echo "Compiling smart contracts..."
npx hardhat compile

echo "Generating Go bindings with abigen..."

# Generate Payments contract bindings
abigen --abi artifacts/contracts/Payments.sol/Payments.abi \
       --bin artifacts/contracts/Payments.sol/Payments.bin \
       --pkg contracts \
       --type Payments \
       --out blockchain/contracts/payments.go

# Generate UserRegistry contract bindings  
abigen --abi artifacts/contracts/UserRegistry.sol/UserRegistry.abi \
       --bin artifacts/contracts/UserRegistry.sol/UserRegistry.bin \
       --pkg contracts \
       --type UserRegistry \
       --out blockchain/contracts/registry.go

# Generate zkSync contract bindings
abigen --abi artifacts/contracts/zkSync/L2Payments.sol/L2Payments.abi \
       --bin artifacts/contracts/zkSync/L2Payments.sol/L2Payments.bin \
       --pkg contracts \
       --type L2Payments \
       --out blockchain/contracts/zksync.go

echo "✅ Go bindings generated successfully!"
echo "📁 Files created in blockchain/contracts/"
```

### Development Workflow
```bash
# 1. Develop smart contracts
vim contracts/Payments.sol

# 2. Compile and generate bindings
./scripts/generate-bindings.sh

# 3. Test contract integration
go test ./blockchain/...

# 4. Deploy contracts (testnet first)
./scripts/deploy-contracts.sh --network sepolia

# 5. Update contract addresses in .env
vim .env

# 6. Test full integration
go run main.go --enable-blockchain
```

## 🎯 Blockchain Development Roadmap

### Phase 1: Foundation (Q3 2025)
- ✅ Set up abigen workflow and Go bindings
- ✅ Implement basic payment contract
- ✅ User registry with verification
- ✅ Event listening and processing

### Phase 2: zkSync Integration (Q4 2025)
- 🔄 Deploy contracts to zkSync Era
- 🔄 Implement L1/L2 bridge functionality
- 🔄 Batch transaction processing
- 🔄 Gas optimization strategies

### Phase 3: Advanced Features (Q1 2026)
- 🔮 Multi-signature wallet support
- 🔮 DeFi integration (yield farming, staking)
- 🔮 Cross-chain payment routing
- 🔮 Advanced analytics and reporting

---

**💰 Building the future of payments with Liora - Now with native blockchain integration via abigen**

*Type-safe smart contract interactions powered by Go and Ethereum*
