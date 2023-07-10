package aabb

import "math/rand"
import "testing"

func TestAugmentedTreeBasics(t *testing.T) {
	// tip to debug when necessary:
	// > fmt.Print(augmentedTree.StringRepresentation(), "\n")

	augmentedTree := NewAugmentedTree()
	box1 := NewBox(3, 4, 3, 4)
	box2 := NewBox(4, 5, 4, 5)
	box3 := NewBox(5, 6, 5, 6)
	
	// null height check
	height, expectedHeight := augmentedTree.dfsHeight(), 0
	if height != expectedHeight {
		t.Fatalf("expected height to be %d, but got %d instead", expectedHeight, height)
	}

	// trivial addition height check
	augmentedTree.Add(box1)
	height, expectedHeight = augmentedTree.dfsHeight(), 1
	if height != expectedHeight {
		t.Fatalf("expected height to be %d, but got %d instead", expectedHeight, height)
	}

	// trivial addition height check 2
	augmentedTree.Add(box2)
	height, expectedHeight = augmentedTree.dfsHeight(), 2
	if height != expectedHeight {
		t.Fatalf("expected height to be %d, but got %d instead", expectedHeight, height)
	}

	// first left rebalancing height check
	augmentedTree.Add(box3)
	height, expectedHeight = augmentedTree.dfsHeight(), 2
	if height != expectedHeight {
		t.Fatalf("expected height to be %d, but got %d instead", expectedHeight, height)
	}
	
	// clear check
	augmentedTree.Clear()
	height, expectedHeight = augmentedTree.dfsHeight(), 0
	if height != expectedHeight {
		t.Fatalf("expected height to be %d, but got %d instead", expectedHeight, height)
	}

	// inverse addition to test right rebalancing
	augmentedTree.Add(box3)
	augmentedTree.Add(box2)
	augmentedTree.Add(box1)
	height, expectedHeight = augmentedTree.dfsHeight(), 2
	if height != expectedHeight {
		t.Fatalf("expected height to be %d, but got %d instead", expectedHeight, height)
	}
}

func TestAugmentedTreeRotateLeftRight(t *testing.T) {
	augmentedTree := NewAugmentedTree()
	box1 := NewBox(1, 1, 0, 0)
	box2 := NewBox(2, 2, 0, 0)
	box3 := NewBox(3, 3, 0, 0)
	box4 := NewBox(4, 4, 0, 0)
	box5 := NewBox(5, 5, 0, 0)
	augmentedTree.Add(box4)
	augmentedTree.Add(box3)
	augmentedTree.Add(box5)
	augmentedTree.Add(box1)
	augmentedTree.Add(box2)
	height, expectedHeight := augmentedTree.dfsHeight(), 3
	if height != expectedHeight {
		t.Fatalf("expected height to be %d, but got %d instead", expectedHeight, height)
	}
}

func TestAugmentedTreeRotateRightLeft(t *testing.T) {
	augmentedTree := NewAugmentedTree()
	box1 := NewBox(1, 1, 0, 0)
	box2 := NewBox(2, 2, 0, 0)
	box3 := NewBox(3, 3, 0, 0)
	box4 := NewBox(4, 4, 0, 0)
	box5 := NewBox(5, 5, 0, 0)
	augmentedTree.Add(box2)
	augmentedTree.Add(box1)
	augmentedTree.Add(box3)
	augmentedTree.Add(box5)
	augmentedTree.Add(box4)
	height, expectedHeight := augmentedTree.dfsHeight(), 3
	if height != expectedHeight {
		t.Fatalf("expected height to be %d, but got %d instead", expectedHeight, height)
	}
}

func TestAugmentedTreeIssue1(t *testing.T) {
	augmentedTree := NewAugmentedTree()
	var newBox = func(x int) Box { return NewBox(x, x, 0, 0) }
	box20 := newBox(20)
	box30 := newBox(30)
	box32 := newBox(32)
	box26 := newBox(26)
	box10 := newBox(10)
	box06 := newBox( 6)
	box08 := newBox( 8)
	box12 := newBox(12)
	box14 := newBox(14)
	boxIssue4 := newBox(4)
	
	augmentedTree.Add(box20)
	augmentedTree.Add(box10)
	augmentedTree.Add(box30)
	augmentedTree.Add(box06)
	augmentedTree.Add(box12)
	if augmentedTree.GetRootBox() != box20 {
		t.Fatalf("expected root %s, got root %s", BoxString(box20), BoxString(augmentedTree.GetRootBox()))
	}
	augmentedTree.Add(box26)
	augmentedTree.Add(box32)
	height, expectedHeight := augmentedTree.dfsHeight(), 3
	if height != expectedHeight {
		t.Fatalf("expected height to be %d, but got %d instead", expectedHeight, height)
	}
	augmentedTree.Add(box08)
	augmentedTree.Add(box14)
	if augmentedTree.GetRootBox() != box20 {
		t.Fatalf("expected root %s, got root %s", BoxString(box20), BoxString(augmentedTree.GetRootBox()))
	}
	
	augmentedTree.Add(boxIssue4)
	if augmentedTree.GetRootBox() != box20 {
		t.Fatalf("expected root %s, got root %s", BoxString(box20), BoxString(augmentedTree.GetRootBox()))
	}
	height, expectedHeight = augmentedTree.dfsHeight(), 4
	if height != expectedHeight {
		t.Fatalf("expected height to be %d, but got %d instead", expectedHeight, height)
	}
}

