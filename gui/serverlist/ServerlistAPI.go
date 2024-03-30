package serverlist

import (
	"fmt"
	"io"
	"net/http"
	protos "scriptprox/gui/serverlist/protos"

	"google.golang.org/protobuf/proto"
)

func FrameReader(reader io.Reader, cb func([]byte), doneCb func()) {
	var frameHeader [4]uint8
	for {
		n, err := io.ReadAtLeast(reader, frameHeader[:], 4)
		if err != nil {
			if err != io.EOF {
				fmt.Println("reader err: " + err.Error())
				fmt.Println(err)
			}
			doneCb()
			return // wenn eof sind wir fertig, sonst scheißen wir drauf
		}
		if n != 4 {
			fmt.Println("n is not 4, bailing")
			return
		}

		frameLength := int(frameHeader[0]) | int(frameHeader[1])<<8 | int(frameHeader[2])<<16 | int(frameHeader[3])<<24
		if frameLength > 65535 {
			fmt.Println("frame too big! aborting")
			return
		}

		frame := make([]byte, frameLength)
		io.ReadAtLeast(reader, frame, frameLength)
		cb(frame)
	}
}

func GotChunk(chunk []byte) {
	var server protos.Server
	proto.Unmarshal(chunk, &server)
}

func RequestServerlist(OnServer func(*protos.Server), OnDone func()) {
	resp, _ := http.Get("https://servers-frontend.fivem.net/api/servers/streamRedir/")

	FrameReader(resp.Body, func(b []byte) {
		server := protos.Server{}
		proto.Unmarshal(b, &server)
		OnServer(&server)
	}, OnDone)
}
