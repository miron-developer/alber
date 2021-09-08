package app

import (
	"errors"
	"net/http"
	"strconv"
	"wnet/pkg/orm"
)

func GetUsers(op orm.SQLOption, userID int, lastOpArgs ...interface{}) (interface{}, error) {
	FLWRSQ := orm.SQLSelectParams{
		Table: "Users",
		What:  "COUNT(id)",
		Options: orm.DoSQLOption(
			"id IN("+FOLLOWERS_USER_ONE_Q+") OR id IN("+FOLLOWERS_USER_BOTH_Q+")",
			"",
			"",
			"u.id", "u.id",
		),
	}
	FLWINGQ := orm.SQLSelectParams{
		Table: "Users",
		What:  "COUNT(id)",
		Options: orm.DoSQLOption(
			"id IN("+FOLLOWING_USER_ONE_Q+") OR id IN("+FOLLOWING_USER_BOTH_Q+")",
			"",
			"",
			"u.id", "u.id",
		),
	}
	EQ := orm.SQLSelectParams{
		Table:   "Events",
		What:    "COUNT(id)",
		Options: orm.DoSQLOption("userID = u.id", "", ""),
	}
	MQ := orm.SQLSelectParams{
		Table:   "Media",
		What:    "COUNT(id)",
		Options: orm.DoSQLOption("userID = u.id", "", ""),
	}
	GQ := orm.SQLSelectParams{
		Table:   "Relations",
		What:    "COUNT(id)",
		Options: orm.DoSQLOption("(senderUserID = u.id AND receiverGroupID IS NOT NULL) OR (receiverUserID = u.id AND senderGroupID IS NOT NULL AND value = 0)", "", ""),
	}
	IRQ := orm.SQLSelectParams{
		Table:   "Relations",
		What:    "value",
		Options: orm.DoSQLOption("senderUserID = u.id AND receiverUserID = ?", "", "", userID),
	}
	ORQ := orm.SQLSelectParams{
		Table:   "Relations",
		What:    "value",
		Options: orm.DoSQLOption("senderUserID = ? AND receiverUserID = u.id", "", "", userID),
	}

	op.Args = append(op.Args, lastOpArgs...)
	mainQ := orm.SQLSelectParams{
		Table:   "Users AS u",
		What:    "u.*",
		Options: op,
	}

	return orm.GetWithSubqueries(
		mainQ,
		[]orm.SQLSelectParams{FLWRSQ, FLWINGQ, EQ, GQ, MQ, IRQ, ORQ},
		[]string{},
		[]string{"followersCount", "followingCount", "eventsCount", "groupsCount", "galleryCount", "InRlshState", "OutRlshState"},
		orm.User{},
	)
}

// SELECT
//     u.*,
//     (SELECT COUNT(id) FROM Users WHERE id IN (SELECT senderUserID FROM Relations WHERE receiverUserID = u.id AND value != -1) OR
//         id IN (SELECT receiverUserID FROM Relations WHERE senderUserID = u.id AND value = 0)) AS flwrs,
//     (SELECT COUNT(id) FROM Users WHERE id IN (SELECT receiverUserID FROM Relations WHERE senderUserID = u.id AND value != -1) OR
//         id IN (SELECT senderUserID FROM Relations WHERE receiverUserID = u.id AND value = 0)) AS flwing,
//     (SELECT COUNT(id) FROM Events WHERE userID = u.id) AS events,
//     (SELECT COUNT(id) FROM Relations WHERE (senderUserID = u.id AND receiverGroupID IS NOT NULL) OR
//         (receiverUserID = u.id AND senderGroupID IS NOT NULL AND value = 0)) AS groups,
//     (SELECT COUNT(id) FROM Media WHERE userID = u.id) AS gallery,
//     (SELECT value FROM Relations WHERE senderUserID = u.id AND receiverUserID = ?) AS inState,
//     (SELECT value FROM Relations WHERE senderUserID = ? AND receiverUserID = u.id) AS outState
// FROM Users AS u

func Users(usersType, flwType string, userID, first, count int) (interface{}, error) {
	op := orm.DoSQLOption(
		"",
		"u.nName ASC",
		"?,?",
		userID, userID, first, count,
	)

	if usersType == "followers" {
		if flwType == "all" {
			op.Where = "u.id IN(" + FOLLOWERS_USER_ONE_Q + ") OR u.id IN(" + FOLLOWERS_USER_BOTH_Q + ")"
		} else if flwType == "online" {
			op.Where = "(u.id IN(" + FOLLOWERS_USER_ONE_Q + ") OR u.id IN(" + FOLLOWERS_USER_BOTH_Q + ")) AND u.status='online'"
		} else {
			if flwType == "request_g" {
				op.Where = "u.id IN(" + REQUEST_GROUP_USER_Q + ")"
			} else {
				op.Where = "u.id IN(" + REQUEST_USER_USER_Q + ")"
			}
			op.Args = op.Args[1:]
		}
	} else if usersType == "following" {
		op.Where = "u.id IN(" + FOLLOWING_USER_ONE_Q + ") OR u.id IN(" + FOLLOWING_USER_BOTH_Q + ")"
	} else {
		op.Where = "u.id IN(" + FOLLOWERS_GROUP_ONE_Q + ") OR u.id IN (" + FOLLOWERS_GROUP_BOTH_Q + ")"
	}

	return GetUsers(op, userID, first, count)
}

func User(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ID, e := strconv.Atoi(r.FormValue("id"))
	if e != nil {
		return nil, errors.New("Wrong id")
	}

	return GetUsers(orm.DoSQLOption("u.id = ?", "", "", ID), ID)
}

func searchUser(userID int, r *http.Request) (interface{}, error) {
	op := orm.DoSQLOption("", "u.status DESC", "?,?")

	searchGetCountFilter(" u.age >=", r.FormValue("agemin"), 0, &op)
	searchGetCountFilter(" u.age <=", r.FormValue("agemax"), 100, &op)
	searchGetCountFilter(" followersCount >=", r.FormValue("subsmin"), 0, &op)
	searchGetCountFilter(" followersCount <=", r.FormValue("subsmax"), 0, &op)

	searchGetBtnFilter(
		r.FormValue("sort"),
		[][]string{
			{"subs", "followersCount"},
		},
		true,
		&op,
	)
	searchGetBtnFilter(
		r.FormValue("gender"),
		[][]string{
			{"Male", " u.gender ="},
			{"Female", " u.gender ="},
		},
		false,
		&op,
	)

	searchGetSwitchFilter(" u.status =", r.FormValue("online"), "online", "", &op)

	if e := searchGetTextFilter(r.FormValue("q"), []string{"u.fName", "u.lName", "u.nName"}, &op); e != nil {
		return nil, e
	}
	first, count := getLimits(r)
	return GetUsers(op, userID, first, count)
}
