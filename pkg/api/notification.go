package app

import (
	"errors"
	"net/http"
	"wnet/pkg/app"
	"wnet/pkg/orm"
)

func Notifications(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ID, e := app.GetUserID(w, r, r.FormValue("id"))
	if e != nil {
		return nil, errors.New("Wrong id")
	}
	first, count := getLimits(r)

	UJ := orm.DoSQLJoin(orm.LOJOINQ, "Users AS u", "u.id = n.senderUserID")
	GJ := orm.DoSQLJoin(orm.LOJOINQ, "Groups AS g", "g.id = n.groupID")
	PJ := orm.DoSQLJoin(orm.LOJOINQ, "Posts AS p", "p.id = n.postID")
	EJ := orm.DoSQLJoin(orm.LOJOINQ, "Events AS e", "e.id = n.eventID")
	CJ := orm.DoSQLJoin(orm.LOJOINQ, "Comments AS c", "c.id = n.commentID")
	MJ := orm.DoSQLJoin(orm.LOJOINQ, "Media AS m", "m.id = n.mediaID")

	mainQ := orm.SQLSelectParams{
		Table: "Notifications AS n",
		What:  "n.*, u.nName, g.title, p.title, e.title, c.body, m.title",
		Options: orm.DoSQLOption(
			`n.receiverUserID = ? OR
			n.senderUserID IN(`+FOLLOWING_USER_ONE_Q+`) OR
			n.senderUserID IN(`+FOLLOWING_USER_BOTH_Q+`)`,
			"n.datetime DESC",
			"?,?",
			ID, ID, ID, first, count,
		),
		Joins: []orm.SQLJoin{UJ, GJ, PJ, EJ, CJ, MJ},
	}

	return orm.GetWithSubqueries(mainQ, nil, []string{"nickname", "groupTitle", "postTitle", "eventTitle", "commentBody", "mediaTitle"}, nil, orm.Notification{})
}

// SELECT
//     n.*, u.nName, g.title, p.title, e.title, c.body, m.title
// FROM Notifications AS n
// LEFT JOIN Users AS u ON u.id = n.senderUserID
// LEFT JOIN Groups AS g ON g.id = n.groupID
// LEFT JOIN Posts AS p ON p.id = n.postID
// LEFT JOIN Events AS e ON e.id = n.eventID
// LEFT JOIN Comments AS c ON c.id = c.commentID
// LEFT JOIN Media AS m ON m.id = n.mediaID
// WHERE
//     n.receiverUserID = ? OR
//     n.senderUserID IN (
//         SELECT receiverUserID FROM Relations WHERE senderUserID = ? AND value != -1
//     ) OR
//     n.senderUserID IN (
//         SELECT senderUserID FROM Relations WHERE receiverUserID = ? AND value = 0
//     )
// ORDER BY n.datetime DESC
// LIMIT ?,?
