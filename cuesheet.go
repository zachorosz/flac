package flac

import "encoding/binary"

// CueSheet represents a cue sheet metadata block data. This block is for
// storing information that can be used in a cue sheet like track and index
// points.
//
// https://xiph.org/flac/format.html#metadata_block_cuesheet
type CueSheet struct {
	CatalogNumber    string
	NumLeadInSamples uint64
	IsCD             bool
	Tracks           []*CueSheetTrack
}

// CueSheetTrack represents a track point in a cue sheet metadata block.
//
// https://xiph.org/flac/format.html#cuesheet_track
type CueSheetTrack struct {
	OffsetSamples uint64 // track offset in samples, relative to the beginning of the FLAC audio stream.
	TrackNumber   uint8
	ISRC          string
	IsAudio       bool
	PreEmphasis   bool
	Indices       []*CueSheetTrackIndex // track index points, except the lead-out track
}

// CueSheetTrackIndex represents a track index point.
//
// https://xiph.org/flac/format.html#cuesheet_track_index
type CueSheetTrackIndex struct {
	OffsetSamples uint64 // offset in samples, relative to the track offset
	PointNumber   uint8
}

func (r *Reader) decodeCueSheet() (*CueSheet, error) {
	cueSheet := new(CueSheet)

	if !r.readFull(r.buf[:128]) {
		return nil, r.err
	}
	cueSheet.CatalogNumber = string(r.buf[:128])

	if err := binary.Read(r.r, binary.BigEndian, &cueSheet.NumLeadInSamples); err != nil {
		return nil, err
	}

	flags, ok := r.nextByte()
	if !ok {
		return nil, r.err
	}

	cueSheet.IsCD = (flags & 0x80) != 0

	if !r.skip(258) { // 258 reserved bytes
		return nil, r.err
	}

	var numTracks uint8
	if err := binary.Read(r.r, binary.BigEndian, &numTracks); err != nil {
		return nil, err
	}

	cueSheet.Tracks = make([]*CueSheetTrack, numTracks)
	for n := range cueSheet.Tracks {
		track, err := r.decodeCueSheetTrack()
		if err != nil {
			return nil, err
		}
		cueSheet.Tracks[n] = track
	}

	return cueSheet, nil
}

func (r *Reader) decodeCueSheetTrack() (*CueSheetTrack, error) {
	track := new(CueSheetTrack)

	if err := binary.Read(r.r, binary.BigEndian, &track.OffsetSamples); err != nil {
		return nil, err
	}

	if err := binary.Read(r.r, binary.BigEndian, &track.TrackNumber); err != nil {
		return nil, err
	}

	if !r.readFull(r.buf[:12]) {
		return nil, r.err
	}
	track.ISRC = string(r.buf[:12])

	flags, ok := r.nextByte()
	if !ok {
		return nil, r.err
	}

	track.IsAudio = (flags & 0x80) != 0

	track.PreEmphasis = (flags & 0x40) != 0

	if !r.skip(13) { // 13 reserved bytes
		return nil, r.err
	}

	var numIndices uint8
	if err := binary.Read(r.r, binary.BigEndian, &numIndices); err != nil {
		return nil, err
	}

	track.Indices = make([]*CueSheetTrackIndex, numIndices)
	for m := range track.Indices {
		index, err := r.decodeCueSheetTrackIndex()
		if err != nil {
			return nil, err
		}
		track.Indices[m] = index
	}

	return track, nil
}

func (r *Reader) decodeCueSheetTrackIndex() (*CueSheetTrackIndex, error) {
	index := new(CueSheetTrackIndex)

	if err := binary.Read(r.r, binary.BigEndian, &index.OffsetSamples); err != nil {
		return nil, err
	}

	if err := binary.Read(r.r, binary.BigEndian, &index.PointNumber); err != nil {
		return nil, err
	}

	if !r.skip(3) { // 3 reserved bytes
		return nil, r.err
	}

	return index, nil
}
