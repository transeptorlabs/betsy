package mempool

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/transeptorlabs/betsy/internal/data"
)

type UserOpMempool struct {
	userOps    map[common.Hash]*data.UserOpV7Hexify
	mutex      sync.Mutex
	ethClient  *ethclient.Client
	epAddress  common.Address
	bundlerUrl string
	ticker     *time.Ticker
	isRunning  bool
	done       chan bool
}

func NewUserOpMempool(epAddress common.Address, ethClient *ethclient.Client, bundlerUrl string) *UserOpMempool {
	return &UserOpMempool{
		userOps:    make(map[common.Hash]*data.UserOpV7Hexify),
		epAddress:  epAddress,
		ethClient:  ethClient,
		bundlerUrl: bundlerUrl,
		isRunning:  false,
		done:       make(chan bool),
	}
}

func (m *UserOpMempool) addUserOp(op *data.UserOpV7Hexify) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	userOpHash, err := op.GetUserOpHash(m.epAddress, m.ethClient)
	if err != nil {
		return err
	}

	if _, ok := m.userOps[userOpHash]; ok {
		return nil
	}

	m.userOps[userOpHash] = op

	return nil
}

func (m *UserOpMempool) GetUserOps() map[common.Hash]*data.UserOpV7Hexify {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.userOps
}

func (m *UserOpMempool) fetchUserOps() error {
	log.Debug().Msg("Fecthing userOps from bundler node...")
	return nil
}

func (m *UserOpMempool) Run() error {
	if m.isRunning {
		return nil
	}

	log.Info().Msg("Starting up Mempool...")
	m.ticker = time.NewTicker(5 * time.Second)
	go func(m *UserOpMempool) {
		for {
			select {
			case <-m.done:
				return
			case <-m.ticker.C:
				err := m.fetchUserOps()
				if err != nil {
					continue
				}
			}
		}
	}(m)

	m.isRunning = true

	return nil
}

func (m *UserOpMempool) Stop() {
	if !m.isRunning {
		return
	}

	m.ticker.Stop()
	log.Info().Msg("Shutting down Mempool...")
	m.isRunning = false
	m.done <- true
}
