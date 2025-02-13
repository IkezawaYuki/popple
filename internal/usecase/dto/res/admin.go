package res

type Admin struct{}

type Admins struct {
	Admins []*Admin
	Paginate
}
