package loadtest

import (
	"fmt"
	"math/rand"
	"strings"
)

// Tile represents a map tile with coordinates
type Tile struct {
	Z int
	X int
	Y int
}

// TileGenerator generates tiles according to the specified pattern
type TileGenerator struct {
	config Config
}

// NewTileGenerator creates a new tile generator
func NewTileGenerator(config Config) *TileGenerator {
	return &TileGenerator{
		config: config,
	}
}

// NextTile returns the next tile to request based on the pattern
func (g *TileGenerator) NextTile() Tile {
	// For random pattern, generate a new random tile
	z := g.config.MinZoom + rand.Intn(g.config.MaxZoom-g.config.MinZoom+1)
	x := g.config.MinX + rand.Intn(g.config.MaxX-g.config.MinX+1)
	y := g.config.MinY + rand.Intn(g.config.MaxY-g.config.MinY+1)

	return Tile{Z: z, X: x, Y: y}
}

// FormatURL formats the URL template with the tile coordinates
func (g *TileGenerator) FormatURL(tile Tile) string {
	url := g.config.URLTemplate
	url = strings.Replace(url, "{z}", fmt.Sprintf("%d", tile.Z), -1)
	url = strings.Replace(url, "{x}", fmt.Sprintf("%d", tile.X), -1)
	url = strings.Replace(url, "{y}", fmt.Sprintf("%d", tile.Y), -1)
	return url
}
