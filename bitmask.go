package utils

import "math/big"

type Bitmask struct {
	mask *big.Int
}

func (b *Bitmask) checkNil() {
	if b.mask == nil {
		b.mask = big.NewInt(0)
	}
}

func (b *Bitmask) Clear(n int) {
	b.checkNil()
	b.mask = b.mask.SetBit(b.mask, n, 0)
}

func (b *Bitmask) Set(n int) {
	b.checkNil()
	b.mask = b.mask.SetBit(b.mask, n, 1)
}

func (b *Bitmask) Get(n int) bool {
	b.checkNil()
	return b.mask.Bit(n) != 0
}

func (b *Bitmask) Copy() *Bitmask {
	nb := new(Bitmask)
	nb.mask.Set(b.mask)
	return nb
}
