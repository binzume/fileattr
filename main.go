package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func dumpAttrs(targetPath string, listPath string) error {
	list, err := os.Create(listPath)
	if err != nil {
		return err
	}
	defer list.Close()
	var errs []error
	err = filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errs = append(errs, err)
			return nil
		}
		name, _ := filepath.Rel(targetPath, path)
		name = filepath.ToSlash(name)
		attrs := getFileAttrs(info)
		log.Println(name)
		fmt.Fprintf(list, "%s\t%d\t%d\t%d\t%d\t%d\n", name, attrs.Mode, attrs.LastModificationTime, attrs.CreationTime, attrs.LastAccessTime, attrs.Size)
		return nil
	})
	if errs != nil {
		for _, err := range errs {
			log.Print(err)
		}
		err = errs[0]
	}
	return err
}

func restoreAttrs(targetPath string, listPath string, restore, checkSize bool) error {
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
		files[row[0]] = fromArray(row)
	}

	var errs []error
	err = filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errs = append(errs, err)
			return nil
		}
		name, _ := filepath.Rel(targetPath, path)
		name = filepath.ToSlash(name)
		attrs := getFileAttrs(info)
		if savedAttrs, ok := files[name]; ok {
			if checkSize && attrs.Size != savedAttrs.Size {
				log.Println("Skip(size cahnged): ", name)
				return nil
			}

			if isModified(attrs, savedAttrs) {
				if attrs.Mode != savedAttrs.Mode {
					_ = os.Chmod(path, os.FileMode(savedAttrs.Mode))
				}
				if restore {
					log.Println("Restore: ", name)
					setFileAttrs(path, info, savedAttrs)
				} else {
					log.Println("Updated: ", name)
				}
			}
		} else {
			log.Println("Unknown file:", name)
		}
		return nil
	})
	if errs != nil {
		for _, err := range errs {
			log.Print(err)
		}
		err = errs[0]
	}
	return err
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -m save|compare|restore -l ATTR_LIST.tsv TARGET_DIR\n", os.Args[0])
		flag.PrintDefaults()
	}
	listPath := flag.String("l", "", "attr list")
	mode := flag.String("m", "", "save|compare|restore")
	checkSize := flag.Bool("s", false, "check file size")
	flag.Parse()

	targetPath := flag.Arg(0)
	if *listPath == "" || targetPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	var err error
	switch *mode {
	case "save":
		err = dumpAttrs(targetPath, *listPath)
		break
	case "compare":
		err = restoreAttrs(targetPath, *listPath, false, *checkSize)
		break
	case "restore":
		err = restoreAttrs(targetPath, *listPath, true, *checkSize)
		break
	default:
		flag.Usage()
		os.Exit(1)
	}
	if err != nil {
		log.Fatal(err)
	}
}
