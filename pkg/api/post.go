package app

import (
	"errors"
	"net/http"
	"strconv"
	"wnet/pkg/app"
	"wnet/pkg/orm"
)

func GetPosts(op orm.SQLOption, userID int, lastOpArgs ...interface{}) ([]map[string]interface{}, error) {
	op.Where += ` AND (
		SELECT 
			userID == ? OR
			CASE postType
				WHEN "public"
					THEN 1
				WHEN "private"
					THEN (
						` + PRIVATE_ACCESS_CHECK + `
					)
				WHEN "almost_private"
					THEN instr(allowedUsers, '|' || ? || ' ')
			END
		FROM Posts WHERE id = p.id
	)`
	op.Args = append(op.Args, userID, userID, userID, userID, userID, lastOpArgs)

	carmaQ := app.CarmaCountQ("postID", "p.id")
	UJ := orm.DoSQLJoin(orm.LOJOINQ, "Users AS u", "u.id = p.userID")
	GJ := orm.DoSQLJoin(orm.LOJOINQ, "Groups AS g", "g.id = p.groupID")
	LJ := orm.DoSQLJoin(orm.LOJOINQ, "Likes AS l", "l.userID = ? AND l.postID = p.id", userID)
	mainQ := orm.SQLSelectParams{
		Table:   "Posts as p",
		What:    "p.*, u.nName, u.ava, u.status, g.title, g.ava, l.id IS NOT NULL",
		Options: op,
		Joins:   []orm.SQLJoin{UJ, GJ, LJ},
	}

	return orm.GetWithSubqueries(
		mainQ,
		[]orm.SQLSelectParams{carmaQ},
		[]string{"nickname", "userAvatar", "status", "groupTitle", "groupAvatar", "isLiked"},
		[]string{"carma"},
		orm.Post{},
	)
}

// SELECT
//     p.*, u.nName, u.ava, u.status, g.title, g.ava, l.id IS NOT NULL,
//     (SELECT COUNT(id) FROM Likes WHERE postID = p.id) AS carma
// FROM Posts AS p
// LEFT OUTER JOIN Users AS u ON u.id = p.userID
// LEFT OUTER JOIN Groups AS g ON g.id = p.groupID
// LEFT OUTER JOIN Likes AS l ON l.postID = p.id AND l.userID = ?
// WHERE (
//     SELECT
//         userID == > OR
//         CASE postType
//             WHEN "public"
//                 THEN 1
//             WHEN "private"
//                 THEN (
//                     SELECT id IS NOT NULL FROM Relations
//                     WHERE (
//                         (userID IS NOT NULL AND ((senderUserID = ? AND receiverUserID = userID) OR (receiverUserID = ? AND senderUserID = userID AND value = 0))) OR
//                         (groupID IS NOT NULL AND (senderUserID = ? AND receiverGroupID = groupID))
//                     )
//                 )
//             WHEN "almost_private"
//                 THEN instr(allowedUsers, '|' || ? || ' ')
//         END
//     FROM Posts WHERE id = p.id
// )

func Post(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ID, e := strconv.Atoi(r.FormValue("id"))
	if e != nil {
		return nil, errors.New("Wrong id")
	}
	userID := app.GetUserIDfromReq(w, r)
	if userID == -1 {
		return nil, errors.New("Not logged")
	}
	return GetPosts(orm.DoSQLOption("p.id = ?", "", "", ID), userID)
}

func searchPost(userID int, r *http.Request) (interface{}, error) {
	op := orm.DoSQLOption("", "p.datetime DESC", "?,?")

	searchGetCountFilter(" carma >=", r.FormValue("carmamin"), 0, &op)
	searchGetCountFilter(" carma <=", r.FormValue("carmamax"), 100, &op)

	searchGetBtnFilter(
		r.FormValue("sort"),
		[][]string{
			{"pop", "carma"},
		},
		true,
		&op,
	)

	if e := searchGetTextFilter(r.FormValue("q"), []string{"p.title", "p.body"}, &op); e != nil {
		return nil, e
	}
	first, count := getLimits(r)
	return GetPosts(op, userID, first, count)
}
