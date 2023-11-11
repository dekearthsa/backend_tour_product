package model

type Products struct {
	Region       string
	ProductName  string
	Objective    []string
	Introduction string
	Price        []PricePerPerson
	Include      []string
	Exclusive    []string
	Content      []ContentDate
	ImagePath    []string
}

type ContentDate struct {
	Title   string
	Content string
}

type PricePerPerson struct {
	Person string
	Price  string
}
