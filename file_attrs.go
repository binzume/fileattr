package main

import (
	"strconv"
)

type fileAttrs struct {
	Mode                 uint32
	LastModificationTime int64
	CreationTime         int64
	LastAccessTime       int64
	Size                 int64
}

func fromArray(row []string) *fileAttrs {
	mode, _ := strconv.ParseUint(row[1], 10, 32)
	mtime, _ := strconv.ParseInt(row[2], 10, 64)
	ctime, _ := strconv.ParseInt(row[3], 10, 64)
	atime, _ := strconv.ParseInt(row[4], 10, 64)
	var size int64 = 0
	if len(row) > 5 {
		size, _ = strconv.ParseInt(row[5], 10, 64)
	}
	return &fileAttrs{
		Mode:                 uint32(mode),
		LastModificationTime: mtime,
		CreationTime:         ctime,
		LastAccessTime:       atime,
		Size:                 size,
	}
}
