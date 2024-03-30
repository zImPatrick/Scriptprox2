package updater

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
)

var hashes map[string][32]byte
var usedUpdaterHost string

func getNewestHashes() error {
	hashReq, err := http.Get(UPDATER_HOST + "/hashes.php")
	if err != nil {
		fmt.Println("[Updater] Hashes konnten nicht abgerufen werden")
		return err
	}
	body, err := io.ReadAll(hashReq.Body)
	if err != nil {
		fmt.Println("[Updater] Hashes konnten nicht gelesen werden")
		return err
	}

	var data map[string]string
	err = json.Unmarshal(body, &data)

	if err != nil {
		fmt.Println("[Updater] Konnte Hashes nicht lesen")
		return err
	}

	hashes = map[string][32]byte{}
	for filename, hash := range data {
		hexed, _ := hex.DecodeString(hash)
		hashes[filename] = [32]byte(hexed)
	}
	data = nil

	return nil
}

func hashAndCheck(file string, hash [32]byte) (bool, error) {
	opened, err := os.Open(file)

	if err != nil {
		return false, err
	}
	defer opened.Close()

	sh := sha256.New()
	_, err = io.Copy(sh, opened)

	if err != nil {
		return false, err
	}

	checksum := sh.Sum(nil)
	return !slices.Equal(hash[:], checksum[:]), nil
}

func checkAllFiles() []string {
	filesWithWrongHash := []string{}

	for fileName, sha256Hash := range hashes {
		shouldUpdate, err := hashAndCheck(fileName, sha256Hash)
		if err != nil && err == os.ErrNotExist {
			shouldUpdate = true
		} else if err != nil {
			fmt.Printf("Fehler beim Checken von %s aufgetreten: %s\n", fileName, err.Error())
		}

		if shouldUpdate {
			filesWithWrongHash = append(filesWithWrongHash, fileName)
		}
	}
	return filesWithWrongHash
}
