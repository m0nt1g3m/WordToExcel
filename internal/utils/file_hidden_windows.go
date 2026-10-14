//go:build windows

package utils

import "golang.org/x/sys/windows"

func hideTempFile(path string) error {
	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	return windows.SetFileAttributes(pathPtr, windows.FILE_ATTRIBUTE_HIDDEN)
}
