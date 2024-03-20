package models

type Product struct {
	ID             int
	OrderID        int
	Name           string
	Amount         int
	MainRack       string
	SecondaryRacks []string
}
