package http

import (
	"io"
	"net/http"
)

// <3 https://stackoverflow.com/a/39261064/10956224
func flushCopy(dst io.Writer, src io.Reader) (written int64, err error) {
	buf := make([]byte, 1024*8)
	flusher, canFlush := dst.(http.Flusher)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				if canFlush {
					flusher.Flush()
				}
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}
