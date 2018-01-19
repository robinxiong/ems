package location

// Location is a struct, you could embedded it into your model to get the Location feature for your model
type Location struct {
	Address   string
	City      string
	Region    string
	Country   string
	Zip       string
	Latitude  float64
	Longitude float64
}
