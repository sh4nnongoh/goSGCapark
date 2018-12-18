package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func GetCarparkAvailInfo() chan CarparkAvailInfo {
	var wg sync.WaitGroup
	c := make(chan CarparkAvailInfo)

	req, err := http.NewRequest("GET", "https://api.data.gov.sg/v1/transport/carpark-availability", nil)
	if err != nil {
		log.Fatalln("error forming HTTP request:", err)
		return nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("error obtaining HTTP response:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalln("error HTTP response code:", errors.New(resp.Status))
		return nil
	}

	var jRsp CarparkAvailResponse
	json.NewDecoder(resp.Body).Decode(&jRsp)
	wg.Add(len(jRsp.Items[0].Carpark_data))
	for _, r := range jRsp.Items[0].Carpark_data {
		LotsTotal, err := strconv.Atoi(r.Carpark_info[0].Total_lots)
		if err != nil {
			log.Fatalln("error converting alphabet to integer:", err)
			return nil
		}
		LotsAvailable, err := strconv.Atoi(r.Carpark_info[0].Lots_available)
		if err != nil {
			log.Fatalln("error converting alphabet to integer:", err)
			return nil
		}
		go func(r Carpark_data) {
			defer wg.Done()
			c <- CarparkAvailInfo{
				r.Update_datetime,
				r.Carpark_number,
				LotsTotal,
				LotsAvailable,
				r.Carpark_info[0].Lot_type}
		}(r)
	}

	go func() {
		wg.Wait()
		close(c)
	}()
	return c
}

func NewCarparkAvailabilityService() CarparkAvailabilityService {
	var report CarparkAvailabilityReport
	for info := range GetCarparkAvailInfo() {
		report.CarparkAvailInfo = append(report.CarparkAvailInfo, info)
	}
	return report
}

type CarparkAvailabilityService interface {
	PrintCarparkAvailabilityCSV()
}

type CarparkAvailabilityReport struct {
	CarparkAvailInfo []CarparkAvailInfo
}

func (i CarparkAvailabilityReport) PrintCarparkAvailabilityCSV() {
	w := csv.NewWriter(os.Stdout)
	headers := CarparkAvailInfo{}.GetHeaders()
	if err := w.Write(headers); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}
	for _, r := range i.CarparkAvailInfo {
		values := r.ToSlice()
		if err := w.Write(values); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

type CarparkAvailInfo struct {
	Timestamp     string
	CarparkNumber string
	LotsTotal     int
	LotsAvailable int
	LotType       string
}

func (r CarparkAvailInfo) GetHeaders() []string {
	return []string{"Timestamp", "CarparkNumber", "LotsTotal", "LotsAvailable", "LotType"}
}

func (r CarparkAvailInfo) ToSlice() []string {
	LotsTotal := strconv.Itoa(r.LotsTotal)
	LotsAvailable := strconv.Itoa(r.LotsAvailable)
	return []string{r.Timestamp, r.CarparkNumber, LotsTotal, LotsAvailable, r.LotType}
}
