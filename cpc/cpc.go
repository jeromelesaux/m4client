package cpc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// CpcHead structure describes the Amsdos header
type CpcHead struct {
	User        byte
	Filename    [15]byte
	BlockNum    byte
	LastBlock   byte
	Type        byte
	Size        int16
	Address     int16
	FirstBlock  byte
	LogicalSize int16
	Exec        int16
	NotUsed     [0x24]byte
	Size2       int16
	BigLength   byte
	Checksum    int16
	NotUsed4    [0x3B]byte
}

func NewCpcHeader(f *os.File) (*CpcHead, error) {
	header := &CpcHead{}
	data := make([]byte, 128)
	_, err := f.Read(data)
	if err != nil {
		return &CpcHead{}, err
	}
	buf := bytes.NewBuffer(data)
	err = binary.Read(buf, binary.LittleEndian, header)
	if err != nil {
		return &CpcHead{}, err
	}

	return header, nil
}

func (c *CpcHead) Checksum16() uint8 {
	var checksum uint8
	checksum += c.User
	return checksum
}

// ToString Will dislay the CpcHead structure content
func (c *CpcHead) ToString() string {
	return fmt.Sprintf("User:%x\nFilename:%s\nType:%d\nSize:&%.2x\nAddress of loading:&%.2x\nAddress of execution:&%.2x\nChecksum:&%.2x\n",
		int(c.User),
		string(c.Filename[:]),
		c.Type,
		c.Size2,
		c.Address,
		c.Exec,
		c.Checksum)
}

func (c *CpcHead) PrettyPrint() {
	fmt.Printf("%x", *c)
	return
}

var (
	xorstream1 = []byte{0xE2, 0x9D, 0xDB, 0x1A, 0x42, 0x29, 0x39, 0xC6, 0xB3, 0xC6, 0x90, 0x45, 0x8A}
	xorstream2 = []byte{0x49, 0xB1, 0x36, 0xF0, 0x2E, 0x1E, 0x06, 0x2A, 0x28, 0x19, 0xEA}
)

// DecryptHash returns value from decryptage of the data
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

// Checksum16 will generate the checksum of the data amsdos header
func Checksum16(data []byte) uint8 {
	var checksum uint8
	for i := 0; i < len(data); i++ {
		checksum += data[i]
	}
	return checksum
}
