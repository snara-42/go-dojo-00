package imgconv

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"path/filepath"
	"os"
)

func Convert(fn string, in, out string) error {
	if err := IsValidExt(in); err != nil {
		return err
	}
	if err := IsValidExt(out); err != nil {
		return err
	}
	if ext := filepath.Ext(fn)[1:]; ext != in {
		return fmt.Errorf("invalid argument: unmatched extensions %s and %s", ext, in)
	}
	fin, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer fin.Close()
	var img image.Image
	switch in {
	case "jpg", "jpeg":
		img, err = jpeg.Decode(fin)
	case "png":
		img, err = png.Decode(fin)
	case "gif":
		img, err = gif.Decode(fin)
	default:
		return fmt.Errorf("invalid argument: unsupported extension %s", in)
	}
	if err != nil {
		return err
	}
	fout, err := os.Create(ConvertExt(fn, out))
	if err != nil {
		return err
	}
	defer fout.Close()
	switch out {
	case "jpg", "jpeg":
		err = jpeg.Encode(fout, img, nil)
	case "png":
		err = png.Encode(fout, img)
	case "gif":
		err = gif.Encode(fout, img, nil)
	default:
		return fmt.Errorf("invalid argument: unsupported extension %s", out)
	}
	if err != nil {
		return err
	}
	return nil
}

func IsValidExt(ext string) (error) {
	switch ext {
	case "jpg", "jpeg", "png", "gif":

	default:
		return fmt.Errorf("invalid argument: unsupported extension %s", ext)
	}
	return nil
}
func ConvertExt(fn, out string) string{
	return fn[0:len(fn)-len(filepath.Ext(fn))]+"."+out
}
