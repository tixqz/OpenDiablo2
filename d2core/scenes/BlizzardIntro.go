package scenes

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2data/video"
	"github.com/hajimehoshi/ebiten"
)

type BlizzardIntro struct {
	fileProvider  d2interface.FileProvider
	sceneProvider d2interface.SceneProvider
	videoDecoder  *video.BinkDecoder
}

func CreateBlizzardIntro(fileProvider d2interface.FileProvider, sceneProvider d2interface.SceneProvider) *BlizzardIntro {
	result := &BlizzardIntro{
		fileProvider:  fileProvider,
		sceneProvider: sceneProvider,
	}

	return result
}

func (v *BlizzardIntro) Load() []func() {
	return []func(){
		func() {
			videoBytes := v.fileProvider.LoadFile("/data/local/video/BlizNorth640x480.bik")
			v.videoDecoder = video.CreateBinkDecoder(videoBytes)
		},
	}
}

func (v *BlizzardIntro) Unload() {

}

func (v *BlizzardIntro) Render(screen *ebiten.Image) {

}

func (v *BlizzardIntro) Update(tickTime float64) {

}
