package providers

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/IllumiKnowLabs/labstore/cli/credentials"
	"github.com/IllumiKnowLabs/labstore/cli/errs"
	"github.com/IllumiKnowLabs/labstore/cli/tui/messages"
	"github.com/IllumiKnowLabs/labstore/cli/types"
	"github.com/IllumiKnowLabs/labstore/client/s3"
	clientTypes "github.com/IllumiKnowLabs/labstore/client/types"
	"github.com/IllumiKnowLabs/labstore/server/config"
	"github.com/IllumiKnowLabs/labstore/server/helper"
)

type S3FSProvider struct {
	Active  bool
	Host    string
	Port    uint16
	Profile string
	Bucket  string
	Key     string

	lastSelected map[string]string
}

type IndexedPath struct {
	FileIndex int
	Path      string
}

type S3IndexedPath struct {
	FileIndex int
	Bucket    string
	Key       string
}

type WalkDirFunc func(path string, d *types.Entry, err error) error

func NewS3FSProvider() *S3FSProvider {
	// TODO: keep state of last selected bucket and key (both might not exist anymore)
	return &S3FSProvider{
		Host: config.App.Server.S3.Address.Host,
		Port: config.App.Server.S3.Address.Port,

		lastSelected: make(map[string]string),
	}
}

func (p *S3FSProvider) Select(args ...string) error {
	if len(args) < 2 {
		return &errs.ErrInsufficientArguments{}
	}

	p.Active = true
	p.Profile = args[0]
	p.Bucket = args[1]

	lastSelKey := p.lastSelKey()

	if len(args) > 2 {
		p.Key = path.Clean(args[2])
		if p.Key == "." {
			p.Key = ""
		} else {
			p.Key += "/"
		}
	}

	p.lastSelected[lastSelKey] = p.Key

	return nil
}

func (p *S3FSProvider) Deselect() {
	p.Active = false
	p.Profile = ""
	p.Bucket = ""
	p.Key = ""
	p.lastSelected = make(map[string]string)
}

func (p *S3FSProvider) Selected() string {
	return p.Key
}

func (p *S3FSProvider) LastSelected() (string, bool) {
	state, ok := p.lastSelected[p.lastSelKey()]
	return state, ok

}

func (p *S3FSProvider) Children() ([]types.Entry, error) {
	if !p.Active {
		return []types.Entry{}, nil
	}

	client, err := p.newClient()
	if err != nil {
		return nil, err
	}

	resCh := client.ListObjects(p.Bucket, p.Key, true)
	var entries []types.Entry

	if keyParts := strings.Split(p.Key, "/"); len(keyParts) > 1 {
		upLevelEntry := types.Entry{
			Name:    "..",
			Path:    path.Clean(keyParts[len(keyParts)-2]) + "/",
			IsDir:   true,
			ModTime: time.Now(),
		}
		entries = append(entries, upLevelEntry)
	}

	for res := range resCh {
		if res.Err != nil {
			return nil, res.Err
		}

		var entry *types.Entry
		if res.CommonPrefix != nil {
			entry = &types.Entry{
				Name:    path.Base(res.CommonPrefix.Prefix) + "/",
				Path:    res.CommonPrefix.Prefix,
				IsDir:   true,
				ModTime: time.Now(),
			}
		} else {
			entry = &types.Entry{
				Name:    path.Base(res.Object.Key),
				Path:    res.Object.Key,
				IsDir:   false,
				ModTime: time.Time(res.Object.LastModified),
				Size:    res.Object.Size,
			}
		}
		entries = append(entries, *entry)
	}

	return entries, nil
}

