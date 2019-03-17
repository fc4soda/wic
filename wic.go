package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"bufio"
	"os"
	"path/filepath" 
	"strings"
	"errors"
)

func main(){
	const defaultDatPath = "./dat"
	const defaultJpgPath = "./jpg"
	
	fmt.Printf("\n\n\t\t## 批量转换微信 dat 文件为 jpg ##\n")
	fmt.Printf("\n\t\t代码改编自 github@Seraphli/wic\n")
	
	fmt.Printf("\n操作步骤:\n")
	fmt.Printf("\n1. 把目录中的 color_sheet.jpg 图片发给一个微信好友;\n")
	fmt.Printf("\n2. 在 WeChat Files\\你的微信id\\Data 目录下按创建时间倒序排列文件, 会看到产生了一个刚刚创建的 dat 文件, 把这个文件复制到当前目录下, 重命名为 color_trans.dat;\n")
	fmt.Printf("\n3. 完成上面两步之后, 按以下提示操作.\n")
		
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n输入地址 或者 拖动微信 Data 文件目录到窗口:\n( 直接回车则使用当前目录下的 dat 目录 )\n")
	datPath, err := reader.ReadString('\n')
	CheckErr(err)
	datPath, err = CheckInputIsPath(datPath, defaultDatPath, false)
	
	fmt.Printf("\n输入地址 或者 拖动图片保存目录到窗口:\n( 直接回车则使用当前目录下的 jpg 目录 )\n")
	jpgPath, err := reader.ReadString('\n')
	CheckErr(err)
	jpgPath, err = CheckInputIsPath(jpgPath, defaultJpgPath, true)
	
	Dat2Jpg(datPath, jpgPath)

	fmt.Printf("\n按返回键退出...\n")
	fmt.Scanln()
}

func Dat2Jpg(datPath string, jpgPath string){
	sampleDat := "color_trans.dat"
	sampleJpg := "color_sheet.jpg"
	bytesSampleDat, err := ioutil.ReadFile(sampleDat)
	CheckErr(err)
	bytesSampleJpg, err := ioutil.ReadFile(sampleJpg)
	CheckErr(err)
	
	var rule [256]byte
	var mapped [256]bool
	var mappedCount int32
	for i, v := range bytesSampleDat{
		if mapped[int(v)] == false {
			rule[int(v)] = bytesSampleJpg[i]
			mapped[int(v)] = true
			mappedCount ++
		}
		if mappedCount == 256 {
			break
		}
	}

	files, err := ioutil.ReadDir(datPath)
	if err != nil {
		log.Fatal(err)
	}
	
	for i, f := range files{
		if f.IsDir() {
			fmt.Printf("( 忽略目录 %v )\n", f.Name())
			continue
		}
		if filepath.Ext(strings.TrimSpace(f.Name())) != ".dat"{
			fmt.Printf("( 忽略非 dat 文件 %v )\n", f.Name())
			continue
		}

		datPath := filepath.Join(datPath, f.Name())
		fmt.Printf("处理文件 %5d: %v\n", i+1, f.Name())
		bytesDat, err := ioutil.ReadFile(datPath)
		CheckErr(err)
		var bytesJpg []byte
		for _, byte := range bytesDat {
			bytesJpg = append(bytesJpg, rule[byte])
		}
		
		jpgPath := filepath.Join(jpgPath, f.Name()+".jpg")
		err = ioutil.WriteFile(jpgPath, bytesJpg, 0644)
		CheckErr(err)
		fmt.Printf("处理完成 %5d: %v\n", i+1, f.Name()+".jpg")
	}
}

func CheckErr(err error){
	if err != nil {
		log.Printf("\n%T\n%s\n%#v\n", err, err, err)
	}
}	
	
func CheckInputIsPath(input string, defaultValue string, createPath bool)(path string, err error){
	input = strings.Trim(input, "\n")
	input = strings.TrimSpace(input)
	if len(input)>2 {
		if c:=input[len(input)-1]; input[0]==c && (c=='"'||c=='\'') {
			input = input[1:len(input)-1]
		}
	}
	if len(input)>0 {
		fmt.Printf("路径符合要求, 使用路径 %v\n", input)
		path = input
		err = nil
	} else {
		fmt.Printf("路径不符合要求, 使用默认路径 %v\n", defaultValue)
		path = defaultValue
		err = errors.New("路径为空")
	}	

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("目录不存在, ")
		if createPath {
			err = os.MkdirAll(path, 0644)
			CheckErr(err)
			if err==nil {
				fmt.Printf("创建目录 %v 成功\n\n", path)
			} else {
				fmt.Printf("创建目录 %v 失败\n\n", path)
				err = errors.New("创建目录 %v 失败")
			}
		} else {
			fmt.Printf("使用默认路径: %v\n\n", defaultValue)
			return defaultValue, errors.New("路径不存在")
		}
	}

	return path, err
}	
