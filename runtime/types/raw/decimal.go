package raw

import (
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

type prismaDecimalValue struct {
	Value decimal.Decimal `json:"prisma__value"`
	Type  string          `json:"prisma__type"`
}

type Decimal struct {
	decimal.Decimal
}

func (r *Decimal) String() string {
	return r.Decimal.String()
}

func (r *Decimal) MarshalJSON() ([]byte, error) {
	return r.Decimal.MarshalJSON()
}

func (r *Decimal) UnmarshalJSON(data []byte) error {
	var v prismaDecimalValue
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v.Type != "decimal" {
		return fmt.Errorf("invalid type %s, expected decimal", v.Type)
	}
	*r = Decimal{v.Value}
	return nil
}
