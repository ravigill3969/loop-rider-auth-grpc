package models

type CapturePaymentRow struct {
	TripID              string
	RiderID             string
	DriverID            string
	PaymentID           string
	StripePaymentIntent string
	StripeAmount        int64
	PaymentStatus       string
}
