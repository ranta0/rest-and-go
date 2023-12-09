package interfaces

type RepositoryInterface interface {
	Create(user *ModelInterface) error

	GetByID(id string) (*ModelInterface, error)

	Update(id string, model *ModelInterface) error

	Delete(id string) error

	GetAll() ([]ModelInterface, error)
}
