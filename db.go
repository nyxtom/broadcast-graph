package bgraph

import "sync"

type DB interface {
	setEdge(from string, to string, weight float64)
	incrEdge(from string, to string, weight float64)
	decrEdge(from string, to string, weight float64)
	setVertex(vertex string, weight float64)
	incrVertex(vertex string, weight float64)
	decrVertex(vertex string, weight float64)
	findEdges(vertex string) map[string]float64
	sumIntersectEdges(vertices []string) map[string]float64
}

type MemoryGraphDb struct {
	sync.Mutex

	vertices      map[string]int64          // set of vertices and their associated map values
	r_vertices    map[int64]string          // reverse lookup of the vertices index to the cooresponding name
	vertexWeights map[int64]float64         // map of vertex weights
	edges         map[int64]map[int64]int64 // map of vertex to the set of vertices edges[a_vertex][b_vertex]edgeNum
	edgeWeights   map[int64]float64         // map of edge weights
	totalVertices int64                     // atomically updated total of the number of vertices
	totalEdges    int64                     // atomically updated total of the number of edges
	allowNegative bool                      // allow for negative weights to occur
}

func NewMemoryGraphDb() (*MemoryGraphDb, error) {
	mem := new(MemoryGraphDb)
	mem.vertices = make(map[string]int64)
	mem.r_vertices = make(map[int64]string)
	mem.vertexWeights = make(map[int64]float64)
	mem.edges = make(map[int64]map[int64]int64)
	mem.edgeWeights = make(map[int64]float64)
	mem.allowNegative = true
	return mem, nil
}

// getEdgeIndex will return the edge index according to the two vertices presented
func (m *MemoryGraphDb) getEdgeIndex(from string, to string) int64 {
	// ensure that both vertices exist in the map
	f, f_ok := m.vertices[from]
	t, t_ok := m.vertices[to]
	if !f_ok {
		m.totalVertices++
		f = m.totalVertices
		m.vertices[from] = f
		m.r_vertices[f] = from
	}

	if !t_ok {
		m.totalVertices++
		t = m.totalVertices
		m.vertices[to] = t
		m.r_vertices[t] = to
	}

	// find the edge map or create it
	ef, ef_ok := m.edges[f]
	if !ef_ok {
		ef = make(map[int64]int64)
		m.edges[f] = ef
	}

	// find the edge appropriately
	ef_t, ok := ef[t]
	if !ok {
		m.totalEdges++
		ef_t = m.totalEdges
		ef[t] = ef_t
	}

	// find the edge weight or set it automatically
	if _, ok = m.edgeWeights[ef_t]; !ok {
		m.edgeWeights[ef_t] = float64(0)
	}

	return ef_t
}

func (m *MemoryGraphDb) setEdge(from string, to string, weight float64) {
	m.Lock()
	defer m.Unlock()

	// set the edge weight now that we have the proper index
	ef_t := m.getEdgeIndex(from, to)
	m.edgeWeights[ef_t] = weight
	if !m.allowNegative && m.edgeWeights[ef_t] < 0 {
		m.edgeWeights[ef_t] = 0
	}
}

func (m *MemoryGraphDb) incrEdge(from string, to string, weight float64) {
	m.Lock()
	defer m.Unlock()

	// set the edge weight now that we have the proper index
	ef_t := m.getEdgeIndex(from, to)
	m.edgeWeights[ef_t] += weight
	if !m.allowNegative && m.edgeWeights[ef_t] < 0 {
		m.edgeWeights[ef_t] = 0
	}
}

func (m *MemoryGraphDb) decrEdge(from string, to string, weight float64) {
	m.Lock()
	defer m.Unlock()

	// set the edge weight now that we have the proper index
	ef_t := m.getEdgeIndex(from, to)
	m.edgeWeights[ef_t] -= weight
	if !m.allowNegative && m.edgeWeights[ef_t] < 0 {
		m.edgeWeights[ef_t] = 0
	}
}

func (m *MemoryGraphDb) setVertex(vertex string, weight float64) {
	m.Lock()
	defer m.Unlock()

	// set the vertex weight now that we have an index
	f, f_ok := m.vertices[vertex]
	if !f_ok {
		m.totalVertices++
		f = m.totalVertices
		m.vertices[vertex] = f
		m.r_vertices[f] = vertex
	}

	m.vertexWeights[f] = weight
	if !m.allowNegative && m.vertexWeights[f] < 0 {
		m.vertexWeights[f] = 0
	}
}

func (m *MemoryGraphDb) incrVertex(vertex string, weight float64) {
	m.Lock()
	defer m.Unlock()

	// set the vertex weight now that we have an index
	f, f_ok := m.vertices[vertex]
	if !f_ok {
		m.totalVertices++
		f = m.totalVertices
		m.vertices[vertex] = f
		m.r_vertices[f] = vertex
	}

	if _, ok := m.vertexWeights[f]; !ok {
		m.vertexWeights[f] = weight
	} else {
		m.vertexWeights[f] += weight
	}

	if !m.allowNegative && m.vertexWeights[f] < 0 {
		m.vertexWeights[f] = 0
	}
}

func (m *MemoryGraphDb) decrVertex(vertex string, weight float64) {
	m.Lock()
	defer m.Unlock()

	// set the vertex weight now that we have an index
	f, f_ok := m.vertices[vertex]
	if !f_ok {
		m.totalVertices++
		f = m.totalVertices
		m.vertices[vertex] = f
		m.r_vertices[f] = vertex
	}

	if _, ok := m.vertexWeights[f]; !ok {
		m.vertexWeights[f] = -1 * weight
	} else {
		m.vertexWeights[f] -= weight
	}

	if !m.allowNegative && m.vertexWeights[f] < 0 {
		m.vertexWeights[f] = 0
	}
}

func (m *MemoryGraphDb) findEdges(vertex string) map[string]float64 {
	m.Lock()
	defer m.Unlock()

	f, f_ok := m.vertices[vertex]
	if !f_ok {
		return nil
	}

	if vertexEdges, ok := m.edges[f]; ok {
		result := make(map[string]float64)
		for vertexIndex, edgeIndex := range vertexEdges {
			to, v_ok := m.r_vertices[vertexIndex]
			weight, w_ok := m.edgeWeights[edgeIndex]
			if v_ok && w_ok {
				result[to] = weight
			}
		}
		return result
	} else {
		return nil
	}
}

func (m *MemoryGraphDb) sumIntersectEdges(vertices []string) map[string]float64 {
	m.Lock()
	defer m.Unlock()

	values := make([]map[int64]int64, len(vertices))
	minimalIndex := 0
	for i, k := range vertices {
		if index, ok := m.vertices[k]; ok {
			e, ok := m.edges[index]
			if !ok {
				return nil
			}
			values[i] = e
			if len(e) < len(values[minimalIndex]) {
				minimalIndex = i
			}
		} else {
			return nil
		}
	}

	minimalSet := values[minimalIndex]
	results := make(map[string]float64)
	for edgeVertex, edgeIndex := range minimalSet {
		value := true
		sum := float64(0)
		weight, ok := m.edgeWeights[edgeIndex]
		if ok {
			sum += weight
		}
		for i, v := range values {
			if i == minimalIndex {
				continue
			}

			e, ok := v[edgeVertex]
			if !ok {
				value = false
				break
			} else {
				sum += m.edgeWeights[e]
			}
		}

		if value {
			name := m.r_vertices[edgeVertex]
			results[name] = sum
		}
	}

	return results
}
