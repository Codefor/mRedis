package main

// redis protocol's special bytes
const (
	cr_byte    byte = byte('\r')
	lf_byte         = byte('\n')
	space_byte      = byte(' ')
	err_byte        = byte('-')
	ok_byte         = byte('+')
	count_byte      = byte('*')
	size_byte       = byte('$')
	num_byte        = byte(':')
	true_byte       = byte('1')
)

var (
	crlf_bytes = []byte{cr_byte, lf_byte}
)

type commandTable map[string]*redisCommand

/* Our command table.
 *
 * Every entry is composed of the following fields:
 *
 * name: a string representing the command name.
 * function: pointer to the C function implementing the command.
 * arity: number of arguments, it is possible to use -N to say >= N
 * sflags: command flags as string. See below for a table of flags.
 * flags: flags as bitmask. Computed by Redis using the 'sflags' field.
 * get_keys_proc: an optional function to get key arguments from a command.
 *                This is only used when the following three fields are not
 *                enough to specify what arguments are keys.
 * first_key_index: first argument that is a key
 * last_key_index: last argument that is a key
 * key_step: step to get all the keys from first to last argument. For instance
 *           in MSET the step is two since arguments are key,val,key,val,...
 * microseconds: microseconds of total execution time for this command.
 * calls: total number of calls of this command.
 *
 * The flags, microseconds and calls fields are computed by Redis and should
 * always be set to zero.
 *
 * Command flags are expressed using strings where every character represents
 * a flag. Later the populateCommandTable() function will take care of
 * populating the real 'flags' field using this characters.
 *
 * This is the meaning of the flags:
 *
 * w: write command (may modify the key space).
 * r: read command  (will never modify the key space).
 * m: may increase memory usage once called. Don't allow if out of memory.
 * a: admin command, like SAVE or SHUTDOWN.
 * p: Pub/Sub related command.
 * f: force replication of this command, regarless of server.dirty.
 * s: command not allowed in scripts.
 * R: random command. Command is not deterministic, that is, the same command
 *    with the same arguments, with the same key space, may have different
 *    results. For instance SPOP and RANDOMKEY are two random commands.
 * S: Sort command output array if called from script, so that the output
 *    is deterministic.
 * l: Allow command while loading the database.
 * t: Allow command while a slave has stale data but is not allowed to
 *    server this data. Normally no command is accepted in this condition
 *    but just a few.
 * M: Do not automatically propagate the command on MONITOR.
 */

