package aabb

import "testing"
import "strconv"
import "math/rand"
import "fmt"

func init() { if false { fmt.Print("making fmt import ok even if unused") } }

func benchmarkBasicSpace500(b *testing.B, space BasicSpace) {
	testingAddBoxes(space, 500)
	for n := 0; n < b.N; n++ {
		count := testingCountAllCollisions(space, referenceTestArea)
		if count != 2018 {
			panic("expected count 2018, got " + strconv.Itoa(count))
		}
	}
}

func benchmarkBasicSpaceStretch500(b *testing.B, space BasicSpace) {
	testingAddStretchedBoxes(space, 500)
	for n := 0; n < b.N; n++ {
		count := testingCountAllCollisions(space, referenceTestArea)
		if count != 6365 {
			panic("expected count 4934, got " + strconv.Itoa(count))
		}
	}
}

func benchmarkBasicSpaceHorz500(b *testing.B, space BasicSpace) {
	testingAddHorzBoxes(space, 500)
	for n := 0; n < b.N; n++ {
		count := testingCountAllCollisions(space, referenceTestAreaHorz)
		if count != 2039 {
			panic("expected count 2039, got " + strconv.Itoa(count))
		}
	}
}

func benchmarkBasicSpace2000(b *testing.B, space BasicSpace) {
	testingAddBoxes(space, 2000)
	for n := 0; n < b.N; n++ {
		count := testingCountAllCollisions(space, referenceTestArea)
		if count != 8114 {
			panic("expected count 8114, got " + strconv.Itoa(count))
		}
	}
}

func benchmarkSpaceQuarterMuts1000(b *testing.B, space Space) {
	boxes := testingAddAndGetBoxes(space, 1000)
	expectedResults := []int{4008, 4100, 4015, 4041, 3998, 3981}

	const seed = 0x7A33149CF
	var rng *rand.Rand
	rng = rand.New(rand.NewSource(seed))

	for n := 0; n < b.N; n++ {
		count := testingCountAllCollisions(space, referenceTestArea)
		if n < len(expectedResults) {
			if count != expectedResults[n] {
				panic("expected count " + strconv.Itoa(expectedResults[n]) + " at update " + strconv.Itoa(n) + ", but got " + strconv.Itoa(count))
			}
		}
		for i := 0; i < 250; i++ {
			targetBox := rng.Intn(1000)
			newXMin, newYMin := rng.Intn(80), rng.Intn(80)
			newXMax, newYMax := newXMin + rng.Intn(21), newYMin + rng.Intn(21)
			space.MutateBox(boxes[targetBox], newXMin, newXMax, newYMin, newYMax)
		}
	}
}

func benchmarkSpaceStabilize2500(b *testing.B, space Space) {
	const NumBoxes = 2500
	boxes := testingAddAndGetBoxes(space, NumBoxes)
	expectedResults := []int{9984, 9915, 9974, 10100, 9977, 9935}

	const seed = 0xBFA491283A
	var rng *rand.Rand
	rng = rand.New(rand.NewSource(seed))

	for n := 0; n < b.N; n++ {
		count := testingCountAllCollisions(space, referenceTestArea)
		if n < len(expectedResults) {
			if count != expectedResults[n] {
				panic("expected count " + strconv.Itoa(expectedResults[n]) + " at update " + strconv.Itoa(n) + ", but got " + strconv.Itoa(count))
			}
		}

		for _, box := range boxes {
			startX := rng.Intn(80)
			startY := rng.Intn(80)
			box.MutateBox(startX, startX + rng.Intn(21), startY, startY + rng.Intn(21))
		}
		space.Stabilize()
	}
}
