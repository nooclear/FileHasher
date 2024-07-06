package FileHasher

import (
	"errors"
	"log"
	"os"
	"time"
)

func Run(fileLocation string) {
	if err := initDB("./database.db"); err != nil {
		panic(err)
	}
	getEntries(fileLocation)

	for _, file := range files {
		e := entry{
			timestamp: time.Now().UnixMilli(),
			filepath:  file,
			hash:      string(hashFile(file)),
		}
		if !find(&e) {
			if _, err := add(&e); err != nil {
				panic(err)
			}
		} else {
			if _, err := update(&e); err != nil {
				panic(err)
			}
		}
	}
}

func Remove() {
	entries, _ := getAll()
	for _, e := range entries {
		if _, err := os.Stat(e.filepath); errors.Is(err, os.ErrNotExist) {
			if !find(e) {
				log.Printf("file %s not found", e.filepath)
			} else {
				if _, err := del(e); err != nil {
					panic(err)
				}
			}
		}
	}
}
