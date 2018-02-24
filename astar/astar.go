package astar

import (
	"fmt"
	"math"
	"sort"
)

// A* 算法设置

type RectType uint8

const (
	RectType_Pass RectType = iota
	RectType_NoPass
)

type Rect struct {
	X    int
	Y    int
	Type RectType
}

var mapData [8][6]Rect

func Init() {

	for i := 0; i < 8; i++ {
		for j := 0; j < 6; j++ {
			mapData[i][j] = Rect{X: i, Y: j}
			if i == 4 && (j >= 1 && j <= 3) {
				mapData[i][j].Type = RectType_NoPass
			}
		}
	}
}

type Mesh struct {
	Pos       *Rect
	FatherPos *Rect
	F         int
	G         int
	H         int
}

type MeshList []*Mesh

func (m MeshList) Len() int {
	return len(m)
}

func (m MeshList) Less(i, j int) bool {
	return m[i].F < m[j].F
}

func (m MeshList) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m MeshList) Find(x, y int) *Mesh {
	for _, v := range m {
		if v.Pos.X == x && v.Pos.Y == y {
			return v
		}
	}
	return nil
}

func Delete(m *MeshList, x, y int) {
	for k, v := range *m {
		if v.Pos.X == x && v.Pos.Y == y {
			*m = append((*m)[:k], (*m)[k+1:]...)
			return
		}
	}
}

func FindPath(src, dst Rect) []Rect {
	srcMesh := &Mesh{
		Pos:       &src,
		FatherPos: &src,
		G:         0,
		H:         0,
		F:         0,
	}
	openMesh, closeMesh := make(MeshList, 0), make(MeshList, 0)
	openMesh = append(openMesh, srcMesh)

	for {
		if len(openMesh) <= 0 {
			break
		}

		// 如果找到目标点，退出
		if m := openMesh.Find(dst.X, dst.Y); m != nil {
			// fmt.Println("open Mesh contain dst", dst.X, dst.Y)
			break
		}

		curMesh := openMesh[0]

		closeMesh = append(closeMesh, curMesh)
		Delete(&openMesh, curMesh.Pos.X, curMesh.Pos.Y)

		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				if i == 0 && j == 0 {
					continue
				}
				x, y := curMesh.Pos.X+i, curMesh.Pos.Y+j

				// 边界判断
				if x < 0 || y < 0 || x >= 8 || y >= 6 {
					continue
				}

				if mapData[x][y].Type != RectType_Pass { // cannot pass
					continue
				}
				// 不能通过的物体，夹角也不能穿行
				if curMesh.Pos.X != x && curMesh.Pos.Y != y {
					if mapData[x][curMesh.Pos.Y].Type != RectType_Pass || mapData[curMesh.Pos.X][y].Type != RectType_Pass {
						continue
					}
				}

				if m := closeMesh.Find(x, y); m != nil {
					continue
				}

				// 如果周围的点，已经在open列表中，比较是否从当前点过去更近
				if m := openMesh.Find(x, y); m != nil {
					gValue := 14
					if curMesh.Pos.Y == y || curMesh.Pos.X == x {
						gValue = 10
					}
					gValue += curMesh.G

					if gValue < m.G {
						m.FatherPos = curMesh.Pos
						m.G = gValue
						m.F = m.G + m.H
					}
					continue
				}

				// 如果周围的点, 不在open列表中，添加到open列表中，把当前节点当作他们的父节点
				mesh := &Mesh{
					Pos: &Rect{
						X: x,
						Y: y,
					},
					FatherPos: curMesh.Pos,
				}

				gValue := 14
				if mesh.FatherPos.Y == y || mesh.FatherPos.X == x {
					gValue = 10
				}
				gValue += curMesh.G

				hValue := math.Abs(float64(dst.X-x)) + math.Abs(float64(dst.Y-y))

				mesh.G = gValue
				mesh.H = int(hValue)
				mesh.F = mesh.G + mesh.H

				openMesh = append(openMesh, mesh)
			}
		}
		sort.Sort(openMesh)
	}

	// 根据目标点，逆向查询父节点，直到找到起始点
	dstMesh := openMesh.Find(dst.X, dst.Y)
	if dstMesh == nil {
		return nil
	}
	path := []Rect{}
	path = append(path, Rect{X: dstMesh.Pos.X, Y: dstMesh.Pos.Y})
	for {
		fPos := dstMesh.FatherPos
		path = append(path, Rect{X: fPos.X, Y: fPos.Y})

		if fPos.X == src.X && fPos.Y == src.Y {
			break
		}
		dstMesh = openMesh.Find(fPos.X, fPos.Y)
		if dstMesh == nil {
			dstMesh = closeMesh.Find(fPos.X, fPos.Y)
			if dstMesh == nil {
				fmt.Println("czx@@@ cannot find in openMesh and closeMesh", fPos.X, fPos.Y)
				break
			}
		}
	}

	return path
}

// [2,2], [3, 3], [3, 4], [4, 4], [5, 4], [6, 3], [6, 2]
func RunAstar() {
	Init()

	path := FindPath(Rect{X: 2, Y: 2}, Rect{X: 6, Y: 2})
	_ = path

	for _, v := range path {
		fmt.Printf("Path [%v, %v] ", v.X, v.Y)
	}
}
