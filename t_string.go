package main

func getGenericCommand(c *redisClient) int {
	//the value must be string
	key := string(c.argv[1])
	redisLog(REDIS_DEBUG, "get command", key, c.db.dict[key])
	if value, present := c.db.dict[key]; present {
		redisLog(REDIS_DEBUG, "get command", value, present)
		if value.(*robj).rtype != REDIS_STRING {
			c.addReply(shared.wrongtypeerr)
			return REDIS_ERR
		}
		c.addReplyBulk(value.(*robj))
	} else {
		c.addReply(shared.nullbulk)
	}
	return REDIS_OK
}