var redisCommandTable = commandTable{
	"get":              &redisCommand{"get", getCommand, 2, "r", 0, nil, 1, 1, 1, 0, 0},
	"set":              &redisCommand{"set", setCommand, 3, "wm", 0, noPreloadGetKeys, 1, 1, 1, 0, 0},
	"setnx":            &redisCommand{"setnx", setnxCommand, 3, "wm", 0, noPreloadGetKeys, 1, 1, 1, 0, 0},
	"setex":            &redisCommand{"setex", setexCommand, 4, "wm", 0, noPreloadGetKeys, 1, 1, 1, 0, 0},
	"psetex":           &redisCommand{"psetex", psetexCommand, 4, "wm", 0, noPreloadGetKeys, 1, 1, 1, 0, 0},
	"append":           &redisCommand{"append", appendCommand, 3, "wm", 0, nil, 1, 1, 1, 0, 0},
	"strlen":           &redisCommand{"strlen", strlenCommand, 2, "r", 0, nil, 1, 1, 1, 0, 0},
	"del":              &redisCommand{"del", delCommand, -2, "w", 0, noPreloadGetKeys, 1, -1, 1, 0, 0},
	"exists":           &redisCommand{"exists", existsCommand, 2, "r", 0, nil, 1, 1, 1, 0, 0},
	"setbit":           &redisCommand{"setbit", setbitCommand, 4, "wm", 0, nil, 1, 1, 1, 0, 0},
	"getbit":           &redisCommand{"getbit", getbitCommand, 3, "r", 0, nil, 1, 1, 1, 0, 0},
	"setrange":         &redisCommand{"setrange", setrangeCommand, 4, "wm", 0, nil, 1, 1, 1, 0, 0},
	"getrange":         &redisCommand{"getrange", getrangeCommand, 4, "r", 0, nil, 1, 1, 1, 0, 0},
	"substr":           &redisCommand{"substr", getrangeCommand, 4, "r", 0, nil, 1, 1, 1, 0, 0},
	"incr":             &redisCommand{"incr", incrCommand, 2, "wm", 0, nil, 1, 1, 1, 0, 0},
	"decr":             &redisCommand{"decr", decrCommand, 2, "wm", 0, nil, 1, 1, 1, 0, 0},
	"mget":             &redisCommand{"mget", mgetCommand, -2, "r", 0, nil, 1, -1, 1, 0, 0},
	"rpush":            &redisCommand{"rpush", rpushCommand, -3, "wm", 0, nil, 1, 1, 1, 0, 0},
	"lpush":            &redisCommand{"lpush", lpushCommand, -3, "wm", 0, nil, 1, 1, 1, 0, 0},
	"rpushx":           &redisCommand{"rpushx", rpushxCommand, 3, "wm", 0, nil, 1, 1, 1, 0, 0},
	"lpushx":           &redisCommand{"lpushx", lpushxCommand, 3, "wm", 0, nil, 1, 1, 1, 0, 0},
	"linsert":          &redisCommand{"linsert", linsertCommand, 5, "wm", 0, nil, 1, 1, 1, 0, 0},
	"rpop":             &redisCommand{"rpop", rpopCommand, 2, "w", 0, nil, 1, 1, 1, 0, 0},
	"lpop":             &redisCommand{"lpop", lpopCommand, 2, "w", 0, nil, 1, 1, 1, 0, 0},
	"brpop":            &redisCommand{"brpop", brpopCommand, -3, "ws", 0, nil, 1, 1, 1, 0, 0},
	"brpoplpush":       &redisCommand{"brpoplpush", brpoplpushCommand, 4, "wms", 0, nil, 1, 2, 1, 0, 0},
	"blpop":            &redisCommand{"blpop", blpopCommand, -3, "ws", 0, nil, 1, -2, 1, 0, 0},
	"llen":             &redisCommand{"llen", llenCommand, 2, "r", 0, nil, 1, 1, 1, 0, 0},
	"lindex":           &redisCommand{"lindex", lindexCommand, 3, "r", 0, nil, 1, 1, 1, 0, 0},
	"lset":             &redisCommand{"lset", lsetCommand, 4, "wm", 0, nil, 1, 1, 1, 0, 0},
	"lrange":           &redisCommand{"lrange", lrangeCommand, 4, "r", 0, nil, 1, 1, 1, 0, 0},
	"ltrim":            &redisCommand{"ltrim", ltrimCommand, 4, "w", 0, nil, 1, 1, 1, 0, 0},
	"lrem":             &redisCommand{"lrem", lremCommand, 4, "w", 0, nil, 1, 1, 1, 0, 0},
	"rpoplpush":        &redisCommand{"rpoplpush", rpoplpushCommand, 3, "wm", 0, nil, 1, 2, 1, 0, 0},
	"sadd":             &redisCommand{"sadd", saddCommand, -3, "wm", 0, nil, 1, 1, 1, 0, 0},
	"srem":             &redisCommand{"srem", sremCommand, -3, "w", 0, nil, 1, 1, 1, 0, 0},
	"smove":            &redisCommand{"smove", smoveCommand, 4, "w", 0, nil, 1, 2, 1, 0, 0},
	"sismember":        &redisCommand{"sismember", sismemberCommand, 3, "r", 0, nil, 1, 1, 1, 0, 0},
	"scard":            &redisCommand{"scard", scardCommand, 2, "r", 0, nil, 1, 1, 1, 0, 0},
	"spop":             &redisCommand{"spop", spopCommand, 2, "wRs", 0, nil, 1, 1, 1, 0, 0},
	"srandmember":      &redisCommand{"srandmember", srandmemberCommand, -2, "rR", 0, nil, 1, 1, 1, 0, 0},
	"sinter":           &redisCommand{"sinter", sinterCommand, -2, "rS", 0, nil, 1, -1, 1, 0, 0},
	"sinterstore":      &redisCommand{"sinterstore", sinterstoreCommand, -3, "wm", 0, nil, 1, -1, 1, 0, 0},
	"sunion":           &redisCommand{"sunion", sunionCommand, -2, "rS", 0, nil, 1, -1, 1, 0, 0},
	"sunionstore":      &redisCommand{"sunionstore", sunionstoreCommand, -3, "wm", 0, nil, 1, -1, 1, 0, 0},
	"sdiff":            &redisCommand{"sdiff", sdiffCommand, -2, "rS", 0, nil, 1, -1, 1, 0, 0},
	"sdiffstore":       &redisCommand{"sdiffstore", sdiffstoreCommand, -3, "wm", 0, nil, 1, -1, 1, 0, 0},
	"smembers":         &redisCommand{"smembers", sinterCommand, 2, "rS", 0, nil, 1, 1, 1, 0, 0},
	"zadd":             &redisCommand{"zadd", zaddCommand, -4, "wm", 0, nil, 1, 1, 1, 0, 0},
	"zincrby":          &redisCommand{"zincrby", zincrbyCommand, 4, "wm", 0, nil, 1, 1, 1, 0, 0},
	"zrem":             &redisCommand{"zrem", zremCommand, -3, "w", 0, nil, 1, 1, 1, 0, 0},
	"zremrangebyscore": &redisCommand{"zremrangebyscore", zremrangebyscoreCommand, 4, "w", 0, nil, 1, 1, 1, 0, 0},
	"zremrangebyrank":  &redisCommand{"zremrangebyrank", zremrangebyrankCommand, 4, "w", 0, nil, 1, 1, 1, 0, 0},
	"zunionstore":      &redisCommand{"zunionstore", zunionstoreCommand, -4, "wm", 0, zunionInterGetKeys, 0, 0, 0, 0, 0},
	"zinterstore":      &redisCommand{"zinterstore", zinterstoreCommand, -4, "wm", 0, zunionInterGetKeys, 0, 0, 0, 0, 0},
	"zrange":           &redisCommand{"zrange", zrangeCommand, -4, "r", 0, nil, 1, 1, 1, 0, 0},
	"zrangebyscore":    &redisCommand{"zrangebyscore", zrangebyscoreCommand, -4, "r", 0, nil, 1, 1, 1, 0, 0},
	"zrevrangebyscore": &redisCommand{"zrevrangebyscore", zrevrangebyscoreCommand, -4, "r", 0, nil, 1, 1, 1, 0, 0},
	"zcount":           &redisCommand{"zcount", zcountCommand, 4, "r", 0, nil, 1, 1, 1, 0, 0},
	"zrevrange":        &redisCommand{"zrevrange", zrevrangeCommand, -4, "r", 0, nil, 1, 1, 1, 0, 0},
	"zcard":            &redisCommand{"zcard", zcardCommand, 2, "r", 0, nil, 1, 1, 1, 0, 0},
	"zscore":           &redisCommand{"zscore", zscoreCommand, 3, "r", 0, nil, 1, 1, 1, 0, 0},
	"zrank":            &redisCommand{"zrank", zrankCommand, 3, "r", 0, nil, 1, 1, 1, 0, 0},
	"zrevrank":         &redisCommand{"zrevrank", zrevrankCommand, 3, "r", 0, nil, 1, 1, 1, 0, 0},
	"hset":             &redisCommand{"hset", hsetCommand, 4, "wm", 0, nil, 1, 1, 1, 0, 0},
	"hsetnx":           &redisCommand{"hsetnx", hsetnxCommand, 4, "wm", 0, nil, 1, 1, 1, 0, 0},
	"hget":             &redisCommand{"hget", hgetCommand, 3, "r", 0, nil, 1, 1, 1, 0, 0},
	"hmset":            &redisCommand{"hmset", hmsetCommand, -4, "wm", 0, nil, 1, 1, 1, 0, 0},
	"hmget":            &redisCommand{"hmget", hmgetCommand, -3, "r", 0, nil, 1, 1, 1, 0, 0},
	"hincrby":          &redisCommand{"hincrby", hincrbyCommand, 4, "wm", 0, nil, 1, 1, 1, 0, 0},
	"hincrbyfloat":     &redisCommand{"hincrbyfloat", hincrbyfloatCommand, 4, "wm", 0, nil, 1, 1, 1, 0, 0},
	"hdel":             &redisCommand{"hdel", hdelCommand, -3, "w", 0, nil, 1, 1, 1, 0, 0},
	"hlen":             &redisCommand{"hlen", hlenCommand, 2, "r", 0, nil, 1, 1, 1, 0, 0},
	"hkeys":            &redisCommand{"hkeys", hkeysCommand, 2, "rS", 0, nil, 1, 1, 1, 0, 0},
	"hvals":            &redisCommand{"hvals", hvalsCommand, 2, "rS", 0, nil, 1, 1, 1, 0, 0},
	"hgetall":          &redisCommand{"hgetall", hgetallCommand, 2, "r", 0, nil, 1, 1, 1, 0, 0},
	"hexists":          &redisCommand{"hexists", hexistsCommand, 3, "r", 0, nil, 1, 1, 1, 0, 0},
	"incrby":           &redisCommand{"incrby", incrbyCommand, 3, "wm", 0, nil, 1, 1, 1, 0, 0},
	"decrby":           &redisCommand{"decrby", decrbyCommand, 3, "wm", 0, nil, 1, 1, 1, 0, 0},
	"incrbyfloat":      &redisCommand{"incrbyfloat", incrbyfloatCommand, 3, "wm", 0, nil, 1, 1, 1, 0, 0},
	"getset":           &redisCommand{"getset", getsetCommand, 3, "wm", 0, nil, 1, 1, 1, 0, 0},
	"mset":             &redisCommand{"mset", msetCommand, -3, "wm", 0, nil, 1, -1, 2, 0, 0},
	"msetnx":           &redisCommand{"msetnx", msetnxCommand, -3, "wm", 0, nil, 1, -1, 2, 0, 0},
	"randomkey":        &redisCommand{"randomkey", randomkeyCommand, 1, "rR", 0, nil, 0, 0, 0, 0, 0},
	"select":           &redisCommand{"select", selectCommand, 2, "r", 0, nil, 0, 0, 0, 0, 0},
	"move":             &redisCommand{"move", moveCommand, 3, "w", 0, nil, 1, 1, 1, 0, 0},
	"rename":           &redisCommand{"rename", renameCommand, 3, "w", 0, renameGetKeys, 1, 2, 1, 0, 0},
	"renamenx":         &redisCommand{"renamenx", renamenxCommand, 3, "w", 0, renameGetKeys, 1, 2, 1, 0, 0},
	"expire":           &redisCommand{"expire", expireCommand, 3, "w", 0, nil, 1, 1, 1, 0, 0},
	"expireat":         &redisCommand{"expireat", expireatCommand, 3, "w", 0, nil, 1, 1, 1, 0, 0},
	"pexpire":          &redisCommand{"pexpire", pexpireCommand, 3, "w", 0, nil, 1, 1, 1, 0, 0},
	"pexpireat":        &redisCommand{"pexpireat", pexpireatCommand, 3, "w", 0, nil, 1, 1, 1, 0, 0},
	"keys":             &redisCommand{"keys", keysCommand, 2, "rS", 0, nil, 0, 0, 0, 0, 0},
	"dbsize":           &redisCommand{"dbsize", dbsizeCommand, 1, "r", 0, nil, 0, 0, 0, 0, 0},
	"auth":             &redisCommand{"auth", authCommand, 2, "rs", 0, nil, 0, 0, 0, 0, 0},
	"ping":             &redisCommand{"ping", pingCommand, 1, "r", 0, nil, 0, 0, 0, 0, 0},
	"echo":             &redisCommand{"echo", echoCommand, 2, "r", 0, nil, 0, 0, 0, 0, 0},
	"save":             &redisCommand{"save", saveCommand, 1, "ars", 0, nil, 0, 0, 0, 0, 0},
	"bgsave":           &redisCommand{"bgsave", bgsaveCommand, 1, "ar", 0, nil, 0, 0, 0, 0, 0},
	"bgrewriteaof":     &redisCommand{"bgrewriteaof", bgrewriteaofCommand, 1, "ar", 0, nil, 0, 0, 0, 0, 0},
	"shutdown":         &redisCommand{"shutdown", shutdownCommand, -1, "ar", 0, nil, 0, 0, 0, 0, 0},
	"lastsave":         &redisCommand{"lastsave", lastsaveCommand, 1, "r", 0, nil, 0, 0, 0, 0, 0},
	"type":             &redisCommand{"type", typeCommand, 2, "r", 0, nil, 1, 1, 1, 0, 0},
	"multi":            &redisCommand{"multi", multiCommand, 1, "rs", 0, nil, 0, 0, 0, 0, 0},
	"exec":             &redisCommand{"exec", execCommand, 1, "sM", 0, nil, 0, 0, 0, 0, 0},
	"discard":          &redisCommand{"discard", discardCommand, 1, "rs", 0, nil, 0, 0, 0, 0, 0},
	"sync":             &redisCommand{"sync", syncCommand, 1, "ars", 0, nil, 0, 0, 0, 0, 0},
	"replconf":         &redisCommand{"replconf", replconfCommand, -1, "ars", 0, nil, 0, 0, 0, 0, 0},
	"flushdb":          &redisCommand{"flushdb", flushdbCommand, 1, "w", 0, nil, 0, 0, 0, 0, 0},
	"flushall":         &redisCommand{"flushall", flushallCommand, 1, "w", 0, nil, 0, 0, 0, 0, 0},
	"sort":             &redisCommand{"sort", sortCommand, -2, "wm", 0, nil, 1, 1, 1, 0, 0},
	"info":             &redisCommand{"info", infoCommand, -1, "rlt", 0, nil, 0, 0, 0, 0, 0},
	"monitor":          &redisCommand{"monitor", monitorCommand, 1, "ars", 0, nil, 0, 0, 0, 0, 0},
	"ttl":              &redisCommand{"ttl", ttlCommand, 2, "r", 0, nil, 1, 1, 1, 0, 0},
	"pttl":             &redisCommand{"pttl", pttlCommand, 2, "r", 0, nil, 1, 1, 1, 0, 0},
	"persist":          &redisCommand{"persist", persistCommand, 2, "w", 0, nil, 1, 1, 1, 0, 0},
	"slaveof":          &redisCommand{"slaveof", slaveofCommand, 3, "ast", 0, nil, 0, 0, 0, 0, 0},
	"debug":            &redisCommand{"debug", debugCommand, -2, "as", 0, nil, 0, 0, 0, 0, 0},
	"config":           &redisCommand{"config", configCommand, -2, "ar", 0, nil, 0, 0, 0, 0, 0},
	"subscribe":        &redisCommand{"subscribe", subscribeCommand, -2, "rpslt", 0, nil, 0, 0, 0, 0, 0},
	"unsubscribe":      &redisCommand{"unsubscribe", unsubscribeCommand, -1, "rpslt", 0, nil, 0, 0, 0, 0, 0},
	"psubscribe":       &redisCommand{"psubscribe", psubscribeCommand, -2, "rpslt", 0, nil, 0, 0, 0, 0, 0},
	"punsubscribe":     &redisCommand{"punsubscribe", punsubscribeCommand, -1, "rpslt", 0, nil, 0, 0, 0, 0, 0},
	"publish":          &redisCommand{"publish", publishCommand, 3, "pflt", 0, nil, 0, 0, 0, 0, 0},
	"watch":            &redisCommand{"watch", watchCommand, -2, "rs", 0, noPreloadGetKeys, 1, -1, 1, 0, 0},
	"unwatch":          &redisCommand{"unwatch", unwatchCommand, 1, "rs", 0, nil, 0, 0, 0, 0, 0},
	"restore":          &redisCommand{"restore", restoreCommand, 4, "awm", 0, nil, 1, 1, 1, 0, 0},
	"migrate":          &redisCommand{"migrate", migrateCommand, 6, "aw", 0, nil, 0, 0, 0, 0, 0},
	"dump":             &redisCommand{"dump", dumpCommand, 2, "ar", 0, nil, 1, 1, 1, 0, 0},
	"object":           &redisCommand{"object", objectCommand, -2, "r", 0, nil, 2, 2, 2, 0, 0},
	"client":           &redisCommand{"client", clientCommand, -2, "ar", 0, nil, 0, 0, 0, 0, 0},
	"eval":             &redisCommand{"eval", evalCommand, -3, "s", 0, zunionInterGetKeys, 0, 0, 0, 0, 0},
	"evalsha":          &redisCommand{"evalsha", evalShaCommand, -3, "s", 0, zunionInterGetKeys, 0, 0, 0, 0, 0},
	"slowlog":          &redisCommand{"slowlog", slowlogCommand, -2, "r", 0, nil, 0, 0, 0, 0, 0},
	"script":           &redisCommand{"script", scriptCommand, -2, "ras", 0, nil, 0, 0, 0, 0, 0},
	"time":             &redisCommand{"time", timeCommand, 1, "rR", 0, nil, 0, 0, 0, 0, 0},
	"bitop":            &redisCommand{"bitop", bitopCommand, -4, "wm", 0, nil, 2, -1, 1, 0, 0},
	"bitcount":         &redisCommand{"bitcount", bitcountCommand, -2, "r", 0, nil, 1, 1, 1, 0, 0},
}
