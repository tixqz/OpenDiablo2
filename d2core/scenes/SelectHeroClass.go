package scenes

import (
	"image"

	"github.com/OpenDiablo2/OpenDiablo2/d2data/datadict"

	"github.com/OpenDiablo2/OpenDiablo2/d2render"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/OpenDiablo2/d2audio"
	"github.com/OpenDiablo2/OpenDiablo2/d2common"
	dh "github.com/OpenDiablo2/OpenDiablo2/d2helper"
	"github.com/OpenDiablo2/OpenDiablo2/d2render/ui"
	"github.com/hajimehoshi/ebiten"
)

type HeroRenderInfo struct {
	Stance                   d2enum.HeroStance
	IdleSprite               *d2render.Sprite
	IdleSelectedSprite       *d2render.Sprite
	ForwardWalkSprite        *d2render.Sprite
	ForwardWalkSpriteOverlay *d2render.Sprite
	SelectedSprite           *d2render.Sprite
	SelectedSpriteOverlay    *d2render.Sprite
	BackWalkSprite           *d2render.Sprite
	BackWalkSpriteOverlay    *d2render.Sprite
	SelectionBounds          image.Rectangle
	SelectSfx                *d2audio.SoundEffect
	DeselectSfx              *d2audio.SoundEffect
}

type SelectHeroClass struct {
	uiManager      *ui.Manager
	soundManager   *d2audio.Manager
	fileProvider   d2interface.FileProvider
	sceneProvider  d2interface.SceneProvider
	bgImage        *d2render.Sprite
	campfire       *d2render.Sprite
	headingLabel   *ui.Label
	heroClassLabel *ui.Label
	heroDesc1Label *ui.Label
	heroDesc2Label *ui.Label
	heroDesc3Label *ui.Label
	heroRenderInfo map[d2enum.Hero]*HeroRenderInfo
	selectedHero   d2enum.Hero
	exitButton     *ui.Button
}

func CreateSelectHeroClass(
	fileProvider d2interface.FileProvider,
	sceneProvider d2interface.SceneProvider,
	uiManager *ui.Manager, soundManager *d2audio.Manager,
) *SelectHeroClass {
	result := &SelectHeroClass{
		uiManager:      uiManager,
		sceneProvider:  sceneProvider,
		fileProvider:   fileProvider,
		soundManager:   soundManager,
		heroRenderInfo: make(map[d2enum.Hero]*HeroRenderInfo),
		selectedHero:   d2enum.HeroNone,
	}
	return result
}

func (v *SelectHeroClass) loadSprite(path string, palette d2enum.PaletteType) *d2render.Sprite {
	return d2render.CreateSprite(v.fileProvider.LoadFile(path), datadict.Palettes[palette])
}

