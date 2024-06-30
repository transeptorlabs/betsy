package wallet

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/transeptorlabs/betsy/internal/utils"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type Wallet struct {
	client                    *ethclient.Client
	CoinbaseAddress           common.Address
	BundlerBeneficiaryAddress common.Address
	defaultDevAccounts        []DefaultDevAccount
	keyStore                  *keystore.KeyStore
}

type DefaultDevAccount struct {
	Address    common.Address
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
}

const keyStorePath = "./wallet/tmp"
const coinbaseKeyStorePath = "./wallet/tmp/coinbase/"
const DefaultSeedPhrase = "test test test test test test test test test test test junk"

// NewWallet creates a new wallet for Besty
func NewWallet(ctx context.Context, ethNodePort string, coinbaseKeystoreFile string) (*Wallet, error) {
	client, err := ethclient.Dial("http://localhost:" + ethNodePort)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Ethereum client: %v", err)
	}

	// Define the path to the keystore
	ks := keystore.NewKeyStore(keyStorePath, keystore.StandardScryptN, keystore.StandardScryptP)

	// Load the coinbase account from the keystore
	coinbaseFile := coinbaseKeyStorePath + coinbaseKeystoreFile
	jsonBytes, err := os.ReadFile(coinbaseFile)
	if err != nil {
		return nil, err
	}

	password := ""
	cbAccount, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		return nil, err
	}

	err = utils.RemoveFile(coinbaseFile)
	if err != nil {
		return nil, err
	}

	// Create a new account for the bundler beneficiary
	bAccount, err := createAccount(ks, password)
	if err != nil {
		return nil, err
	}

	defaultDevAccounts, err := GenerateAccountsFromSeed(DefaultSeedPhrase, 10)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		client:                    client,
		keyStore:                  ks,
		defaultDevAccounts:        defaultDevAccounts,
		CoinbaseAddress:           cbAccount.Address,
		BundlerBeneficiaryAddress: bAccount,
	}, nil
}

// GetAccounts returns all account addresses
func (w *Wallet) GetAccounts() []common.Address {
	am := accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false}, w.keyStore)
	return am.Accounts()
}

// PrintDevAccounts prints the default development accounts
func (w *Wallet) PrintDevAccounts() {
	fmt.Println("_________________________Default Development Accounts:_________________________")
	for i, account := range w.defaultDevAccounts {
		fmt.Printf("Account %d:\n", i+1)
		fmt.Printf("Address: %s\n", account.Address.Hex())
		fmt.Printf("Private Key: 0x%s\n\n", hex.EncodeToString(crypto.FromECDSA(account.PrivateKey)))
	}
	fmt.Println("_______________________________________________________________________________")
}

// GetAccount generates a new key and stores it into the key directory, encrypting it with the passphrase and return the address
func createAccount(ks *keystore.KeyStore, password string) (common.Address, error) {
	account, err := ks.NewAccount(password)
	if err != nil {
		return common.Address{}, err
	}

	return account.Address, nil
}

// generateAccountsFromSeed generates a number of accounts from a seed phrase
func GenerateAccountsFromSeed(seedPhrase string, numAccounts int) ([]DefaultDevAccount, error) {
	var accounts []DefaultDevAccount

	// Generate the seed from the mnemonic and create a master key from the seed
	seed := bip39.NewSeed(seedPhrase, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}

	// Derive the Ethereum account keys (follows the path m/44'/60'/0'/0/x, where x is the account index)
	// This path is used to generate hierarchical deterministic (HD) wallets.
	for i := 0; i < numAccounts; i++ {
		purpose, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
		if err != nil {
			return nil, err
		}
		coinType, err := purpose.NewChildKey(bip32.FirstHardenedChild + 60)
		if err != nil {
			return nil, err
		}
		account, err := coinType.NewChildKey(bip32.FirstHardenedChild + 0)
		if err != nil {
			return nil, err
		}
		change, err := account.NewChildKey(0)
		if err != nil {
			return nil, err
		}
		key, err := change.NewChildKey(uint32(i))
		if err != nil {
			return nil, err
		}

		// Generate the private key, public key and Ethereum address
		privateKey, err := crypto.ToECDSA(key.Key)
		if err != nil {
			return nil, err
		}

		publicKey := &privateKey.PublicKey
		address := crypto.PubkeyToAddress(*publicKey)
		devAccount := DefaultDevAccount{
			Address:    address,
			PublicKey:  publicKey,
			PrivateKey: privateKey,
		}

		accounts = append(accounts, devAccount)
	}

	return accounts, nil
}
