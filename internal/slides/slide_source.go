package slides

type SlideSource interface {
	Load() (*Collection, error)
}
