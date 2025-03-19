package gwh

import (
	"errors"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

type ManifestVersion int

const (
	ManifestVersion1       ManifestVersion = 1
	ManifestVersionCurrent ManifestVersion = ManifestVersion1
)

type ManifestRepositoryKind string

const (
	ManifestRepositoryKindGit ManifestRepositoryKind = "git"
)

type manifestRepository struct {
	LocalID string                 `koanf:"local_id"`
	Kind    ManifestRepositoryKind `koanf:"kind"`
	Path    string                 `koanf:"path"`
}

type manifestContent struct {
	Version      ManifestVersion      `koanf:"version"`
	Repositories []manifestRepository `koanf:"repositories"`
}

type Manifest struct {
	content manifestContent
}

type ManifestOptions struct {
	Version ManifestVersion
}

func NewManifest(options *ManifestOptions) *Manifest {
	return &Manifest{
		content: manifestContent{
			Version:      options.Version,
			Repositories: []manifestRepository{},
		},
	}
}

func (m *Manifest) AddRepository(kind ManifestRepositoryKind, path string) (string, error) {
	localID, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	localIDString := strings.ToLower(localID.String())

	repository := manifestRepository{
		LocalID: localIDString,
		Kind:    kind,
		Path:    path,
	}

	m.content.Repositories = append(m.content.Repositories, repository)

	return localIDString, nil
}

func (m *Manifest) serializeContent() ([]byte, error) {
	conf := koanf.New(".")

	if err := conf.Load(structs.Provider(m.content, "koanf"), nil); err != nil {
		return nil, err
	}

	raw, err := conf.Marshal(yaml.Parser())

	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (m *Manifest) writeContent(name string, buf []byte) (_err error) {
	f, err := os.Create(name)

	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			_err = errors.Join(_err, err)
		}
	}()

	if _, err := f.Write(buf); err != nil {
		return err
	}

	return nil
}

func ReadManifest(path string) (*Manifest, error) {
	manifest := &Manifest{}

	return manifest, nil
}

func WriteManifest(manifest *Manifest, path string) error {
	raw, err := manifest.serializeContent()

	if err != nil {
		return err
	}

	if err := manifest.writeContent(path, raw); err != nil {
		return err
	}

	return nil
}
