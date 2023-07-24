package updater

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var hashes map[string]string
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

	hashes = data

	return nil
}

func hashAndCheck(file string, hash string) (bool, error) {
	opened, err := os.Open(file)

	if err != nil {
		return false, err
	}

	bytes, err := io.ReadAll(opened)

	if err != nil {
		return false, err 
	}

	checksum := sha256.Sum256(bytes)
	return fmt.Sprintf("%x", checksum) != hash, nil
}

func checkAllFiles() []string {
	filesWithWrongHash := make([]string, 0)

	for fileName, sha256Hash := range hashes {
		shouldUpdate, err := hashAndCheck(fileName, sha256Hash)
		if err != nil && err == os.ErrNotExist {
			shouldUpdate = true
		} else if (err != nil) {
			fmt.Printf("Fehler beim Checken von %s aufgetreten: %s\n", fileName, err.Error())
		}

		if shouldUpdate {
			filesWithWrongHash = append(filesWithWrongHash, fileName)
		}
	} 

	return filesWithWrongHash
}