package router

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/IllumiKnowLabs/labstore/server/pkg/config"
)

type ServerDescriptor struct {
	Name    string
	Server  *http.Server
	Healthy atomic.Bool
}

func Start() {
	if err := ensureDirectories(); err != nil {
		slog.Error("could not create directory structure", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	webServerDescriptor := NewWebServerDescriptor(
		config.App.Web.Address.Host,
		config.App.Web.Address.Port,
	)

	s3ServerDescriptor := NewS3ServerDescriptor(
		config.App.Server.S3.Address.Host,
		config.App.Server.S3.Address.Port,
	)

	iamServerDescriptor := NewIAMServerDescriptor(
		config.App.Server.IAM.Address.Host,
		config.App.Server.IAM.Address.Port,
	)

	adminServerDescriptor := NewAdminServerDescriptor(
		config.App.Server.Admin.Address.Host,
		config.App.Server.Admin.Address.Port,
		[]*ServerDescriptor{s3ServerDescriptor, iamServerDescriptor},
	)

	serverDescriptors := []*ServerDescriptor{
		adminServerDescriptor,
		iamServerDescriptor,
		s3ServerDescriptor,
		webServerDescriptor,
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(serverDescriptors))

	for _, sd := range serverDescriptors {
		wg.Add(1)
		go runServer(sd, &wg, errCh)
		sd.Healthy.Store(true)
	}

	go shutdownServers(ctx, serverDescriptors)
	go waitAndClose(errCh, &wg)

	for err := range errCh {
		if err != nil {
			slog.Error("server error", "err", err)
		}
	}

	slog.Info("all servers shut down cleanly")
}

func ensureDirectories() error {
	slog.Debug("ensuring data directories")

	if err := os.MkdirAll(config.App.Server.Storage.ObjectsPath, 0o750); err != nil {
		return err
	}

	if err := os.MkdirAll(config.App.Server.Storage.MetadataPath, 0o750); err != nil {
		return err
	}

	return nil
}

func runServer(sd *ServerDescriptor, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()

	fmt.Printf("🌐 %s listening on http://%s\n", sd.Name, sd.Server.Addr)

	err := sd.Server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		sd.Healthy.Store(false)
		errCh <- err
	}
}

func shutdownServers(ctx context.Context, serverDescriptors []*ServerDescriptor) {
	<-ctx.Done()
	slog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := len(serverDescriptors) - 1; i >= 0; i-- {
		s := serverDescriptors[i]
		slog.Info("shutting down server", "addr", s.Server.Addr)
		_ = s.Server.Shutdown(shutdownCtx)
	}
}

func waitAndClose(errCh chan error, wg *sync.WaitGroup) {
	wg.Wait()
	close(errCh)
}
