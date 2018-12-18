package main

import "log"

func main() {
	log.Println("Hello World")

	CarparkAvailabilityReport := NewCarparkAvailabilityService()
	CarparkAvailabilityReport.PrintCarparkAvailabilityCSV()
}
