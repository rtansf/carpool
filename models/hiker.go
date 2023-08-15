// models/hiker.go

package models

type Hiker struct {
  ID               uint   `json:"id" gorm:"primary_key"`
  FirstName        string `json:"firstName"`
  LastName         string `json:"lastName"`
  Email            string `json:"email"`
  Phone            string `json:"phone"`
  Street           string `json:"street"`
  City             string `json:"city"`
  State            string `json:"state"`
  Zip              string `json:"zip"`
  Consent          string `json:"consent"`
  RideType         string `json:"rideType"`  // give, need, own
  Location         string `json:"location"`  // WholeFoods, Rockridge
  NumPassangers    uint   `json:"numPassangers"`
}