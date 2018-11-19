package main

import (
	"fmt"
)

type CpcHead struct {
	User      byte
	Filename  [8]byte
	Extension [3]byte
	NotUsed   [6]byte
	Type      byte
	NotUsed2  [2]byte
	Address   byte
	Pad       byte
	Size      byte
	Exec      byte
	NotUsed3  [36]byte
	Size2     byte
	Pad2      byte
	Checksum  byte
	NotUsed4  [59]byte
}

func (c *CpcHead) ToString() string {
	return fmt.Sprintf("User:%x\nFilename:%s\nExtension:%s\n",
		int(c.User),
		string(c.Filename[:]),
		string(c.Extension[:]))
}

var (
	xorstream1 = []byte{0xE2, 0x9D, 0xDB, 0x1A, 0x42, 0x29, 0x39, 0xC6, 0xB3, 0xC6, 0x90, 0x45, 0x8A}
	xorstream2 = []byte{0x49, 0xB1, 0x36, 0xF0, 0x2E, 0x1E, 0x06, 0x2A, 0x28, 0x19, 0xEA}
)

func DecryptHash(data []byte) int {
	size := len(data)
	idx1 := 0
	idx2 := 0
	i := 0
	j := 0
	for j < size {
		if i == 0x80 {
			idx1 = 0
			idx2 = 0
			i = 0
		}
		data[j] ^= xorstream1[idx1]
		idx1++
		data[j] ^= xorstream2[idx2]
		idx2++
		if idx1 == 13 {
			idx1 = 0
		}
		if idx2 == 11 {
			idx2 = 0
		}
		i++
		j++
	}
	return 0
}

func Checksum16(data []byte) uint8 {
	var checksum uint8
	for i := 0; i < len(data); i++ {
		checksum += data[i]
	}
	return checksum
}
