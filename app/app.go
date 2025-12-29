// Package app initiate terminal widgets, buttons, adjust terminal view size, wire termdash UI with session logic
package app

import (
	"context"
	"image"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/rinarudhei/pomcli/session"
)

type App struct {
	ctx            context.Context
	size           image.Point
	controller     *termdash.Controller
	redrawCh       chan bool
	errCh          chan error
	term           terminalapi.Terminal
	sessionService *session.SessionService
}

func NewApp(sessionService *session.SessionService) (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())

	redrawCh := make(chan bool)
	errCh := make(chan error)

	infos, err := sessionService.GetInfos()
	if err != nil {
		cancel()
		return nil, err
	}
	w, err := newWidgets(ctx, errCh, infos)
	if err != nil {
		cancel()
		return nil, err
	}

	keyButtons := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
		if k.Key == keyboard.KeyTab {
			resp, err := sessionService.SwitchState()
			if err == nil {
				w.update(resp.Title, "", "", "", []string{resp.DurationString, resp.NextSessionType}, redrawCh)
			}
		}

		if k.Key == '+' {
			update := func(sessionState, timerString, history string) {
				w.update("", "", "", "", []string{timerString, sessionState}, redrawCh)
			}
			errCh <- sessionService.Increment(update)
		}

		if k.Key == '-' {
			update := func(sessionState, timerString, history string) {
				w.update("", "", "", "", []string{timerString, sessionState}, redrawCh)
			}
			errCh <- sessionService.Decrement(update)
		}
	}
	b, err := newButtonSet(ctx, w, redrawCh, sessionService, errCh)
	if err != nil {
		return nil, err
	}

	var t terminalapi.Terminal
	t, err = tcell.New(tcell.ColorMode(terminalapi.ColorMode256))
	if err != nil {
		return nil, err
	}

	c, err := newGrid(b, w, t)
	if err != nil {
		return nil, err
	}
	controller, err := termdash.NewController(t, c, termdash.KeyboardSubscriber(keyButtons))
	if err != nil {
		return nil, err
	}

	return &App{
		ctx:            ctx,
		controller:     controller,
		redrawCh:       redrawCh,
		term:           t,
		errCh:          errCh,
		sessionService: sessionService,
	}, nil
}

func (a *App) resize() error {
	if a.size == a.term.Size() {
		return nil
	}

	a.size = a.term.Size()
	if err := a.term.Clear(); err != nil {
		return err
	}
	return a.controller.Redraw()
}

func (a *App) Run() error {
	defer a.term.Close()
	defer a.controller.Close()

	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-a.redrawCh:
			if err := a.controller.Redraw(); err != nil {
				return err
			}
		case err := <-a.errCh:
			if err != nil {
				return err
			}
		case <-a.ctx.Done():
			return nil
		case <-t.C:
			if err := a.resize(); err != nil {
				return err
			}
		}
	}
}
