package flac

import (
	"encoding/binary"
)

type VorbisComment struct {
	Vendor       string
	UserComments []string
}

func (r *Reader) decodeVorbisComment() (*VorbisComment, error) {
	vc := new(VorbisComment)

	var length uint32
	if err := binary.Read(r.r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	if !r.fill(int(length)) {
		return nil, r.err
	}

	vc.Vendor = string(r.buf[:length])

	if err := binary.Read(r.r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	vc.UserComments = make([]string, length)
	for i := 0; i < len(vc.UserComments); i++ {
		if err := binary.Read(r.r, binary.LittleEndian, &length); err != nil {
			return nil, err
		}

		if !r.fill(int(length)) {
			return nil, r.err
		}
		vc.UserComments[i] = string(r.buf[:length])
	}

	return vc, nil
}
