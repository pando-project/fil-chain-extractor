package config

type StorageRW uint8

const (
	R StorageRW = iota
	W
	RW
)

type Storage struct {
	DBType      string    `yaml:"DbType"`
	DSN         string    `yaml:"DSN"`
	ReadOrWrite StorageRW `yaml:"RW"`
}

func DefaultStorage() Storage {
	return Storage{
		DBType:      "mongo",
		DSN:         "mongodb://127.0.0.1:27017/?directConnection=true",
		ReadOrWrite: RW,
	}
}
