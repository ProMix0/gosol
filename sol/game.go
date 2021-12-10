// Package sol provides a polymorphic solitaire engine
package sol

import (
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gomps5/ui"
)

// Game represents a game state
type Game struct {
}

var (
	// DebugMode is a boolean set by command line flag -debug
	DebugMode bool = false
	// NoDrawing is set when resizing cards to stop the screen flickering
	NoDrawing bool = false
	// NoGameLoad is a boolean set by command line flag -noload
	NoGameLoad bool = false
	// NoGameSave is a boolean set by command line flag -nosave
	NoGameSave bool = false
	// NoShuffle stops the cards from being shuffled
	NoShuffle bool = false
	// NoLerp stops the cards from transitioning
	NoCardLerp = false
	// CardWidth of cards, start with a silly value to force a rescale/refan
	CardWidth int = 9
	// CardHeight of cards, start with a silly value to force a rescale/refan
	CardHeight int = 13
	// Card Corner Radius
	CardCornerRadius float64 = float64(CardWidth) / 15.0
	// PilePaddingX the gap left to the right of the pile
	PilePaddingX int = CardWidth / 10
	// PilePaddingY the gap left underneath each pile
	PilePaddingY int = CardHeight / 10
	// LeftMargin the gap between the left of the screen and the first pile
	LeftMargin int = (CardWidth / 2) + PilePaddingX
	// TopMargin the gap between top pile and top of baize
	TopMargin int = 48 + CardHeight/3
	// CardFaceImageLibrary
	// thirteen suitless cards,
	// one entry for each face card (4 suits * 13 cards),
	// suits are 1-indexed (eg club == 1) so image to be used for a card is (suit * 13) + (ord - 1).
	// can use (ord - 1) as in index to get suitless card
	TheCardFaceImageLibrary [13 * 5]*ebiten.Image
	// CardBackImage applies to all cards so is kept globally as an optimization
	CardBackImage *ebiten.Image
	// CardShadowImage applies to all cards so is kept globally as an optimization
	CardShadowImage *ebiten.Image
)

// ThePreferences holds serialized game progress data
var ThePreferences = &Preferences{Game: "Solitaire", Variant: "Klondike", BaizeColor: "BaizeGreen", PowerMoves: true, CardFaceColor: "Ivory", CardBackColor: "CornflowerBlue", FixedCards: true}

// TheStatistics holds statistics for all variants
// var TheStatistics *Statistics

// TheBaize points to the Baize, so that main can see it
var TheBaize *Baize

// The UI points to the singleton User Interface object
var TheUI *ui.UI

// TheError is a global copy of the last error reported, for optional toasting
// var TheError string

// NewGame generates a new Game object.
func NewGame() (*Game, error) {
	g := &Game{}
	TheUI = ui.New(Execute)
	// TheStatistics = NewStatistics()
	TheBaize = NewBaize()
	TheBaize.NewVariant()
	return g, nil
}

// Layout implements ebiten.Game's Layout.
func (*Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return TheBaize.Layout(outsideWidth, outsideHeight)
}

// Update updates the current game state.
func (*Game) Update() error {
	if err := TheBaize.Update(); err != nil {
		return err
	}
	return nil
}

// Draw draws the current game to the given screen.
func (*Game) Draw(screen *ebiten.Image) {
	TheBaize.Draw(screen)
}