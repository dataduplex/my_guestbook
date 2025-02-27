package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
)

type FileDB struct {
	dir string
}

func NewFileDB(dir string) (*FileDB, error) {
	return &FileDB{
		dir: dir,
	}, nil
}

func (f *FileDB) SaveGuests(ctx context.Context, g *Guests) error {
	fh, err := os.Create(path.Join(f.dir, "db"))
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	for name := range g.Guests { // (1)
		fh.WriteString(name)
		fh.WriteString(" ")
		fh.Write(strconv.AppendBool(nil, g.IsSpecial(name))) // (2)
		fh.WriteString("\n")
	}

	if err := fh.Close(); err != nil {
		return fmt.Errorf("error saving database: %w", err)
	}

	return nil
}
