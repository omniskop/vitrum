package std

import (
	"path/filepath"

	vit "github.com/omniskop/vitrum/vit"
	"github.com/tdewolff/canvas"
)

type img canvas.Image

func (i *Image) reloadImage() {
	path := i.path.String()
	if path == "" {
		i.imageData = nil
		return
	}
	pos, ok := i.path.Position()
	if !ok {
		i.Context().Global.Environment.Logger().Printf("failed to open image file: path string has no position\r\n")
		return
	}

	file, err := pos.FilePath.Dir().Open(path)
	if err != nil {
		i.Context().Global.Environment.Logger().Printf("failed to open image file: %s\r\n", err)
		return
	}
	defer file.Close()

	var loaded canvas.Image
	ext := filepath.Ext(path)
	switch ext {
	case ".png":
		loaded, err = canvas.NewPNGImage(file)
	case ".jpg", ".jpeg":
		loaded, err = canvas.NewJPEGImage(file)
	}
	if err != nil {
		i.Context().Global.Environment.Logger().Printf("failed to parse image %s: %s\r\n", path, err)
		return
	}
	i.imageData = (*img)(&loaded)

	i.SetContentSize(float64(loaded.Bounds().Dx()), float64(loaded.Bounds().Dy()))
}

func (i *Image) Draw(ctx vit.DrawingContext, area vit.Rect) error {
	rect := i.Bounds()

	if i.imageData != nil {
		imageBounds := vit.ImageRect(i.imageData.Image.Bounds())
		finalBounds := i.fill(rect, imageBounds, Image_FillMode(i.fillMode.Int()))

		scaleX := finalBounds.Width() / imageBounds.Width()
		scaleY := finalBounds.Height() / imageBounds.Height()

		ctx.ScaleAbout(scaleX, scaleY, finalBounds.X1, finalBounds.Y1)
		ctx.DrawImage(finalBounds.X1, finalBounds.Y1, i.imageData, canvas.Resolution(1))
		ctx.ScaleAbout(1/scaleX, 1/scaleY, finalBounds.X1, finalBounds.Y1)
	}

	return i.Root.DrawChildren(ctx, rect)
}

func (i *Image) fill(space vit.Rect, img vit.Rect, fillMode Image_FillMode) vit.Rect {
	switch fillMode {
	case Image_FillMode_Fill:
		return space
	case Image_FillMode_PreferUnchanged:
		if img.Width() < space.Width() && img.Height() < space.Height() {
			return img.CenteredIn(space)
		}
		fallthrough
	case Image_FillMode_Fit:
		scaleX := space.Width() / float64(img.Width())
		scaleY := space.Height() / float64(img.Height())
		if scaleX > scaleY {
			scaleX = scaleY
		} else {
			scaleY = scaleX
		}
		return vit.NewRect(0, 0, img.Width()*scaleX, img.Height()*scaleY).CenteredIn(space)
	default:
		i.Context().Global.Environment.Logger().Printf("image.fill called with unknown fill mode %T\r\n", fillMode)
		return img
	}
}
