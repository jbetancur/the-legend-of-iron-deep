package render

// Cell-grid rendering system based on the Screaming Brain Studios technique.
//
// Cells are always SQUARE, sized from ViewportH. In a widescreen viewport the
// extra horizontal space reveals more side columns — it does NOT stretch cells.
//
// Each depth layer is 50 % the size of the previous, centered on the vanishing
// point (viewport center). Layers contain:
//   depth 0: 3 columns  [-1 … +1]
//   depth 1: 5 columns  [-2 … +2]
//   depth 2: 7 columns  [-3 … +3]
//   depth 3: 9 columns  [-4 … +4]
//
// For widescreen we widen the range so off-screen cells still get their side
// walls drawn (they poke into view on 16:9).

type rect struct {
	x, y, w, h float64
}

func (r rect) left() float64   { return r.x }
func (r rect) right() float64  { return r.x + r.w }
func (r rect) top() float64    { return r.y }
func (r rect) bottom() float64 { return r.y + r.h }

// cellRect returns the bounding rect of a cell at (depth, column).
// Cells are always square (side = ViewportH * scale) and centered on the
// vanishing point at the viewport center.
func cellRect(depth, col int) rect {
	scale := 1.0
	for i := 0; i < depth; i++ {
		scale *= 0.5
	}

	side := float64(ViewportH) * scale

	cx := float64(ViewportW) / 2
	cy := float64(ViewportH) / 2

	x := cx - side/2 + float64(col)*side
	y := cy - side/2

	return rect{x: x, y: y, w: side, h: side}
}

// backWallRect returns the back wall of cell (depth, col).
// Per the SBS nesting rule: backWallRect(d, c) == cellRect(d+1, c).
func backWallRect(depth, col int) rect {
	return cellRect(depth+1, col)
}

// leftWallQuad returns 4 corners (CW from top-left) of the left side-wall
// trapezoid: cell's left edge → back wall's left edge.
func leftWallQuad(depth, col int) (x0, y0, x1, y1, x2, y2, x3, y3 float64) {
	c := cellRect(depth, col)
	bw := backWallRect(depth, col)
	return c.left(), c.top(), bw.left(), bw.top(), bw.left(), bw.bottom(), c.left(), c.bottom()
}

// rightWallQuad returns 4 corners of the right side-wall trapezoid.
func rightWallQuad(depth, col int) (x0, y0, x1, y1, x2, y2, x3, y3 float64) {
	c := cellRect(depth, col)
	bw := backWallRect(depth, col)
	return bw.right(), bw.top(), c.right(), c.top(), c.right(), c.bottom(), bw.right(), bw.bottom()
}

// floorQuad returns 4 corners of the floor trapezoid: bottom edge of cell →
// bottom edge of back wall, forming a receding ground plane.
func floorQuad(depth, col int) (x0, y0, x1, y1, x2, y2, x3, y3 float64) {
	c := cellRect(depth, col)
	bw := backWallRect(depth, col)
	return c.left(), c.bottom(), c.right(), c.bottom(), bw.right(), bw.bottom(), bw.left(), bw.bottom()
}

// ceilingQuad returns 4 corners of the ceiling trapezoid: top edge of cell →
// top edge of back wall.
func ceilingQuad(depth, col int) (x0, y0, x1, y1, x2, y2, x3, y3 float64) {
	c := cellRect(depth, col)
	bw := backWallRect(depth, col)
	return c.left(), c.top(), c.right(), c.top(), bw.right(), bw.top(), bw.left(), bw.top()
}

// columnRange returns the min/max column indices visible at a given depth.
// The base count follows SBS (2*depth+3), but we widen it to cover the full
// viewport width so widescreen doesn't show gaps on the sides.
func columnRange(depth int) (int, int) {
	scale := 1.0
	for i := 0; i < depth; i++ {
		scale *= 0.5
	}
	side := float64(ViewportH) * scale

	// How many cells from center to cover half the viewport width, plus one
	// extra to make sure partially-visible cells get drawn.
	halfCols := int(float64(ViewportW)/2/side) + 2

	// Never fewer than the SBS minimum (depth+1).
	if halfCols < depth+1 {
		halfCols = depth + 1
	}
	return -halfCols, halfCols
}
