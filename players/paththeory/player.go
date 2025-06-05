package main

import (
	"math"
	"sort"

	"github.com/brvoisin/planetwars"
)

type starterBot struct {
	source planetwars.Planet
	target planetwars.Planet
	path   Path
}

type Path []planetwars.PlanetID

func NewStarterBot() planetwars.Player {
	return &starterBot{target: planetwars.Planet{ID: -1}}
}

// DoTurn implements Player.
// It implements the same basic strategy as the starter bot from
// https://github.com/xtevenx/planet-wars-starterpackage/.
func (b *starterBot) DoTurn(pwMap planetwars.Map) []planetwars.Order {
	if b.target.ID == -1 || b.target.Owner == planetwars.Myself {
		b.source = b.FindSource(pwMap)
		b.target = b.FindTarget(pwMap)
		b.path = b.FindPath(pwMap, b.source.ID, b.target.ID)
	}
	return b.makeOrders(pwMap)
}

func (b *starterBot) makeOrders(pwMap planetwars.Map) []planetwars.Order {
	orders := make([]planetwars.Order, 0, len(b.path)-1)
	for i := range len(b.path) - 1 {
		planet := pwMap.PlanetByID(b.path[i])
		if planet.Owner != planetwars.Myself {
			continue
		}
		var cumulativeShips planetwars.Ships
		if len(orders) > 0 {
			previousOrder := orders[len(orders)-1]
			cumulativeShips = previousOrder.Ships
		}
		orders = append(orders, planetwars.Order{
			Source: b.path[i],
			Dest:   b.path[i+1],
			Ships: planetwars.Ships(
				math.Min(float64(planetwars.Ships(planet.Growth)-1+cumulativeShips), float64(planet.Ships-1)),
			),
		})
	}
	return orders
}

func (b *starterBot) FindSource(pwMap planetwars.Map) planetwars.Planet {
	planets := pwMap.MyPlanets()
	if len(planets) != 0 {
		return planets[0]
	}
	return planetwars.Planet{}
}

func (b *starterBot) FindTarget(pwMap planetwars.Map) planetwars.Planet {
	planets := pwMap.NotMyPlanets()
	for i := range planets {
		if planets[i].Owner == planetwars.Opponent {
			return planets[i]
		}
	}
	return planetwars.Planet{}
}

func (b *starterBot) FindPath(pwMap planetwars.Map, source, target planetwars.PlanetID) Path {
	path := Path{source}
	planets := pwMap.Planets
	sourcePlanet := pwMap.PlanetByID(source)
	targetPlanet := pwMap.PlanetByID(target)
	distance := planetwars.Distance(sourcePlanet, targetPlanet)
	for _, candidate := range planets {
		if candidate.ID == source || candidate.ID == target {
			continue
		}
		relDist := planetwars.Distance(sourcePlanet, candidate) + planetwars.Distance(candidate, targetPlanet)
		if relDist/distance < 1.25 {
			path = append(path, candidate.ID)
		}
	}
	sort.Slice(path, func(i, j int) bool {
		return planetwars.Distance(
			sourcePlanet,
			pwMap.PlanetByID(path[i]),
		) < planetwars.Distance(
			sourcePlanet,
			pwMap.PlanetByID(path[j]),
		)
	})
	path = append(path, target)
	return path
}
