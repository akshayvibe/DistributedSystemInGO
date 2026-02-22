//ANCHOR
/*
	Offset   →   Position
	0        →   0
	1        →   12
	2        →   24
	3        →   36
*/
/*❌ Problem without index
If you want to read record #1000, what do you do?
👉 You must:
		Start from beginning
		Read length
		Skip data
		Repeat 1000 times
✅ Solution: Index Table
*/
package log

import (
	"io"
	"os"

	"github.com/tysonmote/gommap"
)

var (
	offWidth uint64 = 4
	posWidth uint64 = 8
	entWidth         = offWidth + posWidth
)

type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
}

func newIndex(f *os.File, c Config) (*index, error) {
	idx := &index{
		file: f,
	}

	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	idx.size = uint64(fi.Size())

	// grow file to max size
	if err := os.Truncate(
		f.Name(),
		int64(c.Segment.MaxIndexBytes),
	); err != nil {
		return nil, err
	}

	// memory map the file
	if idx.mmap, err = gommap.Map(
		idx.file.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE,
		gommap.MAP_SHARED,
	); err != nil {
		return nil, err
	}

	return idx, nil
}

func (i *index) Close() error {
	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}

	if err := i.file.Sync(); err != nil {
		return err
	}

	if err := i.file.Truncate(int64(i.size)); err != nil {
		return err
	}

	return i.file.Close()
}

func (i *index) Name() string {
	return i.file.Name()
}

// Read returns offset and position
func (i *index) Read(in int64) (uint32, uint64, error) {
	if i.size == 0 {
		return 0, 0, io.EOF
	}

	// if -1 → read last entry
	if in == -1 {
		in = int64(i.size/entWidth) - 1
	}

	pos := uint64(in) * entWidth

	if i.size < pos+entWidth {
		return 0, 0, io.EOF
	}

	out := enc.Uint32(i.mmap[pos : pos+offWidth])
	p := enc.Uint64(i.mmap[pos+offWidth : pos+entWidth])

	return out, p, nil
}

// Write adds offset + position to index
func (i *index) Write(off uint32, pos uint64) error {
	// check if full
	if uint64(len(i.mmap)) < i.size+entWidth {
		return io.EOF
	}

	// write offset
	enc.PutUint32(i.mmap[i.size:i.size+offWidth], off)

	// write position
	enc.PutUint64(i.mmap[i.size+offWidth:i.size+entWidth], pos)

	i.size += entWidth

	return nil
}