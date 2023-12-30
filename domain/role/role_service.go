package role

type RoleService struct {
	roleRepo *RoleRepository
}

func NewRoleService(roleRepo *RoleRepository) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
	}
}

func (s *RoleService) GetAll(limit, offset int) ([]Role, error) {
	return s.roleRepo.GetAll(limit, offset)
}

func (s *RoleService) CountAll() (int, error) {
	return s.roleRepo.CountAll()
}

func (s *RoleService) GetByID(id string) (*Role, error) {
	return s.roleRepo.GetByID(id)
}

func (s *RoleService) Create(role *Role) error {
	return s.roleRepo.Create(role)
}

func (s *RoleService) Update(id string, role *Role) error {
	return s.roleRepo.Update(id, role)
}

func (s *RoleService) Delete(id string) error {
	return s.roleRepo.Delete(id)
}
