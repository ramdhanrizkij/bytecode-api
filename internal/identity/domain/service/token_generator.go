package service

type TokenGenerator interface {
	Generate() (string, error)
}
