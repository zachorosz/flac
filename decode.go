package flac

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/icza/bitio"
)

var (
	ErrMissingStreamMarker = errors.New("missing fLaC marker at beginning of stream")
)

type bitReader interface {
	io.Reader
	ReadBits(n uint8) (uint64, error)
}

type Reader struct {
	r             bitReader
	err           error
	buf           []byte
	readMarker    bool
	readLastBlock bool
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		r:   bitio.NewReader(r),
		buf: make([]byte, 1024),
	}
}

func (r *Reader) Reset(reader io.Reader) {
	r.r = bitio.NewReader(reader)
	r.err = nil
	r.buf = make([]byte, 1024)
	r.readMarker = false
	r.readLastBlock = false
}

// fill reads n bytes into r.buf.
func (r *Reader) fill(n int) (ok bool) {
	if n > len(r.buf) { // expand buf size if needed
		r.buf = make([]byte, n)
	}

	return r.readFull(r.buf[:n])
}

func (r *Reader) readFull(p []byte) (ok bool) {
	if _, r.err = io.ReadFull(r.r, p); r.err != nil {
		return false
	}
	return true
}

func (r *Reader) readByte(p *byte) (ok bool) {
	if !r.readFull(r.buf[:1]) {
		return false
	}
	*p = r.buf[0]
	return true
}

func (r *Reader) skip(n int) (ok bool) {
	_, r.err = io.CopyN(io.Discard, r.r, int64(n))
	return r.err == nil
}

func (r *Reader) verifyMarker() (ok bool) {
	if !r.readFull(r.buf[:4]) {
		if errors.Is(r.err, io.EOF) {
			r.err = io.ErrUnexpectedEOF
		}
		return false
	}

	if !bytes.Equal(r.buf[:4], []byte("fLaC")) {
		r.err = fmt.Errorf("not a flac stream: %w", ErrMissingStreamMarker)
		return false
	}

	return true
}

func (r *Reader) readBlock() (*MetadataBlock, bool) {
	b := new(MetadataBlock)

	// read metadata block header: 32 bits
	if !r.readFull(r.buf[:4]) {
		return nil, false
	}

	// <1 bit> Last-metadata-block flag
	b.Last = (r.buf[0] & 0b10000000) != 0
	// <7 bits> Block type
	b.Type = MetadataBlockType(r.buf[0] & 0b1111111)
	// <24 bits> Block length in bytes (big endian encoded)
	b.Length = uint32(r.buf[3]) | uint32(r.buf[2])<<8 | uint32(r.buf[1])<<16

	switch b.Type {
	case MetadataBlockTypeStreamInfo:
		b.Data, r.err = r.decodeStreamInfo()
	case MetadataBlockTypeApplication:
		b.Data, r.err = r.decodeApplication(b.Length)
	case MetadataBlockTypeSeekTable:
		b.Data, r.err = r.decodeSeekTable(b.Length)
	case MetadataBlockTypeVorbisComment:
		b.Data, r.err = r.decodeVorbisComment()
	case MetadataBlockTypeCueSheet:
		b.Data, r.err = r.decodeCueSheet()
	case MetadataBlockTypePicture:
		b.Data, r.err = r.decodePicture()
	default:
		r.skip(int(b.Length))
	}

	if r.err != nil {
		return nil, false
	}

	return b, true
}

func (r *Reader) ReadBlock() (*MetadataBlock, error) {
	if r.err != nil {
		return nil, r.err
	}

	if !r.readMarker {
		if !r.verifyMarker() {
			return nil, r.err
		}
		r.readMarker = true
	}

	block, ok := r.readBlock()
	if !ok {
		return nil, r.err
	}

	r.readLastBlock = block.Last

	return block, nil
}
