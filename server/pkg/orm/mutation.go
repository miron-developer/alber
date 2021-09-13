package orm

import (
	"errors"
	"strconv"
)

// ---------------------Create funcs---------------------------

// Create create one user
func (u *User) Create() (int, error) {
	if u.Nickname == "" || u.Password == "" || u.PhoneNumber == "" {
		return -1, errors.New("n/d")
	}

	r, e := insertSQL(SQLInsertParams{
		Table:  "Users",
		Datas:  "null,?,?,?",
		Values: MakeArrFromStruct(*u)[1:],
	})
	if e != nil {
		return -1, e
	}
	ID, e := r.LastInsertId()
	return int(ID), e
}

// Create create new session in db
func (ses *Session) Create() error {
	if ses.ID == "" || ses.UserID == 0 || ses.Expire == "" {
		return errors.New("n/d")
	}

	_, e := insertSQL(SQLInsertParams{
		Table:  "Sessions",
		Datas:  "?,?,?",
		Values: MakeArrFromStruct(*ses),
	})
	return e
}

// I think, country & city & travelType & topType is not necessary now

// Create one parsel and return it's ID
func (p *Parsel) Create() (int, error) {
	if p.Title == "" || p.ContactNumber == "" ||
		p.Weight*p.Price*p.CreationDatetime*p.ExpireDatetime*p.UserID*p.FromID*p.ToID == 0 {
		return -1, errors.New("n/d")
	}

	r, e := insertSQL(SQLInsertParams{
		Table:  "Parsels",
		Datas:  "null,?,?,?,?,?,?,?,?,?,?,?,?",
		Values: MakeArrFromStruct(*p)[1:],
	})
	if e != nil {
		return -1, e
	}
	ID, e := r.LastInsertId()
	return int(ID), e
}

// Create one parsel and return it's ID
func (t *Traveler) Create() (int, error) {
	if t.Weight*t.CreationDatetime*t.DepartureDatetime*t.DepartureDatetime*t.ArrivalDatetime*
		t.UserID*t.FromID*t.ToID == 0 || t.ContactNumber == "" {
		return -1, errors.New("n/d")
	}

	r, e := insertSQL(SQLInsertParams{
		Table:  "Travelers",
		Datas:  "null,?,?,?,?,?,?,?,?,?,?,?,?",
		Values: MakeArrFromStruct(*t)[1:],
	})
	if e != nil {
		return -1, e
	}
	ID, e := r.LastInsertId()
	return int(ID), e
}

// Create create one clipped image
func (i *Image) Create() (int, error) {
	if i.UserID*i.ParselID == 0 || i.Source == "" || i.Name == "" {
		return -1, errors.New("n/d")
	}

	r, e := insertSQL(SQLInsertParams{
		Table:  "Images",
		Datas:  "null,?,?,?,?",
		Values: MakeArrFromStruct(*i)[1:],
	})
	if e != nil {
		return -1, e
	}
	ID, e := r.LastInsertId()
	return int(ID), e
}

// ---------------------Change funcs---------------------------

// Change change user profile
func (u *User) Change() error {
	if u.ID == 0 {
		return errors.New("absent/d")
	}

	params := SQLUpdateParams{
		Table:   "Users",
		Couples: map[string]string{},
		Options: DoSQLOption("id=?", "", "", u.ID),
	}

	if u.PhoneNumber != "" {
		params.Couples["phoneNumber"] = u.PhoneNumber
	}
	if u.Nickname != "" {
		params.Couples["nickname"] = u.Nickname
	}
	if u.Password != "" {
		params.Couples["password"] = u.Password
	}
	_, e := updateSQL(params)
	return e
}

// Change change expiration
func (s *Session) Change() error {
	if s.ID == "" || s.Expire == "" {
		return errors.New("absent/d")
	}

	_, e := updateSQL(SQLUpdateParams{
		Table:   "Sessions",
		Couples: map[string]string{"expire": s.Expire},
		Options: DoSQLOption("id=?", "", "", s.ID),
	})
	return e
}

// I think, country & city & travelType & topType is not necessary now

