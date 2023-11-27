package model

type ProductsWithOut struct {
	Region       string
	ProductName  string
	Objective    []string
	Introduction string
	Price        []PricePerPerson
	Include      []string
	Exclusive    []string
	Content      []ContentDate
}
