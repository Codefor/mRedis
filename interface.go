package main

import (
	"bufio"
	"bytes"
	"container/list"
	"hash/crc64"
	"strconv"
)

type redisCommand struct {
	name   string
	proc   redisCommandProc
	arity  int
	sflags string /* Flags as string represenation, one char per flag. */
	flags  int    /* The actual flags, obtained from the 'sflags' field. */
	/* Use a function to determine keys arguments in a command line. */
	getkeys_proc redisGetKeysProc
	/* What keys should be loaded in background when calling this command? */
	firstkey     int /* The first argument that's a key (0 = no keys) */
	lastkey      int /* THe last argument that's a key */
	keystep      int /* The step between first and last key */
	microseconds int64
	calls        int64
}

type redisCommandProc func(c *redisClient) int

func getCommand(c *redisClient) int {
	return getGenericCommand(c)
}

func setCommand(c *redisClient) int {
	return setGenericCommand(c, false)
}

func setnxCommand(c *redisClient) int {
	return setGenericCommand(c, true)
}

func setexCommand(c *redisClient) int {
	return REDIS_OK
}

func psetexCommand(c *redisClient) int {
	return REDIS_OK
}

func appendCommand(c *redisClient) int {
	key := string(c.argv[1])
	toappendvalue := string(c.argv[2])

	value, present := c.db.dict[key]
	if present {
		value.(*robj).ptr = value.(*robj).ptr.(string) + toappendvalue
	} else {
		value.(*robj).ptr = toappendvalue
	}

	redisLog(REDIS_DEBUG, "append command", value)
	c.addReplyLongLong(int64(len(value.(*robj).ptr.(string))))

	return REDIS_OK
}

func strlenCommand(c *redisClient) int {
	key := string(c.argv[1])
	if value, present := c.db.dict[key]; present {
		redisLog(REDIS_DEBUG, "strlen command", value, present)
		c.addReplyLongLong(int64(len(value.(*robj).ptr.(string))))
	} else {
		c.addReply(shared.czero)
	}

	return REDIS_OK
}

func delCommand(c *redisClient) int {
	delCnt := 0

	for _, key := range c.argv[1:] {
		if _, present := c.db.dict[string(key)]; present {
			delete(c.db.dict, string(key))
			delCnt += 1
		}
	}

	c.addReplyLongLong(int64(delCnt))
	return REDIS_OK
}

func existsCommand(c *redisClient) int {
	key := string(c.argv[1])
	if _, present := c.db.dict[key]; present {
		c.addReply(shared.cone)
	} else {
		c.addReply(shared.czero)
	}
	return REDIS_OK
}

func setbitCommand(c *redisClient) int {
	return REDIS_OK
}
func getbitCommand(c *redisClient) int {
	return REDIS_OK
}
func setrangeCommand(c *redisClient) int {
	return REDIS_OK
}
func getrangeCommand(c *redisClient) int {
	start, err := strconv.ParseInt(string(c.argv[2]), 10, 32)
	if err != nil {
		c.addReplyError("value is not an integer or out of range")
		redisLog(REDIS_DEBUG, err)
		return REDIS_ERR
	}
	end, err := strconv.ParseInt(string(c.argv[3]), 10, 32)

	if err != nil {
		redisLog(REDIS_DEBUG, err)
		c.addReplyError("value is not an integer or out of range")
		return REDIS_ERR
	}
	o, present := c.db.dict[string(c.argv[1])]
	if present {
		llen := int64(len(o.(*robj).ptr.(string)))
		redisLog(REDIS_DEBUG, "lrange :", llen)

		if start < 0 {
			start += int64(llen)
		}
		if end < 0 {
			end += int64(llen)
		}
		if start < 0 {
			start = 0
		}
		if end < 0 {
			end = 0
		}
		if start > end || start >= int64(llen) {
			c.addReply(shared.emptybulk)
			return REDIS_ERR
		}
		if end >= int64(llen) {
			end = int64(llen - 1)
		}
		c.addReply(createStringObject(o.(*robj).ptr.(string)[start:end], 0))
	} else {
		c.addReply(shared.emptybulk)
		return REDIS_OK
	}

	return REDIS_OK
}

func incrCommand(c *redisClient) int {
	return REDIS_OK
}
func decrCommand(c *redisClient) int {
	return REDIS_OK
}
func mgetCommand(c *redisClient) int {
	return REDIS_OK
}
func rpushCommand(c *redisClient) int {
	pushGenericCommand(c, 1)
	return REDIS_OK
}
func lpushCommand(c *redisClient) int {
	pushGenericCommand(c, 0)
	return REDIS_OK
}

