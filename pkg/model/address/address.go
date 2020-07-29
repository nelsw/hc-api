package address

import (
	"fmt"
	"strconv"
	"strings"
)

type Entity struct {
	Id     string `json:"id" xml:"ID,attr"`
	Unit   string `json:"unit,omitempty" xml:"Address1"`
	Street string `json:"street" xml:"Address2"`
	City   string `json:"city" xml:"City"`
	State  string `json:"state" xml:"State"`
	Zip5   string `json:"zip_5" xml:"Zip5"`
	Zip4   string `json:"zip_4,omitempty" xml:"Zip4"`
}

func (e *Entity) Validate() error {
	if _, err := strconv.Atoi(e.Zip5); err != nil || len(e.Zip5) < 5 {
		return err
	} else if len(e.Street) < 5 {
		return fmt.Errorf("bad street\n")
	} else {
		return nil
	}
}

func (e *Entity) String() string {
	var sb strings.Builder
	sb.WriteString(e.Street)
	sb.WriteString(", ")
	if e.Unit != "" {
		sb.WriteString(e.Unit)
		sb.WriteString(", ")
	}
	sb.WriteString(e.City)
	sb.WriteString(", ")
	sb.WriteString(e.State)
	sb.WriteString(", ")
	sb.WriteString(e.Zip5)
	sb.WriteString("-")
	sb.WriteString(e.Zip4)
	sb.WriteString(", ")
	sb.WriteString("United States")
	return sb.String()
}
