package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/tetsuzawa/go-soundlib/spatial"
)

func init() {
	log.SetFlags(0)
	flag.Usage = func() {
		log.Printf("Usage of %s:\n", os.Args[0])
		log.Printf("overlap-add subject sound_file(.DXX) move_width move_velocity end_angle outdir\n")
		flag.PrintDefaults()
	}
}

func main() {
	if err := run(); err != nil {
		log.Println(err)
		flag.Usage()
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	if flag.NArg() != 6 {
		return errors.New("invalid arguments")
	}
	args := flag.Args()
	subject := args[0]
	soundName := args[1]
	moveWidth, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	moveVelocity, err := strconv.Atoi(args[3])
	if err != nil {
		return err
	}
	endAngle, err := strconv.Atoi(args[4])
	if err != nil {
		return err
	}
	outDir := args[5]
	return spatial.OverlapAdd(subject, soundName, moveWidth, moveVelocity, endAngle, outDir)
}
