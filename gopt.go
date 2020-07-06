package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/discordapp/lilliput"
)

func main() {
	var inputWidth int

	inputFilename := os.Args[0]
	flag.IntVar(&inputWidth, "width", 500, "Input width")
	flag.Parse()

	if inputFilename == "" {
		fmt.Printf("No input gif provided, quitting.\n")
		flag.Usage()
		os.Exit(1)
	}

	fileBuf, err := ioutil.ReadFile(inputFilename)
	if err != nil {
		fmt.Printf("failed to read input file, %s\n", err)
		os.Exit(1)
	}

	decoder, err := lilliput.NewDecoder(fileBuf)
	if err != nil {
		fmt.Printf("error decoding image, %s\n", err)
		os.Exit(1)
	}
	defer decoder.Close()

	header, err := decoder.Header()
	// this error is much more comprehensive and reflects
	// format errors
	if err != nil {
		fmt.Printf("error reading image header, %s\n", err)
		os.Exit(1)
	}

	if decoder.Description() != "GIF" {
		fmt.Printf("input is not a gif")
		os.Exit(1)
	}

	ops := lilliput.NewImageOps(8192)
	defer ops.Close()

	height := int(float64(inputWidth) / float64(header.Width()) * float64(header.Height()))
	fmt.Printf("Resizing to %dpx x %dpx\n", inputWidth, height)

	opts := &lilliput.ImageOptions{
		FileType:             ".gif",
		Width:                inputWidth,
		Height:               height,
		ResizeMethod:         lilliput.ImageOpsResize,
		NormalizeOrientation: true,
		EncodeOptions:        map[int]int{lilliput.WebpQuality: 85},
	}

	outputImg := make([]byte, 50*1024*1024)

	outputImg, err = ops.Transform(decoder, opts, outputImg)
	if err != nil {
		fmt.Printf("error transforming gif, %s\n", err)
		os.Exit(1)
	}

	outputFilename := strings.TrimSuffix(inputFilename, filepath.Ext(inputFilename)) + "-opt.gif"

	err = ioutil.WriteFile(outputFilename, outputImg, 0664)
	if err != nil {
		fmt.Printf("error writing out optimized gif, %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("gif written to %s\n", outputFilename)
}
