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

const (
	timeCost    uint32 = 2
	memoryCost  uint32 = 64 * 1024
	parallelism uint8  = 4
	hashLen     uint32 = 32
	saltLen            = 16
)

var (
	ErrInvalidHashFormat = errors.New("invalid argon2 hash format")
	ErrInvalidPassword   = errors.New("invalid password")
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, parallelism, hashLen)

	return encodeHash(salt, hash), nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	salt, expectedHash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	actualHash := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, parallelism, hashLen)
	if subtle.ConstantTimeCompare(actualHash, expectedHash) == 1 {
		return true, nil
	}

	return false, ErrInvalidPassword
}

func encodeHash(salt, hash []byte) string {
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		memoryCost,
		timeCost,
		parallelism,
		b64Salt,
		b64Hash,
	)
}

func decodeHash(encodedHash string) (salt, hash []byte, err error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, ErrInvalidHashFormat
	}

	if parts[1] != "argon2id" {
		return nil, nil, ErrInvalidHashFormat
	}

	var version int
	if _, err = fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, ErrInvalidHashFormat
	}
	if version != argon2.Version {
		return nil, nil, ErrInvalidHashFormat
	}

	var memory, time uint32
	var threads uint8
	if _, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return nil, nil, ErrInvalidHashFormat
	}
	if memory != memoryCost || time != timeCost || threads != parallelism {
		return nil, nil, ErrInvalidHashFormat
	}

	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, ErrInvalidHashFormat
	}

	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, ErrInvalidHashFormat
	}
	if len(hash) != int(hashLen) {
		return nil, nil, ErrInvalidHashFormat
	}

	return salt, hash, nil
}
