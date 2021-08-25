package internal

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rodaine/table"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

const DefaultAlgorithm = "sha1"
const DefaultDigits = 6
const DefaultPeriod = 30

var configDir = filepath.Join(homeDir(), ".totp")
var configFilePath = filepath.Join(configDir, "config.json")

type Json struct {
	Issuer     string `json:"issuer"`
	Identifier string `json:"identifier"`
	Algorithm  string `json:"algorithm"`
	Digits     int    `json:"digits"`
	Period     int64  `json:"period"`
	Secret     string `json:"secret"`
}

type Result struct {
	Issuer     string
	Identifier string
	Totp       string
	Remaining  int64
}

func Add(json *Json) error {
	_, err := base32.StdEncoding.DecodeString(json.Secret)
	if err != nil {
		return fmt.Errorf("secret must be base32 format")
	}
	err = write([]Json{*json}, false)
	if err != nil {
		return err
	}
	return nil
}

func Delete(idx int) error {
	list, err := read()
	if err != nil {
		return err
	}
	deletedList := append(list[:idx], list[idx+1:]...)
	return write(deletedList, true)
}

func Print() error {
	resultTable := table.New("No", "ISSUER", "IDENTIFIER", "TOTP", "REMAINING")
	resultTable.WithPadding(5)

	list, err := read()
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("An error occurred while read entries: %v", err))
	}

	for index, data := range list {
		result, err := generate(&data)
		if err != nil {
			return err
		}
		resultTable.AddRow(index, result.Issuer, result.Identifier, result.Totp, result.Remaining)
	}
	resultTable.Print()
	return nil
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

func read() ([]Json, error) {
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	result := make([]Json, 0)
	err = json.Unmarshal(file, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func write(list []Json, isOverWrite bool) error {
	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}

	if isOverWrite {
		b, err := json.Marshal(list)
		if err != nil {
			return err
		}
		return os.WriteFile(configFilePath, b, 0644)

	} else {
		registeredList, err := read()
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}

		mergedList := append(registeredList, list...)
		b, err := json.Marshal(mergedList)
		if err != nil {
			return err
		}
		return os.WriteFile(configFilePath, b, 0644)
	}
}

func generate(p *Json) (*Result, error) {
	tm := time.Now().Unix()

	counterBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(counterBytes, uint64(tm/p.Period))
	secretKey, err := base32.StdEncoding.DecodeString(p.Secret)
	if err != nil {
		return nil, fmt.Errorf("secret must be base32 format")
	}
	hmacInit := hmac.New(sha1.New, secretKey)
	_, err = hmacInit.Write(counterBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to compute HMAC")
	}

	hash := hmacInit.Sum(nil)
	offset := hash[len(hash)-1] & 0x0F
	hash = hash[offset : offset+4]
	hash[0] = hash[0] & 0x7F
	decimal := binary.BigEndian.Uint32(hash)
	pass := decimal % uint32(math.Pow10(p.Digits))

	totpStr := strconv.Itoa(int(pass))
	for len(totpStr) != p.Digits {
		totpStr = "0" + totpStr
	}

	return &Result{
		Issuer:     p.Issuer,
		Identifier: p.Identifier,
		Totp:       totpStr,
		Remaining:  p.Period - (tm % p.Period),
	}, nil
}
