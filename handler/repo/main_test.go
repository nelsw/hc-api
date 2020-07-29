package main

import (
	"sam-app/pkg/model/credential"
	"sam-app/pkg/model/request"
	"testing"
)

//func TestSave200(t *testing.T) {
//	e := profile.Entity{
//		Id:        test.ProfileId,
//		Email:     "connor@wiesow.com",
//		FirstName: "Connor",
//		LastName:  "Van Helsing",
//		Phone:     "555-555-5555",
//	}
//	r := request.Entity{
//		Id:      e.Id,
//		Type:    util.TypeOf(&e),
//		Table:   e.TableName(),
//		Keyword: "save",
//		Result:  &e,
//	}
//	if out, err := Handle(r); err != nil {
//		t.Fatal(err)
//	} else {
//		t.Log(out)
//	}
//}
//
//func TestSave400(t *testing.T) {
//	e := profile.Entity{
//		Id:        test.ProfileId,
//		Email:     "connor@wiesow.com",
//		FirstName: "Connor",
//		LastName:  "Van Helsing",
//		Phone:     "555-555-5555",
//	}
//	r := request.Entity{
//		Id:      e.Id,
//		Type:    util.TypeOf(&e),
//		Table:   e.TableName(),
//		Keyword: "save",
//		Result:  &e,
//	}
//	if out, err := Handle(r); err != nil {
//		t.Fatal(err)
//	} else {
//		t.Log(out)
//	}
//}

func TestFindOne200(t *testing.T) {
	e := credential.Entity{"hello@gmail.com", "Pass123!", ""}
	r := request.Entity{
		Id:      e.Id,
		Type:    "*credential.Entity",
		Table:   "credential",
		Keyword: "find-one",
		Result:  e,
	}
	if out, err := Handle(r); err != nil {
		t.Fatal(err)
	} else {
		t.Log(out)
	}
}

//func TestFindMany200(t *testing.T) {
//	e := product.Entity{}
//	r := request.Entity{
//		Ids:     []string{test.ProductId, test.ProductId2},
//		Type:    util.TypeOf(&e),
//		Table:   e.TableName(),
//		Keyword: "find-many",
//		Result:  &e,
//	}
//
//	if out, err := Handle(r); err != nil {
//		t.Fatal(err)
//	} else {
//
//		b, err := json.Marshal(out)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		var arr []product.Entity
//		if err := json.Unmarshal(b, &arr); err != nil {
//			t.Fatal(err)
//		}
//		t.Log(arr)
//	}
//}
//
//func TestFindAll200(t *testing.T) {
//	e := product.Entity{}
//	r := request.Entity{
//		Attributes: map[string]string{"brand_id": "cbdrevolution"},
//		Type:       util.TypeOf(&e),
//		Table:      e.TableName(),
//		Keyword:    "find-all",
//		Result:     &e,
//	}
//
//	if out, err := Handle(r); err != nil {
//		t.Fatal(err)
//	} else {
//
//		b, err := json.Marshal(out)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		var arr []product.Entity
//		if err := json.Unmarshal(b, &arr); err != nil {
//			t.Fatal(err)
//		}
//		t.Log(arr)
//	}
//}
//
//func TestFindAll400(t *testing.T) {
//	e := product.Entity{}
//	r := request.Entity{
//		Attributes: map[string]string{"brand_id": "cbdrevolution"},
//		Type:       util.TypeOf(&e),
//		Table:      "",
//		Keyword:    "find-all",
//		Result:     &e,
//	}
//
//	if out, err := Handle(r); err == nil {
//		t.Fatal(err)
//	} else {
//		t.Log(out)
//	}
//}
//
//func TestFindMany400(t *testing.T) {
//	e := product.Entity{}
//	r := request.Entity{
//		Ids:     []string{},
//		Type:    util.TypeOf(&e),
//		Table:   e.TableName(),
//		Keyword: "find-many",
//		Result:  &e,
//	}
//
//	if out, err := Handle(r); err == nil {
//		t.Fatal(err)
//	} else {
//		t.Log(out)
//	}
//}
//
//func TestFind404(t *testing.T) {
//	e := product.Entity{}
//	r := request.Entity{
//		Id:      "",
//		Type:    util.TypeOf(&e),
//		Table:   e.TableName(),
//		Keyword: "find-one",
//		Result:  &e,
//	}
//	if out, err := Handle(r); err == nil {
//		t.Fatal(err)
//	} else {
//		t.Log(out)
//	}
//}
//
//func TestBadKeyword(t *testing.T) {
//	e := product.Entity{}
//	r := request.Entity{
//		Id:      test.ProductId,
//		Type:    util.TypeOf(&e),
//		Table:   e.TableName(),
//		Keyword: "find-one",
//		Result:  &e,
//	}
//	r.Keyword = ""
//	if _, err := Handle(r); err != InvalidKeyword {
//		t.Fail()
//	}
//}
//
//func TestBadType(t *testing.T) {
//	e := product.Entity{}
//	r := request.Entity{
//		Id:      test.ProductId,
//		Type:    "",
//		Table:   e.TableName(),
//		Keyword: "find-one",
//		Result:  &e,
//	}
//	r.Type = ""
//	if _, err := Handle(r); err != InvalidType {
//		t.Fail()
//	}
//}
//
//func TestAdd(t *testing.T) {
//	u := user.Request{}
//	r := request.Entity{
//		Id:      test.UserId,
//		Table:   u.TableName(),
//		Ids:     []string{test.ProductId2},
//		Keyword: "add product_ids",
//	}
//	if out, err := Handle(r); err != nil {
//		t.Fatal(err)
//	} else {
//		t.Log(out)
//	}
//}
//
//func TestDelete(t *testing.T) {
//	u := user.Request{}
//	r := request.Entity{
//		Id:      test.UserId,
//		Table:   u.TableName(),
//		Ids:     []string{test.ProductId2},
//		Keyword: "delete product_ids",
//	}
//	if out, err := Handle(r); err != nil {
//		t.Fatal(err)
//	} else {
//		t.Log(out)
//	}
//}
//
//func TestRemove(t *testing.T) {
//	e := product.Entity{Id: "dd514926-177a-11ea-b04e-ba9247c7c980"}
//	r := request.Entity{
//		Id:      e.Id,
//		Type:    util.TypeOf(e),
//		Table:   e.TableName(),
//		Keyword: "remove",
//		Result:  e,
//	}
//	if out, err := Handle(r); err != nil {
//		t.Fatal(err)
//	} else {
//		t.Log(out)
//	}
//}

// for code coverage purposes only
func TestHandleMain(t *testing.T) {
	go main()
}
