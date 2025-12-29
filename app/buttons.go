package app

import (
	"context"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/button"
	"github.com/rinarudhei/pomcli/session"
)

type buttonSet struct {
	startButton *button.Button
	pauseButton *button.Button
	increment   *button.Button
	decrement   *button.Button
}

func newButtonSet(ctx context.Context, w *widgets, redrawCh chan<- bool, sessionService *session.SessionService, errCh chan<- error) (*buttonSet, error) {
	var err error

	bs := &buttonSet{}
	bs.startButton, err = initStartButton(ctx, redrawCh, w, sessionService, errCh)
	if err != nil {
		return nil, err
	}

	bs.pauseButton, err = initPauseButton(sessionService, errCh)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func initStartButton(ctx context.Context, redrawCh chan<- bool, w *widgets, s *session.SessionService, errCh chan<- error) (*button.Button, error) {
	return button.New("[s]tart", func() error {
		go func() {
			update := func(sessionState string, timerString, history string) {
				w.update("", history, "", "", []string{timerString, sessionState}, redrawCh)
			}
			errCh <- s.Start(ctx, update)
		}()
		return nil
	},
		button.Height(2),
		button.WidthFor("[p]ause"),
		button.FillColor(cell.ColorGreen),
		button.ShadowColor(cell.ColorGray),
		button.GlobalKey('s'),
	)
}

func initPauseButton(s *session.SessionService, errCh chan<- error) (*button.Button, error) {
	return button.New("[p]ause", func() error {
		go func() {
			errCh <- s.Pause()
		}()
		return nil
	},
		button.Height(2),
		button.FillColor(cell.ColorYellow),
		button.ShadowColor(cell.ColorGray),
		button.GlobalKey('p'),
	)
}
