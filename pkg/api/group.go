package app

import (
	"errors"
	"net/http"
	"strconv"
	"wnet/pkg/app"
	"wnet/pkg/orm"
)

func GetGroups(op orm.SQLOption, userID int, lastOpArgs ...interface{}) (interface{}, error) {
	MBQ := orm.SQLSelectParams{
		Table: "Users",
		What:  "COUNT(id)",
		Options: orm.DoSQLOption(
			"id IN ("+FOLLOWERS_GROUP_ONE_Q+") OR id IN ("+FOLLOWERS_GROUP_BOTH_Q+")",
			"",
			"",
			"g.id", "g,id",
		),
	}
	RQQ := orm.SQLSelectParams{
		Table: "Users",
		What:  "COUNT(id)",
		Options: orm.DoSQLOption(
			"id IN ("+REQUEST_GROUP_USER_Q+")",
			"",
			"",
			"g.id",
		),
	}
	EQ := orm.SQLSelectParams{
		Table:   "Events",
		What:    "COUNT(id)",
		Options: orm.DoSQLOption("groupID = g.id", "", ""),
	}
	MQ := orm.SQLSelectParams{
		Table:   "Media",
		What:    "COUNT(id)",
		Options: orm.DoSQLOption("groupID = g.id", "", ""),
	}
	IRQ := orm.SQLSelectParams{
		Table:   "Relations",
		What:    "value",
		Options: orm.DoSQLOption("senderGroupID = g.id AND receiverUserID = ?", "", "", userID),
	}
	ORQ := orm.SQLSelectParams{
		Table:   "Relations",
		What:    "value",
		Options: orm.DoSQLOption("senderUserID = ? AND receiverGroupID = g.id", "", "", userID),
	}

	op.Args = append(op.Args, lastOpArgs...)
	mainQ := orm.SQLSelectParams{
		Table:   "Groups AS g",
		What:    "g.*",
		Options: op,
	}

	return orm.GetWithSubqueries(
		mainQ,
		[]orm.SQLSelectParams{MBQ, EQ, MQ, RQQ, IRQ, ORQ},
		[]string{},
		[]string{"membersCount", "eventsCount", "galleryCount", "requestsCount", "InRlshState", "OutRlshState"},
		orm.Group{},
	)
}

// SELECT
//     g.*,
//     (SELECT COUNT(id) FROM Users WHERE id IN (SELECT senderUserID FROM Relations WHERE receiverGroupID = g.id AND value != -1) OR
//         id IN (SELECT receiverUserID FROM Relations WHERE senderGroupID = g.id AND value = 0)) AS memb,
//     (SELECT COUNT(id) FROM Users WHERE id IN (SELECT senderUserID FROM Relations WHERE receiverGroupID = g.id AND value = -1)) AS rq,
//     (SELECT COUNT(id) FROM Events WHERE groupID = g.id) AS events,
//     (SELECT COUNT(id) FROM Media WHERE groupID = g.id) AS gallery,
//     (SELECT value FROM Relations WHERE senderGroupID = g.id AND receiverUserID = ?) AS inState,
//     (SELECT value FROM Relations WHERE senderUserID = ? AND receiverGroupID = g.id) AS outState
// FROM Groups AS g

func Groups(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ID, e := app.GetUserID(w, r, r.FormValue("id"))
	if e != nil {
		return nil, errors.New("Wrong id")
	}

	first, count := getLimits(r)
	groupType := r.FormValue("type")
	op := orm.DoSQLOption(
		"",
		"g.title ASC",
		"?,?",
		ID, first, count,
	)
	if groupType == "all" {
		op.Where = "g.id IN(" + FOLLOWERS_USER_GROUP_Q + ")"
	} else {
		op.Where = "g.id IN(" + REQUEST_USER_GROUP_Q + ")"
	}

	return GetGroups(op, ID, first, count)
}

func Group(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ID, e := strconv.Atoi(r.FormValue("id"))
	if e != nil {
		return nil, errors.New("Wrong id")
	}

	return GetGroups(orm.DoSQLOption("g.id=?", "", "", ID), ID)
}

func searchGroup(userID int, r *http.Request) (interface{}, error) {
	op := orm.DoSQLOption("", "g.cdate DESC", "?,?")

	searchGetCountFilter(" membersCount >=", r.FormValue("membmin"), 0, &op)
	searchGetCountFilter(" membersCount <=", r.FormValue("membmax"), 100, &op)

	searchGetBtnFilter(
		r.FormValue("sort"),
		[][]string{
			{"member", "membersCount"},
		},
		true,
		&op,
	)

	searchGetSwitchFilter("g.isPrivate =", r.FormValue("private"), "1", "0", &op)

	if e := searchGetTextFilter(r.FormValue("q"), []string{"g.title", "g.about"}, &op); e != nil {
		return nil, e
	}

	first, count := getLimits(r)
	return GetGroups(op, userID, first, count)
}
