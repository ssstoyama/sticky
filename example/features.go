package example

type Repository interface {
	Find(key string) string
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		values: map[string]string{
			// Sample
			"001": "Giorno Giovanna",
			"002": "Bruno Bucciarati",
		},
	}
}

type InMemoryRepository struct {
	values map[string]string
}

func (r *InMemoryRepository) Find(key string) string {
	return r.values[key]
}

func NewAPIRepository() *APIRepository {
	return &APIRepository{
		Endpoint: "localhost",
		values: map[string]string{
			// Sample
			"001": "Diavolo",
			"002": "Doppio",
		},
	}
}

type APIRepository struct {
	Endpoint string
	values   map[string]string
}

func (r *APIRepository) Find(key string) string {
	return r.values[key]
}

func NewUserService(r Repository) *UserService {
	return &UserService{r}
}

type UserService struct {
	r Repository
}

func (s *UserService) FindName(id string) string {
	return s.r.Find(id)
}
