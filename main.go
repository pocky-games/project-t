package main

import (
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	bg     *ebiten.Image
	person *ebiten.Image
	cat    *ebiten.Image
	window *ebiten.Image
)

type game struct {
}

func (g *game) Update() error {
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.DrawImage(bg, nil)
	screen.DrawImage(cat, nil)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(1280, 0)
	screen.DrawImage(person, op)

	// 描画元画像を９分割する
	srcRects := []image.Rectangle{}
	xs := []int{0, 30, 70, 100}
	ys := []int{0, 30, 70, 100}
	for x := range 3 {
		for y := range 3 {
			rect := image.Rect(xs[x], ys[y], xs[x+1], ys[y+1])
			srcRects = append(srcRects, rect)
		}
	}

	// 描画先領域を９分割する
	dstRects := []image.Rectangle{}
	xs = []int{0, 30, 1080 - 30, 1080}
	ys = []int{0, 30, 260 - 30, 260}
	for x := range 3 {
		for y := range 3 {
			rect := image.Rect(xs[x], ys[y], xs[x+1], ys[y+1])
			dstRects = append(dstRects, rect)
		}
	}

	// ９分割した画像を描画する
	for i := range 9 {
		srcRect := srcRects[i]
		dstRect := dstRects[i]
		// 描画元SubImageを作成する
		subImage := window.SubImage(srcRect).(*ebiten.Image)
		// 拡大率を求める
		scaleX := float64(dstRect.Dx()) / float64(srcRect.Dx())
		scaleY := float64(dstRect.Dy()) / float64(srcRect.Dy())
		// 描画する
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(float64(dstRect.Min.X), float64(dstRect.Min.Y))
		op.GeoM.Translate(100, 450)
		screen.DrawImage(subImage, op)
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	img, _, err := ebitenutil.NewImageFromFile("bg.png")
	if err != nil {
		panic(err)
	}
	bg = img

	img, _, err = ebitenutil.NewImageFromFile("person.png")
	if err != nil {
		panic(err)
	}
	person = img

	img, _, err = ebitenutil.NewImageFromFile("cat.png")
	if err != nil {
		panic(err)
	}
	cat = img

	img, _, err = ebitenutil.NewImageFromFile("window.png")
	if err != nil {
		panic(err)
	}
	window = img

	ebiten.SetWindowTitle("ノベルゲーム")
	ebiten.SetWindowSize(1280, 720)
	g := &game{}
	ebiten.RunGame(g)
}
