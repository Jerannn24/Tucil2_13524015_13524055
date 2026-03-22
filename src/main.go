package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type face struct {
	VertexIndices []int
}

type vertice struct {
	x, y, z float64
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run main.go octree.go viewer.go <path_to_obj_file> <max_depth>")
		return
	}

	max_depth := os.Args[4]
	maxDepth, err := strconv.Atoi(max_depth)
	if err != nil {
		fmt.Println("max_depth argument has to be a number!")
		return
	}

	path := os.Args[1]
	faces, vertices, mainBox, err := parseObj("../test/" + path)

	if err != nil {
		fmt.Printf("Error parsing OBJ file: %v\n", err)
		return
	}

	fmt.Printf("Parsed %d vertices and %d faces from the OBJ file.\n", len(vertices), len(faces))
	// for i, vertex := range vertices {
	// 	fmt.Printf("Vertex %d: (%.2f, %.2f, %.2f)\n", i+1, vertex.x, vertex.y, vertex.z)
	// }

	// for i, face := range faces {
	// 	fmt.Printf("Face %d: ", i+1)
	// 	for _, vertexIndex := range face.VertexIndices {
	// 		vertex := vertices[vertexIndex]
	// 		fmt.Printf("(%.2f, %.2f, %.2f) ", vertex.x, vertex.y, vertex.z)
	// 	}
	// 	fmt.Println()
	// }

	// Print boundary box info
	if mainBox != nil {
		fmt.Printf("\nBoundary Box Center: (%.4f, %.4f, %.4f)\n", mainBox.Center.X, mainBox.Center.Y, mainBox.Center.Z)
		fmt.Printf("Half Box Length: %.4f\n", mainBox.Half)
	}

	octree := &Octree{}
	octree.MaxDepth = maxDepth
	octree.NodesCount = make([]int, maxDepth+1)
	octree.NodesSkipped = make([]int, maxDepth+1)
	octree.LeafList = []Boundary{}

	timestart := time.Now()
	octree.Root = octree.Build(*mainBox, faces, vertices, 0)
	timeEnd := time.Now()

	fmt.Println("\nOctree Construction Result : ")
	fmt.Printf("Total Voxels created: %d\n", octree.TotalLeaf)
	fmt.Printf("Total Vertex created: %d\n", len(octree.LeafList)*8)
	fmt.Printf("Total Faces created: %d\n\n", len(octree.LeafList)*12)

	for i := 0; i <= maxDepth; i++ {
		fmt.Printf("Depth %d: Created %d nodes, Skipped %d nodes\n", i, octree.NodesCount[i], octree.NodesSkipped[i])
	}

	fmt.Printf("Octree construction took %v seconds\n", timeEnd.Sub(timestart).Seconds())
	fmt.Printf("Octree Max Depth: %d\n", octree.MaxDepth)

	var outputPath string

	fmt.Printf("Masukkan path output .obj (contoh: hasil.obj atau ../output/hasil.obj): ")
	fmt.Scanf("%s", &outputPath)

	outputPath = strings.TrimSpace(outputPath)
	if !strings.HasSuffix(strings.ToLower(outputPath), ".obj") {
		outputPath += ".obj"
	}

	if filepath.Dir(outputPath) == "." {
		outputPath = filepath.Join("..", "output", outputPath)
	}

	err = octree.ExportToOBJ(outputPath)

	if err != nil {
		fmt.Printf("Error exporting: %v\n", err)
	} else {
		fmt.Printf("Voxels exported to %s successfully!\n", outputPath)
	}

	facess, verticess, _, err := parseObj(outputPath)
	if err != nil {
		fmt.Printf("Error parsing exported OBJ file: %v\n", err)
		return
	}

	mesh := buildMesh(verticess, facess)

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("3D Viewer")
	ebiten.RunGame(&game{distance: 10, meshe: mesh})
}

func parseObj(path string) ([]face, []vertice, *Boundary, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}
	defer file.Close()

	var verticess []vertice
	var facess []face
	var (
		maxX float64 = -1000000
		minX float64 = 1000000
		maxY float64 = -1000000
		minY float64 = 1000000
		maxZ float64 = -1000000
		minZ float64 = 1000000
	)

	scanner := bufio.NewScanner(file)

	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		switch parts[0] {
		case "v":
			if len(parts) != 4 {
				return nil, nil, nil, fmt.Errorf("invalid vertex format at line %d: %s", lineNum, line)
			}

			var vertex vertice
			_, err := fmt.Sscanf(line, "v %f %f %f", &vertex.x, &vertex.y, &vertex.z)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("error parsing vertex at line %d: %v", lineNum, err)
			}

			//Mencari nilai maksimum dan minimum masing masing sumbu
			if vertex.x > maxX {
				maxX = vertex.x
			}
			if vertex.x < minX {
				minX = vertex.x
			}

			if vertex.y > maxY {
				maxY = vertex.y
			}
			if vertex.y < minY {
				minY = vertex.y
			}

			if vertex.z > maxZ {
				maxZ = vertex.z
			}
			if vertex.z < minZ {
				minZ = vertex.z
			}

			verticess = append(verticess, vertex)
		case "f":
			if len(parts) < 4 {
				return nil, nil, nil, fmt.Errorf("invalid face format at line %d: %s", lineNum, line)
			}

			var face face
			for i := 1; i < len(parts); i++ {
				var vertexIndex int
				_, err := fmt.Sscanf(parts[i], "%d", &vertexIndex)
				if err != nil {
					return nil, nil, nil, fmt.Errorf("error parsing face at line %d: %v", lineNum, err)
				}
				if vertexIndex < 1 || vertexIndex > len(verticess) {
					return nil, nil, nil, fmt.Errorf("vertex index out of range at line %d: %d", lineNum, vertexIndex)
				}

				face.VertexIndices = append(face.VertexIndices, vertexIndex-1)
			}
			facess = append(facess, face)
		default:
			fmt.Printf("Warning: Unrecognized line format at line %d: %s\n", lineNum, line)
			continue
		}
	}

	// Print min/max
	// fmt.Printf("minX: %.4f, maxX: %.4f\n", minX, maxX)
	// fmt.Printf("minY: %.4f, maxY: %.4f\n", minY, maxY)
	// fmt.Printf("minZ: %.4f, maxZ: %.4f\n", minZ, maxZ)

	L := float64(math.Max(maxX-minX, math.Max(maxY-minY, maxZ-minZ)))

	center := Point3D{
		X: (minX + maxX) / 2,
		Y: (minY + maxY) / 2,
		Z: (minZ + maxZ) / 2,
	}

	boundary := &Boundary{
		Center: center,
		Half:   L / 2,
	}

	return facess, verticess, boundary, nil
}
