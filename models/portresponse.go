package models

type PortResponse struct {
	RequestId   string
	PhoneNumber string
	Message     string
	CanPort     bool
}
