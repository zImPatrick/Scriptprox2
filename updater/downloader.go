package updater

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
)

func downloadAndReplace(fileName string) error {
	f, err := os.OpenFile(fileName, os.O_CREATE | os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	url, _ := url.JoinPath(UPDATER_HOST, fileName)
	req, err := http.Get(url)
	if err != nil {
		return err
	}
	
	defer req.Body.Close()
	_, err = io.Copy(f, req.Body)
	if err != nil {
		return err
	}

	return nil
}

func InstallUpdate(update *UpdateInfo) {
	update.State = UPDATING
	var waitgroup sync.WaitGroup
	for _, file := range update.RequiredFiles {
		go func() {
			defer waitgroup.Done()
			err := downloadAndReplace(file)
			if err != nil {
				if update.EncounteredErrors == nil {
					update.EncounteredErrors = map[string]error{}
				}
				update.EncounteredErrors[file] = err
			}
		}()

		waitgroup.Add(1)
	}
	waitgroup.Wait()
	update.State = DONE

	if update.DoneCallback != nil {
		update.DoneCallback()
	}
}