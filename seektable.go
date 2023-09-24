package flac

import (
	"encoding/binary"
)

// SeekTable represents seek table metadata block data.
//
// https://xiph.org/flac/format.html#metadata_block_seektable
type SeekTable struct {
	SeekPoints []*SeekPoint
}

// SeekPoint represents a seek point in a seek table.
//
// https://xiph.org/flac/format.html#seekpoint
type SeekPoint struct {
	// sample number of the first sample in target frame, or 0xFFFFFFFFFFFFFFFF
	// for a placeholder point.
	SampleNumber uint64
	// offset in bytes from the first byte of the first frame header to the
	// first byte of the target frame's header
	Offset uint64
	// the number of samples in the target frame
	NumSamples uint16
}

// IsPlaceholder returns true if the seek point is a placeholder.
func (sp *SeekPoint) IsPlaceholder() bool {
	return sp.SampleNumber == 0xFFFFFFFFFFFFFFFF
}

func (r *Reader) decodeSeekTable(blockLength uint32) (*SeekTable, error) {
	// The number of seek points is implied by the metadata header 'length'
	// field, i.e. equal to length / 18.
	nSeekPoints := blockLength / 18

	b := new(SeekTable)
	b.SeekPoints = make([]*SeekPoint, 0, nSeekPoints)

	for i := 0; i < int(nSeekPoints); i++ {
		seekPoint, err := r.decodeSeekPoint()
		if err != nil {
			return nil, err
		}
		b.SeekPoints = append(b.SeekPoints, seekPoint)
	}

	return b, nil
}

func (r *Reader) decodeSeekPoint() (*SeekPoint, error) {
	seekPoint := new(SeekPoint)

	if err := binary.Read(r.r, binary.BigEndian, &seekPoint.SampleNumber); err != nil {
		return nil, err
	}

	if err := binary.Read(r.r, binary.BigEndian, &seekPoint.Offset); err != nil {
		return nil, err
	}

	if err := binary.Read(r.r, binary.BigEndian, &seekPoint.NumSamples); err != nil {
		return nil, err
	}

	return seekPoint, nil
}
