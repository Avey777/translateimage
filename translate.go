package translateimage

import (
	"image"
	"io"

	languagecodes "github.com/spywiree/langcodes"
)

// Path must be absolute.
// Supported file types: .jpg, .jpeg, .png.
func TranslateFile(path string, source, target languagecodes.LanguageCode) (*ImageData, error) {
	ctx, err := NewContext()
	if err != nil {
		return nil, err
	}
	defer ctx.Close()

	return ctx.TranslateFile(path, source, target)
}

func TranslateImage(img image.Image, source, target languagecodes.LanguageCode) (*ImageData, error) {
	ctx, err := NewContext()
	if err != nil {
		return nil, err
	}
	defer ctx.Close()

	return ctx.TranslateImage(img, source, target)
}

// Supported image types: jpg, jpeg, png.
func TranslateReader(r io.Reader, source, target languagecodes.LanguageCode) (*ImageData, error) {
	ctx, err := NewContext()
	if err != nil {
		return nil, err
	}
	defer ctx.Close()

	return ctx.TranslateReader(r, source, target)
}
