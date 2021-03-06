package chunkreader

import (
	"io"
)

type ChunkReader struct {
	r io.Reader

	buf    []byte
	rp, wp int // buf read position and write position

	options Options
}

type Options struct {
	MinBufLen int // Minimum buffer length
}

func NewChunkReader(r io.Reader) *ChunkReader {
	cr, err := NewChunkReaderEx(r, Options{})
	if err != nil {
		panic("default options can't be bad")
	}

	return cr
}

func NewChunkReaderEx(r io.Reader, options Options) (*ChunkReader, error) {
	if options.MinBufLen == 0 {
		// By historical reasons Postgres currently has 8KB send buffer inside,
		// so here we want to have at least the same size buffer.
		// @see https://github.com/postgres/postgres/blob/249d64999615802752940e017ee5166e726bc7cd/src/backend/libpq/pqcomm.c#L134
		// @see https://www.postgresql.org/message-id/0cdc5485-cb3c-5e16-4a46-e3b2f7a41322%40ya.ru
		options.MinBufLen = 8192
	}

	return &ChunkReader{
		r:       r,
		buf:     make([]byte, options.MinBufLen),
		options: options,
	}, nil
}

// Next returns buf filled with the next n bytes. If an error occurs, buf will
// be nil.
func (r *ChunkReader) Next(n int) (buf []byte, err error) {
	// n bytes already in buf
	if (r.wp - r.rp) >= n {
		buf = r.buf[r.rp : r.rp+n]
		r.rp += n
		return buf, err
	}

	// available space in buf is less than n
	if len(r.buf) < n {
		r.copyBufContents(r.newBuf(n))
	}

	// buf is large enough, but need to shift filled area to start to make enough contiguous space
	minReadCount := n - (r.wp - r.rp)
	if (len(r.buf) - r.wp) < minReadCount {
		newBuf := r.newBuf(n)
		r.copyBufContents(newBuf)
	}

	if err := r.appendAtLeast(minReadCount); err != nil {
		return nil, err
	}

	buf = r.buf[r.rp : r.rp+n]
	r.rp += n
	return buf, nil
}

func (r *ChunkReader) appendAtLeast(fillLen int) error {
	n, err := io.ReadAtLeast(r.r, r.buf[r.wp:], fillLen)
	r.wp += n
	return err
}

func (r *ChunkReader) newBuf(size int) []byte {
	if size < r.options.MinBufLen {
		size = r.options.MinBufLen
	}
	return make([]byte, size)
}

func (r *ChunkReader) copyBufContents(dest []byte) {
	r.wp = copy(dest, r.buf[r.rp:r.wp])
	r.rp = 0
	r.buf = dest
}
