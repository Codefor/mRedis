package main

const (
	/* Error codes */
	REDIS_OK  = 0
	REDIS_ERR = -1

	/* Static server configuration */
	REDIS_HZ                        = 100  /* Time interrupt calls/sec. */
	REDIS_SERVERPORT                = 6379 /* TCP port */
	REDIS_MAXIDLETIME               = 0    /* default client timeout: infinite */
	REDIS_DEFAULT_DBNUM             = 16
	REDIS_CONFIGLINE_MAX            = 1024
	REDIS_EXPIRELOOKUPS_PER_CRON    = 10 /* lookup 10 expires per loop */
	REDIS_EXPIRELOOKUPS_TIME_PERC   = 25 /* CPU max % for keys collection */
	REDIS_MAX_WRITE_PER_EVENT       = (1024 * 64)
	REDIS_SHARED_SELECT_CMDS        = 10
	REDIS_SHARED_INTEGERS           = 10000
	REDIS_SHARED_BULKHDR_LEN        = 32
	REDIS_MAX_LOGMSG_LEN            = 1024 /* Default maximum length of syslog messages */
	REDIS_AOF_REWRITE_PERC          = 100
	REDIS_AOF_REWRITE_MIN_SIZE      = (1024 * 1024)
	REDIS_AOF_REWRITE_ITEMS_PER_CMD = 64
	REDIS_SLOWLOG_LOG_SLOWER_THAN   = 10000
	REDIS_SLOWLOG_MAX_LEN           = 128
	REDIS_MAX_CLIENTS               = 10000
	REDIS_AUTHPASS_MAX_LEN          = 512
	REDIS_DEFAULT_SLAVE_PRIORITY    = 100
	REDIS_REPL_TIMEOUT              = 60
	REDIS_REPL_PING_SLAVE_PERIOD    = 10
	REDIS_RUN_ID_SIZE               = 40
	REDIS_OPS_SEC_SAMPLES           = 16

	/* Protocol and I/O related defines */
	REDIS_MAX_QUERYBUF_LEN  = (1024 * 1024 * 1024) /* 1GB max query buffer. */
	REDIS_IOBUF_LEN         = (1024 * 16)          /* Generic I/O buffer size */
	REDIS_REPLY_CHUNK_BYTES = (16 * 1024)          /* 16k output buffer */
	REDIS_INLINE_MAX_SIZE   = (1024 * 64)          /* Max size of inline reads */
	REDIS_MBULK_BIG_ARG     = (1024 * 32)

	/* Hash table parameters */
	REDIS_HT_MINFILL = 10 /* Minimal hash table fill 10% */

	/* Command flags. Please check the command table defined in the redis.c file
	 * for more information about the meaning of every flag. */
	REDIS_CMD_WRITE    = 1 /* "w" flag */
	REDIS_CMD_READONLY = 2 /* "r" flag */
	REDIS_CMD_DENYOOM  = 4 /* "m" flag */
	//    REDIS_CMD_FORCE_REPLICATION 8       /* "f" flag */
	//    REDIS_CMD_ADMIN 16                  /* "a" flag */
	//    REDIS_CMD_PUBSUB 32                 /* "p" flag */
	//    REDIS_CMD_NOSCRIPT  64              /* "s" flag */
	//    REDIS_CMD_RANDOM 128                /* "R" flag */
	//    REDIS_CMD_SORT_FOR_SCRIPT 256       /* "S" flag */
	//    REDIS_CMD_LOADING 512               /* "l" flag */
	//    REDIS_CMD_STALE 1024                /* "t" flag */
	//    REDIS_CMD_SKIP_MONITOR 2048         /* "M" flag */
	//
	/* Object types */
	REDIS_STRING = 0
	REDIS_LIST   = 1
	REDIS_SET    = 2
	REDIS_ZSET   = 3
	REDIS_HASH   = 4
	/* Objects encoding. Some kind of objects like Strings and Hashes can be
	 * internally represented in multiple ways. The 'encoding' field of the object
	 * is set to one of this fields for this object. */
	REDIS_ENCODING_RAW        = 0 /* Raw representation */
	REDIS_ENCODING_INT        = 1 /* Encoded as integer */
	REDIS_ENCODING_HT         = 2 /* Encoded as hash table */
	REDIS_ENCODING_ZIPMAP     = 3 /* Encoded as zipmap */
	REDIS_ENCODING_LINKEDLIST = 4 /* Encoded as regular linked list */
	REDIS_ENCODING_ZIPLIST    = 5 /* Encoded as ziplist */
	REDIS_ENCODING_INTSET     = 6 /* Encoded as intset */
	REDIS_ENCODING_SKIPLIST   = 7 /* Encoded as skiplist */

	/* The current RDB version. When the format changes in a way that is no longer
	 * backward compatible this number gets incremented. */
	/*  REDIS RDB VERSION */
	REDIS_RDB_VERSION = 6

	/* Defines related to the dump file format. To store 32 bits lengths for short
	 * keys requires a lot of space, so we check the most significant 2 bits of
	 * the first byte to interpreter the length:
	 *
	 * 00|000000 => if the two MSB are 00 the len is the 6 bits of this byte
	 * 01|000000 00000000 =>  01, the len is 14 byes, 6 bits + 8 bits of next byte
	 * 10|000000 [32 bit integer] => if it's 01, a full 32 bit len will follow
	 * 11|000000 this means: specially encoded object will follow. The six bits
	 *           number specify the kind of object that follows.
	 *           See the REDIS_RDB_ENC_* defines.
	 *
	 * Lenghts up to 63 are stored using a single byte, most DB keys, and may
	 * values, will fit inside. */
	REDIS_RDB_6BITLEN  = 0
	REDIS_RDB_14BITLEN = 1
	REDIS_RDB_32BITLEN = 2
	REDIS_RDB_ENCVAL   = 3
	REDIS_RDB_LENERR   = 0xffffffff //UINT_MAX

	/* When a length of a string object stored on disk has the first two bits
	 * set, the remaining two bits specify a special encoding for the object
	 * accordingly to the following defines: */
	REDIS_RDB_ENC_INT8  = 0 /* 8 bit signed integer */
	REDIS_RDB_ENC_INT16 = 1 /* 16 bit signed integer */
	REDIS_RDB_ENC_INT32 = 2 /* 32 bit signed integer */
	REDIS_RDB_ENC_LZF   = 3 /* string compressed with FASTLZ */
	//
	//    /* AOF states */
	//    REDIS_AOF_OFF = 0             /* AOF is off */
	//    REDIS_AOF_ON = 1              /* AOF is on */
	//    REDIS_AOF_WAIT_REWRITE = 2    /* AOF waits rewrite to start appending */
	//
	//    /* Client flags */
	//    REDIS_SLAVE = 1       /* This client is a slave server */
	//    REDIS_MASTER = 2      /* This client is a master server */
	//    REDIS_MONITOR = 4     /* This client is a slave monitor, see MONITOR */
	//    REDIS_MULTI = 8       /* This client is in a MULTI context */
	//    REDIS_BLOCKED = 16    /* The client is waiting in a blocking operation */
	//    REDIS_DIRTY_CAS = 64  /* Watched keys modified. EXEC will fail. */
	//    REDIS_CLOSE_AFTER_REPLY = 128 /* Close after writing entire reply. */
	//    REDIS_UNBLOCKED = 256 /* This client was unblocked and is stored in
	//    server.unblocked_clients */
	//    REDIS_LUA_CLIENT = 512 /* This is a non connected client used by Lua */
	//    REDIS_ASKING = 1024   /* Client issued the ASKING command */
	//    REDIS_CLOSE_ASAP = 2048 /* Close this client ASAP */
	//
	//    /* Client request types */
	//    REDIS_REQ_INLINE = 1
	//    REDIS_REQ_MULTIBULK = 2
	//
	//    /* Client classes for client limits, currently used only for
	//    * the max-client-output-buffer limit implementation. */
	//    REDIS_CLIENT_LIMIT_CLASS_NORMAL 0
	//    REDIS_CLIENT_LIMIT_CLASS_SLAVE 1
	//    REDIS_CLIENT_LIMIT_CLASS_PUBSUB 2
	//    REDIS_CLIENT_LIMIT_NUM_CLASSES 3
	//
	//    /* Slave replication state - slave side */
	//    REDIS_REPL_NONE 0 /* No active replication */
	//    REDIS_REPL_CONNECT 1 /* Must connect to master */
	//    REDIS_REPL_CONNECTING 2 /* Connecting to master */
	//    REDIS_REPL_RECEIVE_PONG 3 /* Wait for PING reply */
	//    REDIS_REPL_TRANSFER 4 /* Receiving .rdb from master */
	//    REDIS_REPL_CONNECTED 5 /* Connected to master */
	//
	//    /* Synchronous read timeout - slave side */
	//    REDIS_REPL_SYNCIO_TIMEOUT 5
	//
	//    /* Slave replication state - from the point of view of master
	//    * Note that in SEND_BULK and ONLINE state the slave receives new updates
	//    * in its output queue. In the WAIT_BGSAVE state instead the server is waiting
	//    * to start the next background saving in order to send updates to it. */
	//    REDIS_REPL_WAIT_BGSAVE_START 3 /* master waits bgsave to start feeding it */
	//    REDIS_REPL_WAIT_BGSAVE_END 4 /* master waits bgsave to start bulk DB transmission */
	//    REDIS_REPL_SEND_BULK 5 /* master is sending the bulk DB */
	//    REDIS_REPL_ONLINE 6 /* bulk DB already transmitted, receive updates */
	//
	/* List related stuff */
	REDIS_HEAD = 0
	REDIS_TAIL = 1

	//    /* Sort operations */
	//    REDIS_SORT_GET 0
	//    REDIS_SORT_ASC 1
	//    REDIS_SORT_DESC 2
	//    REDIS_SORTKEY_MAX 1024
	//
	/* Log levels */
	REDIS_DEBUG   = 0
	REDIS_VERBOSE = 1
	REDIS_NOTICE  = 2
	REDIS_WARNING = 3
	REDIS_LOG_RAW = (1 << 10) /* Modifier to log without timestamp */
	//
	//    /* Anti-warning macro... */
	//    REDIS_NOTUSED(V) ((void) V)
	//
	//    ZSKIPLIST_MAXLEVEL 32 /* Should be enough for 2^32 elements */
	//    ZSKIPLIST_P 0.25      /* Skiplist P = 1/4 */
	//
	//    /* Append only defines */
	//    AOF_FSYNC_NO 0
	//    AOF_FSYNC_ALWAYS 1
	//    AOF_FSYNC_EVERYSEC 2
	//
	//    /* Zip structure related defaults */
	//    REDIS_HASH_MAX_ZIPLIST_ENTRIES 512
	//    REDIS_HASH_MAX_ZIPLIST_VALUE 64
	//    REDIS_LIST_MAX_ZIPLIST_ENTRIES 512
	//    REDIS_LIST_MAX_ZIPLIST_VALUE 64
	//    REDIS_SET_MAX_INTSET_ENTRIES 512
	//    REDIS_ZSET_MAX_ZIPLIST_ENTRIES 128
	//    REDIS_ZSET_MAX_ZIPLIST_VALUE 64
	//
	//    /* Sets operations codes */
	//    REDIS_OP_UNION 0
	//    REDIS_OP_DIFF 1
	//    REDIS_OP_INTER 2
	//
	//    /* Redis maxmemory strategies */
	//    REDIS_MAXMEMORY_VOLATILE_LRU 0
	//    REDIS_MAXMEMORY_VOLATILE_TTL 1
	//    REDIS_MAXMEMORY_VOLATILE_RANDOM 2
	//    REDIS_MAXMEMORY_ALLKEYS_LRU 3
	//    REDIS_MAXMEMORY_ALLKEYS_RANDOM 4
	//    REDIS_MAXMEMORY_NO_EVICTION 5
	//
	//    /* Scripting */
	//    REDIS_LUA_TIME_LIMIT 5000 /* milliseconds */
	//
	/* Units */
	UNIT_SECONDS      = 0
	UNIT_MILLISECONDS = 1

//
//    /* SHUTDOWN flags */
//    REDIS_SHUTDOWN_SAVE 1       /* Force SAVE on SHUTDOWN even if no save
//    points are configured. */
//    REDIS_SHUTDOWN_NOSAVE 2     /* Don't SAVE on SHUTDOWN. */
//
//    /* Command call flags, see call() function */
//    REDIS_CALL_NONE 0
//    REDIS_CALL_SLOWLOG 1
//    REDIS_CALL_STATS 2
//    REDIS_CALL_PROPAGATE 4
//    REDIS_CALL_FULL (REDIS_CALL_SLOWLOG | REDIS_CALL_STATS | REDIS_CALL_PROPAGATE)
//
//    /* Command propagation flags, see propagate() function */
//    //REDIS_PROPAGATE_NONE 0
//    //REDIS_PROPAGATE_AOF 1
//    //REDIS_PROPAGATE_REPL 2
//
)
