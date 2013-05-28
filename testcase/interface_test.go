package testcase

//func TestXxx(*testing.T)
import (
    "testing"
    "mRedis/client"
)

var (
    conn *client.RedisClient
)

func init(){
    conn = client.CreateRedisClient("172.16.252.32",8080)
}

func TestSelect(t *testing.T) {
    conn.SELECT(20)
}

func TestSet(t *testing.T) {
    conn.SET("a","b")
}

func TestGet(t *testing.T) {
    conn.GET("a")
}

func TestAppend(t *testing.T) {
    conn.APPEND("a","ABC")
}

func TestExists(t *testing.T) {
    conn.EXISTS("a")
}

func TestRandom(t *testing.T) {
}
