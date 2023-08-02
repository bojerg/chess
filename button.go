package main

type Button struct {
	text     string
	fontSize int
	x        int
	y        int
}

const (
	BtnHeight = 80
	BtnWidth  = 180
)

func (b *Button) TextX() int {
	return b.x + BtnWidth/2 - ((len(b.text) * b.fontSize) / 2)
}

func (b *Button) TextY() int {
	return b.y + BtnHeight/2 + b.fontSize/2
}

func (b *Button) PosInBounds(x, y int) bool {
	return x >= b.x && x <= b.x+BtnWidth && y >= b.y && y <= b.y+BtnHeight
}
