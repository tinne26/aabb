package aabb

import "testing"

func TestHashGridSingleCollision(t *testing.T) {
	space := NewHashGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	testProcSingleCollision(t, space)
}

func TestHashGridMutateVsStabilize1200(t *testing.T) {
	space1 := NewHashGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	space2 := NewHashGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	testProcMutateVsStabilize1200(t, space1, space2)
}

func TestHashGridStabilize600(t *testing.T) {
	space := NewHashGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	testProcStabilizeN(t, space, 600)
}

func TestHashGridNeutralStabilize3000(t *testing.T) {
	space1 := NewHashGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	space2 := NewHashGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	testProcNeutralStabilizeN(t, space1, space2, 3000)
}

// --- benchmarks ---

func BenchmarkHashGrid500(b *testing.B) {
	space := NewHashGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	benchmarkBasicSpace500(b, space)
}

func BenchmarkHashGridStretch500(b *testing.B) {
	space := NewHashGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	benchmarkBasicSpaceStretch500(b, space)
}

func BenchmarkHashGrid2000(b *testing.B) {
	space := NewHashGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	benchmarkBasicSpace2000(b, space)
}

func BenchmarkHashGridQuarterMuts1000(b *testing.B) {
	space := NewHashGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	benchmarkSpaceQuarterMuts1000(b, space)
}

func BenchmarkHashGridStabilize2500(b *testing.B) {
	space := NewHashGrid(referenceTestArea, 5, 5)
	benchmarkSpaceStabilize2500(b, space)
}
