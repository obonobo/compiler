package scanner

import "github.com/obonobo/esac/core/token"

const DEFAULT_CHANNEL_SIZE = 1024

// Use ObservableScanner if you want to consume scanner tokens in multiple
// places. ObservableScanner allows you to subscribe new observers using
// ObservableScanner.Subscribe() or ObservableScanner.SubscribeSize(). These
// methods return channels that will receive tokens whenever Scanner.NextToken
// is called.
//
// Subscriber channels will be closed once the underlying Scanner encounters an
// error, be it an io.EOF error or otherwise. Use ObservableScanner.Err() to
// inspect the Scanner error afterwards.
//
// Note: you should only have one consumer calling
// ObservableScanner.NextToken(), the other consumers should use their
// subscription channels to receive the same token.
//
// Note: if a consumer calls Scanner.NextToken() on the underlying scanner, that
// token will not appear to subscribers
type ObservableScanner struct {
	Scanner
	err    error
	closed bool
	subs   []chan<- token.Token
}

func NewObservableScanner(scanner Scanner) *ObservableScanner {
	return &ObservableScanner{Scanner: scanner}
}

func (s *ObservableScanner) NextToken() (token.Token, error) {
	if s.err != nil {
		return token.Token{}, s.err
	}
	tok, err := s.nextToken()
	if err != nil {
		return token.Token{}, err
	}
	s.notify(tok)
	return tok, err
}

func (s *ObservableScanner) Err() error {
	return s.err
}

func (s *ObservableScanner) Subscribe() <-chan token.Token {
	return s.SubscribeSize(DEFAULT_CHANNEL_SIZE)
}

func (s *ObservableScanner) SubscribeSize(size int) <-chan token.Token {
	c := make(chan token.Token, size)
	s.subs = append(s.subs, c)
	return c
}

// Manually closes all subscription channels of the scanner. This method is only
// to be used if you want to prematurely close subscribers. Otherwise, you can
// let the scanner encounter EOF and it will close the subscribers on its own.
func (s *ObservableScanner) Close() error {
	if s.closed {
		return nil
	}
	for _, sub := range s.subs {
		close(sub)
	}
	s.closed = true
	return nil
}

func (s *ObservableScanner) nextToken() (token.Token, error) {
	tok, err := s.Scanner.NextToken()
	s.err = err
	if err != nil {
		s.Close()
	}
	return tok, err
}

func (s *ObservableScanner) notify(tok token.Token) {
	for _, sub := range s.subs {
		sub <- tok
	}
}
