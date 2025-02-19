package raid

import (
	"errors"
)

type RAID1 struct {
	numDisks int
	disks    [][]byte
}

func NewRAID1(numDisks int) (*RAID1, error) {
	if numDisks < 2 {
		return nil, errors.New("RAID1: number of disks must be greater or equals than 2")
	}
	raid := &RAID1{
		numDisks, make([][]byte, numDisks),
	}
	for i := range raid.disks {
		raid.disks[i] = make([]byte, 0)
	}
	return raid, nil
}

// Mirror data
func (r *RAID1) Write(data []byte, pos int) error {
	if r.numDisks <= 0 {
		return errors.New("RAID1: no disks available")
	}

	// Increase disk size if needed
	diskOffset := pos + len(data)
	for diskIndex := range r.numDisks {
		if diskOffset >= len(r.disks[diskIndex]) {
			needed := diskOffset - len(r.disks[diskIndex]) + 1
			r.disks[diskIndex] = append(r.disks[diskIndex], make([]byte, needed)...)
		}
	}

	for i := 0; i < len(data); i++ {
		logicalPos := pos + i
		for diskIndex := range r.numDisks {
			r.disks[diskIndex][logicalPos] = data[i]
		}
	}
	return nil
}

func (r *RAID1) Read(length int, pos int) ([]byte, error) {
	result := make([]byte, length)
	if r.numDisks <= 0 {
		return result, errors.New("RAID1: no disks available")
	}

	for i := 0; i < length; i++ {
		logicalPos := pos + i
		if logicalPos >= len(r.disks[0]) {
			return nil, errors.New("raid1: logical position out of range")
		}
		for diskIndex := range r.numDisks {
			if r.disks[diskIndex][logicalPos] != 0 {
				result[i] = r.disks[diskIndex][logicalPos]
				break
			}
		}
	}

	return result, nil
}

func (r *RAID1) ClearDisk(diskIndex int) {
	if diskIndex >= 0 && diskIndex < len(r.disks) {
		for i := range r.disks[diskIndex] {
			r.disks[diskIndex][i] = 0
		}
	}
}
