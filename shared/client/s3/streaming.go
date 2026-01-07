package s3

import (
	"bufio"
	"io"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/auth"
)

type SigV4ChunkedEncoder struct {
	Ctx *auth.SigV4Context
	Src io.ReadCloser

	reader *bufio.Reader
	header *auth.SigV4ChunkHeader
	data   []byte
}

func NewSigV4ChunkedEncoder(src io.ReadCloser, ctx *auth.SigV4Context) *SigV4ChunkedEncoder {
	return &SigV4ChunkedEncoder{Ctx: ctx, Src: src}
}

func (r *SigV4ChunkedEncoder) Read(buf []byte) (int, error) {
	// TODO
	return -1, nil
}

func (r *SigV4ChunkedEncoder) Close() error {
	return r.Src.Close()
}
