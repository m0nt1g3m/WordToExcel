//go:build !windows

package utils

func hideTempFile(string) error {
	return nil
}
