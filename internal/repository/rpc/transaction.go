/*
Rpc package implements bridge to Lachesis full node API interface.

We recommend using local IPC for fast and the most efficient inter-process communication between the API server
and an Opera/Lachesis node. Any remote RPC connection will work, but the performance may be significantly degraded
by extra networking overhead of remote RPC calls.

You should also consider security implications of opening Lachesis RPC interface for a remote access.
If you considering it as your deployment strategy, you should establish encrypted channel between the API server
and Lachesis RPC interface with connection limited to specified endpoints.

We strongly discourage opening Lachesis RPC interface for unrestricted Internet access.
*/
package rpc

import (
	"fantom-api-graphql/internal/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Transaction returns information about a blockchain transaction by hash.
func (b *OperaBridge) Transaction(hash *types.Hash) (*types.Transaction, error) {
	// keep track of the operation
	b.log.Debugf("loading transaction %s", hash.String())

	// call for data
	var trx types.Transaction
	err := b.rpc.Call(&trx, "ftm_getTransactionByHash", hash)
	if err != nil {
		b.log.Error("transaction could not be extracted")
		return nil, err
	}

	// is there a block reference already?
	if trx.BlockNumber != nil {
		// get transaction receipt
		var rec struct {
			Index             hexutil.Uint64  `json:"transactionIndex"`
			CumulativeGasUsed hexutil.Uint64  `json:"cumulativeGasUsed"`
			GasUsed           hexutil.Uint64  `json:"gasUsed"`
			ContractAddress   *common.Address `json:"contractAddress"`
			Status            hexutil.Uint64  `json:"status"`
		}

		// call for the transaction receipt data
		err := b.rpc.Call(&rec, "ftm_getTransactionReceipt", hash)
		if err != nil {
			b.log.Errorf("can not get receipt for transaction %s", hash)
			return nil, err
		}

		// copy some data
		trx.Index = &rec.Index
		trx.CumulativeGasUsed = &rec.CumulativeGasUsed
		trx.GasUsed = &rec.GasUsed
		trx.ContractAddress = rec.ContractAddress
		trx.Status = &rec.Status
	}

	// keep track of the operation
	b.log.Debugf("transaction %s loaded", hash.String())
	return &trx, nil
}
