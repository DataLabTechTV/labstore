//go:build !embedassets

package router

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/IllumiKnowLabs/labstore/server/constants"
	"github.com/IllumiKnowLabs/labstore/server/helper"
)

var assetsDir string

func init() {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = ".labstore"
		slog.Warn("could not find user cache dir")
	}

	assetsDir = filepath.Join(cacheDir, "assets")
	slog.Debug("assets dir", "path", assetsDir)
}

func NewWebServerDescriptor(host string, port uint16) (*ServerDescriptor, error) {
	slog.Info("web ui server", "host", host, "port", port)
	slog.Debug("using runtime downloaded web ui assets")

	addr := fmt.Sprintf("%s:%d", host, port)

	assetsFS, err := loadAssets()
	if err != nil {
		return nil, fmt.Errorf("web server descriptor: %w", err)
	}

	contentFS := helper.Must(fs.Sub(assetsFS, "assets"))
	httpFS := http.FS(contentFS)
	handler := http.FileServer(httpFS)

	server := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	serverDescriptor := &ServerDescriptor{
		Name:   "Web UI",
		Server: &server,
	}

	return serverDescriptor, nil
}

func loadAssets() (fs.FS, error) {
	if helper.FileExists(assetsDir) {
		if !helper.IsDir(assetsDir) {
			return nil, fmt.Errorf("not a directory: %s", assetsDir)
		}
	} else {
		slog.Info("downloading web ui assets from github release")
		if err := fetchAssets(); err != nil {
			return nil, err
		}
	}

	return os.DirFS(assetsDir), nil
}

func fetchAssets() error {
	url := constants.GitRepo + "/releases/download/" + constants.GitTag + "/" + constants.GitAssetsFilename

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer helper.CloseWithErr(resp.Body, &err)

	tmpZip := filepath.Join(os.TempDir(), constants.GitAssetsFilename)
	zipWriter, err := os.Create(tmpZip)
	if err != nil {
		return err
	}
	defer helper.CloseWithErr(zipWriter, &err)

	_, err = io.Copy(zipWriter, resp.Body)
	if err != nil {
		return err
	}

	zipReader, err := zip.OpenReader(tmpZip)
	if err != nil {
		return err
	}
	defer helper.CloseWithErr(zipReader, &err)

	for _, file := range zipReader.File {
		outPath := filepath.Join(assetsDir, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(outPath, 0755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}

		inFile, err := file.Open()
		if err != nil {
			return err
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			return err
		}
		_, err = io.Copy(outFile, inFile)

		helper.CloseWithErr(inFile, &err)
		helper.CloseWithErr(outFile, &err)

		if err != nil {
			return err
		}
	}

	return nil
}
