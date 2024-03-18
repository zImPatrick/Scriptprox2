package updater

import (
	"errors"
	"fmt"
	"net/http"
)

type UpdateState int

const (
	NONEXISTANT UpdateState = iota
	AVAILABLE
	UPDATING
	DONE
	HOST_IS_FUCKED
)

type UpdateInfo struct {
	RequiredFiles     []string
	State             UpdateState
	EncounteredErrors map[string]error
	DoneCallback      func()
}

func CheckForUpdates() (*UpdateInfo, error) {
	// Checking for working host
	usedUpdaterHost = ""

	hosts := []string{
		UPDATER_HOST,
		UPDATER_HOST_2,
	}

	for _, v := range hosts {
		response, err := http.Get(UPDATER_HOST + "/hashes.php")
		if err != nil || response.StatusCode != 200 {
			fmt.Printf("[Updater] Host %s is not okay\n", v)
			continue
		} else {
			fmt.Println(response.StatusCode)
			fmt.Printf("[Updater] Using host %s\n", v)
			usedUpdaterHost = v
			break
		}
	}

	if usedUpdaterHost == "" {
		// we dont have a good host?
		LastUpdate = &UpdateInfo{
			State: HOST_IS_FUCKED,
		}
		return nil, errors.New("Keinen funktionierenden Update-Host gefunden. Bitte nerv Patrick nach nem Update, danke")
	}

	err := getNewestHashes()
	if err != nil {
		return nil, errors.New("Konnte Updater-Service nicht erreichen. " + err.Error())
	}

	toUpdate := checkAllFiles()
	update := &UpdateInfo{
		RequiredFiles: toUpdate,
		State:         AVAILABLE,
	}

	if len(update.RequiredFiles) == 0 {
		update.State = NONEXISTANT
	}

	LastUpdate = update
	return update, nil
}
