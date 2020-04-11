package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

// 原始路径
var rawPath string

// 结果路径
var resultPath string

// 项目根路径
var currentWd string

// 自动创建文件
func ensureDir(fileName string) {
	dirName := fileName
	// dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, 0755)
		if merr != nil {
			panic(merr)
		}
	}
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
					return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
					return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
					return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
					return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// 处理图片尺寸(只处理尺寸大于500的)
func changeSize(name string, path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("文件读取失败")
		log.Fatal(err)
	}
	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println("解码图片文件失败")
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio

	// w := img.Bounds().Dx()
	// h := img.Bounds().Dy()

	m := resize.Resize(0, 500, img, resize.Lanczos3)

	//m := resize.Resize(uint(kuan), 0, img, resize.Lanczos3)
	//m := resize.Resize(1000, 0, img, resize.Lanczos3)

	res := filepath.Join(resultPath, name)
	fmt.Println("输出结果: ", res)
	// return
	out, err := os.Create(res)

	if err != nil {
		log.Println("xxx 错误")
		log.Fatal(err)
	}

	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
	fmt.Println("------------------")
}

// 获取图片大小
func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}

// 获取目录文件
func readDir(path string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	return files, nil
}

// 输入文件夹和输出文件夹
var input, out string

// 宽高
// var width, height int

func init() {
	dir, _ := os.Getwd()
	args := os.Args
	if len(args) >= 3 {
		input = args[1]
		out = args[2]
	} else {
		input = "face"
		out = "result"
	}
	resultPath = filepath.Join(dir, out)
	rawPath = filepath.Join(dir, input)
	fmt.Println("原始路径: ", rawPath)
	fmt.Println("结果路径: ", resultPath)
	// 在处理之前, 创建一个文件夹
	ensureDir(resultPath)
}

// 自动处理
func auto() {
	lists, _ := readDir(rawPath)
	fmt.Println("共有", len(lists), "个文件")
	for _, item := range lists {
		name := item.Name()
		// 只判断 `jpg`
		isJpeg := strings.Contains(name, ".jpg")
		if !isJpeg {
			fmt.Println("文件格式不是`.jpg`, 跳过本次操作, ", name)
			continue
		}
		outPath := filepath.Join(rawPath, name)
		w, h := getImageDimension(outPath)
		if h <= 500 {
			fmt.Println("图片尺寸小于500, 不执行操作, 原样复制, ", name)
			rr := filepath.Join(resultPath, name)
			copy(outPath, rr)
			// TODO
			continue
		}
		changeSize(name, outPath)
		fmt.Printf("图片名称: %s 宽: %d, 高: %d\n", name, w, h)
	}
}
