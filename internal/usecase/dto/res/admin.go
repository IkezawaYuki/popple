package res

import (
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"time"
)

type Admin struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetAdmin(m *model.Admin) *Admin {
	return &Admin{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

type Admins struct {
	Admins []*Admin
	Paginate
}

func GetAdmins(admins []*model.Admin, count int) *Admins {
	resp := make([]*Admin, len(admins))
	for i, v := range admins {
		resp[i] = GetAdmin(v)
	}
	return &Admins{
		Admins: resp,
		Paginate: Paginate{
			Count: count,
			Total: len(admins),
		},
	}
}
