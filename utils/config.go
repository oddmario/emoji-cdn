package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tidwall/gjson"
)

var Config gjson.Result
var configHash string = ""

func LoadConfig(firstLoad bool) {
	cfgFilename := "config.debug.json"

	path, _ := filepath.Abs("./" + cfgFilename)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		cfgFilename = "config.json"
	}

	path, _ = filepath.Abs("./" + cfgFilename)
	cfg_content, _ := os.ReadFile(path)
	cfgContentString := string(cfg_content)

	if !gjson.Valid(cfgContentString) {
		if firstLoad {
			fmt.Println("Malformed configuration file")
			os.Exit(1)
		}

		fmt.Println("[WARNING] Malformed configuration file")
	} else {
		hash := md5.Sum([]byte(cfgContentString))
		md5Hash := hex.EncodeToString(hash[:])

		if md5Hash != configHash {
			Config = gjson.Parse(cfgContentString)

			if !firstLoad {
				fmt.Println("[CONFIG] The configuration file was updated and has been reloaded")
			}

			configHash = md5Hash
		}
	}
}
