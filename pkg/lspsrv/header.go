package lspsrv

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Header struct {
	ContentLength *int
	ContentType   *string
}

func (h *Header) Write(l *zap.Logger, w io.Writer) error {
	kv := map[string]string{}
	if h.ContentLength != nil {
		kv["Content-Length"] = strconv.Itoa(*h.ContentLength)
	}
	if h.ContentType != nil {
		kv["Content-Type"] = *h.ContentType
	}

	lines := []string(nil)
	for k, v := range kv {
		lines = append(lines, k+": "+v)
	}

	hdr := strings.Join(lines, "\r\n") + "\r\n\r\n"

	_, err := w.Write([]byte(hdr))
	return err
}

func ReadHeader(l *zap.Logger, r io.Reader) (*Header, error) {
	b := make([]byte, 1)
	ab := []byte(nil)
	for len(ab) < 4 || string(ab[len(ab)-4:]) != "\r\n\r\n" {
		_, err := io.ReadFull(r, b)
		if err != nil {
			if err == io.EOF {
				return nil, err
			}
			return nil, errors.Wrap(err, "read until sep")
		}
		ab = append(ab, b[0])
	}
	str := string(ab[:len(ab)-4])

	lines := strings.Split(str, "\r\n")
	h := Header{}
	for _, line := range lines {
		assignIndex := strings.Index(line, ":")
		if assignIndex == -1 {
			return nil, fmt.Errorf("no separator in header line '%s'", line)
		}
		if assignIndex == len(line)-1 {
			return nil, fmt.Errorf("no value in header line '%s'", line)
		}
		if strings.HasPrefix(line, "Content-Length") {
			i, err := strconv.Atoi(strings.TrimSpace(line[assignIndex+1:]))
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("parse Content-Length from '%s'", strings.TrimSpace(line[assignIndex+1:])))
			}
			h.ContentLength = &i
		} else if strings.HasPrefix(line, "Content-Type") {
			s := strings.TrimSpace(line[assignIndex+1:])
			h.ContentType = &s
		}
	}

	return &h, nil
}
