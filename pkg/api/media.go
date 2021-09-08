package app

import (
	"errors"
	"net/http"
	"strconv"
	"wnet/pkg/app"
	"wnet/pkg/orm"
)

func GetMedia(op orm.SQLOption, userID int, lastOpArgs ...interface{}) ([]map[string]interface{}, error) {
	op.Where += ` AND (
		SELECT
			userID == ? OR (
				` + PRIVATE_ACCESS_CHECK + `
			)
		FROM Media WHERE id = m.id
	)`
	op.Args = append(op.Args, userID, userID, userID, userID, lastOpArgs)

	carmaQ := app.CarmaCountQ("mediaID", "m.id")
	UJ := orm.DoSQLJoin(orm.LOJOINQ, "Users AS u", "u.id = m.userID")
	GJ := orm.DoSQLJoin(orm.LOJOINQ, "Groups AS g", "g.id = m.groupID")
	LJ := orm.DoSQLJoin(orm.LOJOINQ, "Likes AS l", "l.userID = ? AND l.mediaID = m.id", userID)
	mainQ := orm.SQLSelectParams{
		Table:   "Media as m",
		What:    "m.*, u.nName, u.ava, u.status, g.title, g.ava, l.id IS NOT NULL",
		Options: op,
		Joins:   []orm.SQLJoin{UJ, GJ, LJ},
	}

	return orm.GetWithSubqueries(
		mainQ,
		[]orm.SQLSelectParams{carmaQ},
		[]string{"nickname", "userAvatar", "status", "groupTitle", "groupAvatar", "isLiked"},
		[]string{"carma"},
		orm.Media{},
	)
}

// SELECT
//     m.*, u.nName, u.ava, u.status, g.title, g.ava, l.id IS NOT NULL,
//     (SELECT COUNT(id) FROM Likes WHERE mediaID = m.id) AS carma
// FROM Media AS m
// LEFT JOIN Users AS u ON u.id = m.userID
// LEFT JOIN Groups AS g ON g.id = m.groupID
// LEFT JOIN Likes AS l ON l.mediaID = m.id AND l.userID = ?
// WHERE (
//     SELECT
// 		userID == ? OR (
// 			SELECT id IS NOT NULL FROM Relations
//             WHERE (
//                 (
//                     userID IS NOT NULL AND (
//                         (senderUserID = ? AND receiverUserID = userID) OR
//                         (receiverUserID = ? AND senderUserID = userID AND value = 0)
//                     )
//                 ) OR
//                 (groupID IS NOT NULL AND (senderUserID = ? AND receiverGroupID = groupID))
//             )
// 		)
// 	FROM Media WHERE id = m.id
// )

func Gallery(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ID, e := app.GetUserID(w, r, r.FormValue("ownerID"))
	if e != nil {
		return nil, errors.New("Wrong id")
	}

	first, count := getLimits(r)

	op := orm.DoSQLOption("", "m.datetime DESC", "?,?")
	if r.FormValue("type") == "user" {
		op.Where = "m.userID = ?"
	} else {
		op.Where = "m.groupID = ?"
	}

	galleryType := r.FormValue("galleryType")
	if galleryType != "all" && galleryType != "photo" && galleryType != "video" {
		return nil, errors.New("Wrong gallery type")
	}
	if galleryType == "all" {
		op.Where += ` AND (m.type="photo" OR m.type=?)`
		galleryType = "video"
	} else {
		op.Where += " AND m.type=?"
	}

	op.Args = append(op.Args, ID, galleryType)
	return GetMedia(op, ID, first, count)
}

func Media(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ID, e := strconv.Atoi(r.FormValue("id"))
	if e != nil {
		return nil, errors.New("Wrong id")
	}
	userID := app.GetUserIDfromReq(w, r)
	if userID == -1 {
		return nil, errors.New("Not logged")
	}
	mediaType := r.FormValue("type")
	if mediaType != "photo" && mediaType != "video" {
		return nil, errors.New("Wrong media type")
	}
	return GetMedia(orm.DoSQLOption("m.id = ? AND m.type = ?", "", "", ID, mediaType), userID)
}

func searchVideo(userID int, r *http.Request) (interface{}, error) {
	op := orm.DoSQLOption("m.type = 'video' AND ", "m.datetime DESC", "?,?")

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

	if e := searchGetTextFilter(r.FormValue("q"), []string{"m.title"}, &op); e != nil {
		return nil, e
	}

	first, count := getLimits(r)
	return GetMedia(op, userID, first, count)
}
