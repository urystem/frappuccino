package helpers

import (
	"crypto/md5"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
)

type JSONB map[string]any

func CreateMd5Hash(text string) string {
	hasher := md5.New()
	_, err := io.WriteString(hasher, text)
	if err != nil {
		// panic(err)
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &j)
}
