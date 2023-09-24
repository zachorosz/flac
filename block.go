package flac

type MetadataBlockType uint8

// Metadata block types
const (
	MetadataBlockTypeStreamInfo MetadataBlockType = iota
	MetadataBlockTypePadding
	MetadataBlockTypeApplication
	MetadataBlockTypeSeekTable
	MetadataBlockTypeVorbisComment
	MetadataBlockTypeCueSheet
	MetadataBlockTypePicture
)

func (t MetadataBlockType) Reserved() bool {
	return t >= 7 && t.Valid()
}

func (t MetadataBlockType) Valid() bool {
	return t < 127
}

func (t MetadataBlockType) String() string {
	switch t {
	case MetadataBlockTypeStreamInfo:
		return "STREAMINFO"
	case MetadataBlockTypePadding:
		return "PADDING"
	case MetadataBlockTypeApplication:
		return "APPLICATION"
	case MetadataBlockTypeSeekTable:
		return "SEEKTABLE"
	case MetadataBlockTypeVorbisComment:
		return "VORBIS_COMMENT"
	case MetadataBlockTypeCueSheet:
		return "CUESHEET"
	case MetadataBlockTypePicture:
		return "PICTURE"
	}
	if t.Reserved() {
		return "RESERVED"
	}
	return "INVALID"
}

// MetadataBlock represents a metadata block in a FLAC stream.
//
// https://xiph.org/flac/format.html#metadata_block
type MetadataBlock struct {
	MetadataBlockHeader
	// Data points to a StreamInfo, Application, SeekTable, VorbisComment,
	// CueSheet, or Picture respective of the block's Type or is nil if the
	// block is a padding block.
	Data interface{}
}

// MetadataBlockHeader represents a metadata block header.
//
// https://xiph.org/flac/format.html#metadata_block_header
type MetadataBlockHeader struct {
	Last   bool
	Type   MetadataBlockType
	Length uint32
}
