package gray_image

import (
	"github.com/disintegration/imaging"
	"image/jpeg"
	"io"
	"log"
	"os"

	"github.com/jdeng/goheif"
)

// Skip Writer for exif writing
type writerSkipper struct {
	w           io.Writer
	bytesToSkip int
}

func (w *writerSkipper) Write(data []byte) (int, error) {
	if w.bytesToSkip <= 0 {
		return w.w.Write(data)
	}

	if dataLen := len(data); dataLen < w.bytesToSkip {
		w.bytesToSkip -= dataLen
		return dataLen, nil
	}

	if n, err := w.w.Write(data[w.bytesToSkip:]); err == nil {
		n += w.bytesToSkip
		w.bytesToSkip = 0
		return n, nil
	} else {
		return n, err
	}
}

func newWriterExif(w io.Writer, exif []byte) (io.Writer, error) {
	writer := &writerSkipper{w, 2}
	soi := []byte{0xff, 0xd8}
	if _, err := w.Write(soi); err != nil {
		return nil, err
	}

	if exif != nil {
		app1Marker := 0xe1
		markerlen := 2 + len(exif)
		marker := []byte{0xff, uint8(app1Marker), uint8(markerlen >> 8), uint8(markerlen & 0xff)}
		if _, err := w.Write(marker); err != nil {
			return nil, err
		}

		if _, err := w.Write(exif); err != nil {
			return nil, err
		}
	}

	return writer, nil
}

func ResizeHeif(rawFile, dstFile string, width, height int) (err error) {
	fi, err := os.Open(rawFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer fi.Close()

	exif, err := goheif.ExtractExif(fi)
	if err != nil {
		log.Printf("Warning: no EXIF from %s: %v\n", rawFile, err)
		return
	}

	img, err := goheif.Decode(fi)
	if err != nil {
		log.Printf("Failed to parse %s: %v\n", rawFile, err)
		return
	}
	x, y, err := CalcProportionalXY(rawFile, width, height)
	if err != nil {
		log.Printf("CalcProportionalXY error: %v", err)
		return
	}
	canvas := imaging.Resize(img, x, y, imaging.Lanczos)

	fo, err := os.Create(dstFile)
	if err != nil {
		log.Printf("Failed to create output file %s: %v\n", dstFile, err)
		return
	}
	defer fo.Close()

	w, _ := newWriterExif(fo, exif)
	err = jpeg.Encode(w, canvas, nil)
	if err != nil {
		log.Printf("Failed to encode %s: %v\n", dstFile, err)
		return
	}

	log.Printf("Convert %s to %s successfully\n", rawFile, dstFile)
	return
}
