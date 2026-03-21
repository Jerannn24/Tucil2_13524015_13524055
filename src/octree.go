package main

import (
	"math"
	"os"
	"bufio"
	"fmt"
)
type Point3D struct {
	X, Y, Z float64
}

type Boundary struct {
	Center Point3D 
	Half   float64 // Setengah dari panjang sisi kubus
}

type OctreeNode struct {
    Box Boundary
    Children [8]*OctreeNode
    IsLeaf bool
    Depth int
}

type Octree struct {
    Root *OctreeNode
    MaxDepth int
    TotalLeaf int
	LeafList []Boundary //Nyimpan semua boundary, biar gampang di eksport balek
	NodesCount []int
	NodesSkipped []int
}


//Membuat boundaries untuk child 
// Indeks 0: Kiri, Bawah, Belakang
	// Indeks 1: Kanan, Bawah, Belakang
	// Indeks 2: Kiri, Atas, Belakang
	// Indeks 3: Kanan, Atas, Belakang
	// Indeks 4: Kiri, Bawah, Depan
	// Indeks 5: Kanan, Bawah, Depan
	// Indeks 6: Kiri, Atas, Depan
	// Indeks 7: Kanan, Atas, Depan

func (b Boundary) GetChildBoundaries() [8]Boundary {
	newHalf := b.Half / 2
	var children [8]Boundary
	children[0] = Boundary{
		Center: Point3D{b.Center.X - newHalf, b.Center.Y - newHalf, b.Center.Z - newHalf},
		Half:   newHalf,
	}
	children[1] = Boundary{
		Center: Point3D{b.Center.X + newHalf, b.Center.Y - newHalf, b.Center.Z - newHalf},
		Half:   newHalf,
	}
	children[2] = Boundary{
		Center: Point3D{b.Center.X - newHalf, b.Center.Y + newHalf, b.Center.Z - newHalf},
		Half:   newHalf,
	}
	children[3] = Boundary{
		Center: Point3D{b.Center.X + newHalf, b.Center.Y + newHalf, b.Center.Z - newHalf},
		Half:   newHalf,
	}
	children[4] = Boundary{
		Center: Point3D{b.Center.X - newHalf, b.Center.Y - newHalf, b.Center.Z + newHalf},
		Half:   newHalf,
	}
	children[5] = Boundary{
		Center: Point3D{b.Center.X + newHalf, b.Center.Y - newHalf, b.Center.Z + newHalf},
		Half:   newHalf,
	}
	children[6] = Boundary{
		Center: Point3D{b.Center.X - newHalf, b.Center.Y + newHalf, b.Center.Z + newHalf},
		Half:   newHalf,
	}
	children[7] = Boundary{
		Center: Point3D{b.Center.X + newHalf, b.Center.Y + newHalf, b.Center.Z + newHalf},
		Half:   newHalf,
	}

	return children
}


func (o *Octree) Build(currentBox Boundary, faces []faces, depth int) *OctreeNode {
    if !o.IsIntersecting(currentBox, faces) {
        o.NodesSkipped[depth]++
        return nil
    }

    node := &OctreeNode{
        Box:   currentBox,
        Depth: depth,
    }
    o.NodesCount[depth]++

    if depth == o.MaxDepth {
        node.IsLeaf = true
        o.TotalLeaf++
        
        o.LeafList = append(o.LeafList, currentBox)
        return node
    }

    childBoundaries := currentBox.GetChildBoundaries()
    
    for i := 0; i < 8; i++ {
        node.Children[i] = o.Build(childBoundaries[i], faces, depth+1)
    }

    return node
}

func (o *Octree) IsIntersecting(box Boundary, faces []faces) bool {
    for _, f := range faces {
        fMinX := math.Min(f.Vertices[0].x, math.Min(f.Vertices[1].x, f.Vertices[2].x))
        fMaxX := math.Max(f.Vertices[0].x, math.Max(f.Vertices[1].x, f.Vertices[2].x))

        fMinY := math.Min(f.Vertices[0].y, math.Min(f.Vertices[1].y, f.Vertices[2].y))
        fMaxY := math.Max(f.Vertices[0].y, math.Max(f.Vertices[1].y, f.Vertices[2].y))

        fMinZ := math.Min(f.Vertices[0].z, math.Min(f.Vertices[1].z, f.Vertices[2].z))
        fMaxZ := math.Max(f.Vertices[0].z, math.Max(f.Vertices[1].z, f.Vertices[2].z))

        boxMinX := box.Center.X - box.Half
        boxMaxX := box.Center.X + box.Half

        boxMinY := box.Center.Y - box.Half
        boxMaxY := box.Center.Y + box.Half

        boxMinZ := box.Center.Z - box.Half
        boxMaxZ := box.Center.Z + box.Half

        if fMaxX >= boxMinX && fMinX <= boxMaxX && fMaxY >= boxMinY && fMinY <= boxMaxY && fMaxZ >= boxMinZ && fMinZ <= boxMaxZ {
            return true
        }
    }
    return false
}


func (o *Octree) ExportToOBJ(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	vertexOffset := 1
	
	for _, box := range o.LeafList {
		minX, maxX := box.Center.X - box.Half, box.Center.X + box.Half
		minY, maxY := box.Center.Y - box.Half, box.Center.Y + box.Half
		minZ, maxZ := box.Center.Z - box.Half, box.Center.Z + box.Half
		
		vertices := [8]Point3D{
			{minX, minY, minZ}, //Belakang Kiri Bawah
			{maxX, minY, minZ}, //Belakang Kanan Bawah
			{maxX, maxY, minZ}, //Belakang Kanan Atas
			{minX, maxY, minZ}, //Belakang Kiri Atas
			{minX, minY, maxZ}, //Depan Kiri Bawah
			{maxX, minY, maxZ}, //Depan Kanan Bawah
			{maxX, maxY, maxZ}, //Depan Kanan Atas
			{minX, maxY, maxZ}, //Depan Kiri Atas
		}
		for _, v := range vertices {
			fmt.Fprintf(writer, "v %.4f %.4f %.4f\n", v.X, v.Y, v.Z)
		}

		v := vertexOffset
		faces := [][3]int{
			{v, v + 1, v + 2}, {v, v + 2, v + 3},         //belakang
			{v + 4, v + 5, v + 6}, {v + 4, v + 6, v + 7}, //depan
			{v, v + 4, v + 5}, {v, v + 5, v + 1},         //bawah
			{v + 3, v + 2, v + 6}, {v + 3, v + 6, v + 7}, //atas
			{v, v + 3, v + 7}, {v, v + 7, v + 4},         //kiri
			{v + 1, v + 5, v + 6}, {v + 1, v + 6, v + 2}, //kanan
		}

		for _, f := range faces {
			fmt.Fprintf(writer, "f %d %d %d\n", f[0], f[1], f[2])
		}

		vertexOffset += 8
	}

	return writer.Flush()
}