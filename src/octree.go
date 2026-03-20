package main

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


