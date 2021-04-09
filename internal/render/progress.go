package render

import (
	"context"
	"io"
	"strings"
	"time"
)

var (
	progressIndicators = []string{"-", "\\", "|", "/"}
	// HideCursor writes the sequence for hiding cursor
	hideCursor = "\x1b[?25l"
	// ShowCursor writes the sequence for resotring show cursor
	showCursor = "\x1b[?25h"
)

// StartProgress returns a cancel function to stop progress indicator
func StartSimpleProgress(ctx context.Context, w io.Writer, message string) func() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ticker := time.Tick(100 * time.Millisecond)
		w.Write([]byte(message + " "))
		defer func() {
			w.Write([]byte(strings.Repeat(string(rune(KeyCtrlH)), len(message+" "))))
		}()

		w.Write([]byte(hideCursor))
		defer w.Write([]byte(showCursor))

		s := progressIndicators[0]
		w.Write([]byte(s))

		for i := 1; ; i++ {
			w.Write([]byte(strings.Repeat(string(rune(KeyCtrlH)), len(s))))
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				s = progressIndicators[i%len(progressIndicators)]
				w.Write([]byte(s))
			}
		}
	}()

	return cancel
}
