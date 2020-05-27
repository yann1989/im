package file_log

import "os"

func IsDir(dirname string) bool {
	fhandler, err := os.Stat(dirname)
	if !(err == nil || os.IsExist(err)) {
		return false
	} else {
		return fhandler.IsDir()
	}
}

func IsFile(filename string) bool {
	fhandler, err := os.Stat(filename)
	if !(err == nil || os.IsExist(err)) {
		return false
	} else if fhandler.IsDir() {
		return false
	}
	return true
}

func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileByteSize(filename string) (bool, int64) {
	if !IsFile(filename) {
		return false, 0
	}
	fhandler, _ := os.Stat(filename)
	return true, fhandler.Size()
}

func Write(file *os.File, data string) (bool, error) {
	_, err := file.WriteString(data)
	if err != nil {
		return false, err
	}
	return true, nil
}
