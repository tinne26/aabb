package aabb

import "strings"
import "strconv"

type cellIter struct {
	boxIndex int // set to -1 if there's no box in the cell
	nextIter int // set to -1 if there's no next iterator
}

type cellIterList struct {
	list []cellIter
	freeIndex int
}

func newCellIterList(initialCapacity int) cellIterList {
	return cellIterList {
		list: make([]cellIter, 0, initialCapacity),
		freeIndex: -1,
	}
}

func (self *cellIterList) DebugString() string {
	var strBuilder strings.Builder
	strBuilder.WriteString("cellIterList { freeIndex chain:")
	freeIndex := self.freeIndex
	for freeIndex != -1 {
		strBuilder.WriteRune(' ')
		strBuilder.WriteString(strconv.Itoa(freeIndex))
		freeIndex = self.list[freeIndex].nextIter
	}
	strBuilder.WriteString(" -1 } { cellIters:")
	for i, cellIterElem := range self.list {
		if cellIterElem.boxIndex == -1 { continue }
		strBuilder.WriteString(" ( i#" + strconv.Itoa(i))
		strBuilder.WriteString(" box#" + strconv.Itoa(cellIterElem.boxIndex))
		strBuilder.WriteString(" next#" + strconv.Itoa(cellIterElem.nextIter) + " )")
	}
	strBuilder.WriteString(" }")

	return strBuilder.String()
}

func (self *cellIterList) Clear() {
	self.list = self.list[0 : 0]
	self.freeIndex = -1
}

func (self *cellIterList) Next(iterIndex int) (int, int) {
	cellIterElem := self.list[iterIndex]
	return cellIterElem.boxIndex, cellIterElem.nextIter
}

func (self *cellIterList) AddIterTo(boxIndex int, nextIter int) int {
	if self.freeIndex != -1 {
		// reuse existing free index
		newIndex := self.freeIndex
		self.freeIndex = self.list[newIndex].nextIter
		self.list[newIndex] = cellIter{ boxIndex, nextIter }
		return newIndex
	} else {
		// allocate unless we still have capacity left
		newIndex := len(self.list)
		if cap(self.list) > newIndex {
			self.list = self.list[0 : newIndex + 1]
			self.list[newIndex] = cellIter{ boxIndex, nextIter }
		} else {
			self.list = append(self.list, cellIter{ boxIndex, nextIter })
		}
		return newIndex
	}
}

func (self *cellIterList) CutIter(iterIndex int, prevIter int) {
	if prevIter != -1 {
		self.list[prevIter].nextIter = self.list[iterIndex].nextIter
	}

	self.list[iterIndex].boxIndex = -1 // not strictly needed, but helpful on debug
	self.list[iterIndex].nextIter = self.freeIndex
	self.freeIndex = iterIndex
}
