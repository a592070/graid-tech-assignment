package raid

import (
	"errors"
)

type RAID0 struct {
	numDisks   int
	stripeSize int
	disks      [][]byte
}

func NewRAID0(numDisks, stripeSize int) (*RAID0, error) {
	if numDisks < 2 {
		return nil, errors.New("RAID0: number of disks must be greater or equals than 2")
	}
	raid := &RAID0{
		numDisks, stripeSize, make([][]byte, numDisks),
	}
	for i := range raid.disks {
		raid.disks[i] = make([]byte, 0)
	}
	return raid, nil
}

// When writing data, we need to split the input into chunks of stripeSize,
// and distribute them across the disks in order.
// For example, if the stripeSize is 2, and the data is "abcdef",
// then disk0 would get "ab", disk1 "cd", disk0 "ef", etc.
func (r *RAID0) Write(data []byte, pos int) error {
	if r.numDisks <= 0 {
		return errors.New("RAID0: no disks available")
	}
	for i := 0; i < len(data); i++ {
		logicalPos := pos + i
		stripeNumber := logicalPos / r.stripeSize

		diskIndex := stripeNumber % r.numDisks
		offsetInStripe := logicalPos % r.stripeSize
		stripeInDisk := stripeNumber / r.numDisks
		diskOffset := stripeInDisk*r.stripeSize + offsetInStripe

		if diskOffset >= len(r.disks[diskIndex]) {
			needed := diskOffset - len(r.disks[diskIndex]) + 1
			r.disks[diskIndex] = append(r.disks[diskIndex], make([]byte, needed)...)
		}
		r.disks[diskIndex][diskOffset] = data[i]
	}
	return nil
}

// When reading, for each logical byte index, compute the disk and offset,
// then read the byte from there. If the disk's slice is shorter than the offset, return zero.
func (r *RAID0) Read(length int, pos int) ([]byte, error) {
	result := make([]byte, length)
	if r.numDisks <= 0 {
		return result, errors.New("RAID0: no disks available")
	}
	for i := 0; i < length; i++ {
		logicalPos := pos + i
		stripeNumber := logicalPos / r.stripeSize
		diskIndex := stripeNumber % r.numDisks
		if diskIndex >= len(r.disks) {
			result[i] = 0
			continue
		}
		offsetInStripe := logicalPos % r.stripeSize
		stripeInDisk := stripeNumber / r.numDisks
		diskOffset := stripeInDisk*r.stripeSize + offsetInStripe

		disk := r.disks[diskIndex]
		if diskOffset >= len(disk) {
			result[i] = 0
		} else {
			result[i] = disk[diskOffset]
		}
	}
	return result, nil
}

func (r *RAID0) ClearDisk(diskIndex int) {
	if diskIndex >= 0 && diskIndex < len(r.disks) {
		for i := range r.disks[diskIndex] {
			r.disks[diskIndex][i] = 0
		}
	}
}
