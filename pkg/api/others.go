package app

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"wnet/pkg/app"
	"wnet/pkg/orm"
)

func News(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	userID := app.GetUserIDfromReq(w, r)
	if userID == -1 {
		return nil, errors.New("Not logged")
	}

	first, count := getLimits(r)
	// res ex: [[4 post 1588670115000] [3 event 1588670115000]]
	newsIDS, e := orm.GetWithQueryAndArgs(
		`SELECT id, type, datetime FROM Posts 
		WHERE
			userID = ? OR 
			userID IN(`+FOLLOWING_USER_ONE_Q+`) OR
			userID IN(`+FOLLOWING_USER_BOTH_Q+`) OR
			groupID IN(`+FOLLOWING_USER_GROUP_Q+`)
		UNION ALL
		SELECT id, type, datetime FROM Events 
		WHERE 
			userID = ? OR 
			userID IN(`+FOLLOWING_USER_ONE_Q+`) OR
			userID IN(`+FOLLOWING_USER_BOTH_Q+`) OR
			groupID IN(`+FOLLOWING_USER_GROUP_Q+`) 
		ORDER BY datetime DESC LIMIT ?,?`,
		[]interface{}{userID, userID, userID, userID, userID, userID, userID, userID, first, count},
	)
	if e != nil {
		return nil, errors.New("n/d")
	}

	postIDs, eventIDs := []string{}, []string{}
	for _, v := range newsIDS {
		id := "|" + v[1].(string) + " "
		if v[1] == "post" {
			postIDs = append(postIDs, id)
			continue
		}
		eventIDs = append(eventIDs, id)
	}

	posts, e := GetPosts(orm.DoSQLOption("instr(?, '|' || p.id || ' ')", "", "", postIDs), userID)
	if e != nil {
		return nil, e
	}
	evnts, e := GetEvents(orm.DoSQLOption("instr(?, '|' || e.id || ' ')", "", "", eventIDs), userID)
	if e != nil {
		return nil, e
	}
	return append(posts, evnts...), nil
}

func Publications(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ID, e := app.GetUserID(w, r, r.FormValue("id"))
	if e != nil {
		return nil, errors.New("Wrong id")
	}

	publType := r.FormValue("publicationType")
	which := r.FormValue("type")
	if which != "user" && which != "group" {
		return nil, errors.New("Wrong publication type")
	}

	sq := "SELECT id, type, datetime FROM <TABLE> WHERE <WHICH>ID = ?"
	q := ""
	args := []interface{}{ID}
	if publType == "post" {
		q = strings.ReplaceAll(strings.ReplaceAll(sq, "<TABLE>", "Posts"), "<WHICH>", which)
	} else if publType == "event" {
		q = strings.ReplaceAll(strings.ReplaceAll(sq, "<TABLE>", "Events"), "<WHICH>", which)
	} else {
		q = strings.ReplaceAll(strings.ReplaceAll(sq, "<TABLE>", "Posts"), "<WHICH>", which) +
			" UNION ALL " +
			strings.ReplaceAll(strings.ReplaceAll(sq, "<TABLE>", "Events"), "<WHICH>", which)
		args = append(args, ID)
	}
	q += " ORDER BY datetime DESC LIMIT ?,? "
	first, count := getLimits(r)
	args = append(args, first, count)

	// res ex: [[4 post 1588670115000] [3 event 1588670115000]]
	publIDS, e := orm.GetWithQueryAndArgs(q, args)
	if e != nil {
		return nil, errors.New("n/d")
	}

	postIDs, eventIDs := []string{}, []string{}
	for _, v := range publIDS {
		id := "|" + v[1].(string) + " "
		if v[1] == "post" {
			postIDs = append(postIDs, id)
			continue
		}
		eventIDs = append(eventIDs, id)
	}

	posts, events := []map[string]interface{}{}, []map[string]interface{}{}
	if publType == "post" || publType == "all" {
		posts, e = GetPosts(orm.DoSQLOption("instr(?, '|' || p.id || ' ')", "", "", postIDs), ID)
		if e != nil {
			return nil, e
		}
	}
	if publType == "event" || publType == "all" {
		events, e = GetEvents(orm.DoSQLOption("instr(?, '|' || e.id || ' ')", "", "", eventIDs), ID)
		if e != nil {
			return nil, e
		}
	}
	return append(posts, events...), nil
}

func ClippedFiles(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ID, e := strconv.Atoi(r.FormValue("id"))
	if e != nil {
		return nil, errors.New("Wrong id")
	}

	clippedType := r.FormValue("type")
	if clippedType != "comment" && clippedType != "post" && clippedType != "message" {
		return nil, errors.New("Wrong type")
	}

	return orm.GetWithSubqueries(
		orm.SQLSelectParams{
			Table:   "Files",
			What:    "*",
			Options: orm.DoSQLOption(clippedType+"ID=?", "", "", ID),
		},
		nil,
		nil,
		nil,
		orm.ClippedFile{},
	)
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	file, content, e := orm.GetFileFromDrive(strings.Split(r.URL.Path, "/")[2])
	if e != nil {
		return
	}

	ftype := http.DetectContentType(content[:512])
	w.Header().Set("Content-Disposition", "attachment; filename="+file.Name)
	w.Header().Set("Content-Type", ftype)
	io.Copy(w, bytes.NewReader(content))
}
