package aabb

import "testing"
import "math/rand"
import "time"
import crand "crypto/rand"

const RandomizeTests = true
const FixedSeed      = 0x540B982FD054D8E3 // only used if !RandomizeTests
var rng *rand.Rand

func init() {
	if RandomizeTests {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else {
		rng = rand.New(rand.NewSource(FixedSeed))
	}
}

func TestBitGridBox(t *testing.T) {
	bitGridBox := &bitGridBox {} // xBits, yBits, box
	bigBox := NewBox(0, 99, 0, 99)
	bitGridBox.box = bigBox
	bitGridBox.RefreshGridBits(bigBox, 10, 10)
	if bitGridBox.xBits != 0x000003FF {
		t.Fatalf("expected xBits to be 0x3FF, but got %X instead", bitGridBox.xBits)
	}
	if bitGridBox.yBits != 0x000003FF {
		t.Fatalf("expected yBits to be 0x3FF, but got %X instead", bitGridBox.yBits)
	}

	minBox := NewBox(0, 0, 0, 0)
	bitGridBox.box = minBox
	bitGridBox.RefreshGridBits(bigBox, 10, 10)
	if bitGridBox.xBits != 0x00000001 {
		t.Fatalf("expected xBits to be 0x01, but got %X instead", bitGridBox.xBits)
	}
	if bitGridBox.yBits != 0x00000001 {
		t.Fatalf("expected yBits to be 0x01, but got %X instead", bitGridBox.yBits)
	}

	mediumBox := NewBox(10, 19, 10, 19)
	bitGridBox.box = mediumBox
	bitGridBox.RefreshGridBits(bigBox, 10, 10)
	if bitGridBox.xBits != 0x00000002 {
		t.Fatalf("expected xBits to be 0x02, but got %X instead", bitGridBox.xBits)
	}
	if bitGridBox.yBits != 0x00000002 {
		t.Fatalf("expected yBits to be 0x02, but got %X instead", bitGridBox.yBits)
	}

	mediumBox = NewBox(10, 29, 10, 29)
	bitGridBox.box = mediumBox
	bitGridBox.RefreshGridBits(bigBox, 10, 10)
	if bitGridBox.xBits != 0x00000006 {
		t.Fatalf("expected xBits to be 0x06, but got %X instead", bitGridBox.xBits)
	}
	if bitGridBox.yBits != 0x00000006 {
		t.Fatalf("expected yBits to be 0x05, but got %X instead", bitGridBox.yBits)
	}
}

func getSeed() int64 {
	buffer := make([]byte, 4)
	n, err := crand.Read(buffer)
	if err != nil { panic(err.Error()) }
	if n != 4 { panic("getSeed crand.Read() -> n != 4") }
	seed := int64(0)
	seed = (seed << 8) | int64(buffer[0])
	seed = (seed << 8) | int64(buffer[1])
	seed = (seed << 8) | int64(buffer[2])
	seed = (seed << 8) | int64(buffer[3])
	return seed
}

// func TestStabilizeDebug(t *testing.T) {
// 	boxes := make([]MutableBox, 0, 4096)
// 	for numBoxes := 1; numBoxes < 4096; numBoxes += 93 {
// 		refSpace := NewBruteForce()
// 		space := NewBitGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
//
// 		seed := getSeed()
// 		rng := rand.New(rand.NewSource(seed))
//
// 		// add boxes
// 		for i := 0; i < numBoxes; i++ {
// 			box := testingNewBox99(rng)
// 			refSpace.Add(box)
// 			space.Add(box)
// 		}
//
// 		refCount := testingCountAllCollisions(refSpace)
// 		count    := testingCountAllCollisions(space)
// 		if refCount != count {
// 			t.Fatalf("refCount (%d) != count (%d) at numBoxes = %d, seed = %d", refCount, count, numBoxes, seed)
// 		}
//
// 		for n := 0; n < 10; n++ {
// 			// mutate them!
// 			for _, box := range boxes {
// 				startX := rng.Intn(80)
// 				startY := rng.Intn(80)
// 				box.MutateBox(startX, startY, startX + rng.Intn(21), startY + rng.Intn(21))
// 			}
// 			space.Stabilize()
//
// 			refCount = testingCountAllCollisions(refSpace)
// 			count    = testingCountAllCollisions(space)
// 			if refCount != count {
// 				t.Fatalf("refCount (%d) != count (%d) at numBoxes = %d, seed = %d", refCount, count, numBoxes, seed)
// 			}
// 		}
// 	}
// }

func TestBitGridSingleCollision(t *testing.T) {
	space := NewBitGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	testProcSingleCollision(t, space)
}

func TestBitGridMutateVsStabilize1200(t *testing.T) {
	space1 := NewBitGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	space2 := NewBitGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	testProcMutateVsStabilize1200(t, space1, space2)
}

func TestBitGridStabilize600(t *testing.T) {
	space := NewBitGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	testProcStabilizeN(t, space, 600)
}

func TestBitGridNeutralStabilize3000(t *testing.T) {
	space1 := NewBitGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	space2 := NewBitGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	testProcNeutralStabilizeN(t, space1, space2, 3000)
}

// --- benchmarks ---

func BenchmarkBitGrid500(b *testing.B) {
	bitGridSpace := NewBitGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	benchmarkBasicSpace500(b, bitGridSpace)
}

func BenchmarkBitGridStretch500(b *testing.B) {
	bitGridSpace := NewBitGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	benchmarkBasicSpaceStretch500(b, bitGridSpace)
}

func BenchmarkBitGrid2000(b *testing.B) {
	bitGridSpace := NewBitGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	benchmarkBasicSpace2000(b, bitGridSpace)
}

func BenchmarkBitGridQuarterMuts1000(b *testing.B) {
	bitGridSpace := NewBitGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	benchmarkSpaceQuarterMuts1000(b, bitGridSpace)
}

func BenchmarkBitGridStabilize2500(b *testing.B) {
	space := NewBitGrid(referenceTestArea, refTestAxisCells, refTestAxisCells)
	benchmarkSpaceStabilize2500(b, space)
}