func TestAugmentedTreeRemove(t *testing.T) {
	var newBox = func(x int) Box { return NewBox(x, x, 0, 0) }
	
	augmentedTree := NewAugmentedTree()
	box := newBox(5)
	
	// simple add and remove
	augmentedTree.Add(box)
	height, expectedHeight := augmentedTree.dfsHeight(), 1
	if height != expectedHeight {
		t.Fatalf("expected height to be %d, but got %d instead", expectedHeight, height)
	}
	augmentedTree.Remove(box)
	height, expectedHeight = augmentedTree.dfsHeight(), 0
	if height != expectedHeight {
		t.Fatalf("expected height to be %d, but got %d instead", expectedHeight, height)
	}

	// add 3 boxes and remove the leaves
	b1, b2 := newBox(4), newBox(6)
	augmentedTree.Add(box)
	augmentedTree.Add(b1)
	augmentedTree.Add(b2)
	if augmentedTree.GetRootBox() != box { t.Fatal("root failure") }
	
	augmentedTree.Remove(b1)
	if augmentedTree.GetRootBox() != box { t.Fatal("root failure") }
	height, expectedHeight = augmentedTree.dfsHeight(), 2
	if height != expectedHeight {
		t.Fatalf("expected height to be %d, but got %d instead", expectedHeight, height)
	}
	augmentedTree.Remove(b2)
	if augmentedTree.GetRootBox() != box { t.Fatal("root failure") }
	height, expectedHeight = augmentedTree.dfsHeight(), 1
	if height != expectedHeight {
		t.Fatalf("expected height to be %d, but got %d instead", expectedHeight, height)
	}
	if augmentedTree.GetRootBox() != box { t.Fatal("root failure") }

	// first tricky case: remove root
	augmentedTree.Add(b1)
	augmentedTree.Add(b2)
	if !augmentedTree.Remove(box) {
		t.Fatal("failed to remove root")
	}
	height, expectedHeight = augmentedTree.dfsHeight(), 2
	if height != expectedHeight {
		t.Fatalf("expected height to be %d, but got %d instead", expectedHeight, height)
	}
	if augmentedTree.GetRootBox() != b2 { t.Fatal("root failure") }
	
	// complete removals
	if !augmentedTree.Remove(b2) { t.Fatal("failed to remove b2") }
	if augmentedTree.GetRootBox() != b1 { t.Fatal("root failure") }
	augmentedTree.Clear()

	// more tricky cases
	augmentedTree.Add(box)
	augmentedTree.Add(b1)
	b3 := newBox(7)
	augmentedTree.Add(b3)
	augmentedTree.Add(b2)
	b4 := newBox(8)
	augmentedTree.Add(b4)
	// (H2 | MaxX8 [X5] -> L(H0 | MaxX4 [X4]), R(H1 | MaxX8 [X7] -> L(H0 | MaxX6 [X6]), R(H0 | MaxX8 [X8]))
	
	augmentedTree.Remove(b3)
	if augmentedTree.Collision(newBox(b2.XMin()))  != b2  { t.Fatal("b2 disappeared") }
	if augmentedTree.Collision(newBox(box.XMin())) != box { t.Fatal("box disappeared") }
	if augmentedTree.Collision(newBox(b1.XMin()))  != b1  { t.Fatal("b1 disappeared") }
	if augmentedTree.Collision(newBox(b4.XMin()))  != b4  { t.Fatal("b4 disappeared") }

	// general tricky case
	augmentedTree.Clear()
	augmentedTree.Add(box)
	augmentedTree.Add(b1)
	augmentedTree.Add(b3)
	augmentedTree.Add(b2)
	augmentedTree.Add(b4)
	augmentedTree.Remove(box)
	if augmentedTree.GetRootBox() != b2 { t.Fatal("root failure") }
	height, expectedHeight = augmentedTree.dfsHeight(), 3
	if height != expectedHeight {
		t.Fatalf("expected height to be %d, but got %d instead", expectedHeight, height)
	}
	if augmentedTree.Collision(newBox(b1.XMin())) != b1 { t.Fatal("b1 disappeared") }
	if augmentedTree.Collision(newBox(b2.XMin())) != b2 { t.Fatal("b2 disappeared") }
	if augmentedTree.Collision(newBox(b3.XMin())) != b3 { t.Fatal("b3 disappeared") }
	if augmentedTree.Collision(newBox(b4.XMin())) != b4 { t.Fatal("b4 disappeared") }
}

