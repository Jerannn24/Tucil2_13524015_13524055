# Tucil2_13524015_13524055

| Name                    | NIM      |
|-------------------------|----------|
| Mahatma Brahmana        | 13524015 |
| Junior Natra Situmorang | 13524055 |

## 3DOctree Voxelizer

Halo! Ini adalah tugas kecil 2 (Tucil 2) untuk mata kuliah Strategi Algoritma (STIMA) 2024.

### Deskripsi Singkat
Program ini membaca file 3D model berformat OBJ, membangun boundary box, lalu membagi ruang menggunakan Octree hingga kedalaman tertentu. Hasil akhirnya adalah file OBJ baru berisi voxel (kotak-kotak kecil) yang merepresentasikan model asli. Ada juga fitur viewer 3D sederhana untuk melihat hasilnya.

### Cara Pakai
1. **Clone repo ini**
2. **Jalankan program**
   ```bash
   go run src/main.go src/octree.go src/viewer.go <path_to_obj_file> <max_depth>
   ```
   Contoh:
   ```bash
   go run src/main.go src/octree.go src/viewer.go cow.obj 3
   ```
3. **Ikuti instruksi di terminal**
   - Masukkan path output (misal: hasil.obj atau ../output/hasil.obj)
   - Pilih apakah ingin melihat hasilnya di viewer 3D

### Fitur
- Parsing file OBJ (hanya vertex & face)
- Boundary box otomatis
- Octree subdivision (depth bisa diatur)
- Output file OBJ voxelized
- Viewer 3D (pakai Ebiten)
- Error handling untuk path & file

### Struktur Folder
- `src/`   : kode utama (main.go, octree.go, viewer.go)
- `test/`  : file-file untuk testing
- `doc/`   : dokumentasi
- `bin/`   : executable
- `output/`: hasil build


### Author
- 13524015
- 13524055

### License
Pake aja bosq