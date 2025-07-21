package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

// --- Color Palette ---
var (
	primaryPurple = [3]int{108, 63, 197}  // #6C3FC5
	blueAccent    = [3]int{41, 128, 185}  // #2980B9
	greenAccent   = [3]int{39, 174, 96}   // #27AE60
	gray          = [3]int{52, 73, 94}    // #34495E
	lightGray     = [3]int{233, 244, 255} // #E9F4FF
	sectionBg     = [3]int{247, 245, 251} // #F7F5FB
)

func addLogo(pdf *gofpdf.Fpdf) {
	logoPath := "public/logo.png"
	if _, err := os.Stat(logoPath); err == nil {
		pdf.ImageOptions(logoPath, 15, 12, 30, 18, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	}
}

func addHeader(pdf *gofpdf.Fpdf, req ItineraryRequest) {
	pdf.SetFillColor(primaryPurple[0], primaryPurple[1], primaryPurple[2])
	pdf.Rect(0, 0, 210, 36, "F")
	addLogo(pdf)
	pdf.SetFont("Arial", "B", 22)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(50, 10)
	pdf.CellFormat(150, 10, fmt.Sprintf("Hi, %s!", req.GreetingName), "", 1, "R", false, 0, "")
	pdf.SetFont("Arial", "B", 16)
	pdf.SetXY(50, 18)
	pdf.CellFormat(150, 8, fmt.Sprintf("%s Travel Itinerary", req.Destination), "", 1, "R", false, 0, "")
	pdf.SetFont("Arial", "", 13)
	pdf.SetXY(50, 26)
	pdf.CellFormat(150, 7, fmt.Sprintf("%d Days %d Nights", req.DurationDays, req.DurationNights), "", 1, "R", false, 0, "")
	pdf.SetFont("Arial", "", 13)
	pdf.SetXY(50, 32)
	pdf.CellFormat(150, 6, "✈️ 🏨 🚗 🚌", "", 1, "R", false, 0, "")
}

func addTripSummary(pdf *gofpdf.Fpdf, req ItineraryRequest) {
	pdf.SetY(40)
	pdf.SetFont("Arial", "B", 13)
	pdf.SetTextColor(primaryPurple[0], primaryPurple[1], primaryPurple[2])
	pdf.SetFillColor(lightGray[0], lightGray[1], lightGray[2])
	pdf.SetX(15)
	pdf.CellFormat(180, 12, "Trip Summary", "1", 1, "C", true, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.SetTextColor(gray[0], gray[1], gray[2])
	pdf.SetX(15)
	pdf.CellFormat(36, 10, "From", "1", 0, "C", false, 0, "")
	pdf.CellFormat(36, 10, "Departure", "1", 0, "C", false, 0, "")
	pdf.CellFormat(36, 10, "Arrival", "1", 0, "C", false, 0, "")
	pdf.CellFormat(36, 10, "Destination", "1", 0, "C", false, 0, "")
	pdf.CellFormat(36, 10, "Travellers", "1", 1, "C", false, 0, "")
	pdf.SetX(15)
	pdf.CellFormat(36, 10, req.DepartureFrom, "1", 0, "C", false, 0, "")
	pdf.CellFormat(36, 10, req.DepartureDate, "1", 0, "C", false, 0, "")
	pdf.CellFormat(36, 10, req.ArrivalDate, "1", 0, "C", false, 0, "")
	pdf.CellFormat(36, 10, req.Destination, "1", 0, "C", false, 0, "")
	pdf.CellFormat(36, 10, fmt.Sprintf("%d", req.NumTravellers), "1", 1, "C", false, 0, "")
	pdf.Ln(4)
}

func addItinerary(pdf *gofpdf.Fpdf, req ItineraryRequest) {
	dayKeys := map[string]bool{}
	for k := range req.Activities {
		dayKeys[k] = true
	}
	for k := range req.Flights {
		dayKeys[k] = true
	}
	for k := range req.Transfers {
		dayKeys[k] = true
	}
	for k := range req.Hotels {
		dayKeys[k] = true
	}
	var sortedDays []string
	for k := range dayKeys {
		sortedDays = append(sortedDays, k)
	}
	sort.Strings(sortedDays)

	for i, date := range sortedDays {
		startY := pdf.GetY()
		pdf.SetFillColor(sectionBg[0], sectionBg[1], sectionBg[2])
		pdf.Rect(15, startY, 180, 54, "F")
		pdf.SetY(startY + 6)
		pdf.SetX(22)
		pdf.SetFont("Arial", "B", 14)
		pdf.SetTextColor(primaryPurple[0], primaryPurple[1], primaryPurple[2])
		pdf.CellFormat(0, 8, fmt.Sprintf("Day %d: %s", i+1, date), "", 1, "L", false, 0, "")
		pdf.SetFont("Arial", "", 11)
		pdf.SetTextColor(gray[0], gray[1], gray[2])
		pdf.SetX(28)
		times := []string{"Morning", "Afternoon", "Evening"}
		for _, t := range times {
			pdf.SetFont("Arial", "B", 11)
			pdf.SetTextColor(blueAccent[0], blueAccent[1], blueAccent[2])
			pdf.CellFormat(0, 7, t, "", 1, "L", false, 0, "")
			pdf.SetFont("Arial", "", 10)
			pdf.SetTextColor(gray[0], gray[1], gray[2])
			pdf.SetX(34)
			if acts, ok := req.Activities[date]; ok && len(acts) > 0 {
				for _, a := range acts {
					if strings.HasPrefix(a, t+" - ") {
						pdf.CellFormat(0, 6, "• "+strings.TrimPrefix(a, t+" - "), "", 1, "L", false, 0, "")
					}
				}
			}
			if flights, ok := req.Flights[date]; ok && len(flights) > 0 && t == "Morning" {
				for _, f := range flights {
					pdf.SetTextColor(greenAccent[0], greenAccent[1], greenAccent[2])
					pdf.CellFormat(0, 6, "✈️ "+fmt.Sprintf("%s (%s → %s @ %s)", f.FlightNo, f.From, f.To, f.Time), "", 1, "L", false, 0, "")
				}
			}
			if tr, ok := req.Transfers[date]; ok && tr != "" && t == "Morning" {
				pdf.SetTextColor(blueAccent[0], blueAccent[1], blueAccent[2])
				pdf.CellFormat(0, 6, "🚖 "+tr, "", 1, "L", false, 0, "")
			}
			if hotel, ok := req.Hotels[date]; ok && t == "Evening" {
				pdf.SetTextColor(primaryPurple[0], primaryPurple[1], primaryPurple[2])
				pdf.CellFormat(0, 6, "🏨 "+hotel.Name, "", 1, "L", false, 0, "")
				pdf.SetTextColor(gray[0], gray[1], gray[2])
				pdf.CellFormat(0, 6, hotel.Address, "", 1, "L", false, 0, "")
				pdf.CellFormat(0, 6, fmt.Sprintf("Check-in: %s | Check-out: %s", hotel.CheckIn, hotel.CheckOut), "", 1, "L", false, 0, "")
			}
		}
		pdf.Ln(8)
	}
}

func addFlightSummary(pdf *gofpdf.Fpdf, req ItineraryRequest) {
	pdf.AddPage()
	addHeader(pdf, req)
	pdf.SetY(40)
	pdf.SetFont("Arial", "B", 15)
	pdf.SetTextColor(primaryPurple[0], primaryPurple[1], primaryPurple[2])
	pdf.CellFormat(0, 10, "Flight Summary", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(gray[0], gray[1], gray[2])
	pdf.SetFillColor(lightGray[0], lightGray[1], lightGray[2])
	pdf.SetX(15)
	pdf.CellFormat(35, 8, "Date", "1", 0, "C", true, 0, "")
	pdf.CellFormat(45, 8, "Airline", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 8, "From", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 8, "To", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 8, "Flight No.", "1", 1, "C", true, 0, "")
	pdf.SetFont("Arial", "", 11)
	for _, f := range req.FlightSummary {
		pdf.SetX(15)
		pdf.CellFormat(35, 8, f.Date, "1", 0, "C", false, 0, "")
		pdf.CellFormat(45, 8, f.Airline, "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 8, f.From, "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 8, f.To, "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 8, f.FlightNo, "1", 1, "C", false, 0, "")
	}
	pdf.Ln(4)
}

func addHotelBookings(pdf *gofpdf.Fpdf, req ItineraryRequest) {
	pdf.SetFont("Arial", "B", 15)
	pdf.SetTextColor(primaryPurple[0], primaryPurple[1], primaryPurple[2])
	pdf.CellFormat(0, 10, "Hotel Bookings", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(gray[0], gray[1], gray[2])
	pdf.SetFillColor(lightGray[0], lightGray[1], lightGray[2])
	pdf.SetX(15)
	pdf.CellFormat(30, 8, "City", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Check In", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Check Out", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 8, "Nights", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 8, "Hotel Name", "1", 1, "C", true, 0, "")
	pdf.SetFont("Arial", "", 11)
	for _, h := range req.HotelBookings {
		pdf.SetX(15)
		pdf.CellFormat(30, 8, h.City, "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, h.CheckIn, "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, h.CheckOut, "1", 0, "C", false, 0, "")
		pdf.CellFormat(20, 8, fmt.Sprintf("%d", h.Nights), "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, 8, h.HotelName, "1", 1, "C", false, 0, "")
	}
	pdf.Ln(4)
}

func addNotes(pdf *gofpdf.Fpdf, req ItineraryRequest) {
	pdf.SetFont("Arial", "B", 15)
	pdf.SetTextColor(primaryPurple[0], primaryPurple[1], primaryPurple[2])
	pdf.CellFormat(0, 10, "Important Notes", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(gray[0], gray[1], gray[2])
	pdf.SetFillColor(lightGray[0], lightGray[1], lightGray[2])
	pdf.SetX(15)
	pdf.CellFormat(50, 8, "Point", "1", 0, "C", true, 0, "")
	pdf.CellFormat(120, 8, "Details", "1", 1, "C", true, 0, "")
	pdf.SetFont("Arial", "", 11)
	for _, n := range req.ImportantNotes {
		pdf.SetX(15)
		pdf.CellFormat(50, 8, n.Point, "1", 0, "L", false, 0, "")
		pdf.CellFormat(120, 8, n.Detail, "1", 1, "L", false, 0, "")
	}
	pdf.Ln(4)
}

func addPaymentPlan(pdf *gofpdf.Fpdf, req ItineraryRequest) {
	pdf.SetFont("Arial", "B", 15)
	pdf.SetTextColor(primaryPurple[0], primaryPurple[1], primaryPurple[2])
	pdf.CellFormat(0, 10, "Payment Plan", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(gray[0], gray[1], gray[2])
	pdf.SetFillColor(lightGray[0], lightGray[1], lightGray[2])
	pdf.SetX(15)
	pdf.CellFormat(60, 8, "Total Amount", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 8, "TCS", "1", 1, "C", true, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.SetX(15)
	pdf.CellFormat(60, 8, req.PaymentPlan.TotalAmount, "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 8, req.PaymentPlan.TCS, "1", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 11)
	pdf.SetX(15)
	pdf.CellFormat(40, 8, "Installment", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Amount", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Due Date", "1", 1, "C", true, 0, "")
	pdf.SetFont("Arial", "", 11)
	for _, inst := range req.PaymentPlan.Installments {
		pdf.SetX(15)
		pdf.CellFormat(40, 8, inst.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 8, inst.Amount, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 8, inst.DueDate, "1", 1, "C", false, 0, "")
	}
	pdf.Ln(4)
}

func addVisa(pdf *gofpdf.Fpdf, req ItineraryRequest) {
	pdf.SetFont("Arial", "B", 15)
	pdf.SetTextColor(primaryPurple[0], primaryPurple[1], primaryPurple[2])
	pdf.CellFormat(0, 10, "Visa Details", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(gray[0], gray[1], gray[2])
	pdf.SetFillColor(lightGray[0], lightGray[1], lightGray[2])
	pdf.SetX(15)
	pdf.CellFormat(40, 8, "Visa Type", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Validity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Processing Date", "1", 1, "C", true, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.SetX(15)
	pdf.CellFormat(40, 8, req.Visa.Type, "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 8, req.Visa.Validity, "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 8, req.Visa.ProcessingDate, "1", 1, "C", false, 0, "")
	pdf.Ln(4)
}

func addFooter(pdf *gofpdf.Fpdf, req ItineraryRequest) {
	pdf.SetY(287)
	pdf.SetFont("Arial", "I", 11)
	pdf.SetTextColor(primaryPurple[0], primaryPurple[1], primaryPurple[2])
	pdf.CellFormat(0, 10, "PLAN.PACK.GO.   |   Generated by Vigovia   |   Page "+fmt.Sprint(pdf.PageNo()), "", 0, "C", false, 0, "")
	pdf.SetY(295)
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(gray[0], gray[1], gray[2])
	pdf.CellFormat(0, 6, req.Contact.CompanyName+" | "+req.Contact.Address+" | Phone: "+req.Contact.Phone+" | Email: "+req.Contact.Email, "", 0, "C", false, 0, "")
}

func CreatePDF(req ItineraryRequest) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(0, 0, 0)
	pdf.AddPage()
	addHeader(pdf, req)
	addTripSummary(pdf, req)
	addItinerary(pdf, req)
	addFlightSummary(pdf, req)
	addHotelBookings(pdf, req)
	addNotes(pdf, req)
	addPaymentPlan(pdf, req)
	addVisa(pdf, req)
	addFooter(pdf, req)

	outputDir := "./pdfs"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", err
	}
	safeName := strings.ReplaceAll(req.Name, " ", "_")
	filePath := filepath.Join(outputDir, fmt.Sprintf("itinerary_%s_%s.pdf", safeName, req.StartDate))
	if err := pdf.OutputFileAndClose(filePath); err != nil {
		return "", err
	}

	return filePath, nil
}
