package main

type sharedObjectsStruct struct{
    crlf            *robj
    ok              *robj
    err             *robj
    emptybulk       *robj
    czero           *robj
    cone            *robj
    cnegone         *robj
    pong            *robj
    space           *robj
    colon           *robj
    nullbulk        *robj
    nullmultibulk   *robj
    queued          *robj
    emptymultibulk  *robj
    wrongtypeerr    *robj
    nokeyerr        *robj
    syntaxerr       *robj
    sameobjecterr   *robj
    outofrangeerr   *robj
    noscripterr     *robj
    loadingerr      *robj
    slowscripterr   *robj
    bgsaveerr       *robj
    masterdownerr   *robj
    roslaveerr      *robj
    oomerr          *robj
    plus            *robj
    messagebulk     *robj
    pmessagebulk    *robj
    subscribebulk   *robj
    unsubscribebulk *robj
    psubscribebulk  *robj
    punsubscribebulk *robj
    del             *robj
    rpop            *robj
    lpop            *robj
    lpush           *robj
    //selects[REDIS_SHARED_SELECT_CMDS],
    //integers[REDIS_SHARED_INTEGERS],
    //mbulkhdr[REDIS_SHARED_BULKHDR_LEN], / "<value>\r\n" /
    //bulkhdr[REDIS_SHARED_BULKHDR_LEN] *robj;  / "$<value>\r\n" /
}
