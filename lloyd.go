package main

import (
	"flag"
	"fmt"
	"github.com/gonum/matrix/mat64"
	"github.com/pzsz/voronoi"
	vu "github.com/pzsz/voronoi/utils"
	"io"
	"math/rand"
	"os"
	"sort"
)

var seed int64
var intensity int
var iterations int
var matrix bool

func PoissonVoronoi(rng *rand.Rand, n int) []voronoi.Vertex {
	vertices := make([]voronoi.Vertex, 0, n)
	for i := 0; i < n; i++ {
		v := voronoi.Vertex{rng.Float64(), rng.Float64()}
		vertices = append(vertices, v)
	}
	return vertices
}

type Histogram map[int]int
func (h Histogram)Count(n int) {
	c, ok := h[n]
	if ok {
		h[n] = c + 1
	} else {
		h[n] = 1
	}
}

func (h Histogram)Print(w io.Writer, pfx string, n int) {
	var keys []int
	for k := range h {
    keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
    fmt.Fprintf(w, "%s%d\t%f\n", pfx, k, float64(h[k]) / float64(n))
	}
}

type Polygons map[int]int
func (p Polygons)Transitions(q Polygons) mat64.Matrix {
	size := 0
	ct := make(map[int]map[int]int)
	// count up transitions made
	for id, n := range p {
		if n > size {
			size = n
		}
		m, ok := q[id]
		if !ok {
			panic("undefined transition")
		}
		if m > size {
			size = m
		}
		row, ok := ct[n]
		if !ok {
			row = make(map[int]int)
			ct[n] = row
		}
		c, ok := row[m]
		if ok {
			row[m] = c + 1
		} else {
			row[m] = 1
		}
	}

	// normalize to the total
	t := mat64.NewDense(size+1, size+1, nil)
	for n, row := range ct {
		sum := 0
		for _, c := range row {
			sum += c
		}
		for m, c := range row {
			t.Set(n, m, float64(c) / float64(sum))
		}
	}

	// make matrix stochastic
	for i := 0; i < size+1; i++ {
		row := t.RowView(i)
		sum := float64(0)
		for j := 0; j < size+1; j++ {
			if j != i {
				sum += row.At(j, 0)
			}
		}
		if sum > 0 {
			row.SetVec(i, row.At(i,0) - 1)
		}
	}
	return t
}

func init() {
	flag.Int64Var(&seed, "s", 0, "random seed")
	flag.IntVar(&intensity, "i", 1000, "intensity")
	flag.IntVar(&iterations, "n", 10, "iterations")
	flag.BoolVar(&matrix, "m", false, "output transition matrix")
	flag.Usage = func () {
		fmt.Fprintf(os.Stderr, "\nPoisson Voronoi\n===============\n\n")
		fmt.Fprintf(os.Stderr, "Generate distributions of polygons resulting from repeated\napplications of Lloyd's algorithm.\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}
}

func main() {
	flag.Parse()

	src := rand.NewSource(seed)
	rng := rand.New(src)

	vertices := PoissonVoronoi(rng, intensity)
	bbox := voronoi.NewBBox(0, 1, 0, 1)

	diagram := voronoi.ComputeDiagram(vertices, bbox, true)
	polygons := make(Polygons)

	for i := 0; i < iterations; i++ {
		hist := make(Histogram)
		newp := make(Polygons)
		for id, c := range diagram.Cells {
			nedges := len(c.Halfedges)
			hist.Count(nedges)
			newp[id] = nedges
		}

		if matrix && i > 0 {
			trans := polygons.Transitions(newp)
			ft := mat64.Formatted(trans, mat64.Squeeze())
			fmt.Fprintf(os.Stdout, "%v\n\n", ft)
		}
		polygons = newp

		if !matrix {
			pfx := fmt.Sprintf("%d\t", i)
			hist.Print(os.Stdout, pfx, intensity)
			fmt.Fprintf(os.Stdout, "\n")
		}
		vertices = vu.LloydRelaxation(diagram.Cells)
		diagram = voronoi.ComputeDiagram(vertices, bbox, true)
	}
}
