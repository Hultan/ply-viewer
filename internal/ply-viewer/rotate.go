package ply_viewer

import (
	"math"

	"github.com/hultan/ply-viewer/internal/ply"
)

func (m *MainForm) rotateObject() {
	for _, vertex := range object.Vertexes {
		m.rotatePoint(vertex, 0.1, 0.1, 0.1)
	}
}

func (m *MainForm) rotatePoint(v *ply.Vertex, thetaX, thetaY, thetaZ float64) {
	// v = m.rotatePointAxis(v, m.getRotationMatrix(XAxis, thetaX))
	// v = m.rotatePointAxis(v, m.getRotationMatrix(YAxis, thetaY))
	v = m.rotatePointAxis(v, m.getRotationMatrix(ZAxis, thetaZ))
}

func (m *MainForm) rotatePointAxis(p *ply.Vertex, rm RotationMatrix) *ply.Vertex {
	x := rm[0][0]*p.X + rm[0][1]*p.Y + rm[0][2]*p.Z
	y := rm[1][0]*p.X + rm[1][1]*p.Y + rm[1][2]*p.Z
	z := rm[2][0]*p.X + rm[2][1]*p.Y + rm[2][2]*p.Z
	p.X = x
	p.Y = y
	p.Z = z
	return p
}

func (m *MainForm) getRotationMatrix(a Axis, theta float64) RotationMatrix {
	var rm RotationMatrix
	switch a {
	case XAxis:
		rm = m.createRotateX(theta)
	case YAxis:
		rm = m.createRotateY(theta)
	case ZAxis:
		rm = m.createRotateZ(theta)
	}
	return rm
}

func (m *MainForm) createRotateX(theta float64) RotationMatrix {
	var rm RotationMatrix
	rm[0][0] = 1
	rm[0][1] = 0
	rm[0][2] = 0

	rm[1][0] = 0
	rm[1][1] = float32(math.Cos(theta))
	rm[1][2] = float32(-math.Sin(theta))

	rm[2][0] = 0
	rm[2][1] = float32(math.Sin(theta))
	rm[2][2] = float32(math.Cos(theta))

	return rm
}

func (m *MainForm) createRotateY(theta float64) RotationMatrix {
	var rm RotationMatrix
	rm[0][0] = float32(math.Cos(theta))
	rm[0][1] = 0
	rm[0][2] = float32(math.Sin(theta))

	rm[1][0] = 0
	rm[1][1] = 1
	rm[1][2] = 0

	rm[2][0] = float32(-math.Sin(theta))
	rm[2][1] = 0
	rm[2][2] = float32(math.Cos(theta))

	return rm
}

func (m *MainForm) createRotateZ(theta float64) RotationMatrix {
	var rm RotationMatrix
	rm[0][0] = float32(math.Cos(theta))
	rm[0][1] = float32(-math.Sin(theta))
	rm[0][2] = 0

	rm[1][0] = float32(math.Sin(theta))
	rm[1][1] = float32(math.Cos(theta))
	rm[1][2] = 0

	rm[2][0] = 0
	rm[2][1] = 0
	rm[2][2] = 1

	return rm
}
