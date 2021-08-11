package set

import (
	"math/rand"
	"testing"
)

type comp struct {
	n int
}

func (c *comp) Compare(other Comparable) int {
	o, ok := other.(*comp)
	if ok {
		return c.n - o.n
	}
	return 0
}

var _ Comparable = &comp{n: 0}

func TestAddAndGet(t *testing.T) {
	var comps []*comp = []*comp{
		{n: 0},
		{n: 1},
		{n: 2},
	}

	set := New()
	for i := range comps {
		_, ok := set.Get(i)
		if ok {
			t.Errorf("Get Not Add")
		}
		set.Add(i, comps[i])
		_, ok = set.Get(i)
		if !ok {
			t.Errorf("Not Get Add")
		}
	}
	set.Add(0, comps[2])
	if c, ok := set.Get(0); !ok {
		t.Fail()
	} else {
		cc := c.(*comp)
		if cc.n != 2 {
			t.Errorf("Not Updated")
		}
	}
	set.Add(2, comps[0])
	if c, ok := set.Get(2); !ok {
		t.Fail()
	} else {
		cc := c.(*comp)
		if cc.n != 2 {
			t.Errorf("Updated Older Value")
		}
	}
}

func BenchmarkSet(b *testing.B) {
	set := New()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			k := rand.Intn(20)
			rand.Float32()
			if rand.Float32() > 0.5 {
				v := &comp{rand.Int()}
				set.Add(k, v)
			} else {
				set.Get(k)
			}
		}
	})
}
