package bgraph

import "github.com/nyxtom/broadcast/server"

type BGraphBackend struct {
	server.Backend

	app *server.BroadcastServer
}

func (b *BGraphBackend) SetDEdge(data interface{}, client server.ProtocolClient) error {
}

func (b *BGraphBackend) IncrDEdge(data interface{}, client server.ProtocolClient) error {
}

func (b *BGraphBackend) DecrDEdge(data interface{}, client server.ProtocolClient) error {
}

func (b *BGraphBackend) SetEdge(data interface{}, client server.ProtocolClient) error {
}

func (b *BGraphBackend) SetEdge(data interface{}, client server.ProtocolClient) error {
}

func (b *BGraphBackend) IncrEdge(data interface{}, client server.ProtocolClient) error {
}

func (b *BGraphBackend) DecrEdge(data interface{}, client server.ProtocolClient) error {
}

func (b *BGraphBackend) SetVertex(data interface{}, client server.ProtocolClient) error {
}

func (b *BGraphBackend) IncrVertex(data interface{}, client server.ProtocolClient) error {
}

func (b *BGraphBackend) DecrVertex(data interface{}, client server.ProtocolClient) error {
}

func RegisterBackend(app *server.BroadcastServer) (server.Backend, error) {
	backend := new(BGraphBackend)
	app.RegisterCommand(server.Command{"=>", "Sets the directed edge weight", "=> weight from to [to ...]", true}, backend.SetDEdge)
	app.RegisterCommand(server.Command{"+>", "Increments the directed edge weight", "+> weight from to [to ...]", true}, backend.IncrDEdge)
	app.RegisterCommand(server.Command{"->", "Decrements the directed edge weight", "-> weight from to [to ...]", true}, backend.DecrDEdge)
	app.RegisterCommand(server.Command{"<=>", "Sets the symmetric edge weight", "<=> weight vertex vertex [vertex ...]", true}, backend.SetEdge)
	app.RegisterCommand(server.Command{"<+>", "Increments the symmetric edge weight", "<+> weight vertex vertex [vertex ...]", true}, backend.IncrEdge)
	app.RegisterCommand(server.Command{"<->", "Decrements the symmetric edge weight", "<-> weight vertex vertex [vertex ...]", true}, backend.DecrEdge)
	app.RegisterCommand(server.Command{"=", "Sets a given vertex's own weight", "= weight vertex [weight vertex ...]", true}, backend.SetVertex)
	app.RegisterCommand(server.Command{"+", "Increments a given vertex's own weight", "+ weight vertex [weight vertex ...]", true}, backend.IncrVertex)
	app.RegisterCommand(server.Command{"-", "Decrements a given vertex's own weight", "- weight vertex [weight vertex ...]", true}, backend.DecrVertex)
	backend.app = app
	return backend, nil
}

func (b *BGraphBackend) Load() error {
	return nil
}

func (b *BGraphBackend) Unload() error {
	return nil
}
