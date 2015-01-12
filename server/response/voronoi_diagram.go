package response

import (
	"github.com/Vladimiroff/vec2d"
	"github.com/pzsz/voronoi"

	"warcluster/entities"
)

type Edge struct {
	Start *vec2d.Vector
	End   *vec2d.Vector
}

type Polygon struct {
	Edges []Edge
}

type VoronoiDiagram struct {
	baseResponse
	Polygons []Polygon
}

func NewVoronoiDiagram(position *vec2d.Vector, resolution []uint64) *VoronoiDiagram {
	topLeft, bottomRight := calculateCanvasSize(position, resolution)
	planetEntities := entities.Find("planet.*")
	planets := make([]*entities.Planet, 0, len(planetEntities))
	sites := make([]voronoi.Vertex, 0, len(planetEntities))
	x0, xn, y0, yn := 0.0, 0.0, 0.0, 0.0
	for i, planetEntity := range planetEntities {
		planets = append(planets, planetEntity.(*entities.Planet))
		if x0 > planets[i].Position.X {
			x0 = planets[i].Position.X
		}
		if xn < planets[i].Position.X {
			xn = planets[i].Position.X
		}
		if y0 > planets[i].Position.Y {
			y0 = planets[i].Position.Y
		}
		if yn < planets[i].Position.Y {
			yn = planets[i].Position.Y
		}
		sites = append(sites, voronoi.Vertex{planets[i].Position.X, planets[i].Position.Y})
	}

	bbox := voronoi.NewBBox(x0, xn, y0, yn)

	voronoiDiagram := new(VoronoiDiagram)
	voronoiDiagram.Command = "voronoi_diagram"
	diagram := voronoi.ComputeDiagram(sites, bbox, true)
	for _, cell := range diagram.Cells {
		polygon := new(Polygon)
		for _, hedge := range cell.Halfedges {
			edge := Edge{
				Start: &vec2d.Vector{hedge.Edge.Va.X, hedge.Edge.Va.Y},
				End:   &vec2d.Vector{hedge.Edge.Vb.X, hedge.Edge.Vb.Y},
			}
			if isInside(edge.Start, topLeft, bottomRight) && isInside(edge.End, topLeft, bottomRight) {
				polygon.Edges = append(polygon.Edges, edge)
			} else {
				if isInside(edge.Start, topLeft, bottomRight) {
					edge.End = crossPoint(edge, topLeft, bottomRight)
					polygon.Edges = append(polygon.Edges, edge)
				}
				if isInside(edge.End, topLeft, bottomRight) {
					edge.Start = crossPoint(edge, topLeft, bottomRight)
					polygon.Edges = append(polygon.Edges, edge)
				}
			}
		}
		voronoiDiagram.Polygons = append(voronoiDiagram.Polygons, *polygon)
	}
	return voronoiDiagram
}

func isInside(point, topLeft, bottomRight *vec2d.Vector) bool {
	return point.X > topLeft.X && point.X < bottomRight.X && point.Y < topLeft.Y && point.Y > bottomRight.Y
}

func (_ *VoronoiDiagram) Sanitize(*entities.Player) {}

func crossPoint(edge Edge, topLeft, bottomRight *vec2d.Vector) *vec2d.Vector {
	var x, y float64
	p := entities.NewCartesianEquation(edge.Start, edge.End)
	a, b := p.GetA(), p.GetB()
	p = entities.NewCartesianEquation(topLeft, &vec2d.Vector{topLeft.X, bottomRight.Y})
	c1, d1 := p.GetA(), p.GetB()
	x = (d1 - b) / (a - c1)
	y = a*x + b
	if y < topLeft.Y && y > bottomRight.Y {
		return &vec2d.Vector{x, y}
	}
	p = entities.NewCartesianEquation(topLeft, &vec2d.Vector{bottomRight.X, topLeft.Y})
	c2, d2 := p.GetA(), p.GetB()
	x = (d2 - b) / (a - c2)
	y = a*x + b
	if y < topLeft.Y && y > bottomRight.Y {
		return &vec2d.Vector{x, y}
	}
	p = entities.NewCartesianEquation(bottomRight, &vec2d.Vector{topLeft.X, bottomRight.Y})
	c3, d3 := p.GetA(), p.GetB()
	x = (d3 - b) / (a - c3)
	y = a*x + b
	if y < topLeft.Y && y > bottomRight.Y {
		return &vec2d.Vector{x, y}
	}
	p = entities.NewCartesianEquation(bottomRight, &vec2d.Vector{bottomRight.X, topLeft.Y})
	c4, d4 := p.GetA(), p.GetB()
	x = (d4 - b) / (a - c4)
	y = a*x + b
	if y < topLeft.Y && y > bottomRight.Y {
		return &vec2d.Vector{x, y}
	}
	return &vec2d.Vector{1, 1}
}
