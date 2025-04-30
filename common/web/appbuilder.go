package web

type AppBuilder[T Application] struct {
	app T
}

func NewBuilder[T Application](new func() T) AppBuilder[T] {
	return AppBuilder[T]{
		app: new(),
	}
}

func (builder *AppBuilder[T]) WithServer(serverAddress string, routes Routes) *AppBuilder[T] {
	builder.app.SetServer(serverAddress, routes)
	return builder
}

func (builder *AppBuilder[T]) Build() T {
	return builder.app
}
