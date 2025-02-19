package raid

import (
	"errors"
	"fmt"
)

type RAID5 struct {
	numDisks   int
	stripeSize int
	disks      [][]byte
}

func NewRAID5(numDisks, stripeSize int) (*RAID5, error) {
	if numDisks < 3 {
		return nil, errors.New("RAID5: number of disks must be greater or equals than 3")
	}
	raid := &RAID5{
		numDisks, stripeSize, make([][]byte, numDisks),
	}
	for i := range raid.disks {
		raid.disks[i] = make([]byte, 0)
	}
	return raid, nil
}

func (r *RAID5) Read(length int, offset int) ([]byte, error) {
	numDisks := r.numDisks
	stripeDataSize := (numDisks - 1) * r.stripeSize

	startStripe := offset / stripeDataSize
	endStripe := (offset + length + stripeDataSize - 1) / stripeDataSize

	var data []byte

	for s := startStripe; s < endStripe; s++ {
		p := s % numDisks
		dataDisks := make([]int, 0, numDisks-1)
		for d := 0; d < numDisks; d++ {
			if d != p {
				dataDisks = append(dataDisks, d)
			}
		}

		var blocks [][]byte
		var missingDisk int = -1
		for _, disk := range dataDisks {
			if (s+1)*r.stripeSize > len(r.disks[disk]) {
				if missingDisk == -1 {
					missingDisk = disk
				} else {
					return nil, errors.New("RAID5: multiple disks failed")
				}
				blocks = append(blocks, nil)
				continue
			}
			block := r.disks[disk][s*r.stripeSize : (s+1)*r.stripeSize]
			isZero := true
			for _, b := range block {
				if b != 0 {
					isZero = false
					break
				}
			}
			if isZero {
				if missingDisk == -1 {
					missingDisk = disk
				} else {
					return nil, errors.New("RAID5: multiple disks failed")
				}
				blocks = append(blocks, nil)
			} else {
				blocks = append(blocks, block)
			}
		}

		if missingDisk != -1 {
			parityBlock := r.disks[p][s*r.stripeSize : (s+1)*r.stripeSize]
			missingIdx := -1
			for i, d := range dataDisks {
				if d == missingDisk {
					missingIdx = i
					break
				}
			}
			if missingIdx == -1 {
				return nil, errors.New("RAID5: missing disk not in dataDisks")
			}

			reconstructedBlock := make([]byte, r.stripeSize)
			copy(reconstructedBlock, parityBlock)
			for i, block := range blocks {
				if i == missingIdx || block == nil {
					continue
				}
				for j := 0; j < r.stripeSize; j++ {
					reconstructedBlock[j] ^= block[j]
				}
			}
			blocks[missingIdx] = reconstructedBlock
		}

		for _, block := range blocks {
			data = append(data, block...)
		}
	}

	startOffset := offset % stripeDataSize
	endOffset := startOffset + length
	if endOffset > len(data) {
		endOffset = len(data)
	}
	return data[startOffset:endOffset], nil
}

func (r *RAID5) calculateParity(stripeIndex int) byte {
	parity := byte(0)
	numDataDisks := r.numDisks - 1
	for diskIndex := 0; diskIndex < numDataDisks; diskIndex++ {
		stripeOffset := stripeIndex * r.stripeSize
		parity ^= r.disks[diskIndex][stripeOffset]
		parity ^= r.disks[diskIndex][stripeOffset+1]
		parity ^= r.disks[diskIndex][stripeOffset+2]
		parity ^= r.disks[diskIndex][stripeOffset+3]
	}
	return parity
}

func (r *RAID5) Write(data []byte, offset int) error {
	stripeDataSize := (r.numDisks - 1) * r.stripeSize

	if offset%stripeDataSize != 0 {
		return fmt.Errorf("RAID5: offset %d is not aligned to stripe data size %d", offset, stripeDataSize)
	}

	startStripe := offset / stripeDataSize
	numStripes := (len(data) + stripeDataSize - 1) / stripeDataSize

	for s := 0; s < numStripes; s++ {
		stripePos := startStripe + s
		parityDisk := stripePos % r.numDisks

		start := s * stripeDataSize
		end := start + stripeDataSize
		if end > len(data) {
			end = len(data)
		}
		stripeData := data[start:end]
		if len(stripeData) < stripeDataSize {
			padding := make([]byte, stripeDataSize-len(stripeData))
			stripeData = append(stripeData, padding...)
		}

		blocks := make([][]byte, r.numDisks-1)
		for i := 0; i < r.numDisks-1; i++ {
			blockStart := i * r.stripeSize
			blockEnd := blockStart + r.stripeSize
			blocks[i] = stripeData[blockStart:blockEnd]
		}

		parityBlock := make([]byte, r.stripeSize)
		for i := 0; i < r.stripeSize; i++ {
			for _, block := range blocks {
				parityBlock[i] ^= block[i]
			}
		}

		dataDisks := make([]int, 0, r.numDisks-1)
		for d := 0; d < r.numDisks; d++ {
			if d != parityDisk {
				dataDisks = append(dataDisks, d)
			}
		}

		for i, disk := range dataDisks {
			requiredLen := (stripePos + 1) * r.stripeSize
			if len(r.disks[disk]) < requiredLen {
				newDisk := make([]byte, requiredLen)
				copy(newDisk, r.disks[disk])
				r.disks[disk] = newDisk
			}
			copy(r.disks[disk][stripePos*r.stripeSize:(stripePos+1)*r.stripeSize], blocks[i])
		}

		requiredLen := (stripePos + 1) * r.stripeSize
		if len(r.disks[parityDisk]) < requiredLen {
			newDisk := make([]byte, requiredLen)
			copy(newDisk, r.disks[parityDisk])
			r.disks[parityDisk] = newDisk
		}
		copy(r.disks[parityDisk][stripePos*r.stripeSize:(stripePos+1)*r.stripeSize], parityBlock)
	}

	return nil
}

func (r *RAID5) ClearDisk(diskIndex int) {
	if diskIndex >= 0 && diskIndex < len(r.disks) {
		for i := range r.disks[diskIndex] {
			r.disks[diskIndex][i] = 0
		}
	}
}
