package app

import (
	"context"

	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum4k/termdash/widgets/text"
)

type widgets struct {
	timerTitle       *text.Text
	segDis           *segmentdisplay.SegmentDisplay
	historyT         *text.Text
	summaryT         *text.Text
	commitsT         *text.Text
	guide            *text.Text
	updateTimerTitle chan string
	updateSegDis     chan string
	updateHistoryT   chan string
	updateSummaryT   chan string
	updateCommitsT   chan string
}

func (w *widgets) update(timerTitle, segDis, historyT, summaryT, commitsT string, redrawCh chan<- bool) {
	if timerTitle != "" {
		w.updateTimerTitle <- timerTitle
	}

	if segDis != "" {
		w.updateSegDis <- segDis
	}

	if historyT != "" {
		w.updateHistoryT <- historyT
	}

	if summaryT != "" {
		w.updateSummaryT <- summaryT
	}

	if commitsT != "" {
		w.updateCommitsT <- commitsT
	}

	redrawCh <- true
}

func newWidgets(ctx context.Context, errCh chan<- error) (*widgets, error) {
	w := &widgets{}
	w.updateSegDis = make(chan string)
	w.updateHistoryT = make(chan string)
	w.updateSummaryT = make(chan string)
	w.updateCommitsT = make(chan string)
	w.updateTimerTitle = make(chan string)

	var err error
	w.guide, err = text.New()
	if err != nil {
		return nil, err
	}
	w.guide.Write("Press [q] to quit; [+]/[-] increment/decrement", text.WriteCellOpts(&cell.Options{Dim: true}))

	w.timerTitle, err = newTimerTitle(ctx, w.updateTimerTitle, errCh)
	if err != nil {
		return nil, err
	}
	w.segDis, err = newSegDis(ctx, w.updateSegDis, errCh)
	if err != nil {
		return nil, err
	}

	w.historyT, err = newHistoryT(ctx, w.updateHistoryT, errCh)
	if err != nil {
		return nil, err
	}

	w.summaryT, err = newSummaryT(ctx, w.updateSummaryT, errCh)
	if err != nil {
		return nil, err
	}

	w.commitsT, err = newCommitsT(ctx, w.updateCommitsT, errCh)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func newTimerTitle(ctx context.Context, updateText <-chan string, errCh chan<- error) (*text.Text, error) {
	t, err := text.New()
	if err != nil {
		return nil, err
	}
	t.Write("Pomodoro - focus")

	go func() {
		for {
			select {
			case txt := <-updateText:
				t.Reset()
				errCh <- t.Write(txt)
			case <-ctx.Done():
				return
			}
		}
	}()

	return t, nil
}

func newSegDis(ctx context.Context, updateText <-chan string, errCh chan<- error) (*segmentdisplay.SegmentDisplay, error) {
	sd, err := segmentdisplay.New(segmentdisplay.AlignHorizontal(align.HorizontalRight))
	if err != nil {
		return nil, err
	}
	sd.Write([]*segmentdisplay.TextChunk{
		segmentdisplay.NewChunk(
			"50:00",
			segmentdisplay.WriteCellOpts(cell.FgColor(cell.ColorRed)),
		),
	},
	)

	go func() {
		for {
			select {
			case t := <-updateText:
				sd.Reset()
				errCh <- sd.Write([]*segmentdisplay.TextChunk{segmentdisplay.NewChunk(t)})
			case <-ctx.Done():
				return
			}
		}
	}()

	return sd, nil
}

func newHistoryT(ctx context.Context, updateText <-chan string, errCh chan<- error) (*text.Text, error) {
	t, err := text.New()
	if err != nil {
		return nil, err
	}
	t.Write("-")

	go func() {
		for {
			select {
			case updatedText := <-updateText:
				t.Reset()
				errCh <- t.Write(updatedText)
			case <-ctx.Done():
				return
			}
		}
	}()

	return t, nil
}

func newSummaryT(ctx context.Context, updateText <-chan string, errCh chan<- error) (*text.Text, error) {
	t, err := text.New()
	if err != nil {
		return nil, err
	}

	t.Write(`sessions              : 0\n 
		focus time            : 0\n 
		break time            : 0\n 
		commits               : 0\n 
		code insertion (lines): 0\n 
		distractions          : -
		`)

	go func() {
		for {
			select {
			case updatedText := <-updateText:
				t.Reset()
				errCh <- t.Write(updatedText)
			case <-ctx.Done():
				return
			}
		}
	}()

	return t, nil
}

func newCommitsT(ctx context.Context, updateText <-chan string, errCh chan<- error) (*text.Text, error) {
	t, err := text.New()
	if err != nil {
		return nil, err
	}
	t.Write(
		`-`,
	)

	go func() {
		for {
			select {
			case updatedText := <-updateText:
				t.Reset()
				errCh <- t.Write(updatedText)
			case <-ctx.Done():
				return
			}
		}
	}()

	return t, nil
}