func (v *SelectHeroClass) Load() []func() {
	v.soundManager.PlayBGM(d2common.BGMTitle)
	return []func(){
		func() {
			v.bgImage = v.loadSprite(d2common.CharacterSelectBackground, d2enum.Fechar)
			v.bgImage.MoveTo(0, 0)
		},
		func() {
			v.headingLabel = ui.CreateLabel(v.fileProvider, d2common.Font30, d2enum.Units)
			fontWidth, _ := v.headingLabel.GetSize()
			v.headingLabel.MoveTo(400-int(fontWidth/2), 17)
			v.headingLabel.SetText("Select Hero Class")
			v.headingLabel.Alignment = ui.LabelAlignCenter
		},
		func() {
			v.heroClassLabel = ui.CreateLabel(v.fileProvider, d2common.Font30, d2enum.Units)
			v.heroClassLabel.Alignment = ui.LabelAlignCenter
			v.heroClassLabel.MoveTo(400, 65)
		},
		func() {
			v.heroDesc1Label = ui.CreateLabel(v.fileProvider, d2common.Font16, d2enum.Units)
			v.heroDesc1Label.Alignment = ui.LabelAlignCenter
			v.heroDesc1Label.MoveTo(400, 100)
		},
		func() {
			v.heroDesc2Label = ui.CreateLabel(v.fileProvider, d2common.Font16, d2enum.Units)
			v.heroDesc2Label.Alignment = ui.LabelAlignCenter
			v.heroDesc2Label.MoveTo(400, 115)
		},
		func() {
			v.heroDesc3Label = ui.CreateLabel(v.fileProvider, d2common.Font16, d2enum.Units)
			v.heroDesc3Label.Alignment = ui.LabelAlignCenter
			v.heroDesc3Label.MoveTo(400, 130)
		},
		func() {
			v.campfire = v.loadSprite(d2common.CharacterSelectCampfire, d2enum.Fechar)
			v.campfire.MoveTo(380, 335)
			v.campfire.Animate = true
			v.campfire.Blend = true
		},
		func() {
			v.exitButton = ui.CreateButton(ui.ButtonTypeMedium, v.fileProvider, d2common.TranslateString("#970"))
			v.exitButton.MoveTo(33, 537)
			v.exitButton.OnActivated(func() { v.onExitButtonClicked() })
			v.uiManager.AddWidget(v.exitButton)
		},
		func() {
			v.heroRenderInfo[d2enum.HeroBarbarian] = &HeroRenderInfo{
				d2enum.HeroStanceIdle,
				v.loadSprite(d2common.CharacterSelectBarbarianUnselected, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelectBarbarianUnselectedH, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelectBarbarianForwardWalk, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelectBarbarianForwardWalkOverlay, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelectBarbarianSelected, d2enum.Fechar),
				nil,
				v.loadSprite(d2common.CharacterSelectBarbarianBackWalk, d2enum.Fechar),
				nil,
				image.Rectangle{Min: image.Point{364, 201}, Max: image.Point{90, 170}},
				v.soundManager.LoadSoundEffect(d2common.SFXBarbarianSelect),
				v.soundManager.LoadSoundEffect(d2common.SFXBarbarianDeselect),
			}
			v.heroRenderInfo[d2enum.HeroBarbarian].IdleSprite.MoveTo(400, 330)
			v.heroRenderInfo[d2enum.HeroBarbarian].IdleSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroBarbarian].IdleSelectedSprite.MoveTo(400, 330)
			v.heroRenderInfo[d2enum.HeroBarbarian].IdleSelectedSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroBarbarian].ForwardWalkSprite.MoveTo(400, 330)
			v.heroRenderInfo[d2enum.HeroBarbarian].ForwardWalkSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroBarbarian].ForwardWalkSprite.SpecialFrameTime = 2500
			v.heroRenderInfo[d2enum.HeroBarbarian].ForwardWalkSprite.StopOnLastFrame = true
			v.heroRenderInfo[d2enum.HeroBarbarian].ForwardWalkSpriteOverlay.MoveTo(400, 330)
			v.heroRenderInfo[d2enum.HeroBarbarian].ForwardWalkSpriteOverlay.Animate = true
			v.heroRenderInfo[d2enum.HeroBarbarian].ForwardWalkSpriteOverlay.SpecialFrameTime = 2500
			v.heroRenderInfo[d2enum.HeroBarbarian].ForwardWalkSpriteOverlay.StopOnLastFrame = true
			v.heroRenderInfo[d2enum.HeroBarbarian].SelectedSprite.MoveTo(400, 330)
			v.heroRenderInfo[d2enum.HeroBarbarian].SelectedSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroBarbarian].BackWalkSprite.MoveTo(400, 330)
			v.heroRenderInfo[d2enum.HeroBarbarian].BackWalkSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroBarbarian].BackWalkSprite.SpecialFrameTime = 1000
			v.heroRenderInfo[d2enum.HeroBarbarian].BackWalkSprite.StopOnLastFrame = true
		},
		func() {
			v.heroRenderInfo[d2enum.HeroSorceress] = &HeroRenderInfo{
				d2enum.HeroStanceIdle,
				v.loadSprite(d2common.CharacterSelecSorceressUnselected, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecSorceressUnselectedH, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecSorceressForwardWalk, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecSorceressForwardWalkOverlay, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecSorceressSelected, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecSorceressSelectedOverlay, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecSorceressBackWalk, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecSorceressBackWalkOverlay, d2enum.Fechar),
				image.Rectangle{Min: image.Point{580, 240}, Max: image.Point{65, 160}},
				v.soundManager.LoadSoundEffect(d2common.SFXSorceressSelect),
				v.soundManager.LoadSoundEffect(d2common.SFXSorceressDeselect),
			}
			v.heroRenderInfo[d2enum.HeroSorceress].IdleSprite.MoveTo(626, 352)
			v.heroRenderInfo[d2enum.HeroSorceress].IdleSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroSorceress].IdleSelectedSprite.MoveTo(626, 352)
			v.heroRenderInfo[d2enum.HeroSorceress].IdleSelectedSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroSorceress].ForwardWalkSprite.MoveTo(626, 352)
			v.heroRenderInfo[d2enum.HeroSorceress].ForwardWalkSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroSorceress].ForwardWalkSprite.SpecialFrameTime = 2300
			v.heroRenderInfo[d2enum.HeroSorceress].ForwardWalkSprite.StopOnLastFrame = true
			v.heroRenderInfo[d2enum.HeroSorceress].ForwardWalkSpriteOverlay.Blend = true
			v.heroRenderInfo[d2enum.HeroSorceress].ForwardWalkSpriteOverlay.MoveTo(626, 352)
			v.heroRenderInfo[d2enum.HeroSorceress].ForwardWalkSpriteOverlay.Animate = true
			v.heroRenderInfo[d2enum.HeroSorceress].ForwardWalkSpriteOverlay.SpecialFrameTime = 2300
			v.heroRenderInfo[d2enum.HeroSorceress].ForwardWalkSpriteOverlay.StopOnLastFrame = true
			v.heroRenderInfo[d2enum.HeroSorceress].SelectedSprite.MoveTo(626, 352)
			v.heroRenderInfo[d2enum.HeroSorceress].SelectedSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroSorceress].SelectedSpriteOverlay.Blend = true
			v.heroRenderInfo[d2enum.HeroSorceress].SelectedSpriteOverlay.MoveTo(626, 352)
			v.heroRenderInfo[d2enum.HeroSorceress].SelectedSpriteOverlay.Animate = true
			v.heroRenderInfo[d2enum.HeroSorceress].BackWalkSprite.MoveTo(626, 352)
			v.heroRenderInfo[d2enum.HeroSorceress].BackWalkSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroSorceress].BackWalkSprite.SpecialFrameTime = 1200
			v.heroRenderInfo[d2enum.HeroSorceress].BackWalkSprite.StopOnLastFrame = true
			v.heroRenderInfo[d2enum.HeroSorceress].BackWalkSpriteOverlay.Blend = true
			v.heroRenderInfo[d2enum.HeroSorceress].BackWalkSpriteOverlay.MoveTo(626, 352)
			v.heroRenderInfo[d2enum.HeroSorceress].BackWalkSpriteOverlay.Animate = true
			v.heroRenderInfo[d2enum.HeroSorceress].BackWalkSpriteOverlay.SpecialFrameTime = 1200
			v.heroRenderInfo[d2enum.HeroSorceress].BackWalkSpriteOverlay.StopOnLastFrame = true
		},
		func() {
			v.heroRenderInfo[d2enum.HeroNecromancer] = &HeroRenderInfo{
				d2enum.HeroStanceIdle,
				v.loadSprite(d2common.CharacterSelectNecromancerUnselected, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelectNecromancerUnselectedH, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecNecromancerForwardWalk, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecNecromancerForwardWalkOverlay, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecNecromancerSelected, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecNecromancerSelectedOverlay, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecNecromancerBackWalk, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecNecromancerBackWalkOverlay, d2enum.Fechar),
				image.Rectangle{Min: image.Point{265, 220}, Max: image.Point{55, 175}},
				v.soundManager.LoadSoundEffect(d2common.SFXNecromancerSelect),
				v.soundManager.LoadSoundEffect(d2common.SFXNecromancerDeselect),
			}
			v.heroRenderInfo[d2enum.HeroNecromancer].IdleSprite.MoveTo(300, 335)
			v.heroRenderInfo[d2enum.HeroNecromancer].IdleSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroNecromancer].IdleSelectedSprite.MoveTo(300, 335)
			v.heroRenderInfo[d2enum.HeroNecromancer].IdleSelectedSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroNecromancer].ForwardWalkSprite.MoveTo(300, 335)
			v.heroRenderInfo[d2enum.HeroNecromancer].ForwardWalkSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroNecromancer].ForwardWalkSprite.SpecialFrameTime = 2000
			v.heroRenderInfo[d2enum.HeroNecromancer].ForwardWalkSprite.StopOnLastFrame = true
			v.heroRenderInfo[d2enum.HeroNecromancer].ForwardWalkSpriteOverlay.Blend = true
			v.heroRenderInfo[d2enum.HeroNecromancer].ForwardWalkSpriteOverlay.MoveTo(300, 335)
			v.heroRenderInfo[d2enum.HeroNecromancer].ForwardWalkSpriteOverlay.Animate = true
			v.heroRenderInfo[d2enum.HeroNecromancer].ForwardWalkSpriteOverlay.SpecialFrameTime = 2000
			v.heroRenderInfo[d2enum.HeroNecromancer].ForwardWalkSpriteOverlay.StopOnLastFrame = true
			v.heroRenderInfo[d2enum.HeroNecromancer].SelectedSprite.MoveTo(300, 335)
			v.heroRenderInfo[d2enum.HeroNecromancer].SelectedSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroNecromancer].SelectedSpriteOverlay.Blend = true
			v.heroRenderInfo[d2enum.HeroNecromancer].SelectedSpriteOverlay.MoveTo(300, 335)
			v.heroRenderInfo[d2enum.HeroNecromancer].SelectedSpriteOverlay.Animate = true
			v.heroRenderInfo[d2enum.HeroNecromancer].BackWalkSprite.MoveTo(300, 335)
			v.heroRenderInfo[d2enum.HeroNecromancer].BackWalkSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroNecromancer].BackWalkSprite.SpecialFrameTime = 1500
			v.heroRenderInfo[d2enum.HeroNecromancer].BackWalkSprite.StopOnLastFrame = true
			v.heroRenderInfo[d2enum.HeroNecromancer].BackWalkSpriteOverlay.Blend = true
			v.heroRenderInfo[d2enum.HeroNecromancer].BackWalkSpriteOverlay.MoveTo(300, 335)
			v.heroRenderInfo[d2enum.HeroNecromancer].BackWalkSpriteOverlay.Animate = true
			v.heroRenderInfo[d2enum.HeroNecromancer].BackWalkSpriteOverlay.SpecialFrameTime = 1500
			v.heroRenderInfo[d2enum.HeroNecromancer].BackWalkSpriteOverlay.StopOnLastFrame = true
		},
		func() {
			v.heroRenderInfo[d2enum.HeroPaladin] = &HeroRenderInfo{
				d2enum.HeroStanceIdle,
				v.loadSprite(d2common.CharacterSelectPaladinUnselected, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelectPaladinUnselectedH, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecPaladinForwardWalk, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecPaladinForwardWalkOverlay, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecPaladinSelected, d2enum.Fechar),
				nil,
				v.loadSprite(d2common.CharacterSelecPaladinBackWalk, d2enum.Fechar),
				nil,
				image.Rectangle{Min: image.Point{490, 210}, Max: image.Point{65, 180}},
				v.soundManager.LoadSoundEffect(d2common.SFXPaladinSelect),
				v.soundManager.LoadSoundEffect(d2common.SFXPaladinDeselect),
			}
			v.heroRenderInfo[d2enum.HeroPaladin].IdleSprite.MoveTo(521, 338)
			v.heroRenderInfo[d2enum.HeroPaladin].IdleSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroPaladin].IdleSelectedSprite.MoveTo(521, 338)
			v.heroRenderInfo[d2enum.HeroPaladin].IdleSelectedSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroPaladin].ForwardWalkSprite.MoveTo(521, 338)
			v.heroRenderInfo[d2enum.HeroPaladin].ForwardWalkSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroPaladin].ForwardWalkSprite.SpecialFrameTime = 3400
			v.heroRenderInfo[d2enum.HeroPaladin].ForwardWalkSprite.StopOnLastFrame = true
			v.heroRenderInfo[d2enum.HeroPaladin].ForwardWalkSpriteOverlay.MoveTo(521, 338)
			v.heroRenderInfo[d2enum.HeroPaladin].ForwardWalkSpriteOverlay.Animate = true
			v.heroRenderInfo[d2enum.HeroPaladin].ForwardWalkSpriteOverlay.SpecialFrameTime = 3400
			v.heroRenderInfo[d2enum.HeroPaladin].ForwardWalkSpriteOverlay.StopOnLastFrame = true
			v.heroRenderInfo[d2enum.HeroPaladin].SelectedSprite.MoveTo(521, 338)
			v.heroRenderInfo[d2enum.HeroPaladin].SelectedSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroPaladin].BackWalkSprite.MoveTo(521, 338)
			v.heroRenderInfo[d2enum.HeroPaladin].BackWalkSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroPaladin].BackWalkSprite.SpecialFrameTime = 1300
			v.heroRenderInfo[d2enum.HeroPaladin].BackWalkSprite.StopOnLastFrame = true
		},
		func() {
			v.heroRenderInfo[d2enum.HeroAmazon] = &HeroRenderInfo{
				d2enum.HeroStanceIdle,
				v.loadSprite(d2common.CharacterSelectAmazonUnselected, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelectAmazonUnselectedH, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelecAmazonForwardWalk, d2enum.Fechar),
				nil,
				v.loadSprite(d2common.CharacterSelecAmazonSelected, d2enum.Fechar),
				nil,
				v.loadSprite(d2common.CharacterSelecAmazonBackWalk, d2enum.Fechar),
				nil,
				image.Rectangle{Min: image.Point{70, 220}, Max: image.Point{55, 200}},
				v.soundManager.LoadSoundEffect(d2common.SFXAmazonSelect),
				v.soundManager.LoadSoundEffect(d2common.SFXAmazonDeselect),
			}
			v.heroRenderInfo[d2enum.HeroAmazon].IdleSprite.MoveTo(100, 339)
			v.heroRenderInfo[d2enum.HeroAmazon].IdleSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroAmazon].IdleSelectedSprite.MoveTo(100, 339)
			v.heroRenderInfo[d2enum.HeroAmazon].IdleSelectedSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroAmazon].ForwardWalkSprite.MoveTo(100, 339)
			v.heroRenderInfo[d2enum.HeroAmazon].ForwardWalkSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroAmazon].ForwardWalkSprite.SpecialFrameTime = 2200
			v.heroRenderInfo[d2enum.HeroAmazon].ForwardWalkSprite.StopOnLastFrame = true
			v.heroRenderInfo[d2enum.HeroAmazon].SelectedSprite.MoveTo(100, 339)
			v.heroRenderInfo[d2enum.HeroAmazon].SelectedSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroAmazon].BackWalkSprite.MoveTo(100, 339)
			v.heroRenderInfo[d2enum.HeroAmazon].BackWalkSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroAmazon].BackWalkSprite.SpecialFrameTime = 1500
			v.heroRenderInfo[d2enum.HeroAmazon].BackWalkSprite.StopOnLastFrame = true
		},
		func() {
			v.heroRenderInfo[d2enum.HeroAssassin] = &HeroRenderInfo{
				d2enum.HeroStanceIdle,
				v.loadSprite(d2common.CharacterSelectAssassinUnselected, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelectAssassinUnselectedH, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelectAssassinForwardWalk, d2enum.Fechar),
				nil,
				v.loadSprite(d2common.CharacterSelectAssassinSelected, d2enum.Fechar),
				nil,
				v.loadSprite(d2common.CharacterSelectAssassinBackWalk, d2enum.Fechar),
				nil,
				image.Rectangle{Min: image.Point{175, 235}, Max: image.Point{50, 180}},
				v.soundManager.LoadSoundEffect(d2common.SFXAssassinSelect),
				v.soundManager.LoadSoundEffect(d2common.SFXAssassinDeselect),
			}
			v.heroRenderInfo[d2enum.HeroAssassin].IdleSprite.MoveTo(231, 365)
			v.heroRenderInfo[d2enum.HeroAssassin].IdleSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroAssassin].IdleSelectedSprite.MoveTo(231, 365)
			v.heroRenderInfo[d2enum.HeroAssassin].IdleSelectedSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroAssassin].ForwardWalkSprite.MoveTo(231, 365)
			v.heroRenderInfo[d2enum.HeroAssassin].ForwardWalkSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroAssassin].ForwardWalkSprite.SpecialFrameTime = 3800
			v.heroRenderInfo[d2enum.HeroAssassin].ForwardWalkSprite.StopOnLastFrame = true
			v.heroRenderInfo[d2enum.HeroAssassin].SelectedSprite.MoveTo(231, 365)
			v.heroRenderInfo[d2enum.HeroAssassin].SelectedSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroAssassin].BackWalkSprite.MoveTo(231, 365)
			v.heroRenderInfo[d2enum.HeroAssassin].BackWalkSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroAssassin].BackWalkSprite.SpecialFrameTime = 1500
			v.heroRenderInfo[d2enum.HeroAssassin].BackWalkSprite.StopOnLastFrame = true
		},
		func() {
			v.heroRenderInfo[d2enum.HeroDruid] = &HeroRenderInfo{
				d2enum.HeroStanceIdle,
				v.loadSprite(d2common.CharacterSelectDruidUnselected, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelectDruidUnselectedH, d2enum.Fechar),
				v.loadSprite(d2common.CharacterSelectDruidForwardWalk, d2enum.Fechar),
				nil,
				v.loadSprite(d2common.CharacterSelectDruidSelected, d2enum.Fechar),
				nil,
				v.loadSprite(d2common.CharacterSelectDruidBackWalk, d2enum.Fechar),
				nil,
				image.Rectangle{Min: image.Point{680, 220}, Max: image.Point{70, 195}},
				v.soundManager.LoadSoundEffect(d2common.SFXDruidSelect),
				v.soundManager.LoadSoundEffect(d2common.SFXDruidDeselect),
			}
			v.heroRenderInfo[d2enum.HeroDruid].IdleSprite.MoveTo(720, 370)
			v.heroRenderInfo[d2enum.HeroDruid].IdleSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroDruid].IdleSelectedSprite.MoveTo(720, 370)
			v.heroRenderInfo[d2enum.HeroDruid].IdleSelectedSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroDruid].ForwardWalkSprite.MoveTo(720, 370)
			v.heroRenderInfo[d2enum.HeroDruid].ForwardWalkSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroDruid].ForwardWalkSprite.SpecialFrameTime = 4800
			v.heroRenderInfo[d2enum.HeroDruid].ForwardWalkSprite.StopOnLastFrame = true
			v.heroRenderInfo[d2enum.HeroDruid].SelectedSprite.MoveTo(720, 370)
			v.heroRenderInfo[d2enum.HeroDruid].SelectedSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroDruid].BackWalkSprite.MoveTo(720, 370)
			v.heroRenderInfo[d2enum.HeroDruid].BackWalkSprite.Animate = true
			v.heroRenderInfo[d2enum.HeroDruid].BackWalkSprite.SpecialFrameTime = 1500
			v.heroRenderInfo[d2enum.HeroDruid].BackWalkSprite.StopOnLastFrame = true
		},
	}
}

