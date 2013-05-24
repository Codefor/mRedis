package main

type Db struct {
	dict    map[interface{}]interface{}
	expires map[interface{}]interface{}
	id      int
}

func NewDb(id int) *Db {
	return &Db{
		dict:    make(map[interface{}]interface{}),
		expires: make(map[interface{}]interface{}),
		id:      id,
	}
}

func (db *Db) set(key, value interface{}) int {
	db.dict[key] = value
	return REDIS_OK
}

func (db *Db) get(key interface{}) interface{} {
	if value, present := db.dict[key]; present {
		return value
	}
	return nil
}

func (db *Db) randomKey() interface{} {
	//use the random feature of the golang map
	for key, _ := range db.dict {
		return key
	}

	return &robj{}
}

func noPreloadGetKeys(cmd *redisCommand, argv **robj, argc int, numkeys *int, flags int) *int {
	a := 0
	return &a
}
func zunionInterGetKeys(cmd *redisCommand, argv **robj, argc int, numkeys *int, flags int) *int {
	a := 0
	return &a
}
func renameGetKeys(cmd *redisCommand, argv **robj, argc int, numkeys *int, flags int) *int {
	a := 0
	return &a
}
