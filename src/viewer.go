package main

import (
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	SCREEN_WIDTH  = 1200
	SCREEN_HEIGHT = 800
)

type mesh struct {
	vertices []vertice
	faces    []face
}
type game struct {
	whitePixel *ebiten.Image

	meshe mesh
	yaw   float64
	pitch float64

	distance float64

	lastX, lastY float64
	dragging     bool
}

type point3D struct {
	X, Y, Z float64
}

type shadedTriangle struct {
	p1, p2, p3 [2]float64
	depth      float64
	intensity  float64
}

func sub(a, b point3D) point3D {
	return point3D{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func dot(a, b point3D) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func cross(a, b point3D) point3D {
	return point3D{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func normalize(v point3D) point3D {
	length := math.Sqrt(dot(v, v))
	if length == 0 {
		return point3D{}
	}
	return point3D{v.X / length, v.Y / length, v.Z / length}
}

func buildMesh(vertices []vertice, faces []face) mesh {
	m := mesh{}

	for _, v := range vertices {
		m.vertices = append(m.vertices, vertice{v.x, v.y, v.z})
	}
	for _, f := range faces {
		m.faces = append(m.faces, face{
			VertexIndices: f.VertexIndices,
		})
	}

	return m
}

func project(p point3D) (float64, float64) {
	if p.Z <= 0 {
		return 0, 0
	}
	return p.X / p.Z, p.Y / p.Z
}

func toScreen(x, y float64) (float64, float64) {
	screenX := (x + 1) * SCREEN_WIDTH / 2
	screenY := (1 - y) * SCREEN_HEIGHT / 2
	return screenX, screenY
}

func rotateY(p point3D, angle float64) point3D {
	cosA := math.Cos(angle)
	sinA := math.Sin(angle)

	return point3D{
		X: p.X*cosA - p.Z*sinA,
		Y: p.Y,
		Z: p.X*sinA + p.Z*cosA,
	}
}
func rotateX(p point3D, angle float64) point3D {
	cosA := math.Cos(angle)
	sinA := math.Sin(angle)

	return point3D{
		X: p.X,
		Y: p.Y*cosA - p.Z*sinA,
		Z: p.Y*sinA + p.Z*cosA,
	}
}
func (g *game) Update() error {
	x, y := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if g.dragging {
			dx := float64(x) - g.lastX
			dy := float64(y) - g.lastY

			g.yaw += dx * 0.01
			g.pitch += dy * 0.01
		}

		g.dragging = true
	} else {
		g.dragging = false
	}

	g.lastX = float64(x)
	g.lastY = float64(y)

	_, wheelY := ebiten.Wheel()

	if wheelY != 0 {
		g.distance *= 1.0 - float64(wheelY)*0.1
	}

	maxPitch := math.Pi/2 - 0.05
	if g.pitch > maxPitch {
		g.pitch = maxPitch
	} else if g.pitch < -maxPitch {
		g.pitch = -maxPitch
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	if g.whitePixel == nil {
		white := ebiten.NewImage(1, 1)
		white.Fill(color.White)
		g.whitePixel = white
	}

	lightDir := normalize(point3D{0.5, 0.8, -0.6})
	rotated := make([]point3D, len(g.meshe.vertices))

	for i, v := range g.meshe.vertices {
		r := rotateY(point3D{v.x, v.y, v.z}, g.yaw)
		r = rotateX(r, g.pitch)

		r.Z += g.distance

		rotated[i] = r
	}

	projected := make([][2]float64, len(rotated))
	for i, v := range rotated {
		if v.Z <= 0 {
			continue
		}
		x, y := project(point3D{v.X, v.Y, v.Z})
		sx, sy := toScreen(x, y)

		projected[i] = [2]float64{sx, sy}
	}

	triangles := make([]shadedTriangle, 0, len(g.meshe.faces)*2)

	for _, f := range g.meshe.faces {
		if len(f.VertexIndices) < 3 {
			continue
		}

		for i := 1; i < len(f.VertexIndices); i++ {
			a := f.VertexIndices[0]
			b := f.VertexIndices[i]
			c := f.VertexIndices[(i+1)%len(f.VertexIndices)]

			if rotated[a].Z <= 0 || rotated[b].Z <= 0 || rotated[c].Z <= 0 {
				continue
			}

			a3 := rotated[a]
			b3 := rotated[b]
			c3 := rotated[c]

			ab := sub(b3, a3)
			ac := sub(c3, a3)

			normal := normalize(cross(ab, ac))
			toCamera := normalize(point3D{-a3.X, -a3.Y, -a3.Z})
			if dot(normal, toCamera) < 0 {
				normal = point3D{-normal.X, -normal.Y, -normal.Z}
			}

			intensity := dot(normal, lightDir)
			if intensity < 0 {
				intensity = 0
			}
			intensity = 0.2 + 0.8*intensity
			depth := (a3.Z + b3.Z + c3.Z) / 3.0

			p1 := projected[a]
			p2 := projected[b]
			p3 := projected[c]

			triangles = append(triangles, shadedTriangle{
				p1:        p1,
				p2:        p2,
				p3:        p3,
				depth:     depth,
				intensity: intensity,
			})
		}
	}

	sort.Slice(triangles, func(i, j int) bool {
		return triangles[i].depth > triangles[j].depth
	})

	for _, t := range triangles {
		col := float32(t.intensity)

		vertices := []ebiten.Vertex{
			{DstX: float32(t.p1[0]), DstY: float32(t.p1[1]), SrcX: 0, SrcY: 0, ColorR: col, ColorG: col, ColorB: col, ColorA: 1},
			{DstX: float32(t.p2[0]), DstY: float32(t.p2[1]), SrcX: 0, SrcY: 0, ColorR: col, ColorG: col, ColorB: col, ColorA: 1},
			{DstX: float32(t.p3[0]), DstY: float32(t.p3[1]), SrcX: 0, SrcY: 0, ColorR: col, ColorG: col, ColorB: col, ColorA: 1},
		}

		indices := []uint16{0, 1, 2}
		screen.DrawTriangles(vertices, indices, g.whitePixel, nil)
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

// func main() {
// 	faces, vertices, _, err := parseObj("../output/tesr.obj")
// 	if err != nil {
// 		panic(err)
// 	}

// 	mesh := buildMesh(vertices, faces)

// 	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
// 	ebiten.SetWindowTitle("3D Viewer")
// 	ebiten.RunGame(&game{distance: 10, meshe: mesh})
// }
