package aabb

import "strconv"

func spaceBoxChecks(box Box, workingArea Box) {
	BoxPanicIfInvalid(box)
	if box.XMin() < workingArea.XMin() { panic("box.XMin() < workingArea.XMin() [" + strconv.Itoa(box.XMin()) + " < " + strconv.Itoa(workingArea.XMin()) + "]") }
	if box.XMax() > workingArea.XMax() { panic("box.XMax() > workingArea.XMax() [" + strconv.Itoa(box.XMax()) + " > " + strconv.Itoa(workingArea.XMax()) + "]") }
	if box.YMin() < workingArea.YMin() { panic("box.YMin() < workingArea.YMin() [" + strconv.Itoa(box.YMin()) + " < " + strconv.Itoa(workingArea.YMin()) + "]") }
	if box.YMax() > workingArea.YMax() { panic("box.YMax() > workingArea.YMax() [" + strconv.Itoa(box.YMax()) + " > " + strconv.Itoa(workingArea.YMax()) + "]") }
}

func spaceBoxChecksWith(box Box, workAreaMinX, workAreaMaxX, workAreaMinY, workAreaMaxY int) {
	BoxPanicIfInvalid(box)
	if box.XMin() < workAreaMinX { panic("box.XMin() < workingArea.XMin() [" + strconv.Itoa(box.XMin()) + " < " + strconv.Itoa(workAreaMinX) + "]") }
	if box.XMax() > workAreaMaxX { panic("box.XMax() > workingArea.XMax() [" + strconv.Itoa(box.XMax()) + " > " + strconv.Itoa(workAreaMaxX) + "]") }
	if box.YMin() < workAreaMinY { panic("box.YMin() < workingArea.YMin() [" + strconv.Itoa(box.YMin()) + " < " + strconv.Itoa(workAreaMinY) + "]") }
	if box.YMax() > workAreaMaxY { panic("box.YMax() > workingArea.YMax() [" + strconv.Itoa(box.YMax()) + " > " + strconv.Itoa(workAreaMaxY) + "]") }
}
