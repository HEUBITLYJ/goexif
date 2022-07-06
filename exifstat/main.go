package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/HEUBITLYJ/goexif/exif"
	"github.com/HEUBITLYJ/goexif/mknote"
	"github.com/HEUBITLYJ/goexif/tiff"
)

var mnote = flag.Bool("mknote", false, "try to parse makernote data")
var thumb = flag.Bool("thumb", false, "dump thumbail data to stdout (for first listed image file)")

func DecodeRuntimeError() {
	fname := "exif/2af7ee4c984d376e5658a1f6f959a75f.jpg"

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	// Optionally register camera makenote data parsing - currently Nikon and
	// Canon are supported.
	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	camModel, _ := x.Get(exif.Model) // normally, don't ignore errors!
	fmt.Println(camModel.StringVal())

	focal, _ := x.Get(exif.FocalLength)
	numer, denom, _ := focal.Rat2(0) // retrieve first (only) rat. value
	fmt.Printf("%v/%v", numer, denom)

	// Two convenience functions exist for date/time taken and GPS coords:
	tm, _ := x.DateTime()
	fmt.Println("Taken: ", tm)

	lat, long, _ := x.LatLong()
	fmt.Println("lat, long: ", lat, ", ", long)

	fmt.Println("Decode success")
}

func main() {
	DecodeRuntimeError()

	flag.Parse()
	fnames := flag.Args()

	if *mnote {
		exif.RegisterParsers(mknote.All...)
	}

	for _, name := range fnames {
		f, err := os.Open(name)
		if err != nil {
			log.Printf("err on %v: %v", name, err)
			continue
		}

		x, err := exif.Decode(f)
		if err != nil {
			log.Printf("err on %v: %v", name, err)
			continue
		}

		if *thumb {
			data, err := x.JpegThumbnail()
			if err != nil {
				log.Fatal("no thumbnail present")
			}
			if _, err := os.Stdout.Write(data); err != nil {
				log.Fatal(err)
			}
			return
		}

		fmt.Printf("\n---- Image '%v' ----\n", name)
		x.Walk(Walker{})
	}
}

type Walker struct{}

func (_ Walker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	data, _ := tag.MarshalJSON()
	fmt.Printf("    %v: %v\n", name, string(data))
	return nil
}
