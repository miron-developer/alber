package app

import (
	"errors"
	"net/http"
	"strconv"
	"wnet/pkg/app"
	"wnet/pkg/orm"
)

func Comments(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ID, e := strconv.Atoi(r.FormValue("id"))
	if e != nil {
		return nil, errors.New("wrong id")
	}

	userID := app.GetUserIDfromReq(w, r)
	if userID == -1 {
		return nil, errors.New("not logged")
	}

	commentType := r.FormValue("type")
	if commentType != "post" && commentType != "comment" && commentType != "media" {
		return nil, errors.New("wrong comment type")
	}

	op := orm.DoSQLOption("c.id = ?", "", "", ID)
	if r.FormValue("count") != "single" {
		first, count := getLimits(r)
		op = orm.DoSQLOption("c."+commentType+"ID = ?", "c.datetime DESC", "?,?", ID, first, count)
	}

	CQ := app.CarmaCountQ("commentID", ID)
	UJ := orm.DoSQLJoin(orm.LOJOINQ, "Users AS u", "u.id = c.userID")
	LJ := orm.DoSQLJoin(orm.LOJOINQ, "Likes AS l", "l.userID = ? AND l.commentID = c.id", userID)

	mainQ := orm.SQLSelectParams{
		Table:   "Comments as c",
		What:    "c.*, u.nName, u.ava, u.status, l.id IS NOT NULL",
		Options: op,
		Joins:   []orm.SQLJoin{UJ, LJ},
	}

	return orm.GetWithSubqueries(
		mainQ,
		[]orm.SQLSelectParams{CQ},
		[]string{"nickname", "avatar", "status", "isLiked"},
		[]string{"carma"},
		orm.Comment{},
	)
}
