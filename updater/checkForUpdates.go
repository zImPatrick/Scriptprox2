package updater

import "errors"

type UpdateState int

const (
	NONEXISTANT UpdateState = iota
	AVAILABLE
	UPDATING
	DONE
)

type UpdateInfo struct {
	RequiredFiles     []string
	State             UpdateState
	EncounteredErrors map[string]error
	DoneCallback      func()
}

func CheckForUpdates() (*UpdateInfo, error) {
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
