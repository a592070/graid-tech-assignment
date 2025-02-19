Create a program that demonstrates RAID data handling for RAID0, RAID1, RAID10, RAID5, and RAID6.
1. Define the data type for RAID data storage, 
for example, an array(raid) of array(disk) of byte of array(data stripe).
2. Write a string([]byte) into the RAID at position 0 with a length N, 
where N should be greater than the size of one stripe.
3. Clear one of the disks in the RAID, setting it to zero.
4. Read the data in the RAID from 0 to N, convert it back to a string, and print it.