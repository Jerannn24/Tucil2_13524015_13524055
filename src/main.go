package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Obj struct {
	Name     string
	Vertices []vertices
	Faces    []faces
}

type faces struct {
	Vertices []vertices
}

type vertices struct {
	x, y, z float64
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <path_to_obj_file> <max_depth>")
		return
	}

	max_depth := os.Args[2]
	maxDepth, err := strconv.Atoi(max_depth)
	if err != nil {
		fmt.Println("max_depth argument has to be a number!")
		return
	}
	//tambahan aja, nanti dihapus, biar ga eror
	maxDepth = maxDepth + 0
	//sampai sini

	path := os.Args[1]
	faces, vertices, mainBox, err := parseObj("../test/" + path)

	if err != nil {
		fmt.Printf("Error parsing OBJ file: %v\n", err)
		return
	}

	fmt.Printf("Parsed %d vertices and %d faces from the OBJ file.\n", len(vertices), len(faces))
	for i, vertex := range vertices {
		fmt.Printf("Vertex %d: (%.2f, %.2f, %.2f)\n", i+1, vertex.x, vertex.y, vertex.z)
	}

	for i, face := range faces {
		fmt.Printf("Face %d: ", i+1)
		for _, vertex := range face.Vertices {
			fmt.Printf("(%.2f, %.2f, %.2f) ", vertex.x, vertex.y, vertex.z)
		}
		fmt.Println()
	}

	// Print boundary box info
	if mainBox != nil {
		fmt.Printf("\nBoundary Box Center: (%.4f, %.4f, %.4f)\n", mainBox.Center.X, mainBox.Center.Y, mainBox.Center.Z)
		fmt.Printf("Half Box Length: %.4f\n", mainBox.Half)
	}

}

func parseObj(path string) ([]faces, []vertices, *Boundary, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}
	defer file.Close()

	var verticess []vertices
	var facess []faces
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

			var vertex vertices
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

			var face faces
			for i := 1; i < len(parts); i++ {
				var vertexIndex int
				_, err := fmt.Sscanf(parts[i], "%d", &vertexIndex)
				if err != nil {
					return nil, nil, nil, fmt.Errorf("error parsing face at line %d: %v", lineNum, err)
				}
				if vertexIndex < 1 || vertexIndex > len(verticess) {
					return nil, nil, nil, fmt.Errorf("vertex index out of range at line %d: %d", lineNum, vertexIndex)
				}

				face.Vertices = append(face.Vertices, verticess[vertexIndex-1])
			}
			facess = append(facess, face)
		default:
			return nil, nil, nil, fmt.Errorf("unsupported line type at line %d: %s", lineNum, line)
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
