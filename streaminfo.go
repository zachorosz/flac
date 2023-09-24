package flac

// StreamInfo represents stream info metadata block data.
//
// https://xiph.org/flac/format.html#metadata_block_streaminfo
type StreamInfo struct {
	// The minimum block size (in samples) used in the stream.
	MinimumBlockSize uint16
	// The maximum block size (in samples) used in the stream.
	MaximumBlockSize uint16
	// The minimum frame size (in bytes) used in the stream. May be 0 to imply
	// that the value is not known.
	MinimumFrameSize uint32
	// The maximum frame size (in bytes) used in the stream. May be 0 to imply
	// that the value is not known.
	MaximumFrameSize uint32
	// Sample Rate in Hz. A value of 0 is invalid.
	SampleRate uint32
	// Number of channels. FLAC supports 1-8 channels.
	Channels uint8
	// Bits per sample. FLAC supports 4-32 bits per sample.
	BitsPerSample uint8
	// Total samples in stream. May be 0 to imply that the value is unknown.
	TotalSamples uint64
	// MD5 signature of the unencoded audio data.
	MD5 []byte
}

func (r *Reader) decodeStreamInfo() (*StreamInfo, error) {
	streamInfo := new(StreamInfo)

	minBlockSize, err := r.r.ReadBits(16)
	if err != nil {
		return nil, err
	}
	streamInfo.MinimumBlockSize = uint16(minBlockSize)

	maxBlockSize, err := r.r.ReadBits(16)
	if err != nil {
		return nil, err
	}
	streamInfo.MaximumBlockSize = uint16(maxBlockSize)

	minFrameSize, err := r.r.ReadBits(24)
	if err != nil {
		return nil, err
	}
	streamInfo.MinimumFrameSize = uint32(minFrameSize)

	maxFrameSize, err := r.r.ReadBits(24)
	if err != nil {
		return nil, err
	}
	streamInfo.MaximumFrameSize = uint32(maxFrameSize)

	sampleRate, err := r.r.ReadBits(20)
	if err != nil {
		return nil, err
	}
	streamInfo.SampleRate = uint32(sampleRate)

	channels, err := r.r.ReadBits(3)
	if err != nil {
		return nil, err
	}
	streamInfo.Channels = uint8(channels) + 1 // FLAC contains (number of channels)-1

	bps, err := r.r.ReadBits(5)
	if err != nil {
		return nil, err
	}
	streamInfo.BitsPerSample = uint8(bps) + 1 // FLAC contains (bits per sample)-1

	totalSamples, err := r.r.ReadBits(36)
	if err != nil {
		return nil, err
	}
	streamInfo.TotalSamples = uint64(totalSamples)

	if !r.readFull(r.buf[:16]) {
		return nil, r.err
	}
	streamInfo.MD5 = r.buf[:16]

	return streamInfo, nil
}
