package gwh

import (
	"errors"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

type manifestContent struct {
	Version string `koanf:"name"`
}

type Manifest struct {
	content manifestContent
}

type ManifestOptions struct {
	Version string
}

func NewManifest(options *ManifestOptions) *Manifest {
	return &Manifest{
		content: manifestContent{
			Version: options.Version,
		},
	}
}

func (m *Manifest) WriteFile(name string) error {
	conf := koanf.New(".")

	if err := conf.Load(structs.Provider(m.content, "koanf"), nil); err != nil {
		return err
	}

	raw, err := conf.Marshal(yaml.Parser())
	if err != nil {
		return err
	}

	if err := m.writeFile(name, raw); err != nil {
		return err
	}

	return nil
}

func (m *Manifest) writeFile(name string, buf []byte) (_err error) {
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
