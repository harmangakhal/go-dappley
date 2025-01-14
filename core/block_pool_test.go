// Copyright (C) 2018 go-dappley authors
//
// This file is part of the go-dappley library.
//
// the go-dappley library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either pubKeyHash 3 of the License, or
// (at your option) any later pubKeyHash.
//
// the go-dappley library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-dappley library.  If not, see <http://www.gnu.org/licenses/>.
//

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRUCacheWithIntKeyAndValue(t *testing.T) {
	bp := NewBlockPool(5)
	assert.Equal(t, 0, bp.blkCache.Len())
	const addCount = 200
	for i := 0; i < addCount; i++ {
		if bp.blkCache.Len() == ForkCacheLRUCacheLimit {
			bp.blkCache.RemoveOldest()
		}
		bp.blkCache.Add(i, i)
	}
	//test blkCache is full
	assert.Equal(t, ForkCacheLRUCacheLimit, bp.blkCache.Len())
	//test blkCache contains last added key
	assert.Equal(t, true, bp.blkCache.Contains(199))
	//test blkCache oldest key = addcount - BlockPoolLRUCacheLimit
	assert.Equal(t, addCount-ForkCacheLRUCacheLimit, bp.blkCache.Keys()[0])
}
