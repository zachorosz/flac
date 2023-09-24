package flac

import (
	"io"
	"os"
	"testing"
)

var ietfMetadataExtremes = []struct {
	desc string
	path string
}{
	{"has unknown number of samples in STREAMINFO", "./flac-test-files/subset/45 - no total number of samples set.flac"},
	{"has maximum and minimum framesize set to unknown", "./flac-test-files/subset/46 - no min-max framesize set.flac"},
	{"has only a STREAMINFO block", "./flac-test-files/subset/47 - only STREAMINFO.flac"},
	{"has an extremely large SEEKTABLE", "./flac-test-files/subset/48 - Extremely large SEEKTABLE.flac"},
	{"has an extremely large PADDING block", "./flac-test-files/subset/49 - Extremely large PADDING.flac"},
	{"has an extremely large PICTURE block (JPG of 15.8MB)", "./flac-test-files/subset/50 - Extremely large PICTURE.flac"},
	{"has an extremely large VORBISCOMMENT block", "./flac-test-files/subset/51 - Extremely large VORBISCOMMENT.flac"},
	{"has an extremely large APPLICATION block", "./flac-test-files/subset/52 - Extremely large APPLICATION.flac"},
	{"has a CUESHEET block with absurdly many indexes", "./flac-test-files/subset/53 - CUESHEET with very many indexes.flac"},
	{"with the same 20 VORBISCOMMENTs repeated 1000 times", "./flac-test-files/subset/54 - 1000x repeating VORBISCOMMENT.flac"},
	{"has the metadata of track 47-52 combined", "./flac-test-files/subset/55 - file 48-53 combined.flac"},
	{"has a PICTURE with mimetype image/jpeg", "./flac-test-files/subset/56 - JPG PICTURE.flac"},
	{"has a PICTURE with mimetype image/png", "./flac-test-files/subset/57 - PNG PICTURE.flac"},
	{"has a PICTURE with mimetype image/gif", "./flac-test-files/subset/58 - GIF PICTURE.flac"},
	{"has a PICTURE with mimetype image/avif", "./flac-test-files/subset/59 - AVIF PICTURE.flac"},
}

func readBlocks(reader io.Reader) ([]*MetadataBlock, error) {
	r := NewReader(reader)

	var blocks []*MetadataBlock
	for readLast := false; !readLast; {
		b, err := r.ReadBlock()
		if err != nil {
			return nil, err
		}
		readLast = b.Last
		blocks = append(blocks, b)
	}

	return blocks, nil
}

func TestReadBlocks(t *testing.T) {
	t.Run("metadata extremes", func(t *testing.T) {
		for _, tt := range ietfMetadataExtremes {
			t.Run(tt.desc, func(t *testing.T) {
				f, err := os.Open(tt.path)
				if err != nil {
					t.Fatalf("failed to open file: %v", err)
				}
				defer f.Close()

				_, err = readBlocks(f)
				if err != nil {
					t.Error(err)
				}
			})
		}
	})
}

func BenchmarkReadBlocks_metadataextremes(b *testing.B) {
	for _, in := range ietfMetadataExtremes {
		b.Run(in.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				f, err := os.Open(in.path)
				if err != nil {
					b.Fatalf("failed to open file: %v", err)
				}
				_, err = readBlocks(f)
				if err != nil {
					b.Fatal(err)
				}

				_ = f.Close()
			}
		})
	}
}
