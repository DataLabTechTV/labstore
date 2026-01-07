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

type SigV4Chunk struct {
	Ctx    *SigV4Context
	Header *SigV4ChunkHeader
	Data   []byte
}

type SigV4ChunkHeader struct {
	Size      int
	Signature string
}

type SigV4ChunkDecoder struct {
	Src io.ReadCloser

	reader *bufio.Reader
	chunk  *SigV4Chunk
}

func NewSigV4ChunkedDecoder(ctx *SigV4Context, src io.ReadCloser) *SigV4ChunkDecoder {
	return &SigV4ChunkDecoder{
		Src:    src,
		reader: bufio.NewReader(src),
		chunk:  &SigV4Chunk{Ctx: ctx},
	}
}

func (chunk *SigV4Chunk) BuildChunkStringToSign() string {
	var stringToSign strings.Builder

	stringToSign.WriteString("AWS4-HMAC-SHA256-PAYLOAD")
	stringToSign.WriteString("\n")

	stringToSign.WriteString(chunk.Ctx.Timestamp)
	stringToSign.WriteString("\n")

	stringToSign.WriteString(chunk.Ctx.Credential.Scope)
	stringToSign.WriteString("\n")

	stringToSign.WriteString(chunk.Ctx.Signature)
	stringToSign.WriteString("\n")

	emptyHash := sha256.Sum256([]byte(""))
	stringToSign.WriteString(hex.EncodeToString(emptyHash[:]))
	stringToSign.WriteString("\n")

	chunkHash := sha256.Sum256(chunk.Data)
	stringToSign.WriteString(hex.EncodeToString(chunkHash[:]))

	return stringToSign.String()
}

func (dec *SigV4ChunkDecoder) Read(buf []byte) (int, error) {
	if len(dec.chunk.Data) > 0 {
		n := copy(buf, dec.chunk.Data)
		dec.chunk.Data = dec.chunk.Data[n:]
		return n, nil
	}

	if err := dec.readChunkHeader(); err != nil {
		return 0, err
	}

	if dec.chunk.Header.Size == 0 {
		return 0, io.EOF
	}

	if err := dec.readChunkData(); err != nil {
		return 0, err
	}

	if err := dec.readTrailingCRLF(); err != nil {
		return 0, err
	}

	if err := dec.verifyChunkSigV4(); err != nil {
		return 0, err
	}

	dec.chunk.Ctx.Signature = dec.chunk.Header.Signature

	n := copy(buf, dec.chunk.Data)
	dec.chunk.Data = dec.chunk.Data[n:]

	return n, nil
}

func (dec *SigV4ChunkDecoder) Close() error {
	return dec.Src.Close()
}

func (dec *SigV4ChunkDecoder) readChunkHeader() error {
	line, err := dec.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return err
	}

	line = strings.TrimSuffix(line, "\r\n")

	headerParts := strings.SplitN(line, ";", 2)
	if len(headerParts) != 2 {
		return errors.New("invalid chunk header")
	}
	sizeHex, chunkSig := headerParts[0], headerParts[1]

	size, err := strconv.ParseInt(sizeHex, 16, 64)
	if err != nil {
		return err
	}

	sig, ok := strings.CutPrefix(chunkSig, "chunk-signature=")
	if !ok {
		return errors.New("could not find 'chunk-signature=' prefix")
	}

	dec.chunk.Header = &SigV4ChunkHeader{
		Size:      int(size),
		Signature: sig,
	}

	slog.Debug(
		"chunk header",
		"size", dec.chunk.Header.Size,
		"signature", security.Trunc(dec.chunk.Header.Signature),
	)

	return nil
}

func (dec *SigV4ChunkDecoder) readChunkData() error {
	dec.chunk.Data = make([]byte, dec.chunk.Header.Size)

	if _, err := io.ReadFull(dec.reader, dec.chunk.Data); err != nil {
		return err
	}

	slog.Debug("chunk data", "length", len(dec.chunk.Data))

	return nil
}

func (dec *SigV4ChunkDecoder) readTrailingCRLF() error {
	crlf := make([]byte, 2)

	if _, err := io.ReadFull(dec.reader, crlf); err != nil || !bytes.Equal(crlf, []byte{'\r', '\n'}) {
		return errors.New("invalid chunk termination")
	}

	slog.Debug("chunk crlf")

	return nil
}

func (dec *SigV4ChunkDecoder) verifyChunkSigV4() error {
	stringToSign := dec.chunk.BuildChunkStringToSign()
	slog.Debug("string to sign", "string_to_sign", security.TruncLastLines(stringToSign, 3))

	recomputedSignature, err := ComputeSignature(dec.chunk.Ctx.Credential, stringToSign)

	if err != nil {
		return err
	}

	byteSignature, err := hex.DecodeString(dec.chunk.Header.Signature)
	if err != nil {
		return errors.New("could not decode original signature")
	}

	byteRecomputedSignature, err := hex.DecodeString((recomputedSignature))
	if err != nil {
		return errors.New("could not decode recomputed signature")
	}

	slog.Debug(
		"comparing chunk signatures",
		"original", security.Trunc(dec.chunk.Header.Signature),
		"recomputed", security.Trunc(recomputedSignature),
	)

	if hmac.Equal(byteSignature, byteRecomputedSignature) {
		return nil
	}

	return errors.New("chunk signatures differ")
}