func (p *S3FSProvider) Upload(srcRoot string, srcs ...types.Entry) <-chan messages.UploadProgressMsg {
	progressCh := make(chan messages.UploadProgressMsg, 10)

	go func() {
		defer close(progressCh)

		var wg sync.WaitGroup

		if !p.Active {
			progressCh <- messages.UploadProgressMsg{Err: &errs.ErrProviderInactive{}}
			return
		}

		s3Client, err := p.newClient()
		if err != nil {
			progressCh <- messages.UploadProgressMsg{Err: err}
			return
		}

		var indexedPaths []IndexedPath

		fileIndex := 0
		for _, src := range srcs {
			if src.IsDir {
				err := filepath.WalkDir(src.Path, func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if d.Type().IsRegular() {
						indexedPath := IndexedPath{FileIndex: fileIndex, Path: path}
						indexedPaths = append(indexedPaths, indexedPath)
						fileIndex++
					}
					return nil
				})
				if err != nil {
					progressCh <- messages.UploadProgressMsg{Err: err}
					return
				}
			} else {
				indexedPath := IndexedPath{FileIndex: fileIndex, Path: src.Path}
				indexedPaths = append(indexedPaths, indexedPath)
				fileIndex++
			}
		}

		fileCount := len(indexedPaths)

		for _, indexedPath := range indexedPaths {
			relPath, err := filepath.Rel(srcRoot, indexedPath.Path)
			if err != nil {
				progressCh <- messages.UploadProgressMsg{
					FileCount: fileCount,
					FileIndex: indexedPath.FileIndex,
					Err:       err,
				}
				continue
			}
			key := p.Key + filepath.ToSlash(relPath)

			file, err := os.Open(indexedPath.Path)
			if err != nil {
				progressCh <- messages.UploadProgressMsg{
					FileCount: fileCount,
					FileIndex: indexedPath.FileIndex,
					Err:       err,
				}
				continue
			}

			fileProgressCh := make(chan clientTypes.Progress, 10)

			wg.Add(1)
			go func(fileIndex int, fileProgress <-chan clientTypes.Progress) {
				defer wg.Done()

				for msg := range fileProgress {
					progressCh <- messages.UploadProgressMsg{
						FileCount: fileCount,
						FileIndex: fileIndex,
						Uploaded:  int64(msg.Current),
					}
				}
			}(indexedPath.FileIndex, fileProgressCh)

			_, err = s3Client.PutObject(p.Bucket, key, file, fileProgressCh)
			close(fileProgressCh)
			helper.CloseWithErr(file, &err)

			if err != nil {
				progressCh <- messages.UploadProgressMsg{
					FileCount: fileCount,
					FileIndex: indexedPath.FileIndex,
					Err:       err,
				}
			}
		}

		wg.Wait()
	}()

	return progressCh
}

func (p *S3FSProvider) Download(dstRoot string, srcs ...types.Entry) <-chan messages.DownloadProgressMsg {
	progressCh := make(chan messages.DownloadProgressMsg, 10)

	go func() {
		defer close(progressCh)

		var wg sync.WaitGroup

		if !p.Active {
			progressCh <- messages.DownloadProgressMsg{Err: &errs.ErrProviderInactive{}}
			return
		}

		if !helper.IsDir(dstRoot) {
			progressCh <- messages.DownloadProgressMsg{Err: &errs.ErrNotDirectory{}}
			return
		}

		s3Client, err := p.newClient()
		if err != nil {
			progressCh <- messages.DownloadProgressMsg{Err: err}
			return
		}

		var indexedPaths []S3IndexedPath

		fileIndex := 0
		for _, src := range srcs {
			if src.IsDir {
				err := p.WalkDir(src, func(path string, d *types.Entry, err error) error {
					if err != nil {
						return err
					}
					if !d.IsDir {
						indexedPath := S3IndexedPath{FileIndex: fileIndex, Bucket: p.Bucket, Key: d.Path}
						indexedPaths = append(indexedPaths, indexedPath)
						fileIndex++
					}
					return nil
				})
				if err != nil {
					progressCh <- messages.DownloadProgressMsg{Err: err}
					return
				}
			} else {
				indexedPath := S3IndexedPath{FileIndex: fileIndex, Bucket: p.Bucket, Key: src.Path}
				indexedPaths = append(indexedPaths, indexedPath)
				fileIndex++
			}
		}

		fileCount := len(indexedPaths)

		for _, indexedPath := range indexedPaths {
			relKey, err := filepath.Rel(p.Key, indexedPath.Key)
			if err != nil {
				progressCh <- messages.DownloadProgressMsg{
					FileCount: fileCount,
					FileIndex: indexedPath.FileIndex,
					Err:       err,
				}
				continue
			}
			localPath := filepath.Join(dstRoot, relKey)

			dirPath := filepath.Dir(localPath)
			if err := os.MkdirAll(dirPath, 0o755); err != nil {
				progressCh <- messages.DownloadProgressMsg{
					FileCount: fileCount,
					FileIndex: indexedPath.FileIndex,
					Err:       err,
				}
				return
			}

			writer, err := os.Create(localPath)
			if err != nil {
				progressCh <- messages.DownloadProgressMsg{
					FileCount: fileCount,
					FileIndex: indexedPath.FileIndex,
					Err:       err,
				}
				continue
			}

			fileProgressCh := make(chan clientTypes.Progress, 10)

			wg.Add(1)
			go func(fileIndex int, fileProgress <-chan clientTypes.Progress) {
				defer wg.Done()

				for msg := range fileProgress {
					progressCh <- messages.DownloadProgressMsg{
						FileCount:  fileCount,
						FileIndex:  fileIndex,
						Downloaded: int64(msg.Current),
					}
				}
			}(indexedPath.FileIndex, fileProgressCh)

			_, err = s3Client.GetObject(p.Bucket, indexedPath.Key, writer, fileProgressCh)
			close(fileProgressCh)
			helper.CloseWithErr(writer, &err)

			if err != nil {
				progressCh <- messages.DownloadProgressMsg{
					FileCount: fileCount,
					FileIndex: indexedPath.FileIndex,
					Err:       err,
				}
			}
		}

		wg.Wait()
	}()

	return progressCh
}

