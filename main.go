package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type fileAttrs struct {
	Mode                 uint32
	LastModificationTime int64
	CreationTime         int64
	LastAccessTime       int64
}

func fromArray(row []string) *fileAttrs {
	mode, _ := strconv.ParseUint(row[1], 10, 32)
	mtime, _ := strconv.ParseInt(row[2], 10, 64)
	ctime, _ := strconv.ParseInt(row[3], 10, 64)
	atime, _ := strconv.ParseInt(row[4], 10, 64)
	return &fileAttrs{
		Mode:                 uint32(mode),
		LastModificationTime: mtime,
		CreationTime:         ctime,
		LastAccessTime:       atime,
	}
}

func dumpAttrs(targetPath string, listPath string) error {
	list, err := os.Create(listPath)
	if err != nil {
		return err
	}
	defer list.Close()
	err = filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name, _ := filepath.Rel(targetPath, path)
		name = filepath.ToSlash(name)
		attrs := getFileAttrs(info)
		log.Println(name)
		fmt.Fprintf(list, "%s\t%d\t%d\t%d\t%d\n", name, attrs.Mode, attrs.LastModificationTime, attrs.CreationTime, attrs.LastAccessTime)
		return nil
	})

	return err
}

func restoreAttrs(targetPath string, listPath string, restore bool) error {
	listFile, err := os.Open(listPath)
	if err != nil {
		return err
	}
	defer listFile.Close()
	reader := csv.NewReader(listFile)
	reader.Comma = '\t'

	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	files := map[string]*fileAttrs{}
	for _, row := range rows {
		files[filepath.ToSlash(row[0])] = fromArray(row)
	}

	err = filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name, _ := filepath.Rel(targetPath, path)
		name = filepath.ToSlash(name)
		attrs := getFileAttrs(info)
		if savedAttrs, ok := files[name]; ok {

			if isModified(attrs, savedAttrs) {
				if attrs.Mode != savedAttrs.Mode {
					_ = os.Chmod(path, os.FileMode(savedAttrs.Mode))
				}
				log.Println("Updated: ", name)
				if restore {
					setFileAttrs(path, info, savedAttrs)
				}
			}
		} else {
			log.Println("Unknown file:", name)
		}
		return nil
	})
	return err
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -l ATTR_LIST.tsv -m save|compare|restore TARGET_DIR\n", os.Args[0])
		flag.PrintDefaults()
	}
	listPath := flag.String("l", "", "attr list")
	mode := flag.String("m", "", "save|compare|restore")
	flag.Parse()
	targetPath := flag.Arg(0)
	if *mode == "" || *listPath == "" || targetPath == "" {
		flag.Usage()
		return
	}

	switch *mode {
	case "save":
		err := dumpAttrs(targetPath, *listPath)
		if err != nil {
			log.Fatal(err)
		}
		break
	case "compare":
		err := restoreAttrs(targetPath, *listPath, false)
		if err != nil {
			log.Fatal(err)
		}
		break
	case "restore":
		err := restoreAttrs(targetPath, *listPath, true)
		if err != nil {
			log.Fatal(err)
		}
		break
	default:
		flag.Usage()
	}
}