func (v *SelectHeroClass) Unload() {
	v.heroRenderInfo = nil
}

func (v *SelectHeroClass) onExitButtonClicked() {
	v.sceneProvider.SetNextScene(CreateCharacterSelect(v.fileProvider, v.sceneProvider, v.uiManager, v.soundManager))
}

func (v *SelectHeroClass) Render(screen *ebiten.Image) {
	v.bgImage.DrawSegments(screen, 4, 3, 0)
	v.headingLabel.Draw(screen)
	if v.selectedHero != d2enum.HeroNone {
		v.heroClassLabel.Draw(screen)
		v.heroDesc1Label.Draw(screen)
		v.heroDesc2Label.Draw(screen)
		v.heroDesc3Label.Draw(screen)
	}
	for heroClass, heroInfo := range v.heroRenderInfo {
		if heroInfo.Stance == d2enum.HeroStanceIdle || heroInfo.Stance == d2enum.HeroStanceIdleSelected {
			v.renderHero(screen, heroClass)
		}
	}
	for heroClass, heroInfo := range v.heroRenderInfo {
		if heroInfo.Stance != d2enum.HeroStanceIdle && heroInfo.Stance != d2enum.HeroStanceIdleSelected {
			v.renderHero(screen, heroClass)
		}
	}
	v.campfire.Draw(screen)
}

