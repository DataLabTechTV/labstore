package auth

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/DataLabTechTV/labstore/backend/pkg/logger"
)

const StreamingPayload = "STREAMING-AWS4-HMAC-SHA256-PAYLOAD"

type SigV4ChunkedReader struct {
	Body      io.ReadCloser
	PrevSig   string
	SecretKey string
	Timestamp string
	Scope     string

	reader *bufio.Reader
	header *SigV4ChunkHeader
	data   []byte
}

type SigV4ChunkHeader struct {
	size      int
	signature string
}

func (r *SigV4ChunkedReader) Read(buf []byte) (int, error) {
	if r.reader == nil {
		r.reader = bufio.NewReader(r.Body)
	}

	if len(r.data) > 0 {
		n := copy(buf, r.data)
		r.data = r.data[n:]
		return n, nil
	}

	if err := r.readChunkHeader(); err != nil {
		return 0, err
	}

	if r.header.size == 0 {
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

	r.PrevSig = r.header.signature

	n := copy(buf, r.data)
	r.data = r.data[n:]

	return n, nil
}

func (r *SigV4ChunkedReader) Close() error {
	return r.Body.Close()
}

func (r *SigV4ChunkedReader) readChunkHeader() error {
	line, err := r.reader.ReadString('\n')
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

	logger.Log.Debugf("Read chunk header: %d bytes, signature %s", r.header.size, r.header.signature)

	return nil
}

func (r *SigV4ChunkedReader) readChunkData() error {
	r.data = make([]byte, r.header.size)

	if _, err := io.ReadFull(r.reader, r.data); err != nil {
		return err
	}

	logger.Log.Debugf("Read chunk data: %d bytes", len(r.data))

	return nil
}

func (r *SigV4ChunkedReader) readTrailingCRLF() error {
	crlf := make([]byte, 2)

	if _, err := io.ReadFull(r.reader, crlf); err != nil || !bytes.Equal(crlf, []byte{'\r', '\n'}) {
		return errors.New("invalid chunk termination")
	}

	logger.Log.Debug("Read chunk CRLF")

	return nil
}

func (r *SigV4ChunkedReader) verifyChunkSigV4() error {
	stringToSign := r.buildChunkStringToSign()
	recomputedSignature, err := computeSignature(r.SecretKey, r.Scope, stringToSign)

	if err != nil {
		return err
	}

	byteSignature, err := hex.DecodeString(r.header.signature)
	if err != nil {
		return errors.New("could not decode original signature")
	}

	byteRecomputedSignature, err := hex.DecodeString((recomputedSignature))
	if err != nil {
		return errors.New("could not decode recomputed signature")
	}

	logger.Log.Debugf(
		"Comparing chunk signatures: %s (original), %s (recomputed)",
		r.header.signature,
		recomputedSignature,
	)

	if hmac.Equal(byteSignature, byteRecomputedSignature) {
		return nil
	}

	return errors.New("chunk signatures differ")
}

func (r *SigV4ChunkedReader) buildChunkStringToSign() string {
	var stringToSign strings.Builder

	stringToSign.WriteString("AWS4-HMAC-SHA256-PAYLOAD")
	stringToSign.WriteString("\n")

	stringToSign.WriteString(r.Timestamp)
	stringToSign.WriteString("\n")

	stringToSign.WriteString(r.Scope)
	stringToSign.WriteString("\n")

	stringToSign.WriteString(r.PrevSig)
	stringToSign.WriteString("\n")

	emptyHash := sha256.Sum256([]byte(""))
	stringToSign.WriteString(hex.EncodeToString(emptyHash[:]))
	stringToSign.WriteString("\n")

	chunkHash := sha256.Sum256(r.data)
	stringToSign.WriteString(hex.EncodeToString(chunkHash[:]))

	return stringToSign.String()
}
