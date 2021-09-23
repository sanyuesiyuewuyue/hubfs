/*
 * pathkey.go
 *
 * Copyright 2021 Bill Zissimopoulos
 */
/*
 * This file is part of Hubfs.
 *
 * You can redistribute it and/or modify it under the terms of the GNU
 * Affero General Public License version 3 as published by the Free
 * Software Foundation.
 */

package union

import (
	"crypto/sha256"
	"hash"
)

const Pathkeylen = 16

type Pathkey [Pathkeylen]uint8

// Function ComputePathkey computes the path key for a path.
func ComputePathkey(path string) (k Pathkey) {
	sum := sha256.Sum256([]uint8(path))
	copy(k[1:], sum[:])
	return
}

type PathkeyHash struct {
	hash.Hash
}

func NewPathkeyHash() PathkeyHash {
	return PathkeyHash{sha256.New()}
}

func (h PathkeyHash) Write(s string) {
	h.Hash.Write([]uint8(s))
}

func (h PathkeyHash) ComputePathkey() (k Pathkey) {
	copy(k[1:], h.Hash.Sum(nil))
	return
}
