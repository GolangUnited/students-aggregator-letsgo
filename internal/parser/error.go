package parser

import (
	"errors"
	"fmt"
)

type (
	ErrorWebPageCannotBeDelivered struct {
		URL        string
		StatusCode int
	}

	ErrorCannotParseArticleDatetime struct {
		OriginError error
	}

	ErrorUnknown struct {
		OriginError error
	}
)

func (e ErrorWebPageCannotBeDelivered) Error() string {
	return fmt.Sprintf("web page cannot be deleivered, request: %s, status code: %d", e.URL, e.StatusCode)
}

func (e ErrorCannotParseArticleDatetime) Error() string {
	return fmt.Errorf("datatime parsing error: %w", e.OriginError).Error()
}

func (e ErrorUnknown) Error() string {
	return fmt.Errorf("uknown error: %w", e.OriginError).Error()
}

var (
	ErrorArticleTitleNotFound       = errors.New("an article title html element not found on a web page")
	ErrorArticleAuthorNotFound      = errors.New("an article author html element not found on a web page")
	ErrorArticleDescriptionNotFound = errors.New("an article description html element not found on a web page")
	ErrorArticleDatetimeNotFound    = errors.New("an article datetime html element not found on a web page")
	ErrorArticleURLNotFound         = errors.New("an article URL html element not found on a web page")
)

const (
	ErrorMessage = "Not Found"
)
