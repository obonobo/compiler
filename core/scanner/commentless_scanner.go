package scanner

import "github.com/obonobo/esac/core/token"

// A scanner that filters comment tokens. The comment tokens are collected and
// can be obtained via CommentlessScanner.Collected
type CommentlessScanner struct {
	Scanner
	Ignoring  []token.Kind
	Collected []token.Token
}

func IgnoringComments(scnr Scanner, comments ...token.Kind) *CommentlessScanner {
	return &CommentlessScanner{Scanner: scnr, Ignoring: comments}
}

// Extract the next token in the program, ignoring comments
func (s *CommentlessScanner) NextToken() (token.Token, error) {
	for {
		next, err := s.Scanner.NextToken()
		isComment := s.isComment(next.Id)
		if isComment {
			s.Collected = append(s.Collected, next)
		}
		if err != nil || !isComment {
			return next, err
		}
	}
}

func (s *CommentlessScanner) isComment(t token.Kind) bool {
	for _, c := range s.Ignoring {
		if c == t {
			return true
		}
	}
	return false
}
