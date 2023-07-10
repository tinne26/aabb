package aabb

// To run all tests:
// >> go test ./... | grep "^[^?]"
// To run all benchmarks:
// >> go test -bench=. ./... | grep "^[^?]"

type closedBox struct { a, b, c, d int }
func (self closedBox) XMin() int { return self.a }
func (self closedBox) XMax() int { return self.b }
func (self closedBox) YMin() int { return self.c }
func (self closedBox) YMax() int { return self.d }
var referenceTestBox  = closedBox{ 10, 15, 10, 15 }
var referenceTestArea = closedBox{  0, 99,  0, 99 }
var referenceTestAreaHorz = closedBox{  0, 819,  0, 99 }
var refTestAxisCells  = 25

type testBox struct { xmin, xmax, ymin, ymax int }
var singleNoCollisionTests = []testBox {
	{  0,  5,  0,  5 },
	{ 16, 18,  0,  5 },
	{ 16, 18, 20, 99 },
	{  0,  9, 16, 22 },
	{ 10, 15, 16, 16 },
	{  6,  8, 12, 14 },
}

var singleYesCollisionTests = []testBox {
	{  10,  15,  10,  15 },
	{   8,  18,   8,  10 },
	{   0,  99,   0,  99 },
	{  12,  14,  12,  14 },
	{  13,  13,  13,  13 },
	{   7,  12,  13,  19 },
	{   7,  10,   7,  10 },
	{  15,  99,  15,  99 },
	{  15,  15,  10,  10 },
}
