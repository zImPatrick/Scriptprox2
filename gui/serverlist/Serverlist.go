package serverlist

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	protos "scriptprox/gui/serverlist/protos"

	"google.golang.org/protobuf/proto"
)

func FrameReader(reader io.Reader, cb func([]byte)) {
	offset := 0

	for {
		var frameHeader [4]uint8
		n, err := reader.Read(frameHeader[:])
		if err != nil {
			if err != io.EOF {
				fmt.Println("reader err: " + err.Error())
				fmt.Println(err)
			}
			return // wenn eof sind wir fertig, sonst scheißen wir drauf
		}
		if n != 4 {
			fmt.Println("n is not 4, bailing")
			return
		}
		frameLength := int(frameHeader[0]) | int(frameHeader[1])<<8 | int(frameHeader[2])<<16 | int(frameHeader[3])<<24
		offset = offset + 4 + frameLength
		frame := make([]byte, frameLength)
		reader.Read(frame)
		cb(frame)
	}
}

func GotChunk(chunk []byte) {
	var server protos.Server
	proto.Unmarshal(chunk, &server)
}

func RequestServerlist(OnServer func(*protos.Server)) {
	resp, _ := http.Get("https://servers-frontend.fivem.net/api/servers/streamRedir/")
	// es klappt!
	dat, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(dat))

	FrameReader(resp.Body, func(b []byte) {
		var server protos.Server
		proto.Unmarshal(b, &server)
		OnServer(&server)
	})
}
