package cfx_console

import (
	"encoding/binary"
	"fmt"
	"net"
)

type CfxConsole struct {
	connection  net.Conn
	IsConnected bool
}

// todo (maybe): implement reading by registering
// to server:
// https://github.com/citizenfx/fivem/blob/53dd7f8945d42626c39da2345c2ccd51e9f1faa1/code/components/devcon/src/DevConServer.cpp#L364

func (cfx *CfxConsole) Connect() error {
	cfx.Disconnect()

	conn, err := net.Dial("tcp", "127.0.0.1:29200")
	if err != nil {
		return err
	}

	cfx.connection = conn
	cfx.IsConnected = true

	// loop für disconnection stuff
	go func() {
		buf := make([]byte, 1)
		for cfx.IsConnected {
			fmt.Println("loop")
			_, err := conn.Read(buf)
			if err != nil {
				cfx.IsConnected = false
				cfx.Disconnect()
			}
		}
	}()

	return nil
}

func (cfx *CfxConsole) Disconnect() {
	cfx.IsConnected = false
	if cfx.connection != nil {
		cfx.connection.Close()
		cfx.connection = nil
	}
}

func (cfx *CfxConsole) SendCommand(command string) error {
	if !cfx.IsConnected {
		err := cfx.Connect()
		if err != nil {
			return err
		}
	}

	// https://github.com/citizenfx/fivem/blob/53dd7f8945d42626c39da2345c2ccd51e9f1faa1/code/components/devcon/src/DevConServer.cpp#L405
	byteArr := binary.LittleEndian.AppendUint32([]byte{}, 0x444E4D43) // CMND
	byteArr = binary.LittleEndian.AppendUint16(byteArr, 0x0)          // protocol (unused)
	byteArr = binary.LittleEndian.AppendUint32(byteArr, 0x0)          // length (unused??)
	byteArr = binary.LittleEndian.AppendUint16(byteArr, 0x0)          // unused, just read
	byteArr = append(byteArr, []byte(command)...)
	byteArr = append(byteArr, 0x00)

	_, err := cfx.connection.Write(byteArr)
	if err != nil {
		fmt.Println(err)
		// lets disconnect for safety
		cfx.Disconnect()
		return err
	}

	return nil
}
