package wallet

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Wallet struct {
	client                *ethclient.Client
	CoinbaseAddress       common.Address
	BundlerBeneficiary    EthEoaAccount
	BundlerSignerAccounts []EthEoaAccount
	DefaultDevAccounts    []EthEoaAccount
}

type EthEoaAccount struct {
	Address    common.Address
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
}

const keyStorePath = "./wallet/tmp"
const coinbaseKeyStorePath = "./wallet/tmp/coinbase/"

// NewWallet creates a new wallet for Besty
func NewWallet(ctx context.Context, ethNodePort string, coinbaseKeystoreFile string) (*Wallet, error) {
	client, err := ethclient.Dial("http://localhost:" + ethNodePort)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Ethereum client: %v", err)
	}

	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	publicKey := &privateKey.PublicKey
	address := crypto.PubkeyToAddress(*publicKey)

	// Define the path to the keystore
	file := coinbaseKeyStorePath + coinbaseKeystoreFile
	ks := keystore.NewKeyStore(keyStorePath, keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	password := ""
	cbAccount, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		client:          client,
		CoinbaseAddress: cbAccount.Address,
		BundlerBeneficiary: EthEoaAccount{
			Address:    address,
			PublicKey:  publicKey,
			PrivateKey: privateKey,
		},
	}, nil
}
