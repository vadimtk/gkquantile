package gkquantile

import "testing"
import "math/rand"

func Test1(t *testing.T) {
    tt := NewGKSummary(0.025)

    for i := 1; i < 10; i++ {
        tt.Add(float64(i))
    }

    v := tt.Query(0.95)
    if v != 8 {
        t.Error("Expected 8, got ", v)
    }
}

func TestRand(t *testing.T) {
    tt := NewGKSummary(0.025)
    r := rand.New(rand.NewSource(99))

    for i := 1; i < 1000; i++ {
        tt.Add(r.ExpFloat64() * 1000)
    }

    v := tt.Query(0.50)
    if (v - 758.575157) > 0.000001 {
        t.Error("Expected 758.575157, got ", v)
    }
}

func TestRandLen(t *testing.T) {
    tt := NewGKSummary(0.01)
    r := rand.New(rand.NewSource(99))

    for i := 1; i < 1000; i++ {
        tt.Add(r.ExpFloat64() * 1000)
    }

    if len(tt.Items) != 131 {
        t.Error("Expected 131, got ", len(tt.Items))
    }
}

func Benchmark_1mln_025(t *testing.B) {
    tt := NewGKSummary(0.025)
    r := rand.New(rand.NewSource(99))

    for i := 1; i < 1000000; i++ {
        tt.Add(r.ExpFloat64() * 1000)
    }

    _ = tt.Query(0.95)
}

func Benchmark_1mln_010(t *testing.B) {
    tt := NewGKSummary(0.010)
    r := rand.New(rand.NewSource(99))

    for i := 1; i < 1000000; i++ {
        tt.Add(r.ExpFloat64() * 1000)
    }

    _ = tt.Query(0.95)
}

func Benchmark_1mln_001(t *testing.B) {
    tt := NewGKSummary(0.001)
    r := rand.New(rand.NewSource(99))

    for i := 1; i < 1000000; i++ {
        tt.Add(r.ExpFloat64() * 1000)
    }

    _ = tt.Query(0.95)
}