// Change change parsel
func (p *Parsel) Change() error {
	if p.ID == 0 {
		return errors.New("absent/d")
	}

	params := SQLUpdateParams{
		Table:   "Parsels",
		Couples: map[string]string{},
		Options: DoSQLOption("id=?", "", "", p.ID),
	}
	if p.Title != "" {
		params.Couples["title"] = p.Title
	}
	if p.ContactNumber != "" {
		params.Couples["contactNumber"] = p.ContactNumber
	}
	if p.Weight != 0 {
		params.Couples["weight"] = strconv.Itoa(p.Weight)
	}
	if p.Price != 0 {
		params.Couples["price"] = strconv.Itoa(p.Price)
	}
	if p.CreationDatetime != 0 {
		params.Couples["creationDatetime"] = strconv.Itoa(p.CreationDatetime)
	}
	if p.ExpireDatetime != 0 {
		params.Couples["expireDatetime"] = strconv.Itoa(p.ExpireDatetime)
	}
	if p.ExpireOnTopDatetime != 0 {
		params.Couples["expireOnTopDatetime"] = strconv.Itoa(p.ExpireOnTopDatetime)
	}
	if p.ExpireOnTopDatetime == -1 {
		params.Couples["expireOnTopDatetime"] = strconv.Itoa(p.ExpireOnTopDatetime)
	}
	if p.IsHaveWhatsUp != "" {
		params.Couples["isHaveWhatsUp"] = p.IsHaveWhatsUp
	}
	if p.TopTypeID != 0 {
		params.Couples["topTypeID"] = strconv.Itoa(p.TopTypeID)
	}
	if p.TopTypeID == -1 {
		params.Couples["topTypeID"] = "null"
	}
	if p.FromID != 0 {
		params.Couples["fromID"] = strconv.Itoa(p.FromID)
	}
	if p.ToID != 0 {
		params.Couples["toID"] = strconv.Itoa(p.ToID)
	}

	_, e := updateSQL(params)
	return e
}

// Change change post
func (t *Traveler) Change() error {
	if t.ID == 0 {
		return errors.New("absent/d")
	}

	params := SQLUpdateParams{
		Table:   "Travelers",
		Couples: map[string]string{},
		Options: DoSQLOption("id=?", "", "", t.ID),
	}
	if t.ContactNumber != "" {
		params.Couples["contactNumber"] = t.ContactNumber
	}
	if t.Weight != 0 {
		params.Couples["weight"] = strconv.Itoa(t.Weight)
	}
	if t.CreationDatetime != 0 {
		params.Couples["creationDatetime"] = strconv.Itoa(t.CreationDatetime)
	}
	if t.DepartureDatetime != 0 {
		params.Couples["departureDatetime"] = strconv.Itoa(t.DepartureDatetime)
	}
	if t.ArrivalDatetime != 0 {
		params.Couples["arrivalDatetime"] = strconv.Itoa(t.ArrivalDatetime)
	}
	if t.ExpireOnTopDatetime != 0 {
		params.Couples["expireOnTopDatetime"] = strconv.Itoa(t.ExpireOnTopDatetime)
	}
	if t.ExpireOnTopDatetime == -1 {
		params.Couples["expireOnTopDatetime"] = strconv.Itoa(t.ExpireOnTopDatetime)
	}
	if t.IsHaveWhatsUp != "" {
		params.Couples["isHaveWhatsUp"] = t.IsHaveWhatsUp
	}
	if t.TopTypeID != 0 {
		params.Couples["topTypeID"] = strconv.Itoa(t.TopTypeID)
	}
	if t.TopTypeID == -1 {
		params.Couples["topTypeID"] = "null"
	}
	if t.TravelTypeID != 0 {
		params.Couples["travelTypeID"] = strconv.Itoa(t.TravelTypeID)
	}
	if t.FromID != 0 {
		params.Couples["fromID"] = strconv.Itoa(t.FromID)
	}
	if t.ToID != 0 {
		params.Couples["toID"] = strconv.Itoa(t.ToID)
	}

	_, e := updateSQL(params)
	return e
}

// ---------------------Delete funcs---------------------------

// DeleteByParams delete one by id
func DeleteByParams(params SQLDeleteParams) error {
	_, e := deleteSQL(params)
	return e
}
