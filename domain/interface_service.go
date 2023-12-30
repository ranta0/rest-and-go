package domain

type ServiceInterface interface {
	GetAll(limit, offset int) ([]ModelInterface, error)

	GetByID(id string) (*ModelInterface, error)

	Create(user *ModelInterface) error

	Update(id string, user *ModelInterface) error

	Delete(id string) error

	CountAll() (int, error)
}
