package main

import (
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	languagecodes "github.com/spywiree/langcodes"
	"github.com/spywiree/translateimage"
	// "example.com/google_translateimages/translateimage"
)

func TranslateImages(imageURLs []string, downloadDir, outputDir string, sourceLaguage, targetLanguage languagecodes.LanguageCode) (map[string]string, error) {

	// 调用函数下载图片
	downloadedFilenames, err := DownloadImagesFromURLs(imageURLs, downloadDir)
	if err != nil {
		log.Print("downloadedFilenames err", err)
	}
	// 读取下载的图片，并调用TranslateAndSaveImage函数，将翻译后的图片写入到outputDir目录中
	var wg sync.WaitGroup

	for dfilename := range downloadedFilenames {
		wg.Add(1)
		go func() {

			err := TranslateAndSaveImage(downloadDir, downloadedFilenames[dfilename], outputDir, downloadedFilenames[dfilename], sourceLaguage, targetLanguage)
			if err != nil {
				log.Print("TranslateAndSaveImage err", err)
				// return nil, fmt.Errorf("TranslateAndSaveImage err: %w", err)
			}

			wg.Done()
		}()
	}
	wg.Wait()
	log.Printf("Translate images output Sucessfull")
	return downloadedFilenames, nil
}

// DownloadImagesFromURLs 从URL数组下载图片，并返回下载的图片名称数组
func DownloadImagesFromURLs(imageURLs []string, downloadDir string) (map[string]string, []error) {
	// var downloadedFilenames []string // 用于存储下载的图片名称的数组
	var translateFilenamesMap map[string]string = make(map[string]string)

	var errors []error

	for i, imageURL := range imageURLs {
		// 生成输出文件名，这里简单地使用索引和.png后缀
		downloadFilename := "image" + strconv.Itoa(i) + ".png"

		// 调用DownloadImageFromURL下载图片
		DownloadImageFromURL(imageURL, downloadDir, downloadFilename)

		// 将下载的图片文件名添加到数组中
		// downloadedFilenames = append(downloadedFilenames, downloadFilename)
		translateFilenamesMap[imageURL] = downloadFilename

	}
	return translateFilenamesMap, errors // 返回下载的图片名称数组
}

// DownloadImageFromURL 下载图片并保存到指定目录
// DownloadImageFromURL downloads an image from a given URL and saves it to a specified directory.
func DownloadImageFromURL(imageURL, downloadDir, outputFilename string) {
	// 发送HTTP GET请求获取图片数据
	// Send an HTTP GET request to retrieve image data
	response, err := http.Get(imageURL)
	if err != nil {
		log.Printf("Error downloading image: %v", err)
		// return

	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Error downloading image: HTTP status %d", response.StatusCode)
		// return
	}

	// 构建输出图片文件的完整路径
	// Build the complete path for the output image file
	downloaPath := filepath.Join(downloadDir, outputFilename)

	// 创建输出目录，如果它不存在的话
	// Create the output directory if it doesn't exist
	err = os.MkdirAll(filepath.Dir(downloaPath), 0755)
	if err != nil {
		log.Printf("Error creating output directory: %v", err)
		// return
	}

	// 创建输出文件
	// Create the output file
	downloaFile, err := os.Create(downloaPath)
	if err != nil {
		log.Printf("Error creating output image file: %v", err)
		// return
	}
	defer downloaFile.Close()

	// 将图片数据写入输出文件
	// Write the image data to the output file
	_, err = io.Copy(downloaFile, response.Body)
	if err != nil {
		log.Printf("Error saving image to file: %v", err)

	}

	log.Printf("Image downloaded and saved to %s", downloaPath)
}

// TranslateAndSaveImage 翻译并保存指定目录中的图片
// TranslateAndSaveImage translates and saves the image in the specified directory
func TranslateAndSaveImage(downloadDir, downloadFilename, outputDir, outputFilename string, sourceLaguage, targetLanguage languagecodes.LanguageCode) error {
	// 构建输入图片文件的完整路径
	// Build the complete path for the input image file
	inputPath := filepath.Join(downloadDir, downloadFilename)
	inputPathAbs, err := filepath.Abs(inputPath)
	if err != nil {
		log.Printf("Error getting absolute path for input image: %v", err)
		// return err
	}

	// 翻译图片
	// Translate the image
	img, err := translateimage.TranslateFile(
		// inputPathAbs, languagecodes.DETECT_LANGUAGE, languagecodes.ENGLISH,
		inputPathAbs, sourceLaguage, targetLanguage,
	)
	if err != nil {
		log.Printf("Error translating image: %v", err)
		// return err
	}

	// 构建输出图片文件的完整路径
	// Build the complete path for the output image file
	outputPath := filepath.Join(outputDir, outputFilename)

	// 创建输出目录，如果它不存在的话
	// Create the output directory if it doesn't exist
	err = os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		log.Printf("Error creating output directory: %v", err)
		return err
	}

	// 创建输出文件
	// Create the output file
	f, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Error creating output image file: %v", err)
		return err
	}
	defer f.Close()

	// 将翻译后的图片编码为PNG格式并写入文件
	// Encode the translated image in PNG format and write it to the file
	err = png.Encode(f, img)
	if err != nil {
		log.Printf("Error encoding and saving output image: %v", err)
		return err
	}

	log.Printf("Image translated and saved to %s", outputPath)
	return nil
}