func (v *SelectHeroClass) Update(tickTime float64) {
	canSelect := true
	for _, info := range v.heroRenderInfo {
		if info.Stance != d2enum.HeroStanceIdle && info.Stance != d2enum.HeroStanceIdleSelected && info.Stance != d2enum.HeroStanceSelected {
			canSelect = false
			break
		}
	}
	allIdle := true
	for heroType, data := range v.heroRenderInfo {
		if allIdle && data.Stance != d2enum.HeroStanceIdle {
			allIdle = false
		}
		v.updateHeroSelectionHover(heroType, canSelect)
	}
	if v.selectedHero != d2enum.HeroNone && allIdle {
		v.selectedHero = d2enum.HeroNone
	}
}

func (v *SelectHeroClass) updateHeroSelectionHover(hero d2enum.Hero, canSelect bool) {
	renderInfo := v.heroRenderInfo[hero]
	switch renderInfo.Stance {
	case d2enum.HeroStanceApproaching:
		if renderInfo.ForwardWalkSprite.OnLastFrame() {
			renderInfo.Stance = d2enum.HeroStanceSelected
			renderInfo.SelectedSprite.ResetAnimation()
			if renderInfo.SelectedSpriteOverlay != nil {
				renderInfo.SelectedSpriteOverlay.ResetAnimation()
			}
		}
		return
	case d2enum.HeroStanceRetreating:
		if renderInfo.BackWalkSprite.OnLastFrame() {
			renderInfo.Stance = d2enum.HeroStanceIdle
			renderInfo.IdleSprite.ResetAnimation()
		}
		return
	}
	if !canSelect {
		return
	}
	if renderInfo.Stance == d2enum.HeroStanceSelected {
		return
	}
	mouseX := v.uiManager.CursorX
	mouseY := v.uiManager.CursorY
	b := renderInfo.SelectionBounds
	mouseHover := (mouseX >= b.Min.X) && (mouseX <= b.Min.X+b.Max.X) && (mouseY >= b.Min.Y) && (mouseY <= b.Min.Y+b.Max.Y)
	if mouseHover && v.uiManager.CursorButtonPressed(ui.CursorButtonLeft) {
		// showEntryUi = true;
		renderInfo.Stance = d2enum.HeroStanceApproaching
		renderInfo.ForwardWalkSprite.ResetAnimation()
		if renderInfo.ForwardWalkSpriteOverlay != nil {
			renderInfo.ForwardWalkSpriteOverlay.ResetAnimation()
		}
		for _, heroInfo := range v.heroRenderInfo {
			if heroInfo.Stance != d2enum.HeroStanceSelected {
				continue
			}
			heroInfo.SelectSfx.Stop()
			heroInfo.DeselectSfx.Play()
			heroInfo.Stance = d2enum.HeroStanceRetreating
			heroInfo.BackWalkSprite.ResetAnimation()
			if heroInfo.BackWalkSpriteOverlay != nil {
				heroInfo.BackWalkSpriteOverlay.ResetAnimation()
			}
		}
		v.selectedHero = hero
		v.updateHeroText()
		renderInfo.SelectSfx.Play()

		return
	}

	if mouseHover {
		renderInfo.Stance = d2enum.HeroStanceIdleSelected
	} else {
		renderInfo.Stance = d2enum.HeroStanceIdle
	}

	if v.selectedHero == d2enum.HeroNone && mouseHover {
		v.selectedHero = hero
		v.updateHeroText()
	}

}

