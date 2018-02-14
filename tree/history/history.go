// Copyright 2018 BBVA. All rights reserved.
// Use of this source code is governed by a Apache 2 License
// that can be found in the LICENSE file

package history

import (
	"verifiabledata/store"
)

var prefixZero = []byte{0x0}
var prefixOne = []byte{0x1}

type HistoryTree struct {

}

func NewHistoryTree(s *store.Store) *HistoryTree {
	return &HistoryTree{}
}
