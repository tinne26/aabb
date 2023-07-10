package aabb

import "testing"

func TestGridSingleCollision(t *testing.T) {
	space := NewGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	testProcSingleCollision(t, space)
}

func TestGridMutateVsStabilize1200(t *testing.T) {
	space1 := NewGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	space2 := NewGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	testProcMutateVsStabilize1200(t, space1, space2)
}

func TestGridStabilize600(t *testing.T) {
	space := NewGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	testProcStabilizeN(t, space, 600)
}

func TestGridNeutralStabilize3000(t *testing.T) {
	space1 := NewGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	space2 := NewGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	testProcNeutralStabilizeN(t, space1, space2, 3000)
}

// --- benchmarks ---

func BenchmarkGrid500(b *testing.B) {
	space := NewGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	benchmarkBasicSpace500(b, space)
}

func BenchmarkGridStretch500(b *testing.B) {
	space := NewGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	benchmarkBasicSpaceStretch500(b, space)
}

func BenchmarkGrid2000(b *testing.B) {
	space := NewGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	benchmarkBasicSpace2000(b, space)
}

func BenchmarkGridQuarterMuts1000(b *testing.B) {
	space := NewGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	benchmarkSpaceQuarterMuts1000(b, space)
}

func BenchmarkGridStabilize2500(b *testing.B) {
	space := NewGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	benchmarkSpaceStabilize2500(b, space)
}
