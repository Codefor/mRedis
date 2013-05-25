package main

import (
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
	//the value must be string
	key := string(c.argv[1])
	if value, present := c.db.dict[key]; present {
		redisLog(REDIS_DEBUG, "get command", value, present)
		o := createStringObject(value.(string), 0)
		c.addReplyBulk(o)
	} else {
		c.addReply(shared.nullbulk)
	}
	c.argv = c.argv[2:]
	return REDIS_OK
}

func setCommand(c *redisClient) int {
	//the value must be string

	key := c.argv[1]
	value := string(c.argv[2])

	redisLog(REDIS_NOTICE, "Command set:", string(key), string(value), c.argv)
	c.db.set(string(key), value)
	redisLog(REDIS_NOTICE, "key "+string(key)+" set:"+string(value))

	c.addReply(shared.ok)

	c.argv = c.argv[3:]
	return REDIS_OK
}

func setnxCommand(c *redisClient) int {
	return REDIS_OK
}
func setexCommand(c *redisClient) int {
	return REDIS_OK
}
func psetexCommand(c *redisClient) int {
	return REDIS_OK
}
func appendCommand(c *redisClient) int {
	return REDIS_OK
}
func strlenCommand(c *redisClient) int {
	return REDIS_OK
}
func delCommand(c *redisClient) int {
	return REDIS_OK
}
func existsCommand(c *redisClient) int {
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
	return REDIS_OK
}
func lpushCommand(c *redisClient) int {
	return REDIS_OK
}
func rpushxCommand(c *redisClient) int {
	return REDIS_OK
}
func lpushxCommand(c *redisClient) int {
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
	return REDIS_OK
}
func msetCommand(c *redisClient) int {
	return REDIS_OK
}
func msetnxCommand(c *redisClient) int {
	return REDIS_OK
}
func randomkeyCommand(c *redisClient) int {
	key := c.db.randomKey().(*robj)
	if key == nil {
		c.addReply(shared.nullbulk)
		return REDIS_ERR
	}
	c.addReplyBulk(key)
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
	return REDIS_OK
}
func echoCommand(c *redisClient) int {
	return REDIS_OK
}

func saveCommand(c *redisClient) int {
	if mRedisServer.rdbSave("") != REDIS_OK {
		c.addReply(shared.err)
		return REDIS_ERR
	} else {
		c.addReply(shared.ok)
		return REDIS_OK
	}
}

func bgsaveCommand(c *redisClient) int {
	go mRedisServer.rdbSave("")
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
	return REDIS_OK
}
func objectCommand(c *redisClient) int {
	return REDIS_OK
}
func clientCommand(c *redisClient) int {
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
