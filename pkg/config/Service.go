package config

import (
	"github.com/BurntSushi/toml"
	"os"
)

type Driver struct {
	fileName string
	Meta     toml.MetaData
}

func New(fName string) *Driver {
	return &Driver{fileName: fName}
}
func (d *Driver) SetFile(f string) *Driver {
	d.fileName = f
	return d
}

func (d *Driver) Load() (Config, error) {
	var err error
	var config Config
	d.Meta, err = toml.DecodeFile(d.fileName, &config)
	if err != nil {
		return Config{}, err
	}
	return config, err
}

func (d *Driver) Update(config Config) error {

	f, err := os.Create(d.fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	err = toml.NewEncoder(f).Encode(config)
	if err != nil {
		return err
	}
	return nil
}

func CreateNewInitFile(path string) error {
	err := os.WriteFile(path, []byte(InitFileServer), 0666)
	if err != nil {
		return err
	}
	return nil
}
