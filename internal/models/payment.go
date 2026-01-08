package models

type Payment struct {
	Id                  string
	User_id             string
	Site_id             string
	Yookassa_payment_id string
	Status              string
	Amount              float64
	Currency            string
	Description         string
}
