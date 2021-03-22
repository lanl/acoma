package csv

import (
        "bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
	"acoma/oligo"
	"acoma/oligo/long"
)

func Read(fname string, ignoreBad bool) ([]oligo.Oligo, error) {
	var oligos []oligo.Oligo

	err := Parse(fname, func(id, sequence string, quality []byte, reverse bool) error {
		ol, ok := long.FromString(sequence)
		if !ok {
			if ignoreBad {
				// skip
				return nil
			} else {
				return fmt.Errorf("invalid oligo: %s\n", sequence)
			}
		}

		if reverse {
			oligo.Reverse(ol)
			oligo.Invert(ol)
		}

		oligos = append(oligos, ol)
		return nil
	})

	return oligos, err
}

func Parse(fname string, process func(id, sequence string, quality []byte, reverse bool) error) error {
	var r io.Reader

	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	if cf, err := gzip.NewReader(f); err == nil {
		r = cf
	} else {
		f.Seek(0, 0)
		r = f

	}

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		var seq string

		l := sc.Text()
		if strings.Contains(l, ",") {
			ls := strings.Split(l, ",")
			seq = ls[0]
		} else if strings.Contains(l, " ") {
			ls := strings.Split(l, " ")
			seq = ls[0]
		} else {
			seq = l
		}

		if err := process("", seq, nil, false); err != nil {
			return err
		}
	}

	return nil
}
