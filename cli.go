package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	scriptproxHttp "scriptprox/proxies/http"
	"strings"
)

func cli() {
	commands := map[string](func(restOfString string) int){
		"cfx.re/join/": func(restOfString string) int {
			fmt.Println("Cfx.re Link: " + restOfString)
			resp, err := http.DefaultClient.Get("https://cfx.re/join/" + restOfString)
			if err == nil && resp.StatusCode == 200 {
				url := resp.Header.Get("X-Citizenfx-Url")
				scriptproxHttp.ChangeEndpoint(url)

				fmt.Println("Versuche dich zu verbinden...")
				fmt.Println("Falls dies nicht klappt, bitte verwende diesen Befehl in der Konsole:")
				fmt.Println("connect localhost:30120")

				exec.Command("rundll32", "url.dll,FileProtocolHandler", "fivem://connect/localhost:30120")
			} else {
				fmt.Println("Fehler " + err.Error())
				fmt.Println("Fehler " + resp.Status)
			}
			return 0
		},
	}

	returnCode := 0
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%d > ", returnCode)
		input, err := reader.ReadString('\n')
		if err != nil {
			continue
		}

		input = strings.TrimSuffix(input, "\r\n")

		for k, v := range commands {
			if strings.HasPrefix(input, k) {
				returnCode = v(strings.TrimPrefix(input, k))
			}
		}

		fmt.Println(input)
	}
}
