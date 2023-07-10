package aabb

// Spaces and types for AABB (axis-aligned bounding box) collisions.

// Quick API reference for the Space interface:
// >> Clear()
// >> Add(Box)
// >> MutateBox(MutableBox, int, int, int, int) bool
// >> Collision(Box) Box
// >> EachCollision(Box, func(Box) bool)
// >> Remove(Box) bool
// >> Stabilize()

// Auxiliary type used on [BasicSpace].EachCollision()
// that allows the user to stop the iteration process early.
type SearchControl bool
const SearchContinue SearchControl = true
const SearchStop     SearchControl = false

// A subset of the [Space] interface.
type BasicSpace interface {
   // Removes all Boxes without deallocating underlying memory
   // (e.g: slice capacity). Mostly relevant for spaces that
   // are rebuilt frequently or can't be updated directly.
   Clear()

   // Adds the given Box to the Space. Once added, the position
   // and size of the Box *can't be modified* unless it's through
	// the explicitly documented methods (`Transform`, `Update`,
	// or any other specific methods for the underlying type).
   Add(Box)

   // Returns a colliding Box in the Space, or nil if there are no
   // collisions. Collisions where `givenBox == otherBox` are not
	// reported. If multiple collisions are possible, any of them
   // might be returned. See EachCollision if you need more control.
   Collision(Box) Box

   // A more complete version of Collision() that calls the given
   // function for each collision (instead of stopping at one).
   EachCollision(Box, func(Box) SearchControl)
}

// Additional technical notes for Space types:
//  - Space implementations are not concurrent-safe unless indicated.
//  - The Space interface is agnostic with regards to what happens with
//    overlapping Boxes, duplicate Boxes and [Box] coordinate boundaries.
//    For details on those you will have to check the documentation of
//    the underlying types implementing the [Space] interface.
//  - While updates can also be achieved by calling Remove() followed
//    by Add(), many [Space] implementers can implement MutateBox() much
//    more efficiently when it's done as a single operation, so use that.
//  - Following the principle of least surprise, Boxes reported
//    to argument functions will be reported exactly once.
//  - Even if a [Space] wraps Boxes in other types internally during
//    operation, all Boxes exposed to the user will be exposed exactly
//    as they were added, making type assertions reliable.
//  - Higher-level methods may be implemented on top of the current
//    ones, e.g., EachTouch(Box, func(Box) bool), NextCollision(Box,
//    Direction) Box, FreeRoom(Box, Direction) int, CollisionsFromTo(
//    Box, Box, func(Box)).
type Space interface {
	BasicSpace

	// Removes a Box where `givenBox == otherBox`. The returned bool
	// indicates whether any box has been removed. Remove is a slow
	// operation in many Space implementations, so make sure to check
	// the documentation of the underlying type if you need to make
	// heavy use of the operation.
   Remove(Box) bool

	// `MutateBox` searches the Space for a Box where `givenBox == otherBox`,
	// and calls MutateBox() with the given (xMin, xMax, yMin, yMax) values
	// on the Box found in the Space. If no such Box can be found or other
	// incongruences arise, the function will panic.
   MutateBox(MutableBox, int, int, int, int) // xMin, xMax, yMin, yMax

	// `Stabilize` puts the Space in a stable state. We say a Space is
	// "unstable" when at least one of its Box elements have been modified
	// externally outside the `MutateBox` function. In an unstable state,
	// any function other than `Stabilize` might err, panic or hang.
	Stabilize()
}
