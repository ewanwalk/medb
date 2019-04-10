package job

import (
	"bufio"
	"bytes"
	"context"
	"io"
)

var (
	scanQuit = `sigquit`
)

// scan
// read the output of the encode task until it is done or cancelled
func (e *Encode) scan(ctx context.Context, reader io.Reader) {

	scanner := bufio.NewScanner(reader)
	scanner.Split(lines)

	stream := e.report.Stream()
	defer close(stream)

	for scanner.Scan() {

		line := scanner.Text()

		if line == scanQuit || ctx.Err() == context.Canceled {
			return
		}

		select {
		case <-ctx.Done():
			return
		case stream <- line:
			continue
		}
	}

}

// dropCR drops a terminal \r from the data.
func dropCarriageReturn(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// lines
// used to scan lines which may contain multiple returns (\r) in a single line
func lines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// remove \r\n
	if i := bytes.Index(data, []byte{'\r', '\n'}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, dropCarriageReturn(data[0:i]), nil
	}

	// remove any remaining \r
	if i := bytes.Index(data, []byte{'\r'}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCarriageReturn(data[0:i]), nil
	}

	if len(data) == 7 && string(data) == "sigquit" {
		return len(data), data, nil
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCarriageReturn(data), nil
	}
	// Request more data.
	return 0, nil, nil
}
