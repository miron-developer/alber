package orm

// User - table
type User struct {
	ID          int    `json:"id"`
	Nickname    string `json:"nickname"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string
}

// Session - table
type Session struct {
	ID     string
	Expire string
	UserID int
}

// Country - table
type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// City - table
type City struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TravelType - table
type TravelType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TopType - table
type TopType struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	Duration string `json:"duration"`
	Cost     string `json:"cost"`
}

// CountryCode - table
type CountryCode struct {
	ID        int    `json:"id"`
	Code      string `json:"code"`
	CountryID int    `json:"countryID"`
}

// Parsel - table
type Parsel struct {
	ID                  int    `json:"id"`
	Title               string `json:"title"`
	Weight              int    `json:"weight"`
	Price               int    `json:"price"`
	ContactNumber       string `json:"contactNumber"`
	CreationDatetime    int    `json:"creationDatetime"`
	ExpireDatetime      int    `json:"expireDatetime"`
	ExpireOnTopDatetime int    `json:"expireOnTopDatetime"`
	IsHaveWhatsUp       string `json:"isHaveWhatsUp"`
	UserID              int    `json:"userID"`
	TopTypeID           int    `json:"topTypeID"`
	FromID              int    `json:"fromID"`
	ToID                int    `json:"toID"`
}

// Traveler - table
type Traveler struct {
	ID                  int    `json:"id"`
	ContactNumber       string `json:"contactNumber"`
	Weight              int    `json:"weight"`
	CreationDatetime    int    `json:"creationDatetime"`
	DepartureDatetime   int    `json:"departureDatetime"`
	ArrivalDatetime     int    `json:"arrivalDatetime"`
	ExpireOnTopDatetime int    `json:"expireOnTopDatetime"`
	IsHaveWhatsUp       string `json:"isHaveWhatsUp"`
	UserID              int    `json:"userID"`
	TopTypeID           int    `json:"topTypeID"`
	TravelTypeID        int    `json:"travelTypeID"`
	FromID              int    `json:"fromID"`
	ToID                int    `json:"toID"`
}

// Image - table
type Image struct {
	ID       int    `json:"id"`
	Source   string `json:"src"`
	Name     string `json:"filename"`
	UserID   int    `json:"userID"`
	ParselID int    `json:"parselID"`
}
