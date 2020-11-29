package service

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (scv *Service) GetUser(id int) (*User, error) {
	user, err := scv.dao.Get(id)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}
