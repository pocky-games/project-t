package main

import (
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	bg     *ebiten.Image
	person *ebiten.Image
	cat    *ebiten.Image
	window *ebiten.Image
)

var (
	fontFace *text.GoTextFace
	ticks    = 0
)

type game struct {
}

func (g *game) Update() error {
	ticks++
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ticks = 0
	}
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
		op.ColorScale.ScaleAlpha(0.8)
		screen.DrawImage(subImage, op)
	}

	// 名前欄を表示する
	dstRects = []image.Rectangle{}
	xs = []int{0, 30, 200 - 30, 200}
	ys = []int{0, 30, 70 - 30, 70}
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
		op.GeoM.Translate(150, 400)
		screen.DrawImage(subImage, op)
	}

	// 会話文表示
	textop := &text.DrawOptions{}
	textop.GeoM.Translate(200, 480)
	textop.ColorScale.Scale(0, 0, 0, 1)
	textop.LineSpacing = 30 * 1.5

	glyphs := text.AppendGlyphs(nil, "吾輩はsakamotoである。名前はまだない。", fontFace, &textop.LayoutOptions)
	length := ticks / 5 // 5tickに1文字送り
	for i, g := range glyphs {
		if i > length {
			break // 表示文字数を超えたら終了
		}
		textop.GeoM.Reset() // textopを全文字で使い回して、GeoMだけリセット
		textop.GeoM.Translate(200, 480)
		textop.GeoM.Translate(g.X, g.Y)
		screen.DrawImage(g.Image, &textop.DrawImageOptions)
	}
	// 話者表示
	textop = &text.DrawOptions{}
	w, h := text.Measure("sakamoto", fontFace, 30*1.5)
	textop.GeoM.Translate(250-w/2, 430-h/2)
	textop.ColorScale.Scale(0, 0, 0, 1)
	text.Draw(screen, "sakamoto", fontFace, textop)
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

	// ファイルを読み込む
	f, err := os.Open("static/NotoSansJP-Bold.ttf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// フォントを読み込む
	src, err := text.NewGoTextFaceSource(f)
	if err != nil {
		panic(err)
	}

	// フォントサイズを指定し、「フォントフェイス」を作る
	fontFace = &text.GoTextFace{Source: src, Size: 30}

	ebiten.SetWindowTitle("ノベルゲーム")
	ebiten.SetWindowSize(1280, 720)
	g := &game{}
	ebiten.RunGame(g)
}
