package main

import (verify-pack/cmd/verify-pack/main.go
	"bufio"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/git-lfs/git-lfs/git/odb/pack"
	"github.com/pkg/errors"

	"golang.org/x/exp/mmap"
)

func main() {
	if len(os.Args) <= 1 {
		die(usage())
	}

	f, err := mmap.Open(os.Args[1])
	if err != nil {
		die(errors.Wrap(err, "could not open index file").Error())
	}

	idx, err := pack.DecodeIndex(f)
	if err != nil {
		die(errors.Wrap(err, "could not decode index").Error())
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		name := scanner.Text()

		oid, err := hex.DecodeString(name)
		if err != nil {
			die(errors.Wrapf(err, "could not decode name %q", name).Error())
		}

		e, err := idx.Entry(oid)
		if err != nil {
			die(errors.Wrapf(err, "could not find entry %q", name).Error())
		}

		fmt.Printf("%+v\n", e)
	}

	if err := scanner.Err(); err != nil {
		die(err.Error())
	}

	if err := f.Close(); err != nil {
		die(errors.Wrap(err, "could not close file").Error())
	}
}

func die(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func usage() string {
	return fmt.Sprintf("usage: %s <pack>", os.Args[0])
}
