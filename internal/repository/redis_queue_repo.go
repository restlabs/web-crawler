package repository

type RedisRepo struct {
	Repository
}

func NewRedisRepo() *RedisRepo {
	return &RedisRepo{}
}

func (r *RedisRepo) Insert(data string) error {
	return nil
}

func (r *RedisRepo) Remove(data string) error {
	return nil
}

func (r *RedisRepo) Get(data string) (string, error) {
	return "", nil
}
