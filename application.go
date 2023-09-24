package flac

// Application represents an application metadata block. This block is for use
// by third-party applications.
//
// https://xiph.org/flac/format.html#metadata_block_application
type Application struct {
	ID   string // Registered application ID
	Data []byte // Application data
}

func (r *Reader) decodeApplication(n uint32) (*Application, error) {
	b := new(Application)

	if !r.fill(int(n)) {
		return nil, r.err
	}

	i := 0
	b.ID = string(r.buf[i:4])
	i += 4
	b.Data = r.buf[i:n]

	return b, nil
}
