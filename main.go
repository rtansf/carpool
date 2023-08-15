package main


import (
  "fmt"
  "net/http"
  "github.com/gin-gonic/gin"
  "example.com/m/models"
  "example.com/m/persist"
  "log"
  "encoding/json"
)

type CreateHikerInput struct {
  FirstName  string `json:"firstName"`
  LastName   string `json:"lastName"`
  Email      string `json:"email"`
  Phone      string `json:"phone"`
  Street     string `json:"street"` 
  City       string `json:"city"`
  State      string `json:"state"` 
  Zip        string `json:"zip"`
  Consent    string `json:"consent"`
  RideType   string `json:"rideType"`
  Location   string `json:"location"`
  NumPassangers uint `json:"numPassangers"`
}

// POST /hiker
// Create new hiker
func CreateHiker(c *gin.Context) {
  // Validate input
  var input CreateHikerInput
  if err := c.ShouldBindJSON(&input); err != nil {
    fmt.Println("FAILED TO BIND")
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  // Create hiker
  hiker := models.Hiker{
     FirstName: input.FirstName,
     LastName:  input.LastName,
     Email:     input.Email,
     Phone:     input.Phone,
     Street:    input.Street,
     City:      input.City,
     State:     input.State,
     Zip:       input.Zip,
     Consent:   input.Consent,
     RideType:  input.RideType,
     Location:  input.Location,
     NumPassangers: input.NumPassangers,
  }
  models.DB.Create(&hiker)

  c.JSON(http.StatusOK, gin.H{"data": hiker})
}

// GET /hikers
// Get all hikers
func FindHikers(c *gin.Context) {
  var hikers []models.Hiker
  models.DB.Find(&hikers)

  c.JSON(http.StatusOK, gin.H{"data": hikers})
}

// GET /rideassignments
// Return saved ride assignments
func GetRideAssignments(c *gin.Context) {
   var rideAssignments [2]models.AssignmentLocation
   if err := persist.Load("./assignmentLocations.tmp", &rideAssignments); err != nil {
       rideAssignments = PerformInitialRideAssignments()

       // Save assignment locations
       if err := persist.Save("./assignmentLocations.tmp", &rideAssignments); err != nil {
           log.Fatalln(err)
       }
   }
   c.JSON(http.StatusOK, gin.H{"data": rideAssignments})
}

func PerformInitialRideAssignments() [2]models.AssignmentLocation {
  var hikers []models.Hiker
  models.DB.Find(&hikers)

  // --------------------------------------------

  rideGiverMapWF := make(map[string]models.Hiker)
  rideGetterMapWF := make(map[string]models.Hiker)
  rideGiverMapRR := make(map[string]models.Hiker)
  rideGetterMapRR := make(map[string]models.Hiker)
  
  for i := 0; i < len(hikers); i++ {
      hiker := hikers[i]
      if hiker.Location == "Whole Foods" {
         if hiker.RideType == "give" {
            _, ok := rideGiverMapWF[hiker.Email]
   	    if !ok {
	        rideGiverMapWF[hiker.Email] = hiker
	    }
	 } else if hiker.RideType == "need" {
            _, ok := rideGetterMapWF[hiker.Email]
   	    if !ok {
	        rideGetterMapWF[hiker.Email] = hiker
	    }
	 }
      } else if hiker.Location == "Rockridge" {
         if hiker.RideType == "give" {
            _, ok := rideGiverMapRR[hiker.Email]
   	    if !ok {
	        rideGiverMapRR[hiker.Email] = hiker
	    }
	 } else if hiker.RideType == "need" {
            _, ok := rideGetterMapRR[hiker.Email]
   	    if !ok {
	        rideGetterMapRR[hiker.Email] = hiker
	    }
	 }
      }
  }

  fmt.Println("b", rideGetterMapWF)
  fmt.Println("c", rideGiverMapRR)
  fmt.Println("d", rideGetterMapRR)

  var rideAssignments []models.RideAssignment
  for _, element := range rideGiverMapWF {
     rideAssignment := models.RideAssignment {
          RideGiver : element,
	  RideGetters : []models.Hiker {
	  },
     }
     rideAssignments = append(rideAssignments, rideAssignment)     
  }

  getters := []models.Hiker {}
  for _, element := range rideGetterMapWF {
     getters = append(getters, element)
  }

  orphanGetters := []models.Hiker {}

  j := 0
  for j < len(getters) {
     getter := getters[j]
     fmt.Println ("-->", getter.Email)
     added := false
     for i := 0; i < len(rideAssignments); i++ {
	  giver := rideAssignments[i].RideGiver
          if len(rideAssignments[i].RideGetters) < int(giver.NumPassangers) {
	     fmt.Println ("Adding", getter.Email, "to", giver.Email)
	     rideAssignments[i].RideGetters = append (rideAssignments[i].RideGetters, getter)
	     added = true
	     break
          }
      }
      if !added {
          orphanGetters = append (orphanGetters, getter)
      }
      j = j + 1
  }

  fmt.Println ("oprhans", orphanGetters)
  
  assignmentLocationWF := models.AssignmentLocation {
      Location : "Whole Foods",
      RideAssignments : rideAssignments,
      OrphanRiders : orphanGetters,
  }

  assignmentLocationRR := models.AssignmentLocation {
      Location : "Rockridge",
      RideAssignments : []models.RideAssignment {},
      OrphanRiders : []models.Hiker {},
  }

  var assignmentLocations [2]models.AssignmentLocation
  assignmentLocations[0] = assignmentLocationWF
  assignmentLocations[1] = assignmentLocationRR

  return assignmentLocations
}

// GET /initialrideassignments
// Perform and return initial ride assignments
func InitialRideAssignments(c *gin.Context) {
  assignmentLocations := PerformInitialRideAssignments() 

  // Save assignment locations
  if err := persist.Save("./assignmentLocations.tmp", &assignmentLocations); err != nil {
     log.Fatalln(err)
  }
  c.JSON(http.StatusOK, gin.H{"data": assignmentLocations})
}

// POST /rideassignments
// Save Ride Assignments
func SaveRideAssignments(c *gin.Context) {
   jsonData, _ := c.GetRawData()

   var locations = []models.AssignmentLocation{} 
   json.Unmarshal(jsonData, &locations)
   fmt.Println("Saving: ", locations)

   // Save assignment locations
   if err := persist.Save("./assignmentLocations.tmp", &locations); err != nil {
     log.Fatalln(err)
   }
   
   c.JSON(http.StatusOK, gin.H{"data": ""})
}

func main() {

  r := gin.Default()
  r.LoadHTMLFiles("html/index.html", "html/managecarpool.html")

  r.GET("/", func(c *gin.Context) {
     c.HTML(http.StatusOK, "index.html", gin.H{})
  })

  r.GET("/managecarpool", func(c *gin.Context) {
     c.HTML(http.StatusOK, "managecarpool.html", gin.H{})
  })

  r.GET("/hikers", FindHikers)
  r.POST("/hiker", CreateHiker)
  r.GET("/initialrideassignments", InitialRideAssignments)
  r.GET("/rideassignments", GetRideAssignments)
  r.POST("/rideassignments", SaveRideAssignments)
  
  models.ConnectDatabase()
 
  r.Run()
}

