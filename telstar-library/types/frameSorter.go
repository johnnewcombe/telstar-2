package types

import "sort"

func SortFrames(frames []Frame) {
	pageNumber := func(c1, c2 *Frame) bool {
		return c1.PID.PageNumber < c2.PID.PageNumber
	}
	frameId := func(c1, c2 *Frame) bool {
		return c1.PID.FrameId < c2.PID.FrameId
	}

	frameSortOrder(pageNumber, frameId).Sort(frames)
}

// Implementation of sort interface for multiple fields
type lessFunc func(p1, p2 *Frame) bool

// FrameSortOrder returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func frameSortOrder(less ...lessFunc) *frameSorter {
	return &frameSorter{
		less: less,
	}
}

// frameSorter implements the Sort interface, sorting the frames within.
type frameSorter struct {
	changes []Frame
	less    []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to FrameSortOrder.
func (ms *frameSorter) Sort(changes []Frame) {
	ms.changes = changes
	sort.Sort(ms)
}

// Len is part of sort.Interface.
func (ms *frameSorter) Len() int {
	return len(ms.changes)
}

// Swap is part of sort.Interface.
func (ms *frameSorter) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that discriminates between
// the two items (one is less than the other). Note that it can call the
// less functions twice per call. We could change the functions to return
// -1, 0, 1 and reduce the number of calls for greater efficiency: an
// exercise for the reader.
func (ms *frameSorter) Less(i, j int) bool {
	p, q := &ms.changes[i], &ms.changes[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}
