//go:build 386 || arm || mips || mipsle

package hash

import (
	"encoding/binary"
	"math/bits"
)

const (
	prime1 uint32 = 2654435761
	prime2 uint32 = 2246822519
	prime3 uint32 = 3266489917
	prime4 uint32 = 668265263
	prime5 uint32 = 374761393
)

var prime1v = prime1

// xxHash implementation for 32 bit system
func Sum(b []byte) uintptr {
	n := len(b)
	h32 := uint32(n)

	if n < 16 {
		h32 += prime5
	} else {
		v1 := prime1 + prime2
		v2 := prime2
		v3 := uint32(0)
		v4 := -prime1
		p := 0
		for n := n - 16; p <= n; p += 16 {
			sub := b[p:][:16] //BCE hint for compiler
			v1 = rol13(v1+u32(sub[:])*prime2) * prime1
			v2 = rol13(v2+u32(sub[4:])*prime2) * prime1
			v3 = rol13(v3+u32(sub[8:])*prime2) * prime1
			v4 = rol13(v4+u32(sub[12:])*prime2) * prime1
		}
		b = b[p:]
		n -= p
		h32 += rol1(v1) + rol7(v2) + rol12(v3) + rol18(v4)
	}

	p := 0
	for n := n - 4; p <= n; p += 4 {
		h32 += u32(b[p:p+4]) * prime3
		h32 = rol17(h32) * prime4
	}
	for p < n {
		h32 += uint32(b[p]) * prime5
		h32 = rol11(h32) * prime1
		p++
	}

	h32 ^= h32 >> 15
	h32 *= prime2
	h32 ^= h32 >> 13
	h32 *= prime3
	h32 ^= h32 >> 16

	return uintptr(h32)
}

func u32(b []byte) uint32 { return binary.LittleEndian.Uint32(b) }

func rol1(x uint32) uint32  { return bits.RotateLeft32(x, 1) }
func rol7(x uint32) uint32  { return bits.RotateLeft32(x, 7) }
func rol11(x uint32) uint32 { return bits.RotateLeft32(x, 11) }
func rol12(x uint32) uint32 { return bits.RotateLeft32(x, 12) }
func rol13(x uint32) uint32 { return bits.RotateLeft32(x, 13) }
func rol17(x uint32) uint32 { return bits.RotateLeft32(x, 17) }
func rol18(x uint32) uint32 { return bits.RotateLeft32(x, 18) }