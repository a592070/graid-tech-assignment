package raid

import "errors"

type RAID6 struct {
	numDisks   int
	stripeSize int
	disks      [][]byte
	dataDisks  int
}

func NewRAID6(numDisks, stripeSize int) (*RAID6, error) {
	if numDisks < 4 {
		return nil, errors.New("RAID6: number of disks must be greater or equals than 4")
	}
	raid := &RAID6{
		numDisks, stripeSize, make([][]byte, numDisks), numDisks - 2,
	}
	for i := range raid.disks {
		raid.disks[i] = make([]byte, 0)
	}
	return raid, nil
}

func gfMultiply(a, b byte) byte {
	var product byte = 0
	for i := 0; i < 8; i++ {
		if (b & 1) != 0 {
			product ^= a
		}
		highBit := a & 0x80
		a <<= 1
		if highBit != 0 {
			a ^= 0x1D
		}
		b >>= 1
	}
	return product
}

func gfInverse(a byte) byte {
	if a == 0 {
		return 0
	}
	for b := byte(1); b != 0; b++ {
		if gfMultiply(a, b) == 1 {
			return b
		}
	}
	return 0
}

func (r *RAID6) Read(length int, offset int) ([]byte, error) {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		logicalPos := offset + i
		stripe := logicalPos / (r.stripeSize * r.dataDisks)
		stripeByte := logicalPos % (r.stripeSize * r.dataDisks)
		dataBlockIndex := stripeByte / r.stripeSize
		byteInBlock := stripeByte % r.stripeSize

		diskIdx := dataBlockIndex
		stripeOffset := stripe * r.stripeSize
		diskOffset := stripeOffset + byteInBlock

		var dataByte byte
		if diskIdx < 0 || diskIdx >= r.dataDisks {
			dataByte = 0
		} else {
			if diskOffset < len(r.disks[diskIdx]) {
				dataByte = r.disks[diskIdx][diskOffset]
			}
		}

		if dataByte == 0 {
			pDisk := r.dataDisks
			qDisk := r.dataDisks + 1

			var p, q byte
			if diskOffset < len(r.disks[pDisk]) {
				p = r.disks[pDisk][diskOffset]
			}
			if diskOffset < len(r.disks[qDisk]) {
				q = r.disks[qDisk][diskOffset]
			}

			sumP := byte(0)
			missingIndices := make([]int, 0)
			for j := 0; j < r.dataDisks; j++ {
				if j == diskIdx {
					missingIndices = append(missingIndices, j)
					continue
				}
				if stripeOffset+byteInBlock < len(r.disks[j]) {
					sumP ^= r.disks[j][stripeOffset+byteInBlock]
				}
			}

			if len(missingIndices) == 1 {
				dataByte = sumP ^ p
			} else if len(missingIndices) == 2 {
				sumQ := byte(0)
				for j := 0; j < r.dataDisks; j++ {
					if j == missingIndices[0] || j == missingIndices[1] {
						continue
					}
					if stripeOffset+byteInBlock < len(r.disks[j]) {
						b := r.disks[j][stripeOffset+byteInBlock]
						sumQ ^= gfMultiply(b, byte(j+1))
					}
				}
				coeffA := byte(missingIndices[0] + 1)
				coeffB := byte(missingIndices[1] + 1)
				a := gfMultiply(0x01, coeffA) ^ gfMultiply(0x01, coeffB)
				b := p ^ sumP
				c := q ^ sumQ
				dataByte = gfMultiply(c, gfInverse(a)) ^ b
			}
		}

		result[i] = dataByte
	}
	return result, nil
}

func (r *RAID6) Write(data []byte, offset int) error {
	stripeDataSize := r.stripeSize * r.dataDisks

	// Calculate starting stripe and offset within the stripe
	startStripe := offset / stripeDataSize

	// Process each byte in the data
	for i := 0; i < len(data); i++ {
		logicalPos := offset + i
		stripe := logicalPos / stripeDataSize
		byteInStripe := logicalPos % stripeDataSize

		// Determine data block and byte position within the block
		blockIndex := byteInStripe / r.stripeSize
		byteInBlock := byteInStripe % r.stripeSize

		// Calculate disk index (data disk for this block)
		diskIndex := blockIndex
		if diskIndex >= r.dataDisks {
			return errors.New("RAID6: invalid disk index calculation")
		}

		// Calculate stripe offset in disk (same for all disks in stripe)
		stripeOffset := stripe * r.stripeSize
		diskOffset := stripeOffset + byteInBlock

		// Expand disk if needed
		if diskOffset >= len(r.disks[diskIndex]) {
			needed := diskOffset - len(r.disks[diskIndex]) + 1
			r.disks[diskIndex] = append(r.disks[diskIndex], make([]byte, needed)...)
		}

		// Write data byte to data disk
		r.disks[diskIndex][diskOffset] = data[i]
	}

	// Calculate parity for all affected stripes
	for stripe := startStripe; stripe <= (offset+len(data)-1)/stripeDataSize; stripe++ {

		// Prepare data blocks for parity calculation
		dataBlocks := make([][]byte, r.dataDisks)
		for i := 0; i < r.dataDisks; i++ {
			start := stripe * r.stripeSize
			end := start + r.stripeSize
			if start >= len(r.disks[i]) {
				dataBlocks[i] = make([]byte, r.stripeSize)
				continue
			}
			if end > len(r.disks[i]) {
				end = len(r.disks[i])
			}
			block := make([]byte, r.stripeSize)
			copy(block, r.disks[i][start:end])
			dataBlocks[i] = block
		}

		// Calculate P and Q parity
		pParity := make([]byte, r.stripeSize)
		qParity := make([]byte, r.stripeSize)
		for i := 0; i < r.stripeSize; i++ {
			var p, q byte
			for j := 0; j < r.dataDisks; j++ {
				p ^= dataBlocks[j][i]
				q ^= gfMultiply(dataBlocks[j][i], byte(j+1))
			}
			pParity[i] = p
			qParity[i] = q
		}

		// Write parity to disks
		pDisk := r.dataDisks
		qDisk := r.dataDisks + 1
		stripeOffset := stripe * r.stripeSize

		// Expand parity disks if needed
		if stripeOffset+r.stripeSize > len(r.disks[pDisk]) {
			needed := stripeOffset + r.stripeSize - len(r.disks[pDisk])
			r.disks[pDisk] = append(r.disks[pDisk], make([]byte, needed)...)
		}
		copy(r.disks[pDisk][stripeOffset:], pParity)

		if stripeOffset+r.stripeSize > len(r.disks[qDisk]) {
			needed := stripeOffset + r.stripeSize - len(r.disks[qDisk])
			r.disks[qDisk] = append(r.disks[qDisk], make([]byte, needed)...)
		}
		copy(r.disks[qDisk][stripeOffset:], qParity)
	}
	return nil
}

func (r *RAID6) ClearDisk(diskIndex int) {
	if diskIndex >= 0 && diskIndex < len(r.disks) {
		for i := range r.disks[diskIndex] {
			r.disks[diskIndex][i] = 0
		}
	}
}
