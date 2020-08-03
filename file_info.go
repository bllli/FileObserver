package main

import "time"

type FileInfo struct {
	Name    string // base name of the file
	Size    int64  // length in bytes for regular files; system-dependent for others
	IsDir   bool   // abbreviation for Mode().IsDir()
	ModTime time.Time
}

type FileInfoSlice []FileInfo

func (f FileInfoSlice) LessBySize(i, j int) bool {
	if f[i].IsDir && !f[j].IsDir {
		return false
	}
	if !f[i].IsDir && f[j].IsDir {
		return true
	}
	return f[i].Size < f[j].Size
}

func (f FileInfoSlice) LessByModTime(i, j int) bool {
	if f[i].IsDir && !f[j].IsDir {
		return false
	}
	if !f[i].IsDir && f[j].IsDir {
		return true
	}
	return f[i].ModTime.Before(f[j].ModTime)
}