func rpushxCommand(c *redisClient) int {
	pushxGenericCommand(c, 1)
	return REDIS_OK
}

func lpushxCommand(c *redisClient) int {
	pushxGenericCommand(c, 0)
	return REDIS_OK
}

func linsertCommand(c *redisClient) int {
	return REDIS_OK
}
func rpopCommand(c *redisClient) int {
	return REDIS_OK
}
func lpopCommand(c *redisClient) int {
	return REDIS_OK
}
func brpopCommand(c *redisClient) int {
	return REDIS_OK
}
func brpoplpushCommand(c *redisClient) int {
	return REDIS_OK
}
func blpopCommand(c *redisClient) int {
	return REDIS_OK
}
func llenCommand(c *redisClient) int {
	return REDIS_OK
}
func lindexCommand(c *redisClient) int {
	return REDIS_OK
}
func lsetCommand(c *redisClient) int {
	return REDIS_OK
}
func lrangeCommand(c *redisClient) int {
	start, err := strconv.ParseInt(string(c.argv[2]), 10, 32)
	if err != nil {
		c.addReplyError("value is not an integer or out of range")
		redisLog(REDIS_DEBUG, err)
		return REDIS_ERR
	}
	end, err := strconv.ParseInt(string(c.argv[3]), 10, 32)

	if err != nil {
		redisLog(REDIS_DEBUG, err)
		c.addReplyError("value is not an integer or out of range")
		return REDIS_ERR
	}

	o, present := c.db.dict[string(c.argv[1])]
	if present {
		if checkType(c, o.(*robj), REDIS_LIST) == REDIS_ERR {
			return REDIS_ERR
		}
		//container.list linked-double-list implement
		llen := int64(o.(*robj).ptr.(*list.List).Len())
		redisLog(REDIS_DEBUG, "lrange :", llen)
		if start < 0 {
			start += int64(llen)
		}
		if end < 0 {
			end += int64(llen)
		}
		if end < 0 {
			end = 0
		}
		if start > end || start >= int64(llen) {
			c.addReply(shared.emptymultibulk)
			return REDIS_ERR
		}
		if end >= int64(llen) {
			end = int64(llen - 1)
		}
		rangelen := end - start + 1
		redisLog(REDIS_DEBUG, "lrange :", rangelen, start, end)
		c.addReplyMultiBulkLen(rangelen)

		/**a->b->c->d->e->f
		 * if start < llen /2,we search from the Front
		 * else,we search from the Back
		 */
		if start > llen/2 {
			start -= llen
		}
		//find the start node
		ln := listIndex(o.(*robj).ptr.(*list.List), int(start))
		/**
		h := o.(*robj).ptr.(*list.List).Front()
		for h != nil {
			redisLog(REDIS_DEBUG, h.Value.(*robj).ptr.(string))
			h = h.Next()
		}*/
		for rangelen > 0 {
			redisLog(REDIS_DEBUG, "lrange:", ln.Value)
			c.addReplyBulk(ln.Value.(*robj))
			ln = ln.Next()
			rangelen--
		}
	} else {
		c.addReply(shared.emptymultibulk)
		return REDIS_OK
	}

	return REDIS_OK
}
func ltrimCommand(c *redisClient) int {
	return REDIS_OK
}
func lremCommand(c *redisClient) int {
	return REDIS_OK
}
func rpoplpushCommand(c *redisClient) int {
	return REDIS_OK
}
func saddCommand(c *redisClient) int {
	return REDIS_OK
}
func sremCommand(c *redisClient) int {
	return REDIS_OK
}
func smoveCommand(c *redisClient) int {
	return REDIS_OK
}
func sismemberCommand(c *redisClient) int {
	return REDIS_OK
}
func scardCommand(c *redisClient) int {
	return REDIS_OK
}
func spopCommand(c *redisClient) int {
	return REDIS_OK
}
func srandmemberCommand(c *redisClient) int {
	return REDIS_OK
}
func sinterstoreCommand(c *redisClient) int {
	return REDIS_OK
}
func sunionCommand(c *redisClient) int {
	return REDIS_OK
}
func sunionstoreCommand(c *redisClient) int {
	return REDIS_OK
}
func sdiffCommand(c *redisClient) int {
	return REDIS_OK
}
func sdiffstoreCommand(c *redisClient) int {
	return REDIS_OK
}
func sinterCommand(c *redisClient) int {
	return REDIS_OK
}
func zaddCommand(c *redisClient) int {
	return REDIS_OK
}
func zincrbyCommand(c *redisClient) int {
	return REDIS_OK
}
func zremCommand(c *redisClient) int {
	return REDIS_OK
}
func zremrangebyscoreCommand(c *redisClient) int {
	return REDIS_OK
}
func zremrangebyrankCommand(c *redisClient) int {
	return REDIS_OK
}
func zunionstoreCommand(c *redisClient) int {
	return REDIS_OK
}
func zinterstoreCommand(c *redisClient) int {
	return REDIS_OK
}
func zrangeCommand(c *redisClient) int {
	return REDIS_OK
}
func zrangebyscoreCommand(c *redisClient) int {
	return REDIS_OK
}
func zrevrangebyscoreCommand(c *redisClient) int {
	return REDIS_OK
}
func zcountCommand(c *redisClient) int {
	return REDIS_OK
}
func zrevrangeCommand(c *redisClient) int {
	return REDIS_OK
}
func zcardCommand(c *redisClient) int {
	return REDIS_OK
}
func zscoreCommand(c *redisClient) int {
	return REDIS_OK
}
func zrankCommand(c *redisClient) int {
	return REDIS_OK
}
func zrevrankCommand(c *redisClient) int {
	return REDIS_OK
}
func hsetCommand(c *redisClient) int {
	return REDIS_OK
}
func hsetnxCommand(c *redisClient) int {
	return REDIS_OK
}
func hgetCommand(c *redisClient) int {
	return REDIS_OK
}
func hmsetCommand(c *redisClient) int {
	return REDIS_OK
}
func hmgetCommand(c *redisClient) int {
	return REDIS_OK
}
func hincrbyCommand(c *redisClient) int {
	return REDIS_OK
}
func hincrbyfloatCommand(c *redisClient) int {
	return REDIS_OK
}
func hdelCommand(c *redisClient) int {
	return REDIS_OK
}
func hlenCommand(c *redisClient) int {
	return REDIS_OK
}
func hkeysCommand(c *redisClient) int {
	return REDIS_OK
}
func hvalsCommand(c *redisClient) int {
	return REDIS_OK
}
func hgetallCommand(c *redisClient) int {
	return REDIS_OK
}
func hexistsCommand(c *redisClient) int {
	return REDIS_OK
}
func incrbyCommand(c *redisClient) int {
	return REDIS_OK
}
func decrbyCommand(c *redisClient) int {
	return REDIS_OK
}
func incrbyfloatCommand(c *redisClient) int {
	return REDIS_OK
}

