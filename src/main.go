package main

import (
	"bufio"
	"fmt"
	"os"
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
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path_to_obj_file>")
		return
	}

	path := os.Args[1]
	faces, vertices, err := parseObj("../test/" + path)

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

}

func parseObj(path string) ([]faces, []vertices, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var verticess []vertices
	var facess []faces

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
				return nil, nil, fmt.Errorf("invalid vertex format at line %d: %s", lineNum, line)
			}

			var vertex vertices
			_, err := fmt.Sscanf(line, "v %f %f %f", &vertex.x, &vertex.y, &vertex.z)
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing vertex at line %d: %v", lineNum, err)
			}

			verticess = append(verticess, vertex)
		case "f":
			if len(parts) < 4 {
				return nil, nil, fmt.Errorf("invalid face format at line %d: %s", lineNum, line)
			}

			var face faces
			for i := 1; i < len(parts); i++ {
				var vertexIndex int
				_, err := fmt.Sscanf(parts[i], "%d", &vertexIndex)
				if err != nil {
					return nil, nil, fmt.Errorf("error parsing face at line %d: %v", lineNum, err)
				}
				if vertexIndex < 1 || vertexIndex > len(verticess) {
					return nil, nil, fmt.Errorf("vertex index out of range at line %d: %d", lineNum, vertexIndex)
				}

				face.Vertices = append(face.Vertices, verticess[vertexIndex-1])
			}
			facess = append(facess, face)
		default:
			return nil, nil, fmt.Errorf("unsupported line type at line %d: %s", lineNum, line)
		}
	}
	return facess, verticess, nil
}
