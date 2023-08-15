// models/assignmentLocation.go

package models

type AssignmentLocation struct {
  Location         string  `json:"location"`
  RideAssignments  []RideAssignment `json:"rideAssignments"`
  OrphanRiders     []Hiker `json:"orphanRiders"`
}

