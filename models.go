package main

type Flat struct {
	Id           int64
	Number       int64 `pg:",unique"`
	Floor        int8
	Rooms        int8
	SquareTotal  float64
	Section      int8
	Type         string
	PropertyType string
	Price        string
	PriceM2      string
	Status       string
	StatusAlias  string
}
