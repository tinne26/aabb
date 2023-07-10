package aabb

import "testing"
import "math/rand"

func TestQuadtreeSingleCollision(t *testing.T) {
	space := NewQuadtree(referenceTestArea, StdQuadtreeDepth, StdQuadtreeSplitCutoff)
	testProcSingleCollision(t, space)
}

func TestQuadtreeMutateDebug(t *testing.T) {
	space1 := NewQuadtree(referenceTestArea, StdQuadtreeDepth, StdQuadtreeSplitCutoff)
	space2 := NewBruteForce()

	const seed = 0x4728BB12D
	var rng *rand.Rand
	rng = rand.New(rand.NewSource(seed))

	numBoxes := 6
	boxes := make([]MutableBox, numBoxes)
	for i := 0; i < numBoxes; i++ {
		box := testingNewBox99(rng)
		boxes[i] = box
		space1.Add(box)
		space2.Add(box)
	}
	for i := 0; i < numBoxes; i++ {
		newBox := testingNewBox99(rng)
		space1.MutateBox(boxes[i], newBox.XMin(), newBox.XMax(), newBox.YMin(), newBox.YMax())
	}
	space2.Stabilize()

	count1 := testingCountAllCollisions(space1, referenceTestArea)
	count2 := testingCountAllCollisions(space2, referenceTestArea)
	if count1 != count2 {
		t.Fatalf("number of collision for spaces 1 and two differ: %d vs %d", count1, count2)
	}
}

func TestQuadtreeMutateVsStabilize1200(t *testing.T) {
	space1 := NewQuadtree(referenceTestArea, StdQuadtreeDepth, StdQuadtreeSplitCutoff)
	space2 := NewQuadtree(referenceTestArea, StdQuadtreeDepth, StdQuadtreeSplitCutoff)
	testProcMutateVsStabilize1200(t, space1, space2)
}

func TestQuadtreeStabilize600(t *testing.T) {
	space := NewQuadtree(referenceTestArea, StdQuadtreeDepth, StdQuadtreeSplitCutoff)
	testProcStabilizeN(t, space, 600)
}

func TestQuadtreeNeutralStabilize3000(t *testing.T) {
	space1 := NewQuadtree(referenceTestArea, StdQuadtreeDepth, StdQuadtreeSplitCutoff)
	space2 := NewQuadtree(referenceTestArea, StdQuadtreeDepth, StdQuadtreeSplitCutoff)
	testProcNeutralStabilizeN(t, space1, space2, 3000)
}

// --- benchmarks ---

func BenchmarkQuadtree500(b *testing.B) {
	quadtreeSpace := NewQuadtree(referenceTestArea, StdQuadtreeDepth, StdQuadtreeSplitCutoff)
	benchmarkBasicSpace500(b, quadtreeSpace)
}

func BenchmarkQuadtreeStretch500(b *testing.B) {
	quadtreeSpace := NewQuadtree(referenceTestArea, StdQuadtreeDepth, StdQuadtreeSplitCutoff)
	benchmarkBasicSpaceStretch500(b, quadtreeSpace)
}

func BenchmarkQuadtreeHorz500(b *testing.B) {	
	space := NewQuadtree(referenceTestAreaHorz, StdQuadtreeDepth, StdQuadtreeSplitCutoff)
	benchmarkBasicSpaceHorz500(b, space)
}

func BenchmarkQuadtree2000(b *testing.B) {
	quadtreeSpace := NewQuadtree(referenceTestArea, StdQuadtreeDepth, StdQuadtreeSplitCutoff)
	benchmarkBasicSpace2000(b, quadtreeSpace)
}

func BenchmarkQuadtreeQuarterMuts1000(b *testing.B) {
	quadtreeSpace := NewQuadtree(referenceTestArea, StdQuadtreeDepth, StdQuadtreeSplitCutoff)
	benchmarkSpaceQuarterMuts1000(b, quadtreeSpace)
}

func BenchmarkQuadtreeStabilize2500(b *testing.B) {
	space := NewQuadtree(referenceTestArea, StdQuadtreeDepth, StdQuadtreeSplitCutoff)
	benchmarkSpaceStabilize2500(b, space)
}
