package main

import "fmt"

type segment struct {
	length int64
	data   []int64
}

func newSegment(length int64) *segment {
	s := &segment{length, make([]int64, 1, 8)}
	s.data[0] = length
	return s
}

func (s *segment) insert(from int64, size int64) (large bool) {
	large = false
	if size <= 0 {
		panic("invalid argment")
	}
	if from > s.length {
		panic("invalid argment")
	}
	if from < 0 {
		panic("invalid argment")
	}
	end := from + size
	if end > s.length {
		panic("invalid argment")
	}
	var i int
	for i = 0; i < len(s.data); i++ {
		if from <= s.data[i] {
			break
		}
	}
	if i%2 == 0 {
		if end <= s.data[i] {
			if cap(s.data) < len(s.data)+2 {
				t := make([]int64, len(s.data)+2)
				copy(t, s.data)
				s.data = t
			} else {
				s.data = s.data[:len(s.data)+2]
			}
			copy(s.data[i+2:], s.data[i:])
			s.data[i] = from
			s.data[i+1] = end
		} else {
			s.data[i] = from
			if end <= s.data[i+1] {
				return
			} else if end <= s.data[i+2] {
				s.data[i+1] = end
				return
			} else {
				var j int
				for j = i + 3; j < len(s.data); j++ {
					if end <= s.data[j] {
						break
					}
				}
				i++
				if j%2 == 0 {
					s.data[i] = end
					i++
				}
				copy(s.data[i:], s.data[j:])
				s.data = s.data[:len(s.data)-(j-i)]
			}
		}
	} else {
		if end <= s.data[i] {
			return
		}
		if end <= s.data[i+1] {
			s.data[i] = end
			return
		}
		var j int
		for j = i + 2; j < len(s.data); j++ {
			if end <= s.data[j] {
				break
			}
		}
		if j%2 == 0 {
			s.data[i] = end
			i++
		}
		copy(s.data[i:], s.data[j:])
		s.data = s.data[:len(s.data)-(j-i)]
	}

	return
}

func (s *segment) white() int64 {
	t := s.data[0]
	for i := 2; i < len(s.data); i += 2 {
		t += s.data[i] - s.data[i-1]
	}
	return t
}

func (s *segment) dump() {
	fmt.Printf("[0,%d) ", s.data[0])
	for i := 2; i < len(s.data); i += 2 {
		fmt.Printf("[%d,%d) ", s.data[i-1], s.data[i])
	}
	fmt.Printf("%d\n", s.white())
}

func (s *segment) ins(from int64, end int64) {
	fmt.Printf("insert %d,%d ", from, end)
	s.insert(from, end-from)
	s.dump()
}

//func main() {
//	s := newSegment(100)
//	s.ins(0,100)
//	s.ins(20,50)
//	s.ins(30,40)
//	s.ins(0,10)
//	s.ins(45,60)
//	s.ins(70,80)
//	s.ins(5,75)
//	s.ins(15,75)
//}