func getsetCommand(c *redisClient) int {
	if getGenericCommand(c) == REDIS_ERR {
		return REDIS_ERR
	}
	c.db.set(string(c.argv[1]), createStringObject(string(c.argv[2]), 0))
	return REDIS_OK
}

func msetCommand(c *redisClient) int {
	return REDIS_OK
}
func msetnxCommand(c *redisClient) int {
	return REDIS_OK
}

func randomkeyCommand(c *redisClient) int {
	key := c.db.randomKey()
	if key == nil {
		c.addReply(shared.nullbulk)
		return REDIS_ERR
	}
	c.addReplyBulk(key.(*robj))
	return REDIS_OK
}

func selectCommand(c *redisClient) int {
	idx, _ := strconv.ParseInt(string(c.argv[1]), 10, 32)

	if c.selectDB(int(idx)) != REDIS_OK {
		c.addReplyError("invalid DB index")
		return REDIS_ERR
	} else {
		c.addReply(shared.ok)
	}

	return REDIS_OK
}

func moveCommand(c *redisClient) int {
	return REDIS_OK
}
func renameCommand(c *redisClient) int {
	return REDIS_OK
}
func renamenxCommand(c *redisClient) int {
	return REDIS_OK
}
func expireCommand(c *redisClient) int {
	return REDIS_OK
}
func expireatCommand(c *redisClient) int {
	return REDIS_OK
}
func pexpireCommand(c *redisClient) int {
	return REDIS_OK
}
func pexpireatCommand(c *redisClient) int {
	return REDIS_OK
}
func keysCommand(c *redisClient) int {
	return REDIS_OK
}
func dbsizeCommand(c *redisClient) int {
	return REDIS_OK
}
func authCommand(c *redisClient) int {
	return REDIS_OK
}

func pingCommand(c *redisClient) int {
	c.addReply(shared.pong)
	return REDIS_OK
}

