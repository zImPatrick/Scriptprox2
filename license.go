package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"os"
	"scriptprox/utils"
	"strconv"
	"time"

	"golang.org/x/sys/windows/registry"
)

func DoLicenseCheck() {
	doCheckRequest(hwid())
}

const (
	FingerprintCheckFailure = iota
	DialingError
	ClientPostError
	InvalidHttpCode
	RegistryOpenError
	ReadMachineGuidError
	OtherFailure
)

var fingerprint []byte

func doCheckRequest(hwid []byte) {
	fingerprint, _ = hex.DecodeString("c1a2d66ed06cd7d67064e2dff0cf26db202f8c5c326e920641112b5e0440256c")
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, "https://scriptproxlicense.zimpatrick.workers.dev/v", bytes.NewBuffer(
		[]byte(fmt.Sprintf("%x", hwid)),
	))
	req.Header.Set("X-Vers", fmt.Sprintf("2.%s", utils.Commit))

	resp, err := client.Do(req)
	if err != nil {
		throwError(ClientPostError, err.Error())
	}

	if resp.StatusCode != 200 {
		throwError(InvalidHttpCode, fmt.Sprintf("%d", resp.StatusCode))
	}
}

func hwid() []byte {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		throwError(RegistryOpenError, err.Error())
	}
	defer k.Close()
	s, _, err := k.GetStringValue("MachineGuid")
	if err != nil {
		throwError(ReadMachineGuidError, err.Error())
	}

	normal := []byte(s)
	xored := make([]byte, len(s))
	for i, ch := range normal {
		xored[i] = ch ^ byte(0x69+i)
	}

	return xored
}

func dialerTLSPin(network, addr string) (net.Conn, error) {
	c, err := tls.Dial(network, addr, &tls.Config{})
	if err != nil {
		throwError(DialingError, "")
		return c, err
	}

	connstate := c.ConnectionState()
	valid := false
	for _, cert := range connstate.PeerCertificates {
		der, _ := x509.MarshalPKIXPublicKey(cert.PublicKey)
		hash := sha256.Sum256(der)

		if bytes.Equal(hash[0:], fingerprint) {
			valid = true
			break
		}
	}

	if !valid {
		// var subjects []string

		// for _, cert := range connstate.PeerCertificates {
		// 	subjects = append(subjects, cert.Subject.String())
		// }

		throwError(FingerprintCheckFailure, "")
	}

	return c, nil
}

func throwError(err int, additionalInfo string) { // vllt error auch noch an den server reporten? wär cool
	if additionalInfo != "" {
		additionalInfo = " (" + additionalInfo + ")"
	}
	formatted := fmt.Sprintf("LE %d%s", err, additionalInfo)
	os.WriteFile("le.txt", []byte(strconv.Itoa(int(time.Now().Unix()))+" "+formatted), 0600)
	fmt.Println(formatted)

	os.Exit(1)
}
