package flac

import "encoding/binary"

type PictureType uint32

// Picture types according to the ID3v2 APIC frame.
// Others are reserved and should not be used.
const (
	PictureTypeOther PictureType = iota
	PictureTypeFileIcon
	PictureTypeOtherFileIcon
	PictureTypeCoverFront
	PictureTypeCoverBack
	PictureTypeLeafletPage
	PictureTypeMedia
	PictureTypeLeadArtist
	PictureTypeArtistPerformer
	PictureTypeConductor
	PictureTypeBandOrchestra
	PictureTypeComposer
	PictureTypeLyricistTextWriter
	PictureTypeRecordingLocation
	PictureTypeDuringRecording
	PictureTypeDuringPerformance
	PictureTypeMovieVideoScreenCapture
	PictureTypeBrightColoredFish
	PictureTypeIllustration
	PictureTypeBandArtistLogotype
	PictureTypePublisherStudioLogotype
)

func (t PictureType) String() string {
	switch t {
	case PictureTypeOther:
		return "Other"
	case PictureTypeFileIcon:
		return "32x32 pixels file icon"
	case PictureTypeOtherFileIcon:
		return "Other file icon"
	case PictureTypeCoverFront:
		return "Cover (front)"
	case PictureTypeCoverBack:
		return "Cover (back)"
	case PictureTypeLeafletPage:
		return "Leaflet page"
	case PictureTypeMedia:
		return "Media"
	case PictureTypeLeadArtist:
		return "Lead artist/lead performer/soloist"
	case PictureTypeArtistPerformer:
		return "Artist/performer"
	case PictureTypeConductor:
		return "Conductor"
	case PictureTypeBandOrchestra:
		return "Band/Orchestra"
	case PictureTypeComposer:
		return "Composer"
	case PictureTypeLyricistTextWriter:
		return "Lyricist/text writer"
	case PictureTypeRecordingLocation:
		return "Recording location"
	case PictureTypeDuringRecording:
		return "During recording"
	case PictureTypeDuringPerformance:
		return "During performance"
	case PictureTypeMovieVideoScreenCapture:
		return "Movie/video screen capture"
	case PictureTypeBrightColoredFish:
		return "A bright colored fish"
	case PictureTypeIllustration:
		return "Illustration"
	case PictureTypeBandArtistLogotype:
		return "Band/artist logotype"
	case PictureTypePublisherStudioLogotype:
		return "Publisher/Studio logotype"
	default:
		return "Invalid (reserved)"
	}
}

// Picture represents a picture metadata block. This block is for storing
// pictures associated with the file.
//
// https://xiph.org/flac/format.html#metadata_block_picture
type Picture struct {
	Type PictureType
	// The MIME type string, in printable ASCII characters 0x20-0x7e. The MIME
	// type may also be --> to signify that the data part is a URL of the
	// picture instead of the picture data itself.
	MimeType string
	// The description of the picture in UTF-8.
	Description string
	// The width of the picture in pixels.
	Width uint32
	// The height of the picture in pixels.
	Height uint32
	// The color depth of the picture in bits-per-pixel.
	Depth uint32
	// For indexed-color pictures (e.g. GIF), the number of colors used, or 0
	// for non-indexed pictures.
	Colors uint32
	// The binary picture data.
	Data []byte
}

func (r *Reader) decodePicture() (*Picture, error) {
	picture := new(Picture)

	if err := binary.Read(r.r, binary.BigEndian, &picture.Type); err != nil {
		return nil, err
	}

	var length uint32

	if err := binary.Read(r.r, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	if !r.fill(int(length)) {
		return nil, r.err
	}
	picture.MimeType = string(r.buf[:length])

	if err := binary.Read(r.r, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	if !r.fill(int(length)) {
		return nil, r.err
	}
	picture.Description = string(r.buf[:length])

	if err := binary.Read(r.r, binary.BigEndian, &picture.Width); err != nil {
		return nil, err
	}

	if err := binary.Read(r.r, binary.BigEndian, &picture.Height); err != nil {
		return nil, err
	}

	if err := binary.Read(r.r, binary.BigEndian, &picture.Depth); err != nil {
		return nil, err
	}

	if err := binary.Read(r.r, binary.BigEndian, &picture.Colors); err != nil {
		return nil, err
	}

	if err := binary.Read(r.r, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	picture.Data = make([]byte, length)
	if !r.readFull(picture.Data) {
		return nil, r.err
	}

	return picture, nil
}
