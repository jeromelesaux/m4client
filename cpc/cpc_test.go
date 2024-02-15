package cpc_test

import (
	"os"
	"testing"

	"github.com/jeromelesaux/m4client/cpc"
	"github.com/stretchr/testify/assert"
)

// nolint: funlen
func TestHeader(t *testing.T) {
	t.Parallel()
	t.Run("Checksum16", func(t *testing.T) {
		h := cpc.CpcHead{
			User:        1,
			Filename:    setFilename("TEST.BAS"),
			BlockNum:    1,
			LastBlock:   2,
			Type:        cpc.BasicType,
			Size:        218,
			Address:     400,
			FirstBlock:  1,
			LogicalSize: 400,
			Exec:        0,
			Size2:       400,
		}
		h.Checksum = h.ComputedChecksum16()
		assert.Equal(t, h.Checksum, h.ComputedChecksum16())
	})
	t.Run("ByteEncode", func(t *testing.T) {
		h := cpc.CpcHead{
			User:        1,
			Filename:    setFilename("TEST.BAS"),
			BlockNum:    1,
			LastBlock:   2,
			Type:        cpc.BasicType,
			Size:        218,
			Address:     400,
			FirstBlock:  1,
			LogicalSize: 400,
			Exec:        0,
			Size2:       400,
		}
		h.Checksum = h.ComputedChecksum16()
		assert.Equal(t, h.Checksum, h.ComputedChecksum16())
		_, err := h.Bytes()
		assert.Nil(t, err)
	})
	t.Run("DecodeHeader", func(t *testing.T) {
		h := cpc.CpcHead{
			User:        1,
			Filename:    setFilename("TEST.BAS"),
			BlockNum:    1,
			LastBlock:   2,
			Type:        cpc.BasicType,
			Size:        218,
			Address:     400,
			FirstBlock:  1,
			LogicalSize: 400,
			Exec:        0,
			Size2:       400,
		}
		h.Checksum = h.ComputedChecksum16()
		assert.Equal(t, h.Checksum, h.ComputedChecksum16())
		head, err := h.Bytes()
		assert.Nil(t, err)
		got, err := cpc.BytesCpcHeader(head)
		assert.Nil(t, err)
		assert.Equal(t, h, *got)
	})
	t.Run("DecodeFile", func(t *testing.T) {
		filename := "TEST.BAS"
		h := cpc.CpcHead{
			User:        1,
			Filename:    setFilename(filename),
			BlockNum:    1,
			LastBlock:   2,
			Type:        cpc.BasicType,
			Size:        218,
			Address:     400,
			FirstBlock:  1,
			LogicalSize: 400,
			Exec:        0,
			Size2:       400,
		}
		h.Checksum = h.ComputedChecksum16()
		assert.Equal(t, h.Checksum, h.ComputedChecksum16())
		head, err := h.Bytes()
		assert.Nil(t, err)
		err = os.WriteFile(filename, head, 0o666)
		assert.Nil(t, err)
		fo, err := os.Open(filename)
		assert.Nil(t, err)
		got, err := cpc.NewCpcHeader(fo)
		assert.Nil(t, err)
		assert.Equal(t, h, *got)
		// clean file temporary file
		os.Remove(filename)
	})
}

func setFilename(filename string) [15]byte {
	var d [15]byte
	l := len(filename)
	if l > 15 {
		l = 15
	}
	copy(d[:], filename[0:l])
	return d
}
