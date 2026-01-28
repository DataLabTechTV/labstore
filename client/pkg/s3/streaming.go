package s3

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log/slog"
	"strconv"

	client "github.com/IllumiKnowLabs/labstore/client/pkg/types"
	"github.com/IllumiKnowLabs/labstore/server/pkg/auth"
)

type SigV4ChunkEncoder struct {
	Ctx       context.Context
	Src       io.ReadCloser
	DataSize  int
	TotalSize int
	ChunkSize int

	progress chan<- client.Progress
	reader   *bufio.Reader
	chunk    *auth.SigV4Chunk
	chunkBuf []byte
	read     int
	done     bool
}

func NewSigV4ChunkEncoder(
	ctx context.Context,
	src io.ReadCloser,
	dataSize int,
	chunkSize int,
	progress chan<- client.Progress,
) *SigV4ChunkEncoder {
	enc := &SigV4ChunkEncoder{
		Ctx:       ctx,
		Src:       src,
		DataSize:  dataSize,
		ChunkSize: chunkSize,

		progress: progress,
		reader:   bufio.NewReaderSize(src, chunkSize),

		chunk: &auth.SigV4Chunk{
			Ctx: &auth.SigV4Context{},
		},
	}

	enc.TotalSize = enc.calculateTotalSize()

	return enc
}

func (enc *SigV4ChunkEncoder) Read(buf []byte) (int, error) {
	if len(enc.chunkBuf) > 0 {
		n := copy(buf, enc.chunkBuf)
		enc.chunkBuf = enc.chunkBuf[n:]

		enc.read += n
		slog.Debug(
			"sigv4 chunk encoder",
			"enc.read", enc.read,
			"enc.TotalSize", enc.TotalSize,
			"pct", float64(enc.read)/float64(enc.TotalSize)*100,
		)

		if enc.progress != nil {
			msg := client.Progress{Current: enc.read, Total: enc.TotalSize}

			select {
			case <-enc.Ctx.Done():
				return 0, io.ErrClosedPipe
			case enc.progress <- msg:
			}
		}

		return n, nil
	}

	if enc.done {
		return 0, io.EOF
	}

	enc.chunk.Data = make([]byte, enc.ChunkSize)
	n, err := enc.reader.Read(enc.chunk.Data)
	if err != nil && err != io.EOF {
		return 0, err
	}

	if err == io.EOF {
		enc.done = true
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
	chunkBuf.Write(enc.chunk.Data[:n])
	chunkBuf.WriteString("\r\n")

	enc.chunkBuf = chunkBuf.Bytes()
	enc.chunk.Ctx.Signature = enc.chunk.Header.Signature

	n = copy(buf, enc.chunkBuf)
	enc.chunkBuf = enc.chunkBuf[n:]

	enc.read += n
	slog.Debug(
		"sigv4 chunk encoder",
		"enc.read", enc.read,
		"enc.TotalSize", enc.TotalSize,
		"pct", float64(enc.read)/float64(enc.TotalSize)*100,
	)

	if enc.progress != nil {
		msg := client.Progress{Current: enc.read, Total: enc.TotalSize}

		select {
		case <-enc.Ctx.Done():
			return 0, io.ErrClosedPipe
		case enc.progress <- msg:
		}
	}

	return n, nil
}

func (enc *SigV4ChunkEncoder) Close() error {
	return enc.Src.Close()
}

func (enc *SigV4ChunkEncoder) calculateTotalSize() int {
	const crLfSize = 2
	const signatureSize = 64
	const fixedChunkHeaderSize = len(";chunk-signature=") + signatureSize

	nFirstChunks := enc.DataSize / enc.ChunkSize

	chunkHeaderSize := len(strconv.FormatInt(int64(enc.ChunkSize), 16)) + fixedChunkHeaderSize + 2*crLfSize

	lastChunkSize := enc.DataSize % enc.ChunkSize

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
		enc.DataSize

	slog.Debug(
		"sigv4 chunk encoder",
		"enc.Size", enc.DataSize,
		"enc.chunkSize", enc.ChunkSize,
		"lastChunkSize", lastChunkSize,
		"totalSize", totalSize,
	)

	return totalSize
}
