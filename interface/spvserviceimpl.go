package _interface

import (
	"errors"
	"os"
	"os/signal"

	"github.com/ioeXNetwork/ioeX.SPV/log"
	"github.com/ioeXNetwork/ioeX.SPV/sdk"
	"github.com/ioeXNetwork/ioeX.SPV/spvwallet"
	"github.com/ioeXNetwork/ioeX.SPV/spvwallet/db"

	"github.com/ioeXNetwork/ioeX.MainChain/bloom"
	. "github.com/ioeXNetwork/ioeX.MainChain/core"
	. "github.com/ioeXNetwork/ioeX.Utility/common"
)

type SPVServiceImpl struct {
	*spvwallet.SPVWallet
	clientId   uint64
	seeds      []string
	accounts   []*Uint168
	proofs     Proofs
	queue      Queue
	addrFilter *sdk.AddrFilter
	listeners  map[TransactionType][]TransactionListener
}

func newSPVServiceImpl(clientId uint64, seeds []string) *SPVServiceImpl {
	return &SPVServiceImpl{
		clientId:  clientId,
		seeds:     seeds,
		listeners: make(map[TransactionType][]TransactionListener),
	}
}

func (service *SPVServiceImpl) RegisterAccount(address string) error {
	account, err := Uint168FromAddress(address)
	if err != nil {
		return errors.New("Invalid address format")
	}
	service.accounts = append(service.accounts, account)
	return nil
}

func (service *SPVServiceImpl) RegisterTransactionListener(listener TransactionListener) {
	listeners := service.listeners[listener.Type()]
	listeners = append(listeners, listener)
	service.listeners[listener.Type()] = listeners
	log.Debug("Listener registered:", listeners)
}

func (service *SPVServiceImpl) SubmitTransactionReceipt(txHash Uint256) error {
	return service.queue.Delete(&txHash)
}

func (service *SPVServiceImpl) VerifyTransaction(proof bloom.MerkleProof, tx Transaction) error {
	if service.SPVWallet == nil {
		return errors.New("SPV service not started")
	}

	// Get Header from main chain
	header, err := service.Headers().GetHeader(proof.BlockHash)
	if err != nil {
		return errors.New("can not get block from main chain")
	}

	// Check if merkleroot is match
	merkleBlock := bloom.MerkleBlock{
		Header:       header.Header,
		Transactions: proof.Transactions,
		Hashes:       proof.Hashes,
		Flags:        proof.Flags,
	}
	txIds, err := bloom.CheckMerkleBlock(merkleBlock)
	if err != nil {
		return errors.New("check merkle branch failed, " + err.Error())
	}
	if len(txIds) == 0 {
		return errors.New("invalid transaction proof, no transactions found")
	}

	// Check if transaction hash is match
	match := false
	for _, txId := range txIds {
		if *txId == tx.Hash() {
			match = true
			break
		}
	}
	if !match {
		return errors.New("transaction hash not match proof")
	}

	return nil
}

func (service *SPVServiceImpl) SendTransaction(tx Transaction) error {
	if service.SPVWallet == nil {
		return errors.New("SPV service not started")
	}

	return service.SPVWallet.SendTransaction(tx)
}

func (service *SPVServiceImpl) Start() error {
	if service.SPVWallet != nil {
		return errors.New("SPV service already started")
	}

	var err error
	service.SPVWallet, err = spvwallet.Init(service.clientId, service.seeds)
	if err != nil {
		return err
	}

	// Initialize proofs db
	service.proofs, err = NewProofsDB()
	if err != nil {
		return err
	}

	service.queue, err = NewQueueDB()
	if err != nil {
		return err
	}

	// Register accounts
	if len(service.accounts) == 0 {
		return errors.New("No account registered")
	}
	for _, account := range service.accounts {
		service.DataStore().Addrs().Put(account, RegisteredAccountScript, db.TypeNotify)
	}

	// Create address filter by accounts
	service.addrFilter = sdk.NewAddrFilter(service.accounts)

	// Set callback
	service.SPVWallet.Blockchain().AddStateListener(service)

	// Handle interrupt signal
	stop := make(chan int, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		for range signals {
			log.Trace("SPV service shutting down...")
			service.Stop()
			stop <- 1
		}
	}()

	// Start SPV service
	service.SPVWallet.Start()

	<-stop

	return nil
}

func (service *SPVServiceImpl) OnTxCommitted(tx Transaction, height uint32) {}

func (service *SPVServiceImpl) OnChainRollback(height uint32) {
	service.queue.Rollback(height)
	service.notifyRollback(height)
}

func (service *SPVServiceImpl) OnBlockCommitted(block bloom.MerkleBlock, txs []Transaction) {
	header := block.Header

	// Store merkle proof
	service.proofs.Put(&bloom.MerkleProof{
		BlockHash:    header.Hash(),
		Height:       header.Height,
		Transactions: block.Transactions,
		Hashes:       block.Hashes,
		Flags:        block.Flags,
	})

	// If no transactions return
	if len(txs) == 0 {
		return
	}

	// Find transactions matches registered accounts
	var matchedTxs []Transaction
	for _, tx := range txs {
		for _, output := range tx.Outputs {
			if service.addrFilter.ContainAddr(output.ProgramHash) {
				matchedTxs = append(matchedTxs, tx)
			}
		}
	}

	// Queue matched transactions
	for _, tx := range matchedTxs {
		item := &QueueItem{
			TxHash:    tx.Hash(),
			BlockHash: header.Hash(),
			Height:    header.Height,
		}

		// Save to queue db
		service.queue.Put(item)
	}

	// Look up for queued transactions
	items, err := service.queue.GetAll()
	if err != nil {
		return
	}
	for _, item := range items {
		//	Get proof from db
		proof, err := service.proofs.Get(&item.BlockHash)
		if err != nil {
			log.Error("Query merkle proof failed, block hash:", item.BlockHash.String())
			return
		}
		//	Get transaction from db
		storeTx, err := service.DataStore().Txs().Get(&item.TxHash)
		if err != nil {
			log.Error("Query transaction failed, tx hash:", item.TxHash.String())
			return
		}
		// Prune the proof by the given transaction id
		proof = getTransactionProof(proof, storeTx.TxId)

		// Notify listeners
		service.notifyTransaction(*proof, storeTx.Data, header.Height-item.Height)
	}
}

func (service *SPVServiceImpl) notifyTransaction(proof bloom.MerkleProof, tx Transaction, confirmations uint32) {
	listeners := service.listeners[tx.TxType]
	for _, listener := range listeners {
		if listener.Confirmed() {
			if confirmations >= getConfirmations(tx) {
				go listener.Notify(proof, tx)
			}
		} else {
			go listener.Notify(proof, tx)
		}
	}
}

func (service *SPVServiceImpl) notifyRollback(height uint32) {
	for _, group := range service.listeners {
		for _, listener := range group {
			go listener.Rollback(height)
		}
	}
}

func getConfirmations(tx Transaction) uint32 {
	// TODO user can set confirmations attribute in transaction,
	// if the confirmation attribute is set, use it instead of default value
	return DefaultConfirmations
}

func getTransactionProof(proof *bloom.MerkleProof, txHash Uint256) *bloom.MerkleProof {
	// TODO Pick out the merkle proof of the transaction
	return proof
}
