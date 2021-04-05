// Package conv provides image format conversion.
package conv

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

// Options defines conversion options.
type Options struct {
	SrcExt Ext
	DstExt Ext
}

// convFunc defines the function used in DoFile.
var convFunc = Do

// DoFile converts an image file.
func DoFile(srcPath string, opts *Options) (err error) {
	in, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("conv.DoFile %v (%w)", err, ErrFileAccess)
	}
	defer func() {
		if cErr := in.Close(); err == nil && cErr != nil {
			err = fmt.Errorf("DoFile defer %v (%w) <- %v", cErr, ErrFileAccess, err)
		}
	}()

	// e.g., "abc.jpg" -> "abc.png"
	dstPath := strings.TrimSuffix(srcPath, filepath.Ext(srcPath)) + string(opts.DstExt)

	out, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("conv.DoFile %v (%w)", err, ErrFileAccess)
	}
	defer func() {
		if cErr := in.Close(); err == nil && cErr != nil {
			err = fmt.Errorf("DoFile defer %v (%w) <- %v", cErr, ErrFileAccess, err)
		}
	}()

	if err := convFunc(in, out, opts); err != nil {
		return fmt.Errorf("conv.DoFile -> %w", err)
	}
	return
}

// Do converts an image except for the case of:
//   - 1. `opt.SrcExt` and `opt.DstExt` is same.
//   - 2. The image format from `r` is not same with `opt.SrcExt`.
func Do(r io.Reader, w io.Writer, opt *Options) error {
	if opt.SrcExt == opt.DstExt {
		return nil // Case 1
	}
	var buf bytes.Buffer
	r = io.TeeReader(r, &buf)
	srcExt, err := GetExtFromSignature(&buf)
	if err != nil {
		return fmt.Errorf("conv.DoIfNeeded -> %w", err)
	}
	if srcExt != opt.SrcExt {
		return nil // Case 2
	}
	return do(r, w, opt)
}

// do converts an image to `opt.DstExt`.
func do(r io.Reader, w io.Writer, opts *Options) error {
	m, _, err := image.Decode(r)
	if err != nil {
		return fmt.Errorf("conv.Do %v (%w)", err, ErrCouldNotDecode)
	}
	switch opts.DstExt {
	case JPEG:
		err = jpeg.Encode(w, m, nil)
	case PNG:
		err = png.Encode(w, m)
	case BMP:
		err = bmp.Encode(w, m)
	case TIFF:
		err = tiff.Encode(w, m, nil)
	}
	if err != nil {
		return fmt.Errorf("conv.Do %v (%w)", err, ErrCouldNotEncode)
	}
	return err
}
