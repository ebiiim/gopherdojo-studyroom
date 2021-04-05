package conv

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// Ext defines image file extensions.
type Ext string

// Extensions.
const (
	JPEG = ".jpg"
	PNG  = ".png"
	BMP  = ".bmp"
	TIFF = ".tiff"
)

// toExt maps extension string (e.g., ".jpg") to Ext.
var toExt map[string]Ext = map[string]Ext{
	".jpeg": JPEG,
	".jpg":  JPEG,
	".png":  PNG,
	".bmp":  BMP,
	".tiff": TIFF,
	".tif":  TIFF,
}

// GetExt gets Ext from a path string (e.g., "path/to/abc.jpg", "abc.jpg", ".jpg").
func GetExt(imgPath string) (Ext, error) {
	extStr := strings.ToLower(filepath.Ext(imgPath))
	ext, ok := toExt[extStr]
	if !ok {
		return "", fmt.Errorf("conv.GetExt %s (%w)", extStr, ErrExtensionNotSupported)
	}
	return ext, nil
}

// Image file signatures.
var (
	sigJPEG   = []byte{0xff, 0xd8}
	sigPNG    = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	sigBMP    = []byte{0x42, 0x4D}
	sigTIFFle = []byte{0x49, 0x49, 0x2A, 0x00} // little-endian
	sigTIFFbe = []byte{0x4D, 0x4D, 0x00, 0x2A} // big-endian
)

// maxSigLen defines length of the longest signature in bytes.
const maxSigLen = 8

// GetExtFromSignature gets Ext by checking file signature.
// This function reads `r` so please use io.TeeReader to replicate `r`.
func GetExtFromSignature(r io.Reader) (Ext, error) {
	sig := make([]byte, maxSigLen)
	_, err := r.Read(sig)
	if err != nil {
		return "", fmt.Errorf("conv.GetExtFromSignature (%w)", ErrFileAccess)
	}
	if bytes.HasPrefix(sig, sigJPEG) {
		return JPEG, nil
	}
	if bytes.HasPrefix(sig, sigPNG) {
		return PNG, nil
	}
	if bytes.HasPrefix(sig, sigBMP) {
		return BMP, nil
	}
	if bytes.HasPrefix(sig, sigTIFFbe) || bytes.HasPrefix(sig, sigTIFFle) {
		return TIFF, nil
	}
	return "", fmt.Errorf("conv.GetExtFromSignature (%w)", ErrFileNotSupported)
}
