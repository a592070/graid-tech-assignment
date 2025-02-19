package raid

import "errors"

type RAID10 struct {
	numDisks   int
	stripeSize int
	disks      [][]byte
}

func NewRAID10(numDisks, stripeSize int) (*RAID10, error) {
	if numDisks < 4 {
		return nil, errors.New("RAID10: number of disks must be greater or equals than 4")
	}
	if numDisks%2 != 0 {
		return nil, errors.New("RAID10: number of disks must be even")
	}
	raid := &RAID10{
		numDisks, stripeSize, make([][]byte, numDisks),
	}
	for i := range raid.disks {
		raid.disks[i] = make([]byte, 0)
	}
	return raid, nil
}

// When reading, for each logical position, determine which pair and offset,
// then read from either of the disks in the pair (since they are mirrored).
// However, if one disk is cleared (zeroed),
// then the other disk in the pair should still have the data.
func (r *RAID10) Read(length int, pos int) ([]byte, error) {
	result := make([]byte, length)

	numPairs := len(r.disks) / 2
	for i := 0; i < length; i++ {
		logicalPos := pos + i
		stripeNumber := logicalPos / r.stripeSize
		pairIndex := stripeNumber % numPairs
		disk1 := pairIndex * 2
		disk2 := pairIndex*2 + 1

		offset := (stripeNumber/numPairs)*r.stripeSize + (logicalPos % r.stripeSize)
		var value byte

		// Check first disk in pair
		if disk1 < len(r.disks) {
			disk := r.disks[disk1]
			if offset < len(disk) {
				value = disk[offset]
			}
		}

		// If first disk's value is zero, check second disk
		if value == 0 && disk2 < len(r.disks) {
			disk := r.disks[disk2]
			if offset < len(disk) {
				value = disk[offset]
			}
		}

		result[i] = value
	}

	return result, nil
}

// When writing data, we need to split the data into chunks (stripes),
// each stripe's size is determined by the stripe size.
// Then, for each stripe, write it to a pair of disks (mirrored).
// The pairs are selected in a round-robin fashion.
// So stripe 0 goes to pair 0 (disks 0 and 1), stripe 1 to pair 1 (disks 2 and 3), stripe 2 to pair 0 again, etc.
func (r *RAID10) Write(data []byte, pos int) error {
	numPairs := len(r.disks) / 2
	for i := 0; i < len(data); i++ {
		logicalPos := pos + i
		stripeNumber := logicalPos / r.stripeSize
		pairIndex := stripeNumber % numPairs
		disk1 := pairIndex * 2
		disk2 := pairIndex*2 + 1

		offset := (stripeNumber/numPairs)*r.stripeSize + (logicalPos % r.stripeSize)

		// Write to first disk in pair
		if disk1 < len(r.disks) {
			if offset >= len(r.disks[disk1]) {
				needed := offset - len(r.disks[disk1]) + 1
				r.disks[disk1] = append(r.disks[disk1], make([]byte, needed)...)
			}
			r.disks[disk1][offset] = data[i]
		}

		// Write to second disk in pair
		if disk2 < len(r.disks) {
			if offset >= len(r.disks[disk2]) {
				needed := offset - len(r.disks[disk2]) + 1
				r.disks[disk2] = append(r.disks[disk2], make([]byte, needed)...)
			}
			r.disks[disk2][offset] = data[i]
		}
	}

	return nil
}

func (r *RAID10) ClearDisk(diskIndex int) {
	if diskIndex >= 0 && diskIndex < len(r.disks) {
		for i := range r.disks[diskIndex] {
			r.disks[diskIndex][i] = 0
		}
	}
}
