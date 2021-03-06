// Cache package implements bridge to fast in-memory object cache.
package cache

import (
	"fantom-api-graphql/internal/types"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// PullAccount extracts account information from the in-memory cache if available.
func (b *MemBridge) PullBlock(num *hexutil.Uint64) *types.Block {
	// try to get the account data from the cache
	data, err := b.cache.Get(num.String())
	if err != nil {
		// cache returns ErrEntryNotFound if the key does not exist
		return nil
	}

	// do we have the data?
	blk, err := types.UnmarshalBlock(data)
	if err != nil {
		b.log.Criticalf("can not decode block data from in-memory cache; %s", err.Error())
		return nil
	}

	return blk
}

// PushBlockByNumber stores provided block in the in-memory cache.
func (b *MemBridge) PushBlock(blk *types.Block) error {
	// we need valid account
	if nil == blk {
		return fmt.Errorf("undefined block can not be pushed to the in-memory cache")
	}

	// encode account
	data, err := blk.Marshal()
	if err != nil {
		b.log.Criticalf("can not marshal block to JSON; %s", err.Error())
		return err
	}

	// set the data to cache by block number
	return b.cache.Set(blk.Number.String(), data)
}