func (p *S3FSProvider) Stat(key string) (*types.Entry, error) {
	s3Client, err := p.newClient()
	if err != nil {
		return nil, err
	}

	key = path.Clean(key)

	resp, err := s3Client.HeadObject(p.Bucket, key)
	if err != nil {
		return nil, err
	}

	entry := types.Entry{
		Name:        path.Base(key),
		Path:        key,
		IsDir:       false,
		ModTime:     resp.LastModified,
		Size:        resp.ContentLength,
		ContentType: resp.ContentType,
	}

	return &entry, nil
}

func (p *S3FSProvider) WalkDir(root types.Entry, fn WalkDirFunc) error {
	if !p.Active {
		return nil
	}

	if !root.IsDir {
		return &errs.ErrNotDirectory{}
	}

	s3Client, err := p.newClient()
	if err != nil {
		return err
	}

	resCh := s3Client.ListObjects(p.Bucket, root.Path, true)

	for res := range resCh {
		if res.Err != nil {
			if err := fn(root.Path, nil, res.Err); err != nil {
				return err
			}
			continue
		}

		if res.IsCommonPrefix() {
			dirEntry := types.Entry{
				Name:  path.Base(res.CommonPrefix.Prefix),
				Path:  res.CommonPrefix.Prefix,
				IsDir: true,
			}
			if err := fn(res.CommonPrefix.Prefix, &dirEntry, nil); err != nil {
				return err
			}
			if err := p.WalkDir(dirEntry, fn); err != nil {
				return err
			}
			continue
		}

		if res.IsObject() {
			entry := types.Entry{
				Name:    path.Base(res.Object.Key),
				Path:    res.Object.Key,
				IsDir:   false,
				ModTime: time.Time(res.Object.LastModified),
				Size:    res.Object.Size,
			}
			if err := fn(res.Object.Key, &entry, nil); err != nil {
				return err
			}
			continue
		}
	}

	return nil
}

func (p *S3FSProvider) Delete(keys ...string) (int, error) {
	if !p.Active {
		return 0, nil
	}

	s3Client, err := p.newClient()
	if err != nil {
		return 0, err
	}

	res, _, err := s3Client.DeleteObjects(p.Bucket, keys...)
	if err != nil {
		return 0, err
	}

	okCount := len(res.Deleted)

	return okCount, nil
}

func (p *S3FSProvider) newClient() (*s3.Client, error) {
	profile, err := credentials.LoadProfile(p.Profile)
	if err != nil {
		return nil, err
	}

	client := s3.NewS3Client(
		context.Background(),
		p.Host,
		p.Port,
		profile.AccessKey,
		profile.SecretKey,
		false,
		config.App.Server.S3.IO.BufferSize,
	)

	return client, nil
}

func (p *S3FSProvider) lastSelKey() string {
	return fmt.Sprintf("%s@%s/%s", p.Profile, p.Bucket, p.Key)
}
