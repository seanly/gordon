package gordon

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/v29/github"
	"github.com/kyoh86/gordon/internal/archive"
)

func Download(ctx context.Context, ev Env, client *github.Client, release *Release) error {
	path := ReleasePath(ev, release)
	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}
	unarchiver, err := openAsset(ctx, client, release)
	if err != nil {
		return err
	}
	if err := extractAsset(path, unarchiver); err != nil {
		return err
	}
	return nil
}

func openAsset(ctx context.Context, client *github.Client, release *Release) (archive.Unarchiver, error) {
	reader, _, err := client.Repositories.DownloadReleaseAsset(
		ctx,
		release.owner,
		release.name,
		release.asset.ID,
		http.DefaultClient,
	)
	if err != nil {
		return nil, err
	}
	return release.opener(reader)
}

func extractAsset(path string, unarchiver archive.Unarchiver) error {
	return unarchiver.Walk(func(info os.FileInfo, entry archive.Entry) (retErr error) {
		entryReader, err := entry()
		if err != nil {
			return err
		}
		defer func() {
			if err := entryReader.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()

		// TODO: support subdirectory
		file, err := os.OpenFile(
			filepath.Join(path, info.Name()),
			os.O_CREATE|os.O_EXCL|os.O_WRONLY,
			info.Mode(),
		)
		if err != nil {
			return err
		}
		defer func() {
			if err := file.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()

		if _, err := io.Copy(file, entryReader); err != nil {
			return err
		}
		return nil
	})
}
