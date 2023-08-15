// models/rideAssignment.go

package models

type RideAssignment struct {
  RideGiver   Hiker   `json:"rideGiver"`
  RideGetters []Hiker `json:"rideGetters"`
}