package activity

type Handler interface {
	OnMessage() interface{}
}
