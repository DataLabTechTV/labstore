package s3

import (
	"bufio"
	"bytes"
	"io"
	"log/slog"
	"strconv"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/auth"
)

type SigV4ChunkEncoder struct {
	Src       io.ReadCloser
	Size      int
	ChunkSize int

	reader   *bufio.Reader
	chunk    *auth.SigV4Chunk
	chunkBuf []byte
}

func NewSigV4ChunkEncoder(src io.ReadCloser, size int, chunkSize int) *SigV4ChunkEncoder {
	return &SigV4ChunkEncoder{
		Src:       src,
		Size:      size,
		ChunkSize: chunkSize,
		reader:    bufio.NewReaderSize(src, chunkSize),
		chunk: &auth.SigV4Chunk{
			Ctx: &auth.SigV4Context{},
		},
	}
}

func (enc *SigV4ChunkEncoder) ContentLength() int {
	const crLfSize = 2
	const signatureSize = 64
	const fixedChunkHeaderSize = len(";chunk-signature=") + signatureSize

	nFirstChunks := enc.Size / enc.ChunkSize

	chunkHeaderSize := len(strconv.FormatInt(int64(enc.ChunkSize), 16)) + fixedChunkHeaderSize + 2*crLfSize
	lastChunkSize := enc.Size % enc.ChunkSize

	var lastChunkHeaderSize int
	if lastChunkSize == 0 {
		lastChunkHeaderSize = 0
	} else {
		lastChunkHeaderSize = len(strconv.FormatInt(int64(lastChunkSize), 16)) + fixedChunkHeaderSize + 2*crLfSize
	}

	emptyChunkHeaderSize := 1 + fixedChunkHeaderSize + 2*crLfSize

	totalSize := nFirstChunks*chunkHeaderSize +
		lastChunkHeaderSize +
		emptyChunkHeaderSize +
		enc.Size

	slog.Debug(
		"sigv4 chunk encoder",
		"enc.Size", enc.Size,
		"enc.chunkSize", enc.ChunkSize,
		"lastChunkSize", lastChunkSize,
		"totalSize", totalSize,
	)

	return totalSize
}

func (enc *SigV4ChunkEncoder) Read(buf []byte) (int, error) {
	if len(enc.chunkBuf) > 0 {
		n := copy(buf, enc.chunkBuf)
		enc.chunkBuf = enc.chunkBuf[n:]
		return n, nil
	}

	enc.chunk.Data = make([]byte, enc.ChunkSize)
	n, err := enc.reader.Read(enc.chunk.Data)
	if err != nil && err != io.EOF {
		return 0, err
	}

	chunkSignature, err := auth.ComputeSignature(enc.chunk.Ctx.Credential, enc.chunk.BuildChunkStringToSign())
	if err != nil {
		return 0, err
	}

	enc.chunk.Header = &auth.SigV4ChunkHeader{
		Size:      n,
		Signature: chunkSignature,
	}

	var chunkBuf bytes.Buffer
	chunkBuf.WriteString(strconv.FormatInt(int64(enc.chunk.Header.Size), 16))
	chunkBuf.WriteString(";chunk-signature=")
	chunkBuf.WriteString(enc.chunk.Header.Signature)
	chunkBuf.WriteString("\r\n")
	chunkBuf.Write(enc.chunk.Data)
	chunkBuf.WriteString("\r\n")

	enc.chunkBuf = chunkBuf.Bytes()
	enc.chunk.Ctx.Signature = enc.chunk.Header.Signature

	n = copy(buf, enc.chunkBuf)
	enc.chunkBuf = enc.chunkBuf[n:]

	return n, nil
}

func (enc *SigV4ChunkEncoder) Close() error {
	return enc.Src.Close()
}
