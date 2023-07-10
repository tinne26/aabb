package aabb

import "testing"
import "math/rand"

func testProcSingleCollision(t *testing.T, space BasicSpace) {
	space.Add(referenceTestBox)
	for nth, test := range singleNoCollisionTests {
		if space.Collision(NewBox(test.xmin, test.xmax, test.ymin, test.ymax)) != nil {
			t.Fatalf("singleNoCollisionTests#%d (found collision when none expected)", nth)
		}
	}

	for nth, test := range singleYesCollisionTests {
		if space.Collision(NewBox(test.xmin, test.xmax, test.ymin, test.ymax)) == nil {
			t.Fatalf("singleYesCollisionTests#%d (found no collision when one was expected)", nth)
		}
	}
}

func testProcMutateVsStabilize1200(t *testing.T, space1 Space, space2 Space) {
	const seed = 0xC1024A5A83
	var rng *rand.Rand
	rng = rand.New(rand.NewSource(seed))

	boxes := make([]MutableBox, 0, 1200)
	for i := 0; i < 1200; i++ {
		box := testingNewBox99(rng)
		space1.Add(box)
		space2.Add(box)
		boxes = append(boxes, box)
	}

	count1 := testingCountAllCollisions(space1, referenceTestArea)
	count2 := testingCountAllCollisions(space2, referenceTestArea)
	if count1 != count2 {
		t.Fatalf("space1 got %d collisions, but space2 got %d", count1, count2)
	}
	if count1 != 4770 {
		t.Fatalf("Expected 4770 collisions on space1, but got %d", count1)
	}

	for _, box := range boxes {
		newBox := testingNewBox99(rng)
		space1.MutateBox(box, newBox.XMin(), newBox.XMax(), newBox.YMin(), newBox.YMax())
	}
	space2.Stabilize()

	count1 = testingCountAllCollisions(space1, referenceTestArea)
	count2 = testingCountAllCollisions(space2, referenceTestArea)
	if count1 != count2 {
		t.Fatalf("space1 got %d collisions, but space2 got %d", count1, count2)
	}
	if count1 != 4859 {
		t.Fatalf("expected 4859 collisions on space1, but got %d", count1)
	}

	for _, box := range boxes {
		newBox := testingNewBox99(rng)
		space2.MutateBox(box, newBox.XMin(), newBox.XMax(), newBox.YMin(), newBox.YMax())
	}
	space1.Stabilize()

	count1 = testingCountAllCollisions(space1, referenceTestArea)
	count2 = testingCountAllCollisions(space2, referenceTestArea)
	if count1 != 4735 {
		t.Fatalf("expected 4735 collisions on space1, but got %d", count1)
	}
	if count1 != count2 {
		t.Fatalf("space1 got %d collisions, but space2 got %d", count1, count2)
	}
}

func testProcNeutralStabilizeN(t *testing.T, space1 Space, space2 Space, numBoxes int) {
	const seed = 0x6719B80442A
	var rng *rand.Rand
	rng = rand.New(rand.NewSource(seed))

	boxes := make([]MutableBox, 0, numBoxes)
	for i := 0; i < numBoxes; i++ {
		box := testingNewBox99(rng)
		space1.Add(box)
		space2.Add(box)
		boxes = append(boxes, box)
	}

	space1.Stabilize()
	count1 := testingCountAllCollisions(space1, referenceTestArea)
	count2 := testingCountAllCollisions(space2, referenceTestArea)
	if count1 != count2 {
		t.Fatalf("space1 got %d collisions, but space2 got %d", count1, count2)
	}
	if count1 != 12161 {
		t.Fatalf("Expected 12161 collisions on space1, but got %d", count1)
	}

	for i := 0; i < numBoxes/4; i++ {
		removed := space1.Remove(boxes[i])
		if !removed {
			t.Fatalf("space1 failed to remove box %s (at index %d)", BoxString(boxes[i]), i)
		}
		removed = space2.Remove(boxes[i])
		if !removed {
			t.Fatalf("space2 failed to remove box %s (at index %d)", BoxString(boxes[i]), i)
		}
	}

	space1.Stabilize()
	count1 = testingCountAllCollisions(space1, referenceTestArea)
	count2 = testingCountAllCollisions(space2, referenceTestArea)
	if count1 != count2 {
		t.Fatalf("space1 got %d collisions, but space2 got %d", count1, count2)
	}
	if count1 != 9215 {
		t.Fatalf("Expected 9215 collisions on space1, but got %d", count1)
	}
}

