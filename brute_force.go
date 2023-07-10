package aabb

// BruteForce implements the naive, O(n^2) algorithm to check collisions
// between Boxes. One advantage of this implementation, besides being
// very simple, is that since Boxes are not ordered or structured at all,
// making changes doesn't require stabilizing the structure.
type BruteForce struct { boxes []Box }
func NewBruteForce() *BruteForce {
   return &BruteForce { boxes: make([]Box, 0, 64) }
}

func (self *BruteForce) Clear() { self.boxes = self.boxes[0 : 0] }
func (self *BruteForce) Add(box Box) {
   self.boxes = append(self.boxes, box)
}

func (self *BruteForce) MutateBox(box MutableBox, xMin, xMax, yMin, yMax int) {
	for _, otherBox := range self.boxes {
		if box == otherBox {
			otherBox.(MutableBox).MutateBox(xMin, xMax, yMin, yMax)
			return
		}
	}
	panic("box to be updated not found")
}

// Always stable, this method does nothing for this approach.
func (self *BruteForce) Stabilize() { /* always stable */ }

func (self *BruteForce) Collision(box Box) Box {
   for _, otherBox := range self.boxes {
      if BoxesCollide(box, otherBox) && box != otherBox {
         return otherBox
      }
   }
   return nil
}

func (self *BruteForce) EachCollision(box Box, eachFunc func(Box) SearchControl) {
   for _, otherBox := range self.boxes {
      if BoxesCollide(box, otherBox) && box != otherBox {
         searchControl := eachFunc(otherBox)
         if searchControl == SearchStop { return }
      }
   }
}

func (self *BruteForce) Remove(box Box) bool {
   for i, otherBox := range self.boxes {
      if otherBox == box {
         lastIdx := len(self.boxes) - 1 // can't be negative, we have a box
         self.boxes[i], self.boxes[lastIdx] = self.boxes[lastIdx], self.boxes[i]
         self.boxes = self.boxes[0 : lastIdx]
         return true
      }
   }
   return false
}
