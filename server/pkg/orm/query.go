package orm

import (
	"errors"
)

// GetFrom get more than one
func GetFrom(params SQLSelectParams) ([][]interface{}, error) {
	if params.What == "" || params.Table == "" {
		return nil, errors.New("н/д")
	}

	if result, e := selectSQL(prepareGetQueryAndArgs(params)); len(result) != 0 && e == nil {
		return result, nil
	}
	return nil, errors.New("н/д")
}

// GetOneFrom get one value
func GetOneFrom(params SQLSelectParams) ([]interface{}, error) {
	if result, e := GetFrom(params); e == nil {
		return result[0], e
	}
	return nil, errors.New("н/д")
}

// GetWithSubqueries get nested querys
func GetWithSubqueries(mainQ SQLSelectParams, subQuerys []SQLSelectParams, joinAs, qAs []string, sampleStruct interface{}) ([]map[string]interface{}, error) {
	if len(subQuerys) != len(qAs) {
		return nil, errors.New("len(querys) != len(queryAs)")
	}
	for i, v := range subQuerys {
		curQ, curArgs := prepareGetQueryAndArgs(v)
		mainQ.What += ", (" + curQ + ") AS " + qAs[i]
		mainQ.Args = append(mainQ.Args, curArgs...)
	}

	joinAs = append(joinAs, qAs...)

	result, e := GetFrom(mainQ)
	if len(result) == 0 || e != nil {
		return nil, errors.New("н/д")
	}
	return MapFromStructAndMatrix(result, sampleStruct, joinAs...), nil
}

// GetWithQueryAndArgs get with query and args
func GetWithQueryAndArgs(query string, args []interface{}) ([][]interface{}, error) {
	return selectSQL(query, args)
}

// GeneralGet get with one select params
func GeneralGet(selectParams SQLSelectParams, joinAs []string, sampleStruct interface{}) []map[string]interface{} {
	results, e := GetFrom(selectParams)
	if e != nil || len(results) == 0 {
		return nil
	}
	return MapFromStructAndMatrix(results, sampleStruct, joinAs...)
}

// DoSQLOption create new sqloption & return
func DoSQLOption(where, order, limit string, args ...interface{}) SQLOption {
	return SQLOption{Where: where, Order: order, Limit: limit, Args: args}
}

// DoSQLJoin create new sqljoin & return
func DoSQLJoin(jtype, jtable, inter string, args ...interface{}) SQLJoin {
	return SQLJoin{JoinType: jtype, JoinTable: jtable, Intersection: inter, Args: args}
}
