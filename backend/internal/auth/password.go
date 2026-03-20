package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type PasswordHashParams struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

var defaultParams = PasswordHashParams{
	Time:    2,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
	SaltLen: 16,
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("empty_password")
	}

	salt := make([]byte, defaultParams.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, defaultParams.Time, defaultParams.Memory, defaultParams.Threads, defaultParams.KeyLen)
	return fmt.Sprintf(
		"argon2id$%s$%s",
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(hash),
	), nil
}

func VerifyPassword(password string, encodedHash string) (bool, error) {
	// Expected format: argon2id$<salt_b64>$<hash_b64>
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 3 || parts[0] != "argon2id" {
		return false, errors.New("invalid_hash_format")
	}

	salt, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}
	wantHash, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return false, err
	}

	gotHash := argon2.IDKey([]byte(password), salt, defaultParams.Time, defaultParams.Memory, defaultParams.Threads, defaultParams.KeyLen)
	if len(gotHash) != len(wantHash) {
		return false, nil
	}
	if subtle.ConstantTimeCompare(gotHash, wantHash) != 1 {
		return false, nil
	}
	return true, nil
}

func splitPreserveDollar(s string) []string {
	// Split by '$' but keep simple for our fixed format.
	// Format uses $ between fields.
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '$' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