func (v *SelectHeroClass) renderHero(screen *ebiten.Image, hero d2enum.Hero) {
	renderInfo := v.heroRenderInfo[hero]
	switch renderInfo.Stance {
	case d2enum.HeroStanceIdle:
		renderInfo.IdleSprite.Draw(screen)
	case d2enum.HeroStanceIdleSelected:
		renderInfo.IdleSelectedSprite.Draw(screen)
	case d2enum.HeroStanceApproaching:
		renderInfo.ForwardWalkSprite.Draw(screen)
		if renderInfo.ForwardWalkSpriteOverlay != nil {
			renderInfo.ForwardWalkSpriteOverlay.Draw(screen)
		}
	case d2enum.HeroStanceSelected:
		renderInfo.SelectedSprite.Draw(screen)
		if renderInfo.SelectedSpriteOverlay != nil {
			renderInfo.SelectedSpriteOverlay.Draw(screen)
		}
	case d2enum.HeroStanceRetreating:
		renderInfo.BackWalkSprite.Draw(screen)
		if renderInfo.BackWalkSpriteOverlay != nil {
			renderInfo.BackWalkSpriteOverlay.Draw(screen)
		}
	}
}

func (v *SelectHeroClass) updateHeroText() {
	switch v.selectedHero {
	case d2enum.HeroNone:
		return
	case d2enum.HeroBarbarian:
		v.heroClassLabel.SetText(d2common.TranslateString("partycharbar"))
		v.setDescLabels("#1709")
	case d2enum.HeroNecromancer:
		v.heroClassLabel.SetText(d2common.TranslateString("partycharnec"))
		v.setDescLabels("#1704")
	case d2enum.HeroPaladin:
		v.heroClassLabel.SetText(d2common.TranslateString("partycharpal"))
		v.setDescLabels("#1711")
	case d2enum.HeroAssassin:
		v.heroClassLabel.SetText(d2common.TranslateString("partycharass"))
		v.setDescLabels("#305")
	case d2enum.HeroSorceress:
		v.heroClassLabel.SetText(d2common.TranslateString("partycharsor"))
		v.setDescLabels("#1710")
	case d2enum.HeroAmazon:
		v.heroClassLabel.SetText(d2common.TranslateString("partycharama"))
		v.setDescLabels("#1698")
	case d2enum.HeroDruid:
		v.heroClassLabel.SetText(d2common.TranslateString("partychardru"))
		v.setDescLabels("#304")
	}
	/*
	   if (selectedHero == null)
	                   return;

	               switch (selectedHero.Value)
	               {

	               }

	               heroClassLabel.Location = new Point(400 - (heroClassLabel.TextArea.Width / 2), 65);
	               heroDesc1Label.Location = new Point(400 - (heroDesc1Label.TextArea.Width / 2), 100);
	               heroDesc2Label.Location = new Point(400 - (heroDesc2Label.TextArea.Width / 2), 115);
	               heroDesc3Label.Location = new Point(400 - (heroDesc3Label.TextArea.Width / 2), 130);
	*/
}

func (v *SelectHeroClass) setDescLabels(descKey string) {
	heroDesc := d2common.TranslateString(descKey)
	parts := dh.SplitIntoLinesWithMaxWidth(heroDesc, 37)
	if len(parts) > 1 {
		v.heroDesc1Label.SetText(parts[0])
	} else {
		v.heroDesc1Label.SetText("")
	}
	if len(parts) > 1 {
		v.heroDesc2Label.SetText(parts[1])
	} else {
		v.heroDesc2Label.SetText("")
	}
	if len(parts) > 2 {
		v.heroDesc3Label.SetText(parts[2])
	} else {
		v.heroDesc3Label.SetText("")
	}
}
