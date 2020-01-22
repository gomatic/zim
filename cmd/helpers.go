package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/LuminalHQ/zim/git"
	"github.com/LuminalHQ/zim/project"
	"github.com/LuminalHQ/zim/store"
	"github.com/spf13/viper"
)

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func getRepository(dir string) (string, error) {
	if dir == "" {
		dir = "."
	}
	repo, err := git.RepoRoot(dir)
	return repo, err
}

func gitRoot(dir string) (string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	return getRepository(absDir)
}

func getProject(dir string) (*project.Project, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	if repo, err := getRepository(absDir); err == nil {
		absDir = repo
	}
	return project.New(absDir)
}

func downloadSource(ctx context.Context, objStore store.Store, key string) (string, error) {
	if objStore == nil {
		return "", errors.New("S3 bucket is not configured")
	}
	workspace, err := ioutil.TempDir("", "zim")
	if err != nil {
		return "", fmt.Errorf("Failed to create temp dir: %s", err)
	}
	err = git.DownloadExtractArchive(ctx, objStore, workspace, key)
	if err != nil {
		return "", fmt.Errorf("Failed to fetch source from %s: %s", key, err)
	}
	return workspace, nil
}

type zimOptions struct {
	Directory  string
	Bucket     string
	URL        string
	Region     string
	Cache      string
	UseDocker  bool
	Kinds      []string
	Components []string
	Rules      []string
	Debug      bool
	Jobs       int
	CacheMode  string
}

func getZimOptions() zimOptions {
	opts := zimOptions{
		Directory:  viper.GetString("dir"),
		Bucket:     viper.GetString("bucket"),
		URL:        viper.GetString("url"),
		Region:     viper.GetString("region"),
		Cache:      viper.GetString("cache"),
		Kinds:      viper.GetStringSlice("kinds"),
		Components: viper.GetStringSlice("components"),
		Rules:      viper.GetStringSlice("rules"),
		UseDocker:  viper.GetBool("docker"),
		Debug:      viper.GetBool("debug"),
		Jobs:       viper.GetInt("jobs"),
		CacheMode:  viper.GetString("cache-mode"),
	}
	if opts.Cache == "" {
		opts.Cache = XDGCache()
	}
	// Strip paths to components if provided, e.g. src/foo -> foo
	for i, c := range opts.Components {
		opts.Components[i] = filepath.Base(c)
	}
	return opts
}

// XDGCache returns the local cache directory
func XDGCache() string {
	value := os.Getenv("XDG_CACHE_HOME")
	if value != "" {
		return value
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return path.Join(home, ".cache")
}

func fileExists(p string) bool {
	if _, err := os.Stat(p); err == nil {
		return true
	}
	return false
}
