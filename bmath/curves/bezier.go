package curves

import (
	math2 "github.com/wieku/danser-go/bmath"
	"math"
)


type Bezier struct {
	points           []math2.Vector2d
	ApproxLength     float64
	lengthCalculated bool
	lastPos math2.Vector2d
	lastC float64
	lastWidth float64
	lastT float64
	lines []Linear
	sections []float64
}

func NewBezier(points []math2.Vector2d) *Bezier {
	bz := &Bezier{points: points, lastPos: points[0], lines: make([]Linear, 0), sections: make([]float64, 1)}

	pointLength := 0.0
	for i := 1; i < len(points); i++ {
		pointLength += points[i].Dst(points[i-1])
	}

	pointLength = math.Ceil(pointLength*30)

	/*for i := 1; i <= int(pointLength); i++ {
		bz.ApproxLength += bz.NPointAt(float64(i) / pointLength).Dst(bz.NPointAt(float64(i-1) / pointLength))
	}*/

	points1 := NewBezierApproximator(points).CreateBezier()

	println(len(points1))

	for i:= 1; i < len(points1); i++ {
		bz.lines = append(bz.lines, NewLinear(points1[i-1], points1[i]))
		bz.sections = append(bz.sections, points1[i-1].Dst(points1[i]))
		bz.sections[len(bz.sections)-1] += bz.sections[len(bz.sections)-2]
	}

	/*println(bz.NPointAt(0).Dst(bz.NPointAt(1.0/pointLength)))

	previous := bz.NPointAt(0)

	for p := 0.0; p < 1.0; p += 1/pointLength  {
		currentPoint := bz.NPointAt(p)

		println(previous.DstSq(bz.NPointAt(p+1/pointLength)))
		if previous.DstSq(bz.NPointAt(p+1/pointLength)) >= BEZIER_QUANTIZATIONSQ {
			bz.lines = append(bz.lines, NewLinear(previous, currentPoint))
			bz.sections = append(bz.sections, previous.Dst(currentPoint))
			if len(bz.sections) > 1 {
				bz.sections[len(bz.sections)-1] += bz.sections[len(bz.sections)-2]
			}
			previous = currentPoint
		}

	}*/

	bz.ApproxLength = 0.0

	for _, l := range bz.lines  {
		bz.ApproxLength += l.GetLength()
	}

	/*for i := range bz.sections {
		bz.sections[i] /= bz.ApproxLength
	}*/

	return bz
}

func NewBezierNA(points []math2.Vector2d) *Bezier {
	bz := &Bezier{points: points, lastPos: points[0], lines: make([]Linear, 0), sections: make([]float64, 1)}
	bz.ApproxLength = 0.0
	return bz
}

func (bz *Bezier) NPointAt(t float64) math2.Vector2d {
	x := 0.0
	y := 0.0
	n := len(bz.points) - 1
	for i := 0; i <= n; i++ {
		b := bernstein(int64(i), int64(n), t)
		x += bz.points[i].X * b
		y += bz.points[i].Y * b
	}
	return math2.NewVec2d(x, y)
}

//It's not a neat solution, but it works
//This calculates point on bezier with constant velocity
func (bz *Bezier) PointAt(t float64) math2.Vector2d {
	desiredWidth := bz.ApproxLength * t

	lineI := len(bz.sections)-2

	for i, k := range bz.sections[:len(bz.sections)-2]  {
		if k <= desiredWidth {
			lineI = i
		}
	}

	//lineI := sort.SearchFloat64s(bz.sections[:len(bz.sections)-2], desiredWidth)

	//println(lineI, len(bz.sections), len(bz.lines))
	line := bz.lines[lineI]

	point := line.PointAt((desiredWidth-bz.sections[lineI])/(bz.sections[lineI+1]-bz.sections[lineI]))

	//width := b
	//pos := bz.lastPos
	//c := 0.0

	/*if desiredWidth == bz.lastWidth {
		return bz.lastPos
	} else if desiredWidth > bz.lastWidth {
		for bz.lastWidth < desiredWidth {
			pt := bz.NPointAt(bz.lastC)
			lsW := bz.lastWidth + pt.Dst(bz.lastPos)
			if lsW > desiredWidth {
				bz.lastC -= 1.0 / float64(bz.ApproxLength*2-1)
				return bz.lastPos
			}
			bz.lastWidth = lsW
			bz.lastPos = pt
			bz.lastC += 1.0 / float64(bz.ApproxLength*2-1)
		}
	} else {
		for bz.lastWidth > desiredWidth {
			pt := bz.NPointAt(bz.lastC)
			lsW := bz.lastWidth - pt.Dst(bz.lastPos)
			if lsW < desiredWidth {
				bz.lastC += 1.0 / float64(bz.ApproxLength*2-1)
				return bz.lastPos
			}
			bz.lastWidth = lsW
			bz.lastPos = pt
			bz.lastC -= 1.0 / float64(bz.ApproxLength*2-1)
		}
	}*/

	return point
}

func (bz Bezier) GetLength() float64 {
	return bz.ApproxLength
}

func (bz Bezier) GetStartAngle() float64 {
	return bz.lines[0].GetStartAngle()//bz.points[0].AngleRV(bz.NPointAt(1.0 / bz.ApproxLength))
}

func (bz Bezier) GetEndAngle() float64 {
	return bz.lines[len(bz.lines)-1].GetEndAngle()//bz.points[len(bz.points)-1].AngleRV(bz.NPointAt((bz.ApproxLength - 1) / bz.ApproxLength))
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func BinomialCoefficient(n, k int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}
	k = min(k, n-k)
	var c int64 = 1
	var i int64 = 0
	for ; i < k; i++ {
		c = c * (n - i) / (i + 1)
	}

	return c
}

func bernstein(i, n int64, t float64) float64 {
	return float64(BinomialCoefficient(n, i)) * math.Pow(t, float64(i)) * math.Pow(1.0-t, float64(n-i))
}

func calcLength() {

}

func (ln Bezier) GetPoints(num int) []math2.Vector2d {
	t0 := 1 / float64(num-1)

	points := make([]math2.Vector2d, num)
	t := 0.0
	for i := 0; i < num; i += 1 {
		points[i] = ln.PointAt(t)
		t += t0
	}

	return points
}

func (ln *Bezier) GetLines() []Linear {
	return ln.lines
}
