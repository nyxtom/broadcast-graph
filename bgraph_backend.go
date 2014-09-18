package bgraph

import (
	"errors"
	"strconv"

	"github.com/nyxtom/broadcast/server"
)

type BGraphBackend struct {
	server.Backend

	app *server.BroadcastServer
	db  DB
}

// SetDEdge will set the directed edge weight for the data passed into it
func (b *BGraphBackend) SetDEdge(data interface{}, client server.ProtocolClient) error {
	d, _ := data.([][]byte)
	if len(d) >= 3 {
		i := 0
		var weight float64
		var from string
		var to string
		for i < (len(d) - 2) {
			weight, _ = strconv.ParseFloat(string(d[i]), 64)
			from = string(d[i+1])
			to = string(d[i+2])

			b.db.setEdge(from, to, weight)
			i += 3
		}
	}

	return nil
}

func (b *BGraphBackend) IncrDEdge(data interface{}, client server.ProtocolClient) error {
	d, _ := data.([][]byte)
	if len(d) >= 3 {
		i := 0
		var weight float64
		var from string
		var to string
		for i < (len(d) - 2) {
			weight, _ = strconv.ParseFloat(string(d[i]), 64)
			from = string(d[i+1])
			to = string(d[i+2])

			b.db.incrEdge(from, to, weight)
			i += 3
		}
	}

	return nil
}

func (b *BGraphBackend) DecrDEdge(data interface{}, client server.ProtocolClient) error {
	d, _ := data.([][]byte)
	if len(d) >= 3 {
		i := 0
		var weight float64
		var from string
		var to string
		for i < (len(d) - 2) {
			weight, _ = strconv.ParseFloat(string(d[i]), 64)
			from = string(d[i+1])
			to = string(d[i+2])

			b.db.decrEdge(from, to, weight)
			i += 3
		}
	}

	return nil
}

func (b *BGraphBackend) SetEdge(data interface{}, client server.ProtocolClient) error {
	d, _ := data.([][]byte)
	if len(d) >= 3 {
		i := 0
		var weight float64
		var from string
		var to string
		for i < (len(d) - 2) {
			weight, _ = strconv.ParseFloat(string(d[i]), 64)
			from = string(d[i+1])
			to = string(d[i+2])

			b.db.setEdge(from, to, weight)
			b.db.setEdge(to, from, weight)
			i += 3
		}
	}

	return nil
}

func (b *BGraphBackend) IncrEdge(data interface{}, client server.ProtocolClient) error {
	d, _ := data.([][]byte)
	if len(d) >= 3 {
		i := 0
		var weight float64
		var from string
		var to string
		for i < (len(d) - 2) {
			weight, _ = strconv.ParseFloat(string(d[i]), 64)
			from = string(d[i+1])
			to = string(d[i+2])

			b.db.incrEdge(from, to, weight)
			b.db.incrEdge(to, from, weight)
			i += 3
		}
	}

	return nil
}

func (b *BGraphBackend) DecrEdge(data interface{}, client server.ProtocolClient) error {
	d, _ := data.([][]byte)
	if len(d) >= 3 {
		i := 0
		var weight float64
		var from string
		var to string
		for i < (len(d) - 2) {
			weight, _ = strconv.ParseFloat(string(d[i]), 64)
			from = string(d[i+1])
			to = string(d[i+2])

			b.db.decrEdge(from, to, weight)
			b.db.decrEdge(to, from, weight)
			i += 3
		}
	}

	return nil
}

func (b *BGraphBackend) SetVertex(data interface{}, client server.ProtocolClient) error {
	d, _ := data.([][]byte)
	if len(d) >= 2 {
		i := 0
		var weight float64
		var vertex string
		for i < (len(d) - 1) {
			weight, _ = strconv.ParseFloat(string(d[i]), 64)
			vertex = string(d[i+1])

			b.db.setVertex(vertex, weight)
			i += 2
		}
	}

	return nil
}

func (b *BGraphBackend) IncrVertex(data interface{}, client server.ProtocolClient) error {
	d, _ := data.([][]byte)
	if len(d) >= 2 {
		i := 0
		var weight float64
		var vertex string
		for i < (len(d) - 1) {
			weight, _ = strconv.ParseFloat(string(d[i]), 64)
			vertex = string(d[i+1])

			b.db.incrVertex(vertex, weight)
			i += 2
		}
	}

	return nil
}

