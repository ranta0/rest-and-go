package domain

type RepositoryInterface interface {
	Create(user *ModelInterface) error

	GetByID(id string) (*ModelInterface, error)

	Update(id string, model *ModelInterface) error

	Delete(id string) error

	GetAll(limit, offset int) ([]ModelInterface, error)

	CountAll() (int, error)
}
