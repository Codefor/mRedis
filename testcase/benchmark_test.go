package testcase

//func TestXxx(*testing.T)
import (
    "testing"
)

func BenchmarkSet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        conn.SET("a","b")
    }
}

func BenchmarkGet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        conn.GET("a")
    }
}
