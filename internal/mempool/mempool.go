package mempool

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/transeptorlabs/betsy/internal/client"
	"github.com/transeptorlabs/betsy/internal/data"
)

type MempoolEntry struct {
	op     *data.UserOpV7Hexify
	status string
}

type UserOpMempool struct {
	userOps                  map[common.Hash]MempoolEntry
	mutex                    sync.Mutex
	ethClient                *ethclient.Client
	epAddress                common.Address
	bundlerClient            *client.BundlerClient
	ticker                   *time.Ticker
	isRunning                bool
	done                     chan bool
	mempoolRefreshErrorCount int
}

func NewUserOpMempool(epAddress common.Address, ethClient *ethclient.Client, bundlerUrl string) *UserOpMempool {
	return &UserOpMempool{
		userOps:                  make(map[common.Hash]MempoolEntry),
		epAddress:                epAddress,
		ethClient:                ethClient,
		bundlerClient:            client.NewBundlerClient(bundlerUrl),
		isRunning:                false,
		done:                     make(chan bool),
		mempoolRefreshErrorCount: 0,
	}
}

func (m *UserOpMempool) GetUserOps() map[common.Hash]*data.UserOpV7Hexify {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var ops = make(map[common.Hash]*data.UserOpV7Hexify)

	for k, v := range m.userOps {
		ops[k] = v.op
	}

	return ops
}

func (m *UserOpMempool) addUserOp(op *data.UserOpV7Hexify) error {
	log.Debug().Msgf("Attempting to add userOp to mempool: %#v\n", op)
	m.mutex.Lock()
	defer m.mutex.Unlock()

	userOpHash, err := op.GetUserOpHash(m.epAddress, m.ethClient)
	if err != nil {
		return err
	}

	if _, ok := m.userOps[userOpHash]; ok {
		log.Debug().Msgf("Skipping userOp already in mempool(userOpHash): %s\n", userOpHash)
		return nil
	}

	m.userOps[userOpHash] = MempoolEntry{
		op:     op,
		status: "pending",
	}
	log.Debug().Msgf("Successfully added userOp in mempool(userOpHash): %s\n", userOpHash)

	return nil
}

func (m *UserOpMempool) refreshMempool() error {
	log.Debug().Msg("Refreshing mempool...")

	userOps, err := m.bundlerClient.Debug_bundler_dumpMempool()
	if err != nil {
		return err
	}

	log.Debug().Msgf("Total userOps fetched from bundler(count): %d", len(userOps))
	if len(userOps) > 0 {
		for _, op := range userOps {
			err = m.addUserOp(&op)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *UserOpMempool) Run() error {
	if m.isRunning {
		return nil
	}

	log.Info().Msg("Starting up Mempool...")
	m.ticker = time.NewTicker(3 * time.Second)
	go func(m *UserOpMempool) {
		for {
			select {
			case <-m.done:
				return
			case <-m.ticker.C:
				err := m.refreshMempool()
				if err != nil {
					log.Err(err).Msg("Could not refresh mempool")
					m.mempoolRefreshErrorCount = m.mempoolRefreshErrorCount + 1
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
