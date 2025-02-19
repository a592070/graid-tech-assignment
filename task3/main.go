package main

import (
	"fmt"
	"graid-tech-assignment/pkg/task3/raid"
	"log"
)

func main() {
	raids := make(map[string]raid.RAID)
	raid0, err := raid.NewRAID0(2, 8)
	if err != nil {
		panic(err)
	}
	raids["RAID0"] = raid0

	raid1, err := raid.NewRAID1(2)
	if err != nil {
		panic(err)
	}
	raids["RAID1"] = raid1

	raid10, err := raid.NewRAID10(4, 8)
	if err != nil {
		panic(err)
	}
	raids["RAID10"] = raid10

	raid5, err := raid.NewRAID5(3, 8)
	if err != nil {
		panic(err)
	}
	raids["RAID5"] = raid5

	raid6, err := raid.NewRAID6(4, 8)
	if err != nil {
		panic(err)
	}
	raids["RAID6"] = raid6

	for k, r := range raids {
		data := []byte(fmt.Sprintf("Hello, %s!", k))
		log.Printf("Write to %s: %s", k, data)
		if err := r.Write(data, 0); err != nil {
			log.Fatal(err)
		}
		r.ClearDisk(0)

		read, err := r.Read(len(data), 0)
		if err != nil {
			log.Printf("Failed to read from %s: %v", k, err)
			continue
		}
		log.Printf("Read data after clearing disk 0: %q\n", read)
	}
}
