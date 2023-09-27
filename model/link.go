package model

import "errors"

type Link struct {
	Id    int64  `db:"id" json:"id"`
	Short string `db:"short" json:"short"`
	Long  string `db:"long" json:"long"`
}

type CreateLinkRequest struct {
	Long string `json:"long"`
}

var LongValidationError = errors.New("long validation error")

func (p *CreateLinkRequest) Validate() error {
	if len(p.Long) == 0 || len(p.Long) > 512 {
		return LongValidationError
	}

	return nil
}

type CreateLinkResponse struct {
	Short string `json:"short"`
}
