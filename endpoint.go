package main

// curl https://api.data.gov.sg/v1/transport/carpark-availability
type CarparkAvailResponse struct {
	Api_info Api_info
	Items    []Items
}

type Api_info struct {
	Status string
}

type Items struct {
	Timestamp    string
	Carpark_data []Carpark_data
}

type Carpark_data struct {
	Carpark_info    []Carpark_info
	Carpark_number  string
	Update_datetime string
}

type Carpark_info struct {
	Total_lots     string
	Lot_type       string
	Lots_available string
}
