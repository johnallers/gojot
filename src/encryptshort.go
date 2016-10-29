package sdees

import (
	"hash/fnv"
	"math/rand"
	"strconv"
	"strings"
)

// EncryptOTP runs a XOR encryption on the input string using the random bytes
// in the massive key.
// Random bytes are used starting at a position based on the hash of the input string.
// The starting position is saved as a prefix to the encrypted string
func EncryptOTP(input string) string {
	if strings.Contains(input, ".otp") || len(input) == 0 {
		return input
	}
	key := Cryptkey
	inputb := []byte(input)

	// Get integer hash of input, using some of the random bytes in key as salt
	h := fnv.New32a()
	h.Write(append(inputb, []byte(key)[:100]...))
	inputToNum := h.Sum32()

	// Use random integer to seed and generate random start position
	rand.Seed(int64(inputToNum))
	startPos := rand.Intn(999000-1) + 1

	// Do XOR encryption based on that start position
	keyb := []byte(key[startPos : startPos+len(input)])
	b := make([]byte, len(inputb))
	for i := 0; i < len(inputb); i++ {
		b[i] = inputb[i] ^ keyb[i]
	}

	// Return string as [startposition]-==-[encryptedstring]
	startPosString := strconv.Itoa(startPos)
	return startPosString + "." + EncodeToString(b) + ".otp"
}

// DecryptOTP runs a XOR encryption on the input string using the random bytes
// in the massive key.
// Random bytes are used starting at a position based on the prefix in the input
func DecryptOTP(input string) string {
	if !strings.Contains(input, ".otp") {
		return input
	}
	key := Cryptkey
	parts := strings.Split(input, ".")
	inputb := DecodeString(parts[1])
	startPos, _ := strconv.Atoi(parts[0])
	keyb := []byte(key[startPos : startPos+len(inputb)])
	b := make([]byte, len(inputb))
	for i := 0; i < len(inputb); i++ {
		b[i] = inputb[i] ^ keyb[i]
	}
	return string(b)
}
