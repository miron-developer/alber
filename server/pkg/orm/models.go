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
	Description         string `json:"description"`
	Weight              int    `json:"weight"`
	Price               int    `json:"price"`
	ContactNumber       string `json:"contactNumber"`
	CreationDatetime    int    `json:"creationDatetime"`
	ExpireDatetime      int    `json:"expireDatetime"`
	ExpireOnTopDatetime int    `json:"expireOnTopDatetime"`
	IsHaveWhatsUp       string `json:"isHaveWhatsUp"`
	UserID              int    `json:"userID"`
	FromID              int    `json:"fromID"`
	ToID                int    `json:"toID"`
	TopTypeID           int    `json:"topTypeID"`
}

// Traveler - table
type Traveler struct {
	ID                  int    `json:"id"`
	Description         string `json:"description"`
	ContactNumber       string `json:"contactNumber"`
	CreationDatetime    int    `json:"creationDatetime"`
	ExpireOnTopDatetime int    `json:"expireOnTopDatetime"`
	IsHaveWhatsUp       string `json:"isHaveWhatsUp"`
	UserID              int    `json:"userID"`
	TravelTypeID        int    `json:"travelTypeID"`
	FromID              int    `json:"fromID"`
	ToID                int    `json:"toID"`
	TopTypeID           int    `json:"topTypeID"`
}

// Image - table
type Image struct {
	ID       int    `json:"id"`
	Source   string `json:"src"`
	Name     string `json:"filename"`
	UserID   int    `json:"userID"`
	ParselID int    `json:"parselID"`
}
