package gkquantile

import "fmt"
import "math"
import "sort"

/* This work is based on
"Space-Efficient Online Computation of Quantile Summaries"
by M. Greenwald and S. Khanna (known as GK-algorithm)
http://infolab.stanford.edu/~datar/courses/cs361a/papers/quantiles.pdf

Also ideas for the implementation are taken from
http://www.mathcs.emory.edu/~cheung/Courses/584-StreamDB/Syllabus/08-Quantile/Greenwald.html
*/

// === GK implementation ===========================

type GKItem struct {
	value float64
	g     int
	delta int
}

type GKSummary struct {
	Items   []*GKItem
	count   uint
	epsilon float64
}

func NewGKSummary(eps float64) *GKSummary {
	td := new(GKSummary)
	td.count = 0
	td.epsilon = eps
	td.Items = make([]*GKItem, 0)
	return td
}

func (t *GKSummary) Add(val float64) {

	// We maintain t.Items sorted, so can use a binary search
	idx := sort.Search(len(t.Items), func(i int) bool { return t.Items[i].value >= val })

	if idx > 0 {
		for idx < len(t.Items) && t.Items[idx].value == val {
			//  find first index where value > t.Items[i]
			idx++
		}
	}

	delta := 0
	// for first and last index - delta is always 0
	if idx != 0 && idx != len(t.Items) {
		delta = t.Items[idx].g + t.Items[idx].delta - 1
		// sanity check, if condition is true something went wrong. this should never happen
		if delta > int(math.Floor(t.epsilon*float64(t.count*2))) {
			panic("Delta is greater than allowable error. This never should happen. Really.")
		}
	}

	nItem := &GKItem{value: val, g: 1, delta: delta}
	t.Items = append(t.Items, nItem)
	// insert item into position idx, unless it was the last idx
	if idx != len(t.Items)-1 {
		copy(t.Items[idx+1:], t.Items[idx:])
		t.Items[idx] = nItem
	}

	t.count += 1

	// Periodically perform compress
	if int(math.Mod(float64(t.count), math.Floor(1/(2*t.epsilon)))) == 0 {
		t.Compress()
	}
}

func (t *GKSummary) Compress() {
	/* original logic from
	   http://www.mathcs.emory.edu/~cheung/Courses/584-StreamDB/Syllabus/08-Quantile/Greenwald.html
	       for ( i = s-1; i = 2; i = j - 1 )
	      {
	         j = i-1;

	         while ( j = 1 && gj + ... + gi + ?i < 2eN )
	     {
	        j--;
	     }

	     j++;        // We went one index too far in the while...

	         if ( j < i )
	     {
	        replace entries j, .., i with the entry (vi, gj+ ... + gi, ?i);
	     }
	      }
	*/
	j := 0
	for i := len(t.Items) - 1; i >= 2; i = j - 1 {
		j = i - 1
		rollSum := t.Items[i].g
		for {
			if j < 1 {
				break
			}
			if rollSum+t.Items[j].g+t.Items[i].delta > int(math.Floor(t.epsilon*float64(t.count*2))) {
				// error too big, we stop now
				break
			}
			//fmt.Println("Merge items", i, j)
			rollSum = rollSum + t.Items[j].g
			t.Items[j] = nil
			j = j - 1
		}
		j = j + 1 // We went one index too far in the while...

		if j < i {
			// remove all items from j to i-1
			// and use i as a new item
			t.Items[i].g = rollSum
			copy(t.Items[j:], t.Items[i:])
			t.Items = t.Items[:len(t.Items)-(i-j)]
		}
	}
	return
}

func (t *GKSummary) Query(q float64) float64 {
	rankMin := 0
	if q == 0 {
		return t.Items[0].value
	}
	if q == 1 {
		return t.Items[len(t.Items)-1].value
	}
	desired := int(q * (float64(t.count)))
	for i := 1; i < len(t.Items); i++ {
		prev := t.Items[i-1]
		curr := t.Items[i]
		rankMin += prev.g
		if rankMin+curr.g+curr.delta > desired+int(t.epsilon*float64(t.count)) {
			return prev.value
		}
	}
	return t.Items[len(t.Items)-1].value
}

func (t *GKSummary) QueryRank(q float64) (float64, int, int) {
	rankMin := 0
	if q == 0 {
		return t.Items[0].value, t.Items[0].g, t.Items[0].g + t.Items[0].delta
	}

	desired := int(q * (float64(t.count)))

	for i := 1; i < len(t.Items); i++ {
		prev := t.Items[i-1]
		curr := t.Items[i]
		rankMin += prev.g
		if rankMin+curr.g+curr.delta > desired+int(t.epsilon*float64(t.count)) {
			return prev.value, rankMin, rankMin + prev.delta
		}
	}
	return t.Items[len(t.Items)-1].value, rankMin + t.Items[len(t.Items)-1].g, rankMin + t.Items[len(t.Items)-1].g + t.Items[len(t.Items)-1].delta
}

func (t *GKSummary) Print() {
	rank := 0
	for i, v := range t.Items { // loop through all values
		rank = rank + v.g
		fmt.Printf("el: %d v: %f, range (%d-%d), g,d: %d, %d, maxerr: %d\n", i, v.value, rank, rank+v.delta, v.g, v.delta, (v.g+v.delta)/2)
	}
	fmt.Printf("Error range: %d\n", int(math.Floor(t.epsilon*float64(t.count))))
	return
}
