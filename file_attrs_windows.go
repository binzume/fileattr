package main

import (
	"os"
	"syscall"
)

func getFileAttrs(fi os.FileInfo) *fileAttrs {
	stat := fi.Sys().(*syscall.Win32FileAttributeData)
	return &fileAttrs{
		Mode:                 uint32(fi.Mode() & ^os.ModeDir),
		LastModificationTime: fi.ModTime().UnixNano(),
		CreationTime:         stat.CreationTime.Nanoseconds(),
		LastAccessTime:       stat.LastAccessTime.Nanoseconds(),
		Size:                 fi.Size(),
	}
}

func setFileAttrs(path string, info os.FileInfo, attrs *fileAttrs) error {
	pathp, e := syscall.UTF16PtrFromString(path)
	if e != nil {
		return e
	}
	h, e := syscall.CreateFile(pathp,
		syscall.FILE_WRITE_ATTRIBUTES, syscall.FILE_SHARE_WRITE, nil,
		syscall.OPEN_EXISTING, syscall.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if e != nil {
		return e
	}
	defer syscall.Close(h)
	c := syscall.NsecToFiletime(attrs.CreationTime)
	a := syscall.NsecToFiletime(attrs.LastAccessTime)
	w := syscall.NsecToFiletime(attrs.LastModificationTime)
	return syscall.SetFileTime(h, &c, &a, &w)
}

func isModified(attrs1 *fileAttrs, attrs2 *fileAttrs) bool {
	return *attrs1 != *attrs2
}
