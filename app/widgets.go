package app

import (
	"context"

	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/rinarudhei/pomcli/session"
)

type Tab string

type widgets struct {
	timerTitle       *text.Text
	segDis           *segmentdisplay.SegmentDisplay
	historyT         *text.Text
	summaryT         *text.Text
	commitsT         *text.Text
	guide            *text.Text
	updateTimerTitle chan string
	updateSegDis     chan []string
	updateHistoryT   chan string
	updateSummaryT   chan string
	updateCommitsT   chan string
}

func (w *widgets) update(timerTitle, historyT, summaryT, commitsT string, segDis []string, redrawCh chan<- bool) {
	if timerTitle != "" {
		w.updateTimerTitle <- timerTitle
	}

	if len(segDis) > 0 && segDis[0] != "" && segDis[1] != "" {
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

func newWidgets(ctx context.Context, errCh chan<- error, history string) (*widgets, error) {
	w := &widgets{}
	w.updateSegDis = make(chan []string)
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

	w.historyT, err = newHistoryT(ctx, w.updateHistoryT, errCh, history)
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

func newSegDis(ctx context.Context, updateText <-chan []string, errCh chan<- error) (*segmentdisplay.SegmentDisplay, error) {
	sd, err := segmentdisplay.New(segmentdisplay.AlignHorizontal(align.HorizontalCenter))
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
				var wOpt segmentdisplay.WriteOption
				switch t[1] {
				case session.PomodoroSession:
					wOpt = segmentdisplay.WriteCellOpts(cell.FgColor(cell.ColorRed))
				case session.ShortBreakSession:
					wOpt = segmentdisplay.WriteCellOpts(cell.FgColor(cell.ColorYellow))
				case session.LongBreakSession:
					wOpt = segmentdisplay.WriteCellOpts(cell.FgColor(cell.ColorLime))
				}
				errCh <- sd.Write([]*segmentdisplay.TextChunk{segmentdisplay.NewChunk(t[0], wOpt)})
			case <-ctx.Done():
				return
			}
		}
	}()

	return sd, nil
}

func newHistoryT(ctx context.Context, updateText <-chan string, errCh chan<- error, initial string) (*text.Text, error) {
	t, err := text.New()
	if err != nil {
		return nil, err
	}
	t.Write(initial)

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

	t.Write("- MAMEN")
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
