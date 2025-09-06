package thebus

import (
	"crypto/rand"
	"sync"
	"time"
)

// IDGenerator -
type IDGenerator func() string

func init() {
	_, err := rand.Read(stateIDGen.entropy[:])
	if err != nil {
		panic(err)
	}
}

// crockfordAlphabet see https://www.baeldung.com/cs/crockfords-base32-encoding
// for a simple article about it
const crockfordAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// stateIDGenerator - allows me to keep track of the lastTs and
// generate the entropy again in case we change millisecond
type stateIDGenerator struct {
	mu      sync.Mutex
	lastTs  int64
	entropy [10]byte
}

// next generate again the entropy if the ms increased.
// Get the latest if decrease and increase it if we are in the same ms
func (s *stateIDGenerator) next() (ts int64, e [10]byte) {
	now := time.Now().UnixMilli()
	if now > s.lastTs {
		s.lastTs = now
		_, err := rand.Read(s.entropy[:])
		if err != nil {
			panic(err)
		}
	} else if now == s.lastTs {
		incrementEntropyBE(&s.entropy)
	} else {
		now = s.lastTs
		incrementEntropyBE(&s.entropy)
	}
	return now, s.entropy
}

// stateIDGen - global var for the state
// see init for the default entropy
var stateIDGen = &stateIDGenerator{
	lastTs: time.Now().UnixMilli(),
}

// DefaultIDGenerator - generate a new uniq id.
// See WithIDGenerator Option if you want to include your own generator (ex: uuid.UUID from Google)
func DefaultIDGenerator() string {
	stateIDGen.mu.Lock()
	ts, e := stateIDGen.next()
	stateIDGen.mu.Unlock()
	return encodeULIDLike(ts, e)
}

// incrementEntropyBE increment the entropy.
// See it as a km/miles counter on a car. The latest number on the counter is increased
// each km/miles. If the increase operation result is 0, it means we have to increase the next number (index - 1)
func incrementEntropyBE(e *[10]byte) {
	for i := 9; i >= 0; i-- {
		e[i]++
		if e[i] != 0 {
			return
		}
	}
}

func encodeULIDLike(tsMillis int64, entropy [10]byte) string {
	buf := [26]byte{}
	writeTimestamp48ToBase32(uint64(tsMillis), buf[0:10])
	writeEntropy80ToBase32(entropy, buf[10:26])
	return string(buf[:])
}

// writeTimestamp48ToBase32
func writeTimestamp48ToBase32(ts uint64, out []byte) {
	if len(out) != 10 {
		panic("writeTimestamp48ToBase32: out must be len 10")
	}
	shifts := [...]uint{45, 40, 35, 30, 25, 20, 15, 10, 5, 0}
	for i, sh := range shifts {
		idx := byte((ts >> sh) & 0x1F) // 0x1F = 5 bits
		out[i] = crockfordAlphabet[idx]
	}
}

func writeEntropy80ToBase32(e [10]byte, out []byte) {
	if len(out) != 16 {
		panic("writeEntropy80ToBase32: out must be len 16")
	}
	var acc uint64
	bits := 0
	pos := 0

	for i := 0; i < len(e); i++ {
		acc = (acc << 8) | uint64(e[i]) // on pousse 8 bits
		bits += 8
		// tant qu'on a au moins 5 bits en stock, on sort 1 char
		for bits >= 5 {
			idx := byte((acc >> uint(bits-5)) & 0x1F) // prend les 5 bits de tête
			out[pos] = crockfordAlphabet[idx]
			pos++
			bits -= 5
		}
	}
}
