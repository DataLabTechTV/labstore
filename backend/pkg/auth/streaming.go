package auth

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"strconv"
	"strings"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/security"
)

type sigV4ChunkedDecoder struct {
	ctx *SigV4ChunkContext
	src io.ReadCloser

	reader *bufio.Reader
	header *SigV4ChunkHeader
	data   []byte
}

type SigV4ChunkContext struct {
	PrevSig    string
	Credential *SigV4Credential
	Timestamp  string
}

type SigV4ChunkHeader struct {
	Size      int
	Signature string
}

func NewSigV4ChunkedDecoder(src io.ReadCloser, res *SigV4Result) *sigV4ChunkedDecoder {
	return &sigV4ChunkedDecoder{
		ctx: &SigV4ChunkContext{
			PrevSig:    res.Signature,
			Credential: res.Credential,
			Timestamp:  res.Timestamp,
		},
		src: src,
	}
}

func (r *sigV4ChunkedDecoder) Read(buf []byte) (int, error) {
	if r.reader == nil {
		r.reader = bufio.NewReader(r.src)
	}

	if len(r.data) > 0 {
		n := copy(buf, r.data)
		r.data = r.data[n:]
		return n, nil
	}

	if err := r.readChunkHeader(); err != nil {
		return 0, err
	}

	if r.header.Size == 0 {
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

	r.ctx.PrevSig = r.header.Signature

	n := copy(buf, r.data)
	r.data = r.data[n:]

	return n, nil
}

func (r *sigV4ChunkedDecoder) Close() error {
	return r.src.Close()
}

func (r *sigV4ChunkedDecoder) readChunkHeader() error {
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
		Size:      int(size),
		Signature: sig,
	}

	slog.Debug("chunk header", "size", r.header.Size, "signature", security.Trunc(r.header.Signature))

	return nil
}

func (r *sigV4ChunkedDecoder) readChunkData() error {
	r.data = make([]byte, r.header.Size)

	if _, err := io.ReadFull(r.reader, r.data); err != nil {
		return err
	}

	slog.Debug("chunk data", "length", len(r.data))

	return nil
}

func (r *sigV4ChunkedDecoder) readTrailingCRLF() error {
	crlf := make([]byte, 2)

	if _, err := io.ReadFull(r.reader, crlf); err != nil || !bytes.Equal(crlf, []byte{'\r', '\n'}) {
		return errors.New("invalid chunk termination")
	}

	slog.Debug("chunk crlf")

	return nil
}

func (r *sigV4ChunkedDecoder) verifyChunkSigV4() error {
	stringToSign := r.buildChunkStringToSign()
	slog.Debug("string to sign", "string_to_sign", security.TruncLastLines(stringToSign, 3))

	recomputedSignature, err := ComputeSignature(r.ctx.Credential, stringToSign)

	if err != nil {
		return err
	}

	byteSignature, err := hex.DecodeString(r.header.Signature)
	if err != nil {
		return errors.New("could not decode original signature")
	}

	byteRecomputedSignature, err := hex.DecodeString((recomputedSignature))
	if err != nil {
		return errors.New("could not decode recomputed signature")
	}

	slog.Debug(
		"comparing chunk signatures",
		"original", security.Trunc(r.header.Signature),
		"recomputed", security.Trunc(recomputedSignature),
	)

	if hmac.Equal(byteSignature, byteRecomputedSignature) {
		return nil
	}

	return errors.New("chunk signatures differ")
}

func (r *sigV4ChunkedDecoder) buildChunkStringToSign() string {
	var stringToSign strings.Builder

	stringToSign.WriteString("AWS4-HMAC-SHA256-PAYLOAD")
	stringToSign.WriteString("\n")

	stringToSign.WriteString(r.ctx.Timestamp)
	stringToSign.WriteString("\n")

	stringToSign.WriteString(r.ctx.Credential.Scope)
	stringToSign.WriteString("\n")

	stringToSign.WriteString(r.ctx.PrevSig)
	stringToSign.WriteString("\n")

	emptyHash := sha256.Sum256([]byte(""))
	stringToSign.WriteString(hex.EncodeToString(emptyHash[:]))
	stringToSign.WriteString("\n")

	chunkHash := sha256.Sum256(r.data)
	stringToSign.WriteString(hex.EncodeToString(chunkHash[:]))

	return stringToSign.String()
}
