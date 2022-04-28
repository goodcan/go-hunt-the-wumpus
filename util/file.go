/*
 * @Time     : 2022/4/9 18:22
 * @Author   : cancan
 * @File     : file.go
 * @Function :
 */

package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func CheckDir(dir string) {
	if !PathExists(dir) {
		_ = os.Mkdir(dir, os.ModePerm)
	}
}

func JoinPath(p1, p2 string) string {
	return filepath.Join(p1, p2)
}

func WriteFile(dir, filename string, context []byte) {
	CheckDir(dir)
	_ = ioutil.WriteFile(JoinPath(dir, filename), context, 0660)
}

func ReadFile(dir, filename string) []byte {
	path := JoinPath(dir, filename)
	if !PathExists(path) {
		return []byte{}
	}

	content, _ := ioutil.ReadFile(path)
	return content
}
