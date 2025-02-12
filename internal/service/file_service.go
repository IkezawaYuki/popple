package service

import (
	"context"
	"fmt"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type fileService struct {
	httpClient *infrastructure.HttpClient
}

type FileService interface {
	DownloadMediaFiles(ctx context.Context, customerID int, post entity.InstagramPost) ([]string, error)
	MakeTempDirectory(customerID int) error
	RemoveTempDirectory(customerID int) error
}

func NewFileService(httpClient *infrastructure.HttpClient) FileService {
	return &fileService{
		httpClient: httpClient,
	}
}

func (f *fileService) DownloadMediaFiles(ctx context.Context, customerID int, post entity.InstagramPost) ([]string, error) {
	var fileList []string
	if len(post.Children.Data) == 0 {
		mediaPath, err := f.downloadMedia(ctx, customerID, post.MediaURL)
		if err != nil {
			return nil, err
		}
		fileList = append(fileList, mediaPath)
		return fileList, nil
	}
	for _, child := range post.Children.Data {
		mediaPath, err := f.downloadMedia(ctx, customerID, child.MediaURL)
		if err != nil {
			return nil, err
		}
		fileList = append(fileList, mediaPath)
	}
	return fileList, nil
}

func (f *fileService) downloadMedia(ctx context.Context, customerID int, mediaUrl string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", mediaUrl, nil)
	if err != nil {
		return "", err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	filename := strings.Split(filepath.Base(mediaUrl), "?")[0]
	filePath := filepath.Join(fmt.Sprintf(tempDirectory, customerID), filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = out.Close()
	}()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

const tempDirectory = "./tmp_%d"

func (f *fileService) MakeTempDirectory(customerID int) error {
	err := os.Mkdir(fmt.Sprintf(tempDirectory, customerID), 0777)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func (f *fileService) RemoveTempDirectory(customerID int) error {
	return os.RemoveAll(fmt.Sprintf(tempDirectory, customerID))
}
