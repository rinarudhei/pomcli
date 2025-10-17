package app

import (
	"context"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/button"
)

type buttonSet struct {
	startButton *button.Button
	pauseButton *button.Button
	increment   *button.Button
	decrement   *button.Button
}

func newButtonSet(ctx context.Context, w *widgets, redrawCh chan<- bool) (*buttonSet, error) {
	var err error

	bs := &buttonSet{}
	bs.startButton, err = initStartButton(ctx, redrawCh)
	if err != nil {
		return nil, err
	}

	bs.pauseButton, err = initPauseButton(ctx, redrawCh)
	if err != nil {
		return nil, err
	}

	bs.increment, err = initIncrementButton(ctx, redrawCh)
	if err != nil {
		return nil, err
	}

	bs.decrement, err = initDecrementButton(ctx, redrawCh)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func initStartButton(ctx context.Context, redrawCh chan<- bool) (*button.Button, error) {
	return button.New("[s]tart", func() error {
		redrawCh <- true
		return nil
	},
		button.Height(2),
		button.WidthFor("[p]ause"),
		button.FillColor(cell.ColorGreen),
		button.ShadowColor(cell.ColorGray),
		button.GlobalKey('s'),
	)
}

func initPauseButton(ctx context.Context, redrawCh chan<- bool) (*button.Button, error) {
	return button.New("[p]ause", func() error {
		redrawCh <- true
		return nil
	},
		button.Height(2),
		button.FillColor(cell.ColorYellow),
		button.ShadowColor(cell.ColorGray),
		button.GlobalKey('p'),
	)
}

func initIncrementButton(ctx context.Context, redrawCh chan<- bool) (*button.Button, error) {
	return button.New("[+]", func() error {
		redrawCh <- true
		return nil
	},
		button.Height(1),
		button.WidthFor("[+]"),
		button.FillColor(cell.ColorYellow),
		button.GlobalKey('+'),
	)
}

func initDecrementButton(ctx context.Context, redrawCh chan<- bool) (*button.Button, error) {
	return button.New("[-]", func() error {
		redrawCh <- true
		return nil
	},
		button.Height(1),
		button.WidthFor("[+]"),
		button.FillColor(cell.ColorYellow),
		button.GlobalKey('-'),
	)
}
