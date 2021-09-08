package app

import (
	"errors"
	"net/http"
	"strconv"
	"wnet/pkg/app"
	"wnet/pkg/orm"
)

func GetEvents(op orm.SQLOption, userID int, lastOpArgs ...interface{}) ([]map[string]interface{}, error) {
	op.Where += ` AND (
		SELECT
			userID == ? OR (
				` + PRIVATE_ACCESS_CHECK + `
			)
		FROM Events WHERE id = e.id
	)`
	op.Args = append(op.Args, userID, userID, userID, userID, lastOpArgs)

	UJ := orm.DoSQLJoin(orm.LOJOINQ, "Users AS u", "u.id = e.userID")
	GJ := orm.DoSQLJoin(orm.LOJOINQ, "Groups AS g", "g.id = e.groupID")
	EAJ := orm.DoSQLJoin(orm.LOJOINQ, "EventAnswers AS ea", "ea.userID = ? AND ea.eventID = e.id", userID)
	eventAnswersGoingQ := app.EventAnswerQ(0, "e.id")
	eventAnswersNotGoingQ := app.EventAnswerQ(1, "e.id")
	eventAnswersIDKQ := app.EventAnswerQ(2, "e.id")

	mainQ := orm.SQLSelectParams{
		Table:   "Events as e",
		What:    "e.*, u.nName, u.ava, u.status, g.title, g.ava, ea.answer",
		Options: op,
		Joins:   []orm.SQLJoin{UJ, GJ, EAJ},
	}

	return orm.GetWithSubqueries(
		mainQ,
		[]orm.SQLSelectParams{eventAnswersGoingQ, eventAnswersNotGoingQ, eventAnswersIDKQ},
		[]string{"nickname", "userAvatar", "status", "groupTitle", "groupAvatar", "myVote"},
		[]string{"votes0", "votes1", "votes2"},
		orm.Event{},
	)
}

// SELECT
//     e.*, u.nName, u.ava, u.status, g.title, g.ava, ea.answer,
//     (SELECT COUNT(id) FROM EventAnswers WHERE answer = 0 AND eventID = e.id),
//     (SELECT COUNT(id) FROM EventAnswers WHERE answer = 1 AND eventID = e.id),
//     (SELECT COUNT(id) FROM EventAnswers WHERE answer = 2 AND eventID = e.id)
// FROM Events AS e
// LEFT OUTER JOIN Users AS u ON u.id = e.userID
// LEFT OUTER JOIN Groups AS g ON g.id = e.groupID
// LEFT OUTER JOIN EventAnswers AS ea ON ea.userID = ? AND ea.eventID = e.id
// WHERE (
//     SELECT
// 		userID == ? OR
// 		(
// 			SELECT id IS NOT NULL FROM Relations
// 			WHERE (
// 				(userID IS NOT NULL AND ((senderUserID = ? AND receiverUserID = userID) OR (receiverUserID = ? AND senderUserID = userID AND value = 0))) OR
// 				(groupID IS NOT NULL AND (senderUserID = ? AND receiverGroupID = groupID))
// 			)
// 		)
// 	FROM Events WHERE id = e.id
// )

func Events(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ID, e := app.GetUserID(w, r, r.FormValue("id"))
	if e != nil {
		return nil, errors.New("Wrong id")
	}

	which := r.FormValue("which")
	if which != "user" && which != "group" {
		return nil, errors.New("Wrong event type")
	}

	first, count := getLimits(r)
	return GetEvents(orm.DoSQLOption("e."+which+"ID = ?", "e.datetime DESC", "?,?", ID), ID, first, count)
}

func Event(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ID, e := strconv.Atoi(r.FormValue("id"))
	if e != nil {
		return nil, errors.New("Wrong id")
	}
	userID := app.GetUserIDfromReq(w, r)
	if userID == -1 {
		return nil, errors.New("Not logged")
	}
	return GetEvents(orm.DoSQLOption("e.id = ?", "", "", ID), userID)
}
