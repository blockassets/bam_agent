package tool

import "os"

func IsDirectory(dir string) (bool, error) {
	dirFile, err := os.Open(dir)
	defer dirFile.Close()
	if err != nil {
		return false, err
	}
	stat, err := dirFile.Stat()
	if os.IsNotExist(err) || err != nil {
		return false, err
	}

	return stat.IsDir(), nil
}

func IsExecutable(mode os.FileMode) bool {
	return !mode.IsDir() && (mode.Perm()&0111) > 0
}