func (b *BGraphBackend) DecrVertex(data interface{}, client server.ProtocolClient) error {
	d, _ := data.([][]byte)
	if len(d) >= 2 {
		i := 0
		var weight float64
		var vertex string
		for i < (len(d) - 1) {
			weight, _ = strconv.ParseFloat(string(d[i]), 64)
			vertex = string(d[i+1])

			b.db.decrVertex(vertex, weight)
			i += 2
		}
	}

	return nil
}

func (b *BGraphBackend) FindEdges(data interface{}, client server.ProtocolClient) error {
	d, _ := data.([][]byte)
	if len(d) < 1 {
		client.WriteError(errors.New("*e takes at least 1 parameter (*e vertex [vertex ...])"))
		client.Flush()
		return nil
	}

	vertexEdges := make(map[string]map[string]float64)
	for _, k := range d {
		key := string(k)
		edges := b.db.findEdges(key)
		if edges != nil {
			vertexEdges[key] = edges
		}
	}

	if len(vertexEdges) > 0 {
		client.WriteJson(vertexEdges)
		client.Flush()
	} else {
		client.WriteNull()
		client.Flush()
	}
	return nil
}

func (b *BGraphBackend) IntersectEdges(data interface{}, client server.ProtocolClient) error {
	d, _ := data.([][]byte)
	if len(d) < 2 {
		client.WriteError(errors.New("&e takes at least 2 parameters (&e vertex [vertex ...])"))
		client.Flush()
		return nil
	}

	keys := make([]string, len(d))
	for i, k := range d {
		keys[i] = string(k)
	}

	results := b.db.sumIntersectEdges(keys)
	if results != nil && len(results) > 0 {
		client.WriteJson(results)
		client.Flush()
	} else {
		client.WriteNull()
		client.Flush()
	}

	return nil
}

func RegisterBackend(app *server.BroadcastServer) (server.Backend, error) {
	backend := new(BGraphBackend)
	db, _ := NewMemoryGraphDb()
	backend.db = db

	app.RegisterCommand(server.Command{"=>", "Sets the directed edge weight", "=> weight from to [from to ...]", true}, backend.SetDEdge)
	app.RegisterCommand(server.Command{"+>", "Increments the directed edge weight", "+> weight from to [from to ...]", true}, backend.IncrDEdge)
	app.RegisterCommand(server.Command{"->", "Decrements the directed edge weight", "-> weight from to [from to ...]", true}, backend.DecrDEdge)
	app.RegisterCommand(server.Command{"<=>", "Sets the symmetric edge weight", "<=> weight from to [from to ...]", true}, backend.SetEdge)
	app.RegisterCommand(server.Command{"<+>", "Increments the symmetric edge weight", "<+> weight from to [from to ...]", true}, backend.IncrEdge)
	app.RegisterCommand(server.Command{"<->", "Decrements the symmetric edge weight", "<-> weight from to [from to ...]", true}, backend.DecrEdge)
	app.RegisterCommand(server.Command{"=", "Sets a given vertex's own weight", "= weight vertex [weight vertex ...]", true}, backend.SetVertex)
	app.RegisterCommand(server.Command{"+", "Increments a given vertex's own weight", "+ weight vertex [weight vertex ...]", true}, backend.IncrVertex)
	app.RegisterCommand(server.Command{"-", "Decrements a given vertex's own weight", "- weight vertex [weight vertex ...]", true}, backend.DecrVertex)
	app.RegisterCommand(server.Command{"*e", "Returns a list of all edges from the specified vertices", "*e vertex [vertex ...]", false}, backend.FindEdges)
	app.RegisterCommand(server.Command{"&e", "Returns the intersection of all edges between the set of vertices with the sum of the weights", "&e vertex [vertex ...]", false}, backend.IntersectEdges)
	backend.app = app

	return backend, nil
}

func (b *BGraphBackend) Load() error {
	return nil
}

func (b *BGraphBackend) Unload() error {
	return nil
}