func echoCommand(c *redisClient) int {
	s := string(c.argv[1])
	o := createStringObject(s, 0)
	c.addReplyBulk(o)
	return REDIS_OK
}

func saveCommand(c *redisClient) int {
	if mRedisServer.rdbSave("dump.rdb") != REDIS_OK {
		c.addReply(shared.err)
		return REDIS_ERR
	} else {
		c.addReply(shared.ok)
		return REDIS_OK
	}
}

func bgsaveCommand(c *redisClient) int {
	go mRedisServer.rdbSave("dump.rdb")
	c.addReply(shared.ok)
	return REDIS_OK
}

func bgrewriteaofCommand(c *redisClient) int {
	return REDIS_OK
}
func shutdownCommand(c *redisClient) int {
	return REDIS_OK
}
func lastsaveCommand(c *redisClient) int {
	return REDIS_OK
}
func typeCommand(c *redisClient) int {
	return REDIS_OK
}
func multiCommand(c *redisClient) int {
	return REDIS_OK
}
func execCommand(c *redisClient) int {
	return REDIS_OK
}
func discardCommand(c *redisClient) int {
	return REDIS_OK
}
func syncCommand(c *redisClient) int {
	return REDIS_OK
}
func replconfCommand(c *redisClient) int {
	return REDIS_OK
}
func flushdbCommand(c *redisClient) int {
	return REDIS_OK
}
func flushallCommand(c *redisClient) int {
	return REDIS_OK
}
func sortCommand(c *redisClient) int {
	return REDIS_OK
}
func infoCommand(c *redisClient) int {
	return REDIS_OK
}
func monitorCommand(c *redisClient) int {
	return REDIS_OK
}
func ttlCommand(c *redisClient) int {
	return REDIS_OK
}
func pttlCommand(c *redisClient) int {
	return REDIS_OK
}
func persistCommand(c *redisClient) int {
	return REDIS_OK
}
func slaveofCommand(c *redisClient) int {
	return REDIS_OK
}
func debugCommand(c *redisClient) int {
	return REDIS_OK
}
func configCommand(c *redisClient) int {
	return REDIS_OK
}
func subscribeCommand(c *redisClient) int {
	return REDIS_OK
}
func unsubscribeCommand(c *redisClient) int {
	return REDIS_OK
}
func psubscribeCommand(c *redisClient) int {
	return REDIS_OK
}
func punsubscribeCommand(c *redisClient) int {
	return REDIS_OK
}
func publishCommand(c *redisClient) int {
	return REDIS_OK
}
func watchCommand(c *redisClient) int {
	return REDIS_OK
}
func unwatchCommand(c *redisClient) int {
	return REDIS_OK
}
func restoreCommand(c *redisClient) int {
	return REDIS_OK
}
func migrateCommand(c *redisClient) int {
	return REDIS_OK
}
func dumpCommand(c *redisClient) int {
	key := string(c.argv[1])
	bbuf := new(bytes.Buffer)
	buf := bufio.NewWriter(bbuf)

	o, present := c.db.dict[key]
	if present {
		rdbSaveObjectType(buf, o.(*robj))
		rdbSaveObject(buf, o.(*robj))

		bbuf.Write([]byte{
			byte(REDIS_RDB_VERSION & 0xff),
			byte((REDIS_RDB_VERSION >> 8) & 0xff),
		})
		t := crc64.MakeTable(crc64.ISO)
		h := crc64.New(t)
		bbuf.Write(h.Sum(bbuf.Bytes()))
		c.addReplyBulk(createStringObject(bbuf.String(), 0))
	} else {
		c.addReply(shared.nullbulk)
		return REDIS_ERR
	}
	return REDIS_OK
}

func objectCommand(c *redisClient) int {
	return REDIS_OK
}

func clientCommand(c *redisClient) int {
	if string(c.argv[1]) == "list" && len(c.argv) == 2 {
		o := getAllClientsInfoString()
		c.addReplyBulkCBuffer(o)
	}
	return REDIS_OK
}

func evalCommand(c *redisClient) int {
	return REDIS_OK
}
func evalShaCommand(c *redisClient) int {
	return REDIS_OK
}
func slowlogCommand(c *redisClient) int {
	return REDIS_OK
}
func scriptCommand(c *redisClient) int {
	return REDIS_OK
}
func timeCommand(c *redisClient) int {
	return REDIS_OK
}
func bitopCommand(c *redisClient) int {
	return REDIS_OK
}
func bitcountCommand(c *redisClient) int {
	return REDIS_OK
}
