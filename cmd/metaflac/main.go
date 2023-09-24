package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/zachorosz/flac"
)

func help(w io.Writer) {
	fmt.Fprintln(w, `Usage:
    metaflac FLACfile [FLACfile ...]

List metadata in one or more FLAC files.`)
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func main() {
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if len(os.Args) < 2 {
		help(out)
		os.Exit(1)
	}
	flacFiles := os.Args[1:]

	list(out, flacFiles)
}

func list(w io.Writer, files []string) {
	for _, f := range files {
		var prefix string
		if len(files) > 1 {
			prefix = fmt.Sprintf("%s:", f)
		}
		listMetadata(w, prefix, f)
	}
}

func listMetadata(w io.Writer, prefix, file string) {
	f, err := os.Open(file)
	if err != nil {
		fatalf("%s: failed to open FLAC file", file)
		os.Exit(1)
	}
	defer f.Close()

	r := flac.NewReader(f)

	i := 0
	for readLast := false; !readLast; {
		b, err := r.ReadBlock()
		if err != nil {
			fatalf("%s: failed to read block: %v", file, err)
		}
		readLast = b.Last

		fmt.Fprintf(w, "%sMETADATA block #%d\n", prefix, i+1)
		fmt.Fprintf(w, "%s\ttype: %d (%s)\n", prefix, b.Type, b.Type)
		fmt.Fprintf(w, "%s\tis last: %t\n", prefix, b.Last)
		fmt.Fprintf(w, "%s\tlength: %d\n", prefix, b.Length)

		switch b := b.Data.(type) {
		case *flac.StreamInfo:
			fmt.Fprintf(w, "%s\tminimum block size: %d samples\n", prefix, b.MinimumBlockSize)
			fmt.Fprintf(w, "%s\tmaximum block size: %d samples\n", prefix, b.MaximumBlockSize)
			fmt.Fprintf(w, "%s\tminimum frame size: %d bytes\n", prefix, b.MinimumFrameSize)
			fmt.Fprintf(w, "%s\tmaximum frame size: %d bytes\n", prefix, b.MaximumFrameSize)
			fmt.Fprintf(w, "%s\tsample_rate: %d Hz\n", prefix, b.SampleRate)
			fmt.Fprintf(w, "%s\tchannels: %d\n", prefix, b.Channels)
			fmt.Fprintf(w, "%s\tbits-per-sample: %d\n", prefix, b.BitsPerSample)
			fmt.Fprintf(w, "%s\ttotal samples: %d\n", prefix, b.TotalSamples)
			fmt.Fprintf(w, "%s\tMD5 signature: %x\n", prefix, b.MD5)
		case *flac.Application:
			fmt.Fprintf(w, "%s\tapplication id: %s\n", prefix, b.ID)
			fmt.Fprintf(w, "%s\tapplication data: %x\n", prefix, b.Data)
		case *flac.SeekTable:
			fmt.Fprintf(w, "%s\tseek points: %d\n", prefix, len(b.SeekPoints))
			for i, p := range b.SeekPoints {
				fmt.Fprintf(w, "%s\t\tpoint %d: sample_number=%d, stream_offset=%d, frame_samples=%d\n", prefix, i, p.SampleNumber, p.Offset, p.NumSamples)
			}
		case *flac.VorbisComment:
			fmt.Fprintf(w, "%s\tvendor string: %s\n", prefix, b.Vendor)
			fmt.Fprintf(w, "%s\tcomments: %d\n", prefix, len(b.UserComments))
			for i, c := range b.UserComments {
				fmt.Fprintf(w, "%s\t\tcomment[%d]: %s\n", prefix, i, c)
			}
		case *flac.CueSheet:
			fmt.Fprintf(w, "%s\tmedia catalog number: %s\n", prefix, b.CatalogNumber)
			fmt.Fprintf(w, "%s\tlead-in: %d\n", prefix, b.NumLeadInSamples)
			fmt.Fprintf(w, "%s\tis CD: %t\n", prefix, b.IsCD)
			fmt.Fprintf(w, "%s\tnumber of tracks: %d\n", prefix, len(b.Tracks))
			for i, t := range b.Tracks {
				fmt.Fprintf(w, "%s\t\ttrack[%d]\n", prefix, i)
				fmt.Fprintf(w, "%s\t\t\toffset: %d\n", prefix, t.OffsetSamples)
				fmt.Fprintf(w, "%s\t\t\tnumber: %d\n", prefix, t.TrackNumber)
				fmt.Fprintf(w, "%s\t\t\tISRC: %s\n", prefix, t.ISRC)
				if t.IsAudio {
					fmt.Fprintf(w, "%s\t\t\ttype: AUDIO\n", prefix)
				} else {
					fmt.Fprintf(w, "%s\t\t\ttype: NON-AUDIO\n", prefix)
				}
				fmt.Fprintf(w, "%s\t\t\tpre-emphasis: %t\n", prefix, t.PreEmphasis)
				fmt.Fprintf(w, "%s\t\t\tnumber of index points: %d\n", prefix, len(t.Indices))

				for j, p := range t.Indices {
					fmt.Fprintf(w, "%s\t\t\t\tindex[%d]\n", prefix, j)
					fmt.Fprintf(w, "%s\t\t\t\t\toffset: %d\n", prefix, p.OffsetSamples)
					fmt.Fprintf(w, "%s\t\t\t\t\tnumber: %d\n", prefix, p.PointNumber)
				}
			}
		case *flac.Picture:
			fmt.Fprintf(w, "%s\ttype: %d (%s)\n", prefix, b.Type, b.Type)
			fmt.Fprintf(w, "%s\tMIME type: %s\n", prefix, b.MimeType)
			fmt.Fprintf(w, "%s\tdescription: %s\n", prefix, b.Description)
			fmt.Fprintf(w, "%s\twidth: %d\n", prefix, b.Width)
			fmt.Fprintf(w, "%s\theight: %d\n", prefix, b.Height)
			fmt.Fprintf(w, "%s\tdepth: %d\n", prefix, b.Depth)
			if b.Colors == 0 {
				fmt.Fprintf(w, "%s\tcolors: 0 (unindexed)\n", prefix)
			} else {
				fmt.Fprintf(w, "%s\tcolors: %d\n", prefix, b.Colors)
			}
			fmt.Fprintf(w, "%s\tdata length: %d\n", prefix, len(b.Data))
		}

		i++
	}
}
