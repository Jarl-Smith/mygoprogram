package fileencrypt

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Start() {
	var inputStr string

	fmt.Println("选择功能:\n1.后缀名加密\n2.base64加密\n3.后缀名解密\n4.base64解密")
	//接收用户的选项
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		inputStr = scanner.Text()
	}
	//对用户输入的选项进行检验
	inputInt, err := strconv.Atoi(inputStr) //判断输入是否为数字
	if err != nil {
		log.Fatal(err)
	} else if inputInt < 0 || inputInt > 4 { //判断输入选项是否在选项内
		log.Fatal("输入无效选项")
	}
	//接收用户输入的目标路径
	fmt.Println("输入绝对路径：")
	if scanner.Scan() {
		inputStr = scanner.Text()
	}
	//对用户输入的路径进行检验
	if !filepath.IsAbs(inputStr) { //判断是否为路径
		log.Fatal("输入路径无效")
	}
	_, err = os.Stat(inputStr) //判断路径是否存在
	if err != nil {
		log.Fatal(err)
	}
	//执行加密或解密操作
	Execute(inputStr, inputInt)
}

// 执行加密或解密文件的函数，directoryStr为文件夹，mode为用户的选择
func Execute(directoryStr string, mode int) {
	var action func(input string) string
	chooseEncrypt := false
	switch mode {
	case 1:
		action = encryptWithSuffix
		chooseEncrypt = true
	case 2:
		action = encryptWithBase64
		chooseEncrypt = true
	case 3:
		action = decryptWithSuffix
	case 4:
		action = decryptWithBase64
	default:
		action = func(input string) string { return input }
	}
	fileEntrys, err := os.ReadDir(directoryStr) //获取给定文件夹内的所有文件
	if err != nil {
		fmt.Println(err)
		return
	}
	//依次对文件名称进行加密或解密
	for _, entry := range fileEntrys {
		info, err := entry.Info()
		if err != nil {
			fmt.Println(err)
		}
		if info.IsDir() { //如果是文件夹，则不进行任何处理
			continue
		}
		//给定的目录存在已加密和未加密的文件,根据用户的选项应该选择性执行加密或解密操作
		if chooseEncrypt { //若用户选择的加密
			if isFileNoEncrypt(entry.Name()) { //只对未加密的文件进行加密操作
				newFileName := action(entry.Name())
				oldPath := filepath.Join(directoryStr, entry.Name())
				newPath := filepath.Join(directoryStr, newFileName)
				os.Rename(oldPath, newPath)
			}
		} else { //若用户选择解密
			if !isFileNoEncrypt(entry.Name()) { //只对已加密的文件进行解密
				newFileName := action(entry.Name())
				oldPath := filepath.Join(directoryStr, entry.Name())
				newPath := filepath.Join(directoryStr, newFileName)
				os.Rename(oldPath, newPath)
			}
		}
	}
	fmt.Println("执行完毕")
}

// 对文件夹里的文件判断是否已经加密了
func isFileNoEncrypt(fileName string) bool {
	//如果文件名称不以.bin结尾且有.,视该文件未加密
	if !strings.HasSuffix(fileName, ".bin") && strings.LastIndex(fileName, ".") > 0 {
		return true
	}
	return false
}

// 将文件名称进行base64加密，得到的密文作为新的文件名
func encryptWithBase64(fileName string) string {
	fileName = base64.StdEncoding.EncodeToString([]byte(fileName))
	fileName = strings.ReplaceAll(fileName, "/", "_") //密文可能会有'/'字符，这个字符会导致重命名失败，因此需要将'/'替换
	return fileName
}

// 将密文解密，解密后的明文即为原文件名称
func decryptWithBase64(fileName string) string {
	fileName = strings.ReplaceAll(fileName, "_", "/")
	bytes, err := base64.StdEncoding.DecodeString(fileName)
	if err != nil {
		fmt.Println(err)
		return fileName
	}
	return string(bytes)
}

// 通过添加.bin达到加密效果
func encryptWithSuffix(fileName string) string {
	return fileName + ".bin"
}

// 通过去掉.bin后缀达到解密效果
func decryptWithSuffix(fileName string) string {
	if strings.HasSuffix(fileName, ".bin") {
		end := strings.LastIndex(fileName, ".bin")
		return fileName[:end]
	}
	return fileName
}
