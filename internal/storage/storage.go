package storage

type Storage interface {
	Save(url string, code string) error
	Get(code string) (string, error)
	FindByURL(url string) (string, error)
	Exists(code string) (bool, error)
}
