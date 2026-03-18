package repository

type ListOptions struct {
	Page   int
	Limit  int
	Search string
	Sort   string
	Order  string
}
