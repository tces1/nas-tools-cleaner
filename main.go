package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func checkSymlink(path string) {
	// 获取软链接的实际路径
	dest, err := os.Readlink(path)
	if err != nil {
		fmt.Println("Error reading symlink:", err)
		return
	}

	// 检查软链接的目标是否存在
	_, err = os.Stat(dest)
	if err != nil {
		// 如果目标不存在，则输出无效软链接并删除它
		fmt.Printf("Invalid symlink: %s -> %s (Destination does not exist). Deleting symlink...\n", path, dest)
		err = os.Remove(path)
		if err != nil {
			fmt.Printf("Error deleting symlink %s: %v\n", path, err)
		} else {
			fmt.Printf("Symlink %s deleted.\n", path)
		}
		return
	}
}

func traverseDirectory(dir string) {
	// 标记是否找到任何软链接或子目录
	foundSymlink := false
	hasSubDir := false

	// 遍历目录
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// 如果目录已删除，跳过该路径
			if os.IsNotExist(err) {
				return nil
			}
			fmt.Println("Error walking the path:", err)
			return err
		}

		// 如果是软链接，处理它
		if info.Mode()&os.ModeSymlink != 0 {
			foundSymlink = true
			checkSymlink(path)
		}

		// 如果是子目录，设置标志
		if info.IsDir() && path != dir {
			hasSubDir = true
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error traversing directory:", err)
		return
	}

	// 如果没有软链接也没有子目录，输出提示并删除目录
	if !foundSymlink && !hasSubDir {
		fmt.Printf("Directory %s does not contain symlinks or subdirectories. Deleting...\n", dir)
		if err:=os.RemoveAll(dir); err != nil{
			fmt.Print(err)
		}
	}
}

// 递归遍历目录并处理所有子目录
func recursiveTraverse(dir string) {
	// 首先遍历当前目录
	traverseDirectory(dir)

	// 遍历当前目录的所有子目录
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// 如果目录已删除，跳过该路径
			if os.IsNotExist(err) {
				return nil
			}
			fmt.Println("Error walking the path:", err)
			return err
		}
		// 如果是目录，递归遍历
		if info.IsDir() && path != dir {
			recursiveTraverse(path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error traversing directory:", err)
	}
}

func main() {
	// 监视的目录路径
	dirPath := "/rclone/emby"

	// 遍历并处理软链接
	recursiveTraverse(dirPath)
}
