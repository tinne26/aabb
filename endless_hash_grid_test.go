package aabb

import "testing"

func TestEndlessHashGridSingleCollision(t *testing.T) {
	space := NewEndlessHashGrid(refTestAxisCells, refTestAxisCells)
	testProcSingleCollision(t, space)
}

func TestEndlessHashGridMutateVsStabilize1200(t *testing.T) {
	space1 := NewEndlessHashGrid(refTestAxisCells, refTestAxisCells)
	space2 := NewEndlessHashGrid(refTestAxisCells, refTestAxisCells)
	testProcMutateVsStabilize1200(t, space1, space2)
}

func TestEndlessHashGridStabilize600(t *testing.T) {
	space := NewEndlessHashGrid(refTestAxisCells, refTestAxisCells)
	testProcStabilizeN(t, space, 600)
}

func TestEndlessHashGridNeutralStabilize3000(t *testing.T) {
	space1 := NewEndlessHashGrid(refTestAxisCells, refTestAxisCells)
	space2 := NewEndlessHashGrid(refTestAxisCells, refTestAxisCells)
	testProcNeutralStabilizeN(t, space1, space2, 3000)
}

// --- benchmarks ---

func BenchmarkEndlessHashGrid500(b *testing.B) {
	space := NewEndlessHashGrid(refTestAxisCells, refTestAxisCells)
	benchmarkBasicSpace500(b, space)
}

func BenchmarkEndlessHashGridStretch500(b *testing.B) {
	space := NewEndlessHashGrid(refTestAxisCells, refTestAxisCells)
	benchmarkBasicSpaceStretch500(b, space)
}

func BenchmarkEndlessHashGridHorz500(b *testing.B) {	
	benchmarkBasicSpaceHorz500(b, NewEndlessHashGrid(refTestAxisCells/3, refTestAxisCells*2))
}

func BenchmarkEndlessHashGrid2000(b *testing.B) {
	space := NewEndlessHashGrid(refTestAxisCells, refTestAxisCells)
	benchmarkBasicSpace2000(b, space)
}

func BenchmarkEndlessHashGridQuarterMuts1000(b *testing.B) {
	space := NewEndlessHashGrid(refTestAxisCells, refTestAxisCells)
	benchmarkSpaceQuarterMuts1000(b, space)
}

func BenchmarkEndlessHashGridStabilize2500(b *testing.B) {
	space := NewEndlessHashGrid(refTestAxisCells, refTestAxisCells)
	benchmarkSpaceStabilize2500(b, space)
}
