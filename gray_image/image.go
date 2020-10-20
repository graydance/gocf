package gray_image

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"log"
	"os"
)

/**
等比压缩
宽或高以最大的边为准等比压缩
*/
func ResizeProportional(rawFile, dstFile string, width, height int) (err error) {
	src, err := imaging.Open(rawFile)
	if err != nil {
		log.Printf("failed to open image: %v", err)
		return
	}
	x, y, err := CalcProportionalXY(rawFile, width, height)
	if err != nil {
		log.Printf("CalcProportionalXY error: %v", err)
		return
	}
	src = imaging.Resize(src, x, y, imaging.Lanczos)
	err = imaging.Save(src, dstFile)
	if err != nil {
		log.Printf("image save error:%v", err)
		return
	}

	return
}

/**
等比压缩,计算最终的宽高
*/
func CalcProportionalXY(rawFile string, width, height int) (x, y int, err error) {
	imgReader, err := os.Open(rawFile)
	if err != nil {
		log.Printf("rawFile open error:%v", err)
		return
	}
	defer func() {
		_ = imgReader.Close()
	}()
	//原始图片的信息
	img, _, err := image.Decode(imgReader)
	if err != nil {
		log.Printf("image decode error:%v", err)
		return
	}
	rawX := img.Bounds().Max.X
	rawY := img.Bounds().Max.Y
	if rawX/width == rawY/height {
		x = width
		y = height
		return
	} else if rawX/width > rawY/height {
		x = width
		y = 0
		return
	} else {
		x = 0
		y = width
		return
	}
}

/**
等比压缩,未完成
若给的比例和原比例不符,则生成一个背景,留边
*/
func test(file, whiteBacground string) (err error) {
	imgReader, err := os.Open(file)
	if err != nil {
		log.Printf("rawFile open error:%v", err)
		return
	}
	img, _, err := image.Decode(imgReader)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	b := img.Bounds()

	width := b.Max.X
	height := b.Max.Y
	h := width * 9 / 16
	if h > height { //高不够，填充
		backgroundImageFile, _ := os.Open(whiteBacground)
		backgroundImage, _, _ := image.Decode(backgroundImageFile)
		backgroundImage = imaging.CropCenter(backgroundImage, width, h) //理想宽高
		img = imaging.OverlayCenter(backgroundImage, img, 1)            //最后生成的图片
	} else {
		img = imaging.Resize(img, width, width*9/16, imaging.Lanczos) //按比例缩
	}
	return
}
