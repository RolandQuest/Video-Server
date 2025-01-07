package localfiles

import (
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

type MediaType int

const (
	Movies MediaType = iota
	Images MediaType = iota
	Music MediaType = iota
)

const (
	pathMovies = "C:/Users/Paul/Videos"
	pathImages = "C:/Users/Paul/Pictures"
	pathMusic  = "C:/Users/Paul/Music"
)

var extensionMovies = [...]string { ".mp4" }
var extensionImages = [...]string { ".png", ".jpg" }
var extensionMusic = [...]string { ".mp3" }

func GetMediaDirectory(mediaType MediaType, subPath string) []fs.DirEntry {
	switch mediaType {
		case Movies:
			return GetMovies(subPath)
		case Images:
			return GetImages(subPath)
		case Music:
			return GetMusic(subPath)
	}
	return []fs.DirEntry{}
}




func GetMovies(subPath string) []fs.DirEntry {
	return getDirectory(pathMovies + subPath, extensionMovies[:])
}

func GetImages(subPath string) []fs.DirEntry {
	return getDirectory(pathImages + subPath, extensionImages[:])
}

func GetMusic(subPath string) []fs.DirEntry {
	return getDirectory(pathMusic + subPath, extensionMusic[:])
}

func FindMovie(title string) (movie_path string, err error) {
	err = filepath.WalkDir(pathMovies, func(path string, d fs.DirEntry, _ error) error {
		if !d.IsDir() && title == d.Name() {
			movie_path = path
			return io.EOF
		}
		return nil
	})
	
	if err == io.EOF {
		err = nil
	}
	return
}

func FindImage(title string) (image_path string, err error) {
	return pathImages + "/" + title, nil
}

func filterEntries(entries []fs.DirEntry, test func(fs.DirEntry) bool) []fs.DirEntry {
	ret := []fs.DirEntry{}
	for _, e := range entries {
		if(test(e)) {
			ret = append(ret, e)
		}
	}
	return ret
}

func containsValidFile(entryToTest fs.DirEntry, thePath string, extensions []string) bool {
	entries, _ := os.ReadDir(thePath + "/" + entryToTest.Name())
	for _, entry := range(entries) {
		if entry.IsDir() && containsValidFile(entry, thePath + "/" + entryToTest.Name(), extensions) {
			return true
		}
		if slices.Contains(extensions, path.Ext(entry.Name())) {
			return true
		}
	}
	return false
}

func checkValidEntry(entry fs.DirEntry, thePath string, extensions []string) bool {
	if entry.IsDir() {
		if entry.Name() == "Temp Explorer" {
			return false
		}
		if !containsValidFile(entry, thePath, extensions) {
			return false
		}
	} else {
		if !slices.Contains(extensions, path.Ext(entry.Name())) {
			return false
		}
	}
	
	return true
}

func getEntries(thePath string, extensions []string) []fs.DirEntry {
	entries, err := os.ReadDir(thePath)
	if err != nil {
		return []fs.DirEntry{}
	}
	return filterEntries(entries, func (entry fs.DirEntry) bool {
		return checkValidEntry(entry, thePath, extensions)
	})
}

func getDirectory(path string, extensions []string) []fs.DirEntry {
	entries := getEntries(path, extensions)
	slices.SortStableFunc(entries, func (a fs.DirEntry, b fs.DirEntry) int {
		if a.IsDir() && !b.IsDir() {
			return -1
		}
		if !b.IsDir() && b.IsDir() {
			return 1
		}
		return strings.Compare(a.Name(), b.Name()) 
	})
	return entries
}

func getSubDirectories(path string, extensions []string) []fs.DirEntry {
	entries := getEntries(path, extensions)
	return filterEntries(entries, func (entry fs.DirEntry) bool {
		return entry.IsDir()
	})
}

func getFiles(path string, extensions []string) []fs.DirEntry {
	entries := getEntries(path, extensions)
	return filterEntries(entries, func (entry fs.DirEntry) bool {
		return !entry.IsDir()
	})
}
