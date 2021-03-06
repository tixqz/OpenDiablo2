package tests

import (
	"testing"

	"github.com/hajimehoshi/ebiten"

	_map "github.com/OpenDiablo2/OpenDiablo2/d2render/mapengine"

	"github.com/OpenDiablo2/OpenDiablo2/d2common"
	"github.com/OpenDiablo2/OpenDiablo2/d2core"
	"github.com/OpenDiablo2/OpenDiablo2/d2data/mpq"
)

func TestMapGenerationPerformance(t *testing.T) {
	mpq.InitializeCryptoBuffer()
	d2common.ConfigBasePath = "../"

	engine := d2core.CreateEngine()
	gameState := d2core.CreateGameState()
	mapEngine := _map.CreateMapEngine(gameState, engine.SoundManager, engine)
	mapEngine.GenerateAct1Overworld()
	surface, _ := ebiten.NewImage(800, 600, ebiten.FilterNearest)
	for y := 0; y < 1000; y++ {
		mapEngine.Render(surface)
		mapEngine.OffsetY = float64(-y)
	}

}
