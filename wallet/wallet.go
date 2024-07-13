package wallet

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/transeptorlabs/betsy/contracts/entrypoint"
	"github.com/transeptorlabs/betsy/contracts/examples"
	"github.com/transeptorlabs/betsy/contracts/factory"
	"github.com/transeptorlabs/betsy/internal/utils"

	"github.com/rs/zerolog/log"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

const keyStorePath = "./wallet/tmp"
const coinbaseKeyStorePath = "./wallet/tmp/coinbase/"
const DefaultSeedPhrase = "test test test test test test test test test test test junk"
const cfEntrypointV7Address = "0x5FbDB2315678afecb367f032d93F642f64180aa3"

// BundlerWalletDetails contains the details of the bundler wallet
type BundlerWalletDetails struct {
	Beneficiary       common.Address
	Mnemonic          string
	EntryPointAddress common.Address
}

// Wallet contains the details of the wallet for Betsy
type Wallet struct {
	client                      *ethclient.Client
	coinbaseAddress             common.Address
	bundlerBeneficiaryAddress   common.Address
	devAccounts                 []DevAccount
	keyStore                    *keystore.KeyStore
	password                    string
	entryPointAddress           common.Address
	simpleAccountFactoryAddress common.Address
	globalCounterAddress        common.Address
	chainID                     *big.Int
}

// DevAccount contains the details of the default development account
type DevAccount struct {
	Address       common.Address
	PublicKey     *ecdsa.PublicKey
	PrivateKey    *ecdsa.PrivateKey
	PrivateKeyHex string
	Balance       *big.Int
}

type PreDeployedContracts struct {
	EntryPointAddress           common.Address
	SimpleAccountFactoryAddress common.Address
	GlobalCounterAddress        common.Address
}

// NewWallet creates a new wallet for Betsy
func NewWallet(ctx context.Context, ethNodePort string, coinbaseKeystoreFile string) (*Wallet, error) {
	client, err := ethclient.Dial("http://localhost:" + ethNodePort)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Ethereum client: %v", err)
	}

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, err
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

	// Create a new account for the bundler beneficiary in the keystore
	bAccount, err := createAccount(ks, password)
	if err != nil {
		return nil, err
	}

	devAccounts, err := GenerateAccountsFromSeed(DefaultSeedPhrase, 10)
	if err != nil {
		return nil, err
	}

	// Create the wallet
	wallet := &Wallet{
		client:                      client,
		keyStore:                    ks,
		devAccounts:                 devAccounts,
		coinbaseAddress:             cbAccount.Address,
		bundlerBeneficiaryAddress:   bAccount,
		password:                    password,
		entryPointAddress:           common.HexToAddress(""),
		simpleAccountFactoryAddress: common.HexToAddress(""),
		globalCounterAddress:        common.HexToAddress(""),
		chainID:                     chainID,
	}

	// Fund the default development accounts
	for _, account := range wallet.devAccounts {
		err = wallet.fundAccountWithEth(ctx, account.Address)
		if err != nil {
			return nil, err
		}
	}

	//  Deploy the pre-compiled contracts
	err = wallet.deployPreCompiledContracts(ctx)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// GetKeyStoreAccounts returns all the accounts in the keystore
func (w *Wallet) GetKeyStoreAccounts() []common.Address {
	am := accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false}, w.keyStore)
	return am.Accounts()
}

// GetDevAccounts returns the default development accounts
func (w *Wallet) GetDevAccounts(ctx context.Context) ([]DevAccount, error) {
	for i, account := range w.devAccounts {

		balance, err := w.client.BalanceAt(ctx, account.Address, nil)
		if err != nil {
			return nil, err
		}
		w.devAccounts[i].Balance = balance
	}

	return w.devAccounts, nil
}

// GetBundlerWalletDetails returns the details of the bundler wallet
func (w *Wallet) GetBundlerWalletDetails() BundlerWalletDetails {
	return BundlerWalletDetails{
		Beneficiary:       w.bundlerBeneficiaryAddress,
		Mnemonic:          DefaultSeedPhrase,
		EntryPointAddress: w.entryPointAddress,
	}
}

