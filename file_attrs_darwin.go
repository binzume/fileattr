package main

import (
	"os"
	"syscall"
	"time"
)

func getFileAttrs(fi os.FileInfo) *fileAttrs {
	stat := fi.Sys().(*syscall.Stat_t)
	return &fileAttrs{
		Mode:                 uint32(fi.Mode() & ^os.ModeDir),
		LastModificationTime: fi.ModTime().UnixNano(),
		CreationTime:         syscall.TimespecToNsec(stat.Ctimespec),
		LastAccessTime:       syscall.TimespecToNsec(stat.Atimespec),
	}
}

func setFileAttrs(path string, fi os.FileInfo, attrs *fileAttrs) error {
	return os.Chtimes(path,
		time.Unix(0, attrs.LastAccessTime),
		time.Unix(0, attrs.LastModificationTime))
}

func isModified(attrs1 *fileAttrs, attrs2 *fileAttrs) bool {
	return attrs1.Mode != attrs2.Mode ||
		attrs1.LastModificationTime != attrs2.LastModificationTime ||
		attrs1.LastAccessTime != attrs2.LastAccessTime
}
