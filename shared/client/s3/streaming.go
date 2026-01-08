package s3

import (
	"bufio"
	"bytes"
	"io"
	"math"
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
	nChunks := int(math.Ceil(float64(enc.Size) / float64(enc.ChunkSize)))

	lastChunkSize := enc.Size % enc.ChunkSize
	if lastChunkSize == 0 {
		lastChunkSize = enc.ChunkSize
	}

	totalSize := enc.Size + (nChunks-1)*enc.ChunkSize + lastChunkSize

	return totalSize
}

func (r *SigV4ChunkEncoder) Read(buf []byte) (int, error) {
	if len(r.chunk.Data) > 0 {
		n := copy(buf, r.chunk.Data)
		r.chunk.Data = r.chunk.Data[n:]
		return n, nil
	}

	r.chunk.Data = make([]byte, r.ChunkSize)
	n, err := r.reader.Read(r.chunk.Data)
	if err != nil && err != io.EOF {
		return 0, err
	}

	chunkSignature, err := auth.ComputeSignature(r.chunk.Ctx.Credential, r.chunk.BuildChunkStringToSign())
	if err != nil {
		return 0, err
	}

	r.chunk.Header = &auth.SigV4ChunkHeader{
		Size:      n,
		Signature: chunkSignature,
	}

	var chunkBuf bytes.Buffer
	chunkBuf.WriteString(strconv.FormatInt(int64(r.chunk.Header.Size), 16))
	chunkBuf.WriteString(";chunk-signature=")
	chunkBuf.WriteString(r.chunk.Header.Signature)
	chunkBuf.WriteString("\r\n")
	chunkBuf.Write(r.chunk.Data)
	chunkBuf.WriteString("\r\n")

	r.chunkBuf = chunkBuf.Bytes()
	r.chunk.Ctx.Signature = r.chunk.Header.Signature

	n = copy(buf, r.chunk.Data)
	r.chunk.Data = r.chunk.Data[n:]

	return n, nil
}

func (r *SigV4ChunkEncoder) Close() error {
	return r.Src.Close()
}
