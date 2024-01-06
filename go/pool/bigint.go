package pool

import (
	"math/big"
	"sync"
)

// bigIntPool is a pool of big.Ints.
// It is used to reduce allocation overhead when creating new big.Ints.
var bigIntPool = sync.Pool{
	New: func() interface{} {
		return new(big.Int)
	},
}

// FetchBigIntFromPool gets a big.Int from the pool, reducing allocation overhead.
func FetchBigIntFromPool() *big.Int {
	bi := bigIntPool.Get().(*big.Int)
	bi.SetInt64(0) // Reset to zero
	return bi
}

// ReleaseBigIntToPool releases a big.Int back to the pool.
func ReleaseBigIntToPool(bi *big.Int) {
	bigIntPool.Put(bi)
}
