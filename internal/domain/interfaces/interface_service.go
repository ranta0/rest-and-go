package interfaces

type ServiceInterface interface {
	GetAll() ([]ModelInterface, error)

	GetByID(id string) (*ModelInterface, error)

	Create(user *ModelInterface) error

	Update(id string, user *ModelInterface) error

	Delete(id string) error
}
