package main

import (
	"container/list"
)

func pushxGenericCommand(c *redisClient, where int) {
	key := string(c.argv[0])
	value := string(c.argv[1])

	o, present := c.db.dict[key]
	if present {
		olist := o.(*robj).ptr.(*list.List)
		if where == 1 {
			//REDIS_TAIL = 1
			olist.PushBack(value)
		} else if where == 0 {
			//REDIS_HEAD = 0
			olist.PushFront(value)
		}
		c.addReplyLongLong(int64(olist.Len()))
	} else {
		c.addReply(shared.czero)
	}
}

func pushGenericCommand(c *redisClient, where int) {
	key := string(c.argv[1])
	var o *robj
	value, present := c.db.dict[key]
	redisLog(REDIS_NOTICE, "pushGenericCommand:", value, present)

	if present {
		o = value.(*robj)
		if o.rtype != REDIS_LIST {
			c.addReply(shared.wrongtypeerr)
			return
		}
	} else {
		o = createListObject()
		c.db.set(key, o)
		redisLog(REDIS_NOTICE, "lpush:", c.db.dict)
	}

	cnt := 0
	olist := o.ptr.(*list.List)
	for _, v := range c.argv[2:] {
		if where == 1 {
			//REDIS_TAIL = 1
			olist.PushBack(createStringObject(string(v), 0))
		} else if where == 0 {
			//REDIS_HEAD = 0
			olist.PushFront(createStringObject(string(v), 0))

		}
		cnt += 1
	}
	c.addReplyLongLong(int64(cnt))
}

func listTypePush(c *robj, value *robj, where int) {
}

func listIndex(l *list.List, index int) (n *list.Element) {
	redisLog(REDIS_DEBUG, "listIndex:", index)
	if index < 0 {
		index = (-index) - 1
		n = l.Back()
		for index > 0 && n != nil {
			n = n.Prev()
			index--
		}
	} else {
		n = l.Front()
		for index > 0 && n != nil {
			n = n.Next()
			index--
		}
	}
	return n

}
