package curves

import (
	"github.com/wieku/danser-go/bmath"
)

const minPartWidth = 0.0001

type MultiCurve struct {
	sections []float32
	length   float32
	lines    []Linear
}

func NewMultiCurve(typ string, points []bmath.Vector2f, desiredLength float64) *MultiCurve {
	lines := make([]Linear, 0)

	if len(points) < 3 {
		typ = "L"
	}

	switch typ {
	case "P":
		lines = append(lines, ApproximateCircularArc(points[0], points[1], points[2], 0.125)...)
		break
	case "L":
		for i := 0; i < len(points)-1; i++ {
			lines = append(lines, NewLinear(points[i], points[i+1]))
		}
		break
	case "B":
		lastIndex := 0
		for i, p := range points {
			if (i == len(points)-1 && p != points[i-1]) || (i < len(points)-1 && points[i+1] == p) {
				pts := points[lastIndex : i+1]

				if len(pts) > 2 {
					lines = append(lines, ApproximateBezier(pts)...)
				} else if len(pts) == 1 {
					lines = append(lines, NewLinear(pts[0], pts[0]))
				} else {
					lines = append(lines, NewLinear(pts[0], pts[1]))
				}

				lastIndex = i + 1
			}
		}
		break
	case "C":

		if points[0] != points[1] {
			points = append([]bmath.Vector2f{points[0]}, points...)
		}

		if points[len(points)-1] != points[len(points)-2] {
			points = append(points, points[len(points)-1])
		}

		for i := 0; i < len(points)-3; i++ {
			lines = append(lines, ApproximateCatmullRom(points[i:i+4], 50)...)
		}
		break
	}

	length := float32(0.0)

	for _, l := range lines {
		length += l.GetLength()
	}

	if desiredLength >= 0 {
		diff := float64(length) - desiredLength

		for len(lines) > 0 {
			line := lines[len(lines)-1]

			if float64(line.GetLength()) > diff+minPartWidth {
				pt := line.PointAt((line.GetLength() - float32(diff)) / line.GetLength())
				lines[len(lines)-1] = NewLinear(line.Point1, pt)
				break
			}

			diff -= float64(line.GetLength())
			lines = lines[:len(lines)-1]
		}
	}

	length = 0.0

	for _, l := range lines {
		length += l.GetLength()
	}

	sections := make([]float32, len(lines)+1)
	sections[0] = 0.0
	prev := float32(0.0)

	for i := 0; i < len(lines); i++ {
		prev += lines[i].GetLength()
		sections[i+1] = prev
	}

	return &MultiCurve{sections, length, lines}
}

func (sa *MultiCurve) PointAt(t float32) bmath.Vector2f {

	desiredWidth := sa.length * t

	lineI := len(sa.sections) - 2

	for i, k := range sa.sections[:len(sa.sections)-1] {
		if k <= desiredWidth {
			lineI = i
		}
	}

	line := sa.lines[lineI]

	return line.PointAt((desiredWidth - sa.sections[lineI]) / (sa.sections[lineI+1] - sa.sections[lineI]))
}

func (sa *MultiCurve) GetLength() float32 {
	return sa.length
}

func (sa *MultiCurve) GetStartAngle() float32 {
	return sa.lines[0].GetStartAngle()
}

func (sa *MultiCurve) GetEndAngle() float32 {
	return sa.lines[len(sa.lines)-1].GetEndAngle()
}

func (ln *MultiCurve) GetLines() []Linear {
	return ln.lines
}
