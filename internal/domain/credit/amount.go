package credit

import "fmt"

type Amount int64

func NewAmount(value int64) (Amount, error) {
	if value < 0 {
		return 0, ErrNegativeAmount
	}

	return Amount(value), nil
}

func (a Amount) Value() int64 { return int64(a) }

func (a Amount) IsNegative() bool { return a < 0 }
func (a Amount) IsPositive() bool { return a > 0 }
func (a Amount) IsZero() bool     { return a == 0 }
func (a Amount) String() string   { return fmt.Sprintf("%d", a) }
