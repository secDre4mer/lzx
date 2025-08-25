package slidingwindow

import "errors"

// SlidingWindow implements a window of n bytes where old bytes are replaced with new ones. Old bytes within the last
// window size can be looked up again.
type SlidingWindow struct {
	windowData []byte

	position int
}

func New(size int) *SlidingWindow {
	return &SlidingWindow{
		windowData: make([]byte, size),
		position:   0,
	}
}

var ErrInvalidOffset = errors.New("invalid offset")

func (s *SlidingWindow) Lookback(offset int) (byte, error) {
	if offset < 0 || offset > len(s.windowData) {
		return 0, ErrInvalidOffset
	}
	return s.windowData[(s.position-offset+len(s.windowData))%len(s.windowData)], nil
}

func (s *SlidingWindow) Add(b byte) {
	s.windowData[s.position] = b
	s.position = (s.position + 1) % len(s.windowData)
}

func (s *SlidingWindow) Size() int {
	return len(s.windowData)
}
