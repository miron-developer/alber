package app

import (
	"errors"
	"net/http"
	"strconv"
	"wnet/pkg/app"
	"wnet/pkg/orm"
)

func Messages(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ID, e := app.GetUserID(w, r, "")
	if e != nil {
		return nil, errors.New("wrong id")
	}

	first, count := getLimits(r)
	receiverID := r.FormValue("id")
	chatType := r.FormValue("type")

	op := orm.DoSQLOption(
		"senderUserID=? AND receiverGroupID=?",
		"datetime(datetime) DESC",
		"?,?",
		ID, receiverID,
	)
	if chatType == "user" {
		op.Where = "senderUserID = ? AND receiverUserID = ? OR senderUserID = ? AND receiverUserID = ?"
		op.Args = append(op.Args, receiverID, ID)
	}
	op.Args = append(op.Args, first, count)

	UJ := orm.DoSQLJoin(orm.LOJOINQ, "Users AS u", "m.senderUserID = u.id")
	FJ := orm.DoSQLJoin(orm.LOJOINQ, "Files AS f", "f.messageID = m.id")
	return orm.GetWithSubqueries(
		orm.SQLSelectParams{
			Table:   "Messages AS m",
			What:    "m.*, u.ava, u.nName, u.status, f.src",
			Options: op,
			Joins:   []orm.SQLJoin{UJ, FJ},
		},
		nil,
		[]string{"avatar", "nickname", "status", "src"},
		nil,
		orm.Message{},
	)
}

func Chats(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ID, e := app.GetUserID(w, r, r.FormValue("id"))
	if e != nil {
		return nil, errors.New("wrong id")
	}

	first, count := getLimits(r)
	messageSelectOp := orm.DoSQLOption(
		`(c.senderUserID = m.senderUserID AND c.receiverUserID = m.receiverUserID) OR
		(c.senderUserID = m.receiverUserID AND c.receiverUserID = m.senderUserID) OR
		(c.senderUserID = m.senderUserID AND c.receiverGroupID = m.receiverGroupID)`,
		"m.datetime DESC",
		"?",
		1,
	)
	messageBodyQ := orm.SQLSelectParams{
		Table:   "Messages AS m",
		What:    "m.body",
		Options: messageSelectOp,
	}
	messageDatetimeQ := orm.SQLSelectParams{
		Table:   "Messages AS m",
		What:    "m.datetime",
		Options: messageSelectOp,
	}

	UJ := orm.DoSQLJoin(orm.LOJOINQ, "Users AS u", "c.receiverUserID=u.id")
	GJ := orm.DoSQLJoin(orm.LOJOINQ, "Groups AS g", "c.receiverGroupID=g.id")
	mainQ := orm.SQLSelectParams{
		Table: "Chats c",
		What:  "c.*, u.ava, u.nName, u.status, g.ava, g.title",
		Options: orm.DoSQLOption(
			"(c.users LIKE '%|"+strconv.Itoa(ID)+" %' AND c.closed NOT LIKE '%|"+strconv.Itoa(ID)+" %') ",
			"",
			"?,?",
			first, count,
		),
		Joins: []orm.SQLJoin{UJ, GJ},
	}

	return orm.GetWithSubqueries(
		mainQ,
		[]orm.SQLSelectParams{messageBodyQ, messageDatetimeQ},
		[]string{"userAvatar", "nickname", "status", "groupAvatar", "groupTitle"},
		[]string{"msgBody", "msgDatetime"},
		orm.Chat{},
	)
}
