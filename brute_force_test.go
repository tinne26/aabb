package aabb

import "testing"

func TestBruteForceSingleCollision(t *testing.T) {
	space := NewBruteForce()
	testProcSingleCollision(t, space)
}

func TestBruteForceMutateVsStabilize1200(t *testing.T) {
	space1 := NewBruteForce()
	space2 := NewBruteForce()
	testProcMutateVsStabilize1200(t, space1, space2)
}

func TestBruteForceStabilize600(t *testing.T) {
	space := NewBruteForce()
	testProcStabilizeN(t, space, 600)
}

func TestBruteForceNeutralStabilize3000(t *testing.T) {
	space1 := NewBruteForce()
	space2 := NewBruteForce()
	testProcNeutralStabilizeN(t, space1, space2, 3000)
}


// --- benchmarks ---

func BenchmarkBruteForce500(b *testing.B) {
	bruteSpace := NewBruteForce()
	benchmarkBasicSpace500(b, bruteSpace)
}

func BenchmarkBruteForceStretch500(b *testing.B) {
	bruteSpace := NewBruteForce()
	benchmarkBasicSpaceStretch500(b, bruteSpace)
}

func BenchmarkBruteForceHorz500(b *testing.B) {	
	benchmarkBasicSpaceHorz500(b, NewBruteForce())
}

func BenchmarkBruteForce2000(b *testing.B) {
	bruteSpace := NewBruteForce()
	benchmarkBasicSpace2000(b, bruteSpace)
}

func BenchmarkBruteForceQuarterMuts1000(b *testing.B) {
	bruteSpace := NewBruteForce()
	benchmarkSpaceQuarterMuts1000(b, bruteSpace)
}

func BenchmarkBruteForceStabilize2500(b *testing.B) {
	space := NewBruteForce()
	benchmarkSpaceStabilize2500(b, space)
}
