package raid

type RAID interface {
	Read(length int, pos int) ([]byte, error)
	Write(data []byte, pos int) error
	ClearDisk(diskIndex int)
}
