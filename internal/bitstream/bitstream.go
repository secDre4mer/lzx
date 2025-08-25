package bitstream

import (
	"bufio"
	"errors"
	"io"
)

type BitStream struct {
	Internal io.Reader

	// Currently cached bits
	cachedData int64
	// Number of bits cached in cachedData
	cacheSize int
}

func New(reader io.Reader) *BitStream {
	// Wrap reader in a buffered reader to optimize small reads
	return &BitStream{Internal: bufio.NewReader(reader)}
}

// ReadBits reads c bits from the stream and advances the read position.
func (b *BitStream) ReadBits(c int) (int64, error) {
	bits, err := b.PeekBits(c)
	if err != nil {
		return 0, err
	}
	b.cacheSize -= c
	b.cachedData = b.cachedData & (1<<b.cacheSize - 1)
	return bits, nil
}

// PeekBits reads c bits from the stream without advancing the read position.
func (b *BitStream) PeekBits(c int) (int64, error) {
	if c > 32 {
		return 0, errors.New("invalid bit read")
	}
	for c > b.cacheSize {
		// Read new data
		var nextData [2]byte
		if _, err := io.ReadFull(b.Internal, nextData[:]); err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return 0, err
		}
		b.cachedData = b.cachedData<<16 | int64(nextData[1])<<8 | int64(nextData[0])
		b.cacheSize += 16
	}
	result := b.cachedData >> (b.cacheSize - c)
	return result, nil
}

// Align the stream to the next 16-bit boundary.
func (b *BitStream) Align() {
	b.cacheSize -= b.cacheSize % 16 // Align to 16 bits
	b.cachedData = b.cachedData & (1<<b.cacheSize - 1)
}

// BitsLeft returns the number of bits currently cached in the stream.
func (b *BitStream) BitsLeft() int {
	return b.cacheSize
}
