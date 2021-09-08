package app

import (
	"net/http"
	"strconv"
)

const (
	FOLLOWING_USER_ONE_Q   string = "SELECT receiverUserID FROM Relations WHERE senderUserID = ? AND value != -1"
	FOLLOWING_USER_BOTH_Q  string = "SELECT senderUserID FROM Relations WHERE receiverUserID = ? AND value = 0"
	FOLLOWING_USER_GROUP_Q string = "SELECT receiverGroupID FROM Relations WHERE senderUserID = ? AND value = 1"
	FOLLOWERS_USER_ONE_Q   string = "SELECT senderUserID FROM Relations WHERE receiverUserID = ? AND value != -1"
	FOLLOWERS_USER_BOTH_Q  string = "SELECT receiverUserID FROM Relations WHERE senderUserID = ? AND value = 0"
	FOLLOWERS_USER_GROUP_Q string = "SELECT receiverGroupID FROM Relations WHERE senderUserID = ? AND value = 1"
	FOLLOWERS_GROUP_ONE_Q  string = "SELECT senderUserID FROM Relations WHERE receiverGroupID = ? AND value != -1"
	FOLLOWERS_GROUP_BOTH_Q string = "SELECT receiverUserID FROM Relations WHERE senderGroupID = ? AND value = 0"
	REQUEST_USER_USER_Q    string = "SELECT senderUserID FROM Relations WHERE receiverUserID = ? AND value = -1"
	REQUEST_USER_GROUP_Q   string = "SELECT senderGroupID FROM Relations WHERE receiverUserID = ? AND value = -1"
	REQUEST_GROUP_USER_Q   string = "SELECT senderUserID FROM Relations WHERE receiverGroupID = ? AND value = -1"
	PRIVATE_ACCESS_CHECK   string = `
		SELECT id IS NOT NULL FROM Relations 
		WHERE (
			(
				userID IS NOT NULL AND (
					(senderUserID = ? AND receiverUserID = userID) OR 
					(receiverUserID = ? AND senderUserID = userID AND value = 0)
				)
			) OR
			(groupID IS NOT NULL AND (senderUserID = ? AND receiverGroupID = groupID))
		)`
)

// Gets int value from string with setted default value
func getIntFromString(src string, def int) int {
	val := def
	if v, e := strconv.Atoi(src); e == nil {
		val = v
	}
	return val
}

func getLimits(r *http.Request) (int, int) {
	return getIntFromString(r.FormValue("from"), 0), getIntFromString(r.FormValue("step"), 10)
}
