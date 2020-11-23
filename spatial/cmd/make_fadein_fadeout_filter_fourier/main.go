package main

import (
	"errors"
	"flag"
	"github.com/tetsuzawa/go-soundlib/dxx"
	"github.com/tetsuzawa/go-soundlib/spatial"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func init() {
	log.SetFlags(0)
	flag.Usage = func() {
		log.Printf("Usage: %s signal_length(sample) filename(fadein filter) filename(fadeout filter)\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	if err := run(); err != nil {
		log.Printf("error: %+v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	if flag.NArg() != 3 {
		return errors.New("invalid arguments")
	}
	args := flag.Args()
	samples, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	fadeinFiltName := args[1]
	fadeoutFiltName := args[2]
	fadeinFilt, fadeoutFilt := spatial.GenerateFadeinFadeoutFilt(samples)
	if err := dxx.WriteToFile(fadeinFiltName, fadeinFilt); err != nil {
		return err
	}
	if err := dxx.WriteToFile(fadeoutFiltName, fadeoutFilt); err != nil {
		return err
	}
	return nil
}
