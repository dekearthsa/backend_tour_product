package model

type ProductOut struct {
	Region       string
	ProductName  string
	Objective    []string
	Introduction string
	Price        []PricePerPerson
	Include      []string
	Exclusive    []string
	Content      []ContentDate
	ImagePath    []string
	ArrayBase64  []string
}
