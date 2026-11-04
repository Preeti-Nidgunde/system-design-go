package parkinglot

type Vehiclesize int

const (
	SMALL Vehiclesize = iota
	MEDIUM
	LARGE
)

type Vehicle struct {
	number string
	size   Vehiclesize
}

func (v *Vehicle) GetSize() Vehiclesize {
	return v.size
}

func (v *Vehicle) GetLicense() string {
	return v.number
}
