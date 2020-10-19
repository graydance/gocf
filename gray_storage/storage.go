/**
google cloud storage相关方法
 */
package gray_storage

import (
	"bufio"
	"cloud.google.com/go/storage"
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
)

/**
从google cloud storage下载文件
*/
func DownloadObj(ctx context.Context, bucket, objName, destFile string) (err error) {
	//gcs handler
	client, err := storage.NewClient(ctx)
	defer func() {
		_ = client.Close()
	}()
	if err != nil {
		log.Printf("gcs newclient error: %v", err)
		return
	}
	handler := client.Bucket(bucket).Object(objName)
	//download handler
	f, err := os.Create(destFile)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		log.Printf("dest file create error: %v", err)
		return
	}
	dstWriter := bufio.NewWriter(f)
	objReader, err := handler.NewReader(ctx)
	if err != nil {
		log.Printf("gcs newreader error: %v", err)
		return
	}
	//copy file
	n, err := io.Copy(dstWriter, objReader)
	log.Printf("io copy size: %v\n", n)
	if err != nil {
		log.Printf("iocopy error: %v", err)
		return
	}
	_ = dstWriter.Flush()
	return
}

/**
向google cloud storage 上传文件
*/
func UploadObj(ctx context.Context, sourceFile, bucket, objName string) (err error) {
	//gcs handler
	client, err := storage.NewClient(ctx)
	defer func() {
		_ = client.Close()
	}()
	if err != nil {
		log.Printf("gcs newclient error: %v", err)
		return
	}
	whandler := client.Bucket(bucket).Object(objName)
	wc := whandler.NewWriter(ctx)
	defer func() {
		_ = wc.Close()
	}()
	wc.CacheControl = "public, max-age=31536000"
	thumbBytes, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		log.Printf("outil.ReadFile error: %+v\n", err)
		return
	}
	n, err := wc.Write(thumbBytes)
	log.Printf("thumb upload size: %d\n", n)
	if err != nil {
		log.Printf("wc.Write error: %+v\n", err)
		return
	}
	return
}