func testProcStabilizeN(t *testing.T, space Space, numBoxes int) {
	const seed = 0x049C22B184
	var rng *rand.Rand
	rng = rand.New(rand.NewSource(seed))

	brute := NewBruteForce()
	boxes := make([]MutableBox, 0, numBoxes)
	for i := 0; i < numBoxes; i++ {
		box := testingNewBox99(rng)
		space.Add(box)
		brute.Add(box)
		boxes = append(boxes, box)
	}

	for _, box := range boxes {
		newBox := testingNewBox99(rng)
		box.MutateBox(newBox.XMin(), newBox.XMax(), newBox.YMin(), newBox.YMax())
		//space.MutateBox(box, newBox.XMin(), newBox.XMax(), newBox.YMin(), newBox.YMax())
	}
	brute.Stabilize()
	space.Stabilize()
	collisionBoxes := make(map[Box]int)

	for y := 0; y < 100; y += 10 {
		for x := 0; x < 100; x += 10 {
			brute.EachCollision(
				NewBox(x, x + 9, y, y + 9),
				func(box Box) SearchControl {
					value, found := collisionBoxes[box]
					if found { value += 1 } else { value = 1 }
					collisionBoxes[box] = value
					return SearchContinue
				},
			)
		}
	}

	for y := 0; y < 100; y += 10 {
		for x := 0; x < 100; x += 10 {
			space.EachCollision(
				NewBox(x, x + 9, y, y + 9),
				func(box Box) SearchControl {
					count, found := collisionBoxes[box]
					if !found {
						t.Fatalf("Found collision between %s and space box %s (not in the result set)", BoxString(NewBox(x, x + 9, y, y + 9)), BoxString(box))
					} else if count == 0 {
						t.Fatalf("Found collision between %s and space box %s more times than in result set", BoxString(NewBox(x, x + 9, y, y + 9)), BoxString(box))
					} else {
						collisionBoxes[box] = count - 1
					}
					return SearchContinue
				},
			)
		}
	}


	for box, count := range collisionBoxes {
		if count > 0 {
			t.Fatalf("Missed at least one collision with %s", BoxString(box))
		}
	}
}

// ----- helpers -----

func testingAddBoxes(space BasicSpace, n int) {
	const seed = 0x4F22A80E3652CA75
	var rng *rand.Rand
	rng = rand.New(rand.NewSource(seed - int64(n)))

	for i := 0; i < n; i++ {
		box := testingNewBox99(rng)
		space.Add(box)
	}
}

func testingAddHorzBoxes(space BasicSpace, n int) {
	const seed = 0x352E2E64639361AE
	var rng *rand.Rand
	rng = rand.New(rand.NewSource(seed - int64(n)))

	for i := 0; i < n; i++ {
		box := testingNewHorzBox(rng)
		space.Add(box)
	}
}

func testingAddAndGetBoxes(space BasicSpace, n int) []MutableBox {
	const seed = 0x81182930EABC
	var rng *rand.Rand
	rng = rand.New(rand.NewSource(seed - int64(n)))
	boxes := make([]MutableBox, n)

	for i := 0; i < n; i++ {
		box := testingNewBox99(rng)
		space.Add(box)
		boxes[i] = box
	}
	return boxes
}

func testingAddStretchedBoxes(space BasicSpace, n int) {
	const seed = 0x02AD771296C458EF
	var rng *rand.Rand
	rng = rand.New(rand.NewSource(seed - int64(n)))

	for i := 0; i < n; i++ {
		box := testingNewStretchedBox99(rng)
		space.Add(box)
	}
}

func testingCountAllCollisions(space BasicSpace, testArea closedBox) int {
	count := 0
	for y := testArea.YMin(); y <= testArea.YMax(); y += 10 {
		for x := testArea.XMin(); x < testArea.XMax(); x += 10 {
			space.EachCollision(
				NewBox(x, x + 9, y, y + 9),
				func(Box) SearchControl { count += 1; return SearchContinue },
			)
		}
	}
	return count
}

func testingNewBox99(rng *rand.Rand) MutableBox {
	x := rng.Intn(80)
	y := rng.Intn(80)
	return NewBox(x, x + rng.Intn(21), y, y + rng.Intn(21))
}

func testingNewHorzBox(rng *rand.Rand) MutableBox {
	x := rng.Intn(800)
	y := rng.Intn(80)
	return NewBox(x, x + rng.Intn(21), y, y + rng.Intn(21))
}

func testingNewStretchedBox99(rng *rand.Rand) MutableBox {
	if rng.Intn(2) == 0 { // wide
		x := rng.Intn(20)
		y := rng.Intn(80)
		return NewBox(x, x + 30 + rng.Intn(51), y, y + rng.Intn(21))
	} else { // tall
		x := rng.Intn(80)
		y := rng.Intn(20)
		return NewBox(x, x + rng.Intn(21), y, y + 30 + rng.Intn(51))
	}
}
