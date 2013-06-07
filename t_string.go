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

func setGenericCommand(c *redisClient, nx bool) int {
	key := string(c.argv[1])
	value := string(c.argv[2])

	_, present := c.db.dict[key]
	if present && nx {
		//setnx do set only if key does not exist
		c.addReply(shared.czero)
		return REDIS_ERR
	}

	c.db.set(string(key), createStringObject(value, 0))
	if nx {
		c.addReply(shared.cone)
	} else {
		c.addReply(shared.ok)
	}

	return REDIS_OK
}
