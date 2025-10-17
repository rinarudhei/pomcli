package app

import (
	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

func newGrid(bs *buttonSet, w *widgets, t terminalapi.Terminal) (*container.Container, error) {
	builder := grid.New()
	builder.Add(
		grid.RowHeightPercWithOpts(50, []container.Option{container.BorderTitle("Pomodoro <Press tab to switch>"), container.Border(linestyle.Round)},
			grid.ColWidthPerc(20),
			grid.ColWidthPerc(60,
				grid.RowHeightPerc(20),
				grid.RowHeightPerc(50,
					grid.ColWidthPerc(10),
					grid.ColWidthPerc(80, grid.Widget(w.segDis, container.AlignHorizontal(align.HorizontalCenter))),
					grid.ColWidthPerc(10),
				),
				grid.RowHeightPerc(30,
					grid.ColWidthPerc(50, grid.Widget(bs.startButton)),
					grid.ColWidthPerc(50, grid.Widget(bs.pauseButton)),
				),
			),
			grid.ColWidthPerc(20),
		),
		grid.RowHeightPerc(48,
			grid.ColWidthPercWithOpts(50, []container.Option{container.BorderTitle("Git commits"), container.Border(linestyle.Round)},
				grid.Widget(w.commitsT),
			),
			grid.ColWidthPercWithOpts(25, []container.Option{container.BorderTitle("History"), container.Border(linestyle.Round)},
				grid.Widget(w.historyT),
			),
			grid.ColWidthPercWithOpts(25, []container.Option{container.BorderTitle("Summary"), container.Border(linestyle.Round)},
				grid.Widget(w.summaryT),
			),
		),
		grid.RowHeightPerc(2, grid.Widget(w.guide)),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}

	c, err := container.New(t, gridOpts...)
	if err != nil {
		return nil, err
	}

	return c, nil
}
