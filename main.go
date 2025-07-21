package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ItineraryRequest struct {
	// Header & Trip Info
	Name           string   `json:"name"`
	GreetingName   string   `json:"greeting_name"`
	DepartureFrom  string   `json:"departure_from"`
	DepartureDate  string   `json:"departure_date"`
	ArrivalDate    string   `json:"arrival_date"`
	Destination    string   `json:"destination"`
	NumTravellers  int      `json:"num_travellers"`
	StartDate      string   `json:"start_date"`
	EndDate        string   `json:"end_date"`
	DurationNights int      `json:"duration_nights"`
	DurationDays   int      `json:"duration_days"`
	Destinations   []string `json:"destinations"`

	// Day-wise Details
	Activities map[string][]string            `json:"activities"`
	Flights    map[string][]FlightInfo        `json:"flights"`
	Transfers  map[string]string              `json:"transfers"`
	Hotels     map[string]HotelBookingDetails `json:"hotels"`

	// Flight Summary Table
	FlightSummary []FlightSummaryInfo `json:"flight_summary"`

	// Hotel Bookings Table
	HotelBookings []HotelBookingTableInfo `json:"hotel_bookings"`

	// Important Notes
	ImportantNotes []NotePoint `json:"important_notes"`

	// Scope of Service
	ScopeOfService []ServicePoint `json:"scope_of_service"`

	// Inclusion Summary
	InclusionSummary []InclusionPoint `json:"inclusion_summary"`

	// Payment Plan
	PaymentPlan PaymentPlanDetails `json:"payment_plan"`

	// Visa Details
	Visa VisaDetails `json:"visa"`

	// Contact Info
	Contact ContactInfo `json:"contact"`
}

// Supporting Structs

type FlightInfo struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Time     string `json:"time"`
	FlightNo string `json:"flight_no"`
}

type HotelBookingDetails struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	CheckIn  string `json:"check_in"`
	CheckOut string `json:"check_out"`
}

type FlightSummaryInfo struct {
	Date     string `json:"date"`
	Airline  string `json:"airline"`
	From     string `json:"from"`
	To       string `json:"to"`
	FlightNo string `json:"flight_no"`
}

type HotelBookingTableInfo struct {
	City      string `json:"city"`
	CheckIn   string `json:"check_in"`
	CheckOut  string `json:"check_out"`
	Nights    int    `json:"nights"`
	HotelName string `json:"hotel_name"`
}

type NotePoint struct {
	Point  string `json:"point"`
	Detail string `json:"detail"`
}

type ServicePoint struct {
	Service string `json:"service"`
	Detail  string `json:"detail"`
}

type InclusionPoint struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
	Detail   string `json:"detail"`
	Status   string `json:"status"`
}

type PaymentPlanDetails struct {
	TotalAmount  string             `json:"total_amount"`
	TCS          string             `json:"tcs"`
	Installments []InstallmentPoint `json:"installments"`
}

type InstallmentPoint struct {
	Name    string `json:"name"`
	Amount  string `json:"amount"`
	DueDate string `json:"due_date"`
}

type VisaDetails struct {
	Type           string `json:"type"`
	Validity       string `json:"validity"`
	ProcessingDate string `json:"processing_date"`
}

type ContactInfo struct {
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
}

type ItineraryResponse struct {
	Message  string `json:"message"`
	FilePath string `json:"file_path"`
}

func main() {
	r := gin.Default()
	r.Static("/pdfs", "./pdfs")
	r.POST("/generate-itinerary", generateItinerary)
	r.Run(":8080")
}

func generateItinerary(c *gin.Context) {
	var req ItineraryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	pdfPath, err := CreatePDF(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	c.JSON(http.StatusOK, ItineraryResponse{
		Message:  "Itinerary generated successfully",
		FilePath: fmt.Sprintf("http://localhost:8080/%s", pdfPath),
	})
}
