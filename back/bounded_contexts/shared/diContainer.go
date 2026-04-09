package shared

var container *DIContainer

type DIContainer struct {
}

func NewDIContainer() *DIContainer {
	if container == nil {
		container = &DIContainer{}
	}
	return container
}

