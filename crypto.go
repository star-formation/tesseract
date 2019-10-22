/*  Copyright 2019 The tesseract Authors

    This file is part of tesseract.

    tesseract is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    tesseract is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package tesseract

import (
	"crypto/sha256"
	"encoding/binary"

	xrand "golang.org/x/exp/rand"

	nacl "github.com/kevinburke/nacl"
	//naclbox "github.com/kevinburke/nacl/box"
)

// TODO: move to separate file / package

type Tx struct {
	Addr   []byte
	PubKey nacl.Key
	Action []byte
	SeqNum uint64
}

// NOTE: Secure Merkle trees must use different hash functions for leaves and internal nodes to avoid type confusion based attacks [X]

// [X] https://bitslog.com/2018/06/09/leaf-node-weakness-in-bitcoin-merkle-tree-design/

/* TODO: this is for testing

   This will be replaced with a on-chain random beacon for source of
   deterministic _and_ unpredictable entropy.

   See https://dfinity.org/static/dfinity-consensus-0325c35128c72b42df7dd30c22c41208.pdf
   and https://github.com/ethereum/eth2.0-specs/blob/master/specs/core/0_beacon-chain.md
*/
func RandBytes() ([32]byte, error) {
	return sha256.Sum256([]byte("hello world\n")), nil
}

func NewRand() (*xrand.Rand, error) {
	randBytes, err := RandBytes()
	if err != nil {
		return nil, err
	}

	u64, _ := binary.Uvarint([]byte(randBytes[:]))

	// We use https://www.godoc.org/golang.org/x/exp/rand#PCGSource
	// as the math/rand RNG algo is planned to be deprecated.
	// See https://github.com/golang/go/issues/21835
	//
	// x/exp/rand.NewSource defaults to PCGSource
	src := xrand.NewSource(u64)
	return xrand.New(src), nil
}