// fundAccountWithEth send 4337 ETH to the account using the coinbase account
func (w *Wallet) fundAccountWithEth(ctx context.Context, toAddress common.Address) error {
	// Unlock the account (in the context of the keystore is necessary because the private key is encrypted for security reasons)
	account, err := w.keyStore.Find(accounts.Account{Address: w.coinbaseAddress})
	if err != nil {
		return err
	}

	err = w.keyStore.Unlock(account, w.password)
	if err != nil {
		return err
	}

	nonce, err := w.client.PendingNonceAt(ctx, w.coinbaseAddress)
	if err != nil {
		return err
	}

	value := new(big.Int)
	value.SetString("4337000000000000000000", 10) // 4337 ETH in Wei
	gasLimit := uint64(21000)                     // The gas limit for a standard ETH transfer is 21000 units.
	gasPrice, err := w.client.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// Sign the transaction and send it
	signedTx, err := w.keyStore.SignTx(account, tx, w.chainID)
	if err != nil {
		return err
	}

	err = w.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return err
	}

	log.Info().Msgf("tx sent to fund (%s): %s", toAddress, signedTx.Hash().Hex())

	return nil
}

func (w *Wallet) GetPreDeployedContracts() PreDeployedContracts {
	return PreDeployedContracts{
		EntryPointAddress:           w.entryPointAddress,
		SimpleAccountFactoryAddress: w.simpleAccountFactoryAddress,
		GlobalCounterAddress:        w.globalCounterAddress,
	}
}

// deployPreCompiledContracts deploys the pre-compiled contracts
func (w *Wallet) deployPreCompiledContracts(ctx context.Context) error {
	log.Info().Msg("Deploying the 4337 EntryPointV7 contract...")
	auth, err := bind.NewKeyedTransactorWithChainID(w.devAccounts[0].PrivateKey, w.chainID)
	if err != nil {
		return err
	}

	entryPointAddress, tx1, _, err := entrypoint.DeployEntryPointV7(auth, w.client)
	time.Sleep(300 * time.Millisecond) // Allow it to be processed by the local node

	receipt, err := bind.WaitMined(ctx, w.client, tx1)
	if err != nil {
		return err
	} else if receipt.Status == types.ReceiptStatusFailed {
		return err
	}

	exists, err := checkContractExistence(ctx, entryPointAddress, w.client)
	if err != nil {
		return err
	}
	if !exists {
		return err
	}

	w.entryPointAddress = entryPointAddress

	log.Info().Msg("Deploying the 4337 SimpleAccountFactory contract...")
	simpleAFAddress, tx2, _, err := factory.DeploySimpleAccountFactoryV7(auth, w.client, entryPointAddress)
	time.Sleep(300 * time.Millisecond) // Allow it to be processed by the local node

	receipt2, err := bind.WaitMined(ctx, w.client, tx2)
	if err != nil {
		return err
	} else if receipt2.Status == types.ReceiptStatusFailed {
		return err
	}

	exists2, err := checkContractExistence(ctx, simpleAFAddress, w.client)
	if err != nil {
		return err
	}
	if !exists2 {
		return err
	}

	w.simpleAccountFactoryAddress = simpleAFAddress

	log.Info().Msg("Deploying the GlobalCounter contract...")
	globalCounterAddress, tx3, _, err := examples.DeployGlobalCounter(auth, w.client)
	time.Sleep(300 * time.Millisecond) // Allow it to be processed by the local node

	receipt3, err := bind.WaitMined(ctx, w.client, tx3)
	if err != nil {
		return err
	} else if receipt3.Status == types.ReceiptStatusFailed {
		return err
	}

	exists3, err := checkContractExistence(ctx, globalCounterAddress, w.client)
	if err != nil {
		return err
	}
	if !exists3 {
		return err
	}

	w.globalCounterAddress = globalCounterAddress

	return nil
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
func GenerateAccountsFromSeed(seedPhrase string, numAccounts int) ([]DevAccount, error) {
	var accounts []DevAccount

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
		devAccount := DevAccount{
			Address:       address,
			PublicKey:     publicKey,
			PrivateKey:    privateKey,
			PrivateKeyHex: "0x" + hex.EncodeToString(crypto.FromECDSA(privateKey)),
			Balance:       big.NewInt(0),
		}

		accounts = append(accounts, devAccount)
	}

	return accounts, nil
}

// checkContractExistence checks if a contract exists at the given address
func checkContractExistence(ctx context.Context, contractAddress common.Address, client *ethclient.Client) (bool, error) {
	stopRequestCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Query the blockchain to see if a contract exists at the given address
	code, err := client.CodeAt(stopRequestCtx, contractAddress, nil)
	if err != nil {
		return false, err
	}

	return len(code) > 0, nil
}

func (w *Wallet) GetEthClient() *ethclient.Client {
	return w.client
}

func (w *Wallet) GetGethClient() *gethclient.Client {
	return gethclient.New(w.client.Client())
}
