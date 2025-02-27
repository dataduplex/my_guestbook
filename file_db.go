package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
)

// Needs a method to read and load guests from file DB.
// Currently, any restart of the application will erase the data saved to DB.
type FileDB struct {
	dir string
}

func NewFileDB(dir string) (*FileDB, error) {
	// perhaps validate if path exists, if not, create it here
	return &FileDB{
		dir: dir,
	}, nil
}

func (f *FileDB) SaveGuests(ctx context.Context, g *Guests) error {
	// There are no functions to load existing "db" file to map
	// The below line will erase the contents of any existing "db" file - we lose data.
	// If for any reason, the application crashes while writing to DB - we lose data.
	// Perhaps a safer option is to create a new file each time with current timestamp
	// After finish writing, save it with an extension (For example - "db.<timestamp>.FULL" to indicate it has all data saved)
	// when reading from DB, use the file with latest timestamp
	// old files can be purged in a separate thread or a separate cron job.
	fh, err := os.Create(path.Join(f.dir, "db"))
	if err != nil {
		// return error to caller instead of panic
		panic(err)
	}
	defer fh.Close()

	// as mentioned in another comment, possible deadlock scenario with locking twice here at (1) and (2)
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
