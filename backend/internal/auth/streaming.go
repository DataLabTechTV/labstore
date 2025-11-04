package auth

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

const StreamingPayload = "STREAMING-AWS4-HMAC-SHA256-PAYLOAD"

type SigV4ChunkedReader struct {
	Body    io.ReadCloser
	PrevSig string
	header  *SigV4ChunkHeader
	data    []byte
}

type SigV4ChunkHeader struct {
	size      int
	signature string
}

func (r *SigV4ChunkedReader) readChunkHeader() error {
	reader := bufio.NewReader(r.Body)

	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return err
	}

	line = strings.TrimSuffix(line, "\r\n")

	headerParts := strings.SplitN(line, ";", 2)
	sizeHex, chunkSig := headerParts[0], headerParts[1]

	size, err := strconv.ParseInt(sizeHex, 16, 64)
	if err != nil {
		return err
	}

	sig, ok := strings.CutPrefix(chunkSig, "chunk-signature=")
	if !ok {
		return errors.New("could not find 'chunk-signature=' prefix")
	}

	r.header = &SigV4ChunkHeader{
		size:      int(size),
		signature: sig,
	}

	return nil
}

func (r *SigV4ChunkedReader) readChunkData() error {
	r.data = make([]byte, r.header.size)

	if _, err := io.ReadFull(r.Body, r.data); err != nil {
		return err
	}

	return nil
}

func (r *SigV4ChunkedReader) readTrailingCRLF() error {
	crlf := make([]byte, 2)

	if _, err := io.ReadFull(r.Body, crlf); err != nil || string(crlf) != "\r\n" {
		return errors.New("invalid chunk termination")
	}

	return nil
}

func (r *SigV4ChunkedReader) verifyChunkSigV4() error {
	// TODO
	return nil
}

func (r *SigV4ChunkedReader) Read(buf []byte) (int, error) {
	if err := r.readChunkHeader(); err != nil {
		return 0, err
	}

	if r.header.size == 0 {
		// done
		return 0, io.EOF
	}

	if err := r.readChunkData(); err != nil {
		return 0, err
	}

	if err := r.readTrailingCRLF(); err != nil {
		return 0, err
	}

	if err := r.verifyChunkSigV4(); err != nil {
		return 0, err
	}

	n := 10000 // !FIXME: replace with read bytes

	return n, nil
}

func (r *SigV4ChunkedReader) Close() error {
	return r.Body.Close()
}