func TestAugmentedTreeRng(t *testing.T) {
	var newBox = func(x int) Box { return NewBox(x, x, 0, 0) }
	
	const NumIters = 300
	const NumBoxes = 7
	const Seed = 0x8795812C457D2

	rng := rand.New(rand.NewSource(Seed))
	augmentedTree := NewAugmentedTree()
	boxes := make([]Box, NumBoxes)
	for i := 0; i < NumIters; i++ {
		augmentedTree.Clear()
		for n := 0; n < NumBoxes; n++ {
			boxes[n] = newBox(1 + rng.Intn(99))
		}

		for n := 0; n < NumBoxes; n++ {
			augmentedTree.Add(boxes[n])
		}
		treeInfo := augmentedTree.StringRepresentation()

		augmentedTree.Remove(boxes[0])
		for n := 1; n < NumBoxes; n++ {
			if augmentedTree.Collision(newBox(boxes[n].XMin())) == nil {
				t.Fatalf("iter#%d, removed box %d, but missing box %d on %s\n(after removal, tree = %s)",
					i, boxes[0].XMin(), boxes[n].XMin(), treeInfo, augmentedTree.StringRepresentation(),
				)
			}
		}
	}
}

func TestAugmentedTreeRemoveIssue1(t *testing.T) {
	var newBox = func(x int) Box { return NewBox(x, x, 0, 0) }
	augmentedTree := NewAugmentedTree()
	boxH0N0 := newBox(53)
	boxH1N0 := newBox(21)
	boxH1N1 := newBox(65)
	boxH2N0 := newBox( 9)
	boxH2N2 := newBox(55)
	boxH2N3 := newBox(86)
	boxH2N4 := newBox(85)
	augmentedTree.Add(boxH0N0)
	augmentedTree.Add(boxH1N0)
	augmentedTree.Add(boxH1N1)
	augmentedTree.Add(boxH2N0)
	augmentedTree.Add(boxH2N2)
	augmentedTree.Add(boxH2N3)
	augmentedTree.Add(boxH2N4)
	augmentedTree.Remove(boxH0N0)
	if augmentedTree.Collision(newBox(boxH2N4.XMin())) == nil {
		t.Fatal("broken removal . issue1")
	}
}

func TestAugmentedTreeRemoveIssue2(t *testing.T) {
	var newBox = func(x int) Box { return NewBox(x, x, 0, 0) }
	augmentedTree := NewAugmentedTree()
	boxH0N0 := newBox(6)
	boxH1N0 := newBox(1)
	boxH1N1 := newBox(36)
	boxH2N1 := newBox(5)
	boxH2N2 := newBox(20)
	boxH2N3 := newBox(73)
	boxH3N4 := newBox(14)
	augmentedTree.Add(boxH0N0)
	augmentedTree.Add(boxH1N0)
	augmentedTree.Add(boxH1N1)
	augmentedTree.Add(boxH2N1)
	augmentedTree.Add(boxH2N2)
	augmentedTree.Add(boxH2N3)
	augmentedTree.Add(boxH3N4)
	augmentedTree.Remove(boxH0N0)
	if augmentedTree.Collision(newBox(boxH2N3.XMin())) == nil {
		t.Fatal("broken removal . issue2")
	}
}

func TestAugmentedTreeSingleCollision(t *testing.T) {
	space := NewAugmentedTree()
	testProcSingleCollision(t, space)
}

func TestAugmentedTreeStabilize600(t *testing.T) {
	space := NewAugmentedTree()
	testProcStabilizeN(t, space, 600)
}

func TestAugmentedTreeNeutralStabilize3000(t *testing.T) {
	space1 := NewAugmentedTree()
	space2 := NewAugmentedTree()
	testProcNeutralStabilizeN(t, space1, space2, 3000)
}

func TestAugmentedTreeMutateVsStabilize1200(t *testing.T) {
	space1 := NewAugmentedTree()
	space2 := NewAugmentedTree()
	testProcMutateVsStabilize1200(t, space1, space2)
}

// --- benchmarks ---

func BenchmarkAugmentedTree500(b *testing.B) {
	space := NewAugmentedTree()
	benchmarkBasicSpace500(b, space)
}

func BenchmarkAugmentedTreeStretch500(b *testing.B) {
	space := NewAugmentedTree()
	benchmarkBasicSpaceStretch500(b, space)
}

func BenchmarkAugmentedTreeHorz500(b *testing.B) {	
	benchmarkBasicSpaceHorz500(b, NewAugmentedTree())
}

func BenchmarkAugmentedTree2000(b *testing.B) {
	space := NewAugmentedTree()
	benchmarkBasicSpace2000(b, space)
}

func BenchmarkAugmentedTreeQuarterMuts1000(b *testing.B) {
	space := NewAugmentedTree()
	benchmarkSpaceQuarterMuts1000(b, space)
}

func BenchmarkAugmentedTreeStabilize2500(b *testing.B) {
	space := NewAugmentedTree()
	benchmarkSpaceStabilize2500(b, space)
}
