package res

import (
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"time"
)

type Customer struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	WordpressURL   string     `json:"wordpress_url"`
	FacebookToken  *string    `json:"facebook_token"`
	StartDate      *time.Time `json:"start_date"`
	InstagramID    *string    `json:"instagram_id"`
	InstagramName  *string    `json:"instagram_name"`
	DeleteHashFlag int        `json:"delete_hash_flag"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type Customers struct {
	Customers []*Customer `json:"customers"`
	Paginate
}

func GetCustomer(c *model.Customer) *Customer {
	return &Customer{
		ID:             c.ID,
		Name:           c.Name,
		Email:          c.Email,
		WordpressURL:   c.WordpressURL,
		FacebookToken:  c.FacebookToken,
		StartDate:      c.StartDate,
		InstagramID:    c.InstagramID,
		InstagramName:  c.InstagramName,
		DeleteHashFlag: c.DeleteHashFlag,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}

func GetCustomers(customers []*model.Customer, count int) *Customers {
	resp := make([]*Customer, len(customers))
	for i, post := range customers {
		resp[i] = GetCustomer(post)
	}
	return &Customers{
		Customers: resp,
		Paginate: Paginate{
			Count: count,
			Total: len(customers),
		},
	}
}
