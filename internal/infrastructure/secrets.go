package infrastructure

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CredentialFiles are the only files owned by infrastructure secret generation.
var CredentialFiles = [...]string{"redis.password", "redis.acl"}

var passwordValue = regexp.MustCompile(`^[a-f0-9]{64}$`)

// Password reads the generated credential without exposing values in errors.
func Password(cfg Config) (string, error) {
	path := filepath.Join(cfg.Secrets, "redis.password")
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return "", errors.New("Redis credential must be a regular file")
	}
	data, err := os.ReadFile(path)
	if err != nil || !passwordValue.MatchString(strings.TrimSpace(string(data))) {
		return "", errors.New("Redis credential is unavailable or malformed")
	}
	return strings.TrimSpace(string(data)), nil
}

// PasswordHash is the Redis ACL SHA-256 verifier, not a plaintext password.
func PasswordHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// GenerateSecrets preserves credentials and rejects an inconsistent ACL file.
func GenerateSecrets(cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	if err := os.MkdirAll(cfg.Secrets, 0o700); err != nil {
		return errors.New("cannot create infrastructure credentials directory")
	}
	info, err := os.Lstat(cfg.Secrets)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return errors.New("credentials directory must be a regular directory")
	}
	if err = os.Chmod(cfg.Secrets, 0o700); err != nil {
		return errors.New("cannot protect infrastructure credentials directory")
	}
	path := filepath.Join(cfg.Secrets, "redis.password")
	if _, err = os.Lstat(path); errors.Is(err, os.ErrNotExist) {
		var value [32]byte
		if _, err = rand.Read(value[:]); err != nil {
			return errors.New("cannot generate Redis credential")
		}
		if err = preserveFile(path, []byte(hex.EncodeToString(value[:])+"\n")); err != nil {
			return err
		}
	}
	password, err := Password(cfg)
	if err != nil {
		return err
	}
	if err = os.Chmod(path, 0o644); err != nil {
		return errors.New("cannot set Redis credential permissions")
	}
	acl := "user default off\nuser " + cfg.RedisUser + " on #" + PasswordHash(password) + " ~" + cfg.KeyPrefix + "* -@all +ping +get +set +del +pttl +hello\n"
	return preserveFile(filepath.Join(cfg.Secrets, "redis.acl"), []byte(acl))
}

func preserveFile(path string, data []byte) error {
	if info, err := os.Lstat(path); err == nil {
		if !info.Mode().IsRegular() {
			return errors.New("infrastructure credential must be a regular file")
		}
		current, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(current, data) {
			return errors.New("infrastructure credential file differs from its expected configuration")
		}
	} else if errors.Is(err, os.ErrNotExist) {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			return errors.New("cannot create infrastructure credential")
		}
		_, writeErr := file.Write(data)
		if errors.Join(writeErr, file.Close()) != nil {
			return errors.New("cannot write infrastructure credential")
		}
	} else {
		return errors.New("cannot inspect infrastructure credential")
	}
	if err := os.Chmod(path, 0o644); err != nil {
		return errors.New("cannot set infrastructure credential permissions")
	}
	return nil
}
