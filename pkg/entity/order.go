package entity

import "hc-api/pkg/value"

type Order struct {
	value.Request
	value.Person
	value.Authorization
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	ProfileId string `json:"profile_id"`
	OrderSum  int64  `json:"order_sum"`
	// Package Id (ie Product Id) -> Vendor Id -> Service Type -> Rate.
	Rates    map[string]map[string]map[string]string `json:"rates,omitempty"`
	Packages []value.Package                         `json:"packages"`
	// transient data (when outside of this layer)
	PackageIds []string `json:"package_ids,omitempty"`
	Vendor     string   `json:"-"`
}

func (e *Order) Ids() []string {
	if e.Request.Ids == nil {
		return []string{e.Id}
	}
	return e.Request.Ids
}

func (e *Order) Table() *string {
	panic("implement me")
}

func (e *Order) Function() *string {
	panic("implement me")
}

func (e *Order) Payload() []byte {
	panic("implement me")
}

func (e *Order) Validate() error {
	return nil
}
