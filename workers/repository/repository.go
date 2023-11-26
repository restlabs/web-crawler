package repository

type Repository interface {
	Insert(data string) error
	Remove(data string) error
	Get(data string) (string, error)
}
