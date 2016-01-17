package main

import (
	"math"
	"sort"
)

type P2 struct {
	Count int
	Q     []float64
	Dn    []float64
	Np    []float64
	N     []int
}

func createP2() *P2 {
	p2 := P2{
		Count: 0,
	}
	p2.addEndMarkers()
	return &p2
}

func (p *P2) addEndMarkers() {
	p.Q = make([]float64, 2, 5)
	p.Dn = make([]float64, 2, 5)
	p.Np = make([]float64, 2, 5)
	p.N = make([]int, 2, 5)
	p.Dn[0] = 0
	p.Dn[1] = 1

	p.updateMarkers()
}

func (p *P2) updateMarkers() {
	sort.Float64s(p.Dn)

	for i, v := range p.Dn {
		p.Np[i] = float64(len(p.Dn)-1)*v + 1
	}
}

func (p *P2) appendDn(data ...float64) {
	m := len(p.Dn)
	n := m + len(data)
	if n > cap(p.Dn) {
		newDn := make([]float64, n, n*2)
		copy(newDn, p.Dn)
		p.Dn = newDn
		p.Q = make([]float64, n, n*2)
		p.Np = make([]float64, n, n*2)
		p.N = make([]int, n, n*2)
	} else {
		p.Q = p.Q[0:n]
		p.Dn = p.Dn[0:n]
		p.Np = p.Np[0:n]
		p.N = p.N[0:n]
	}
	copy(p.Dn[m:n], data)
}

func (p *P2) addQuantile(quant float64) {
	p.appendDn(quant, quant/2.0, (1.0+quant)/2.0)
	p.updateMarkers()
}

func (p *P2) addEqualSpacing(count int) {
	for i := 1; i < count; i++ {
		p.appendDn(float64(i) / float64(count))
	}
	p.updateMarkers()
}

func sign(d float64) int {
	if d >= 0.0 {
		return 1
	}
	return -1
}

func (p *P2) parabolic(i int, d int) float64 {
	return p.Q[i] +
		float64(d)/float64(p.N[i+1]-p.N[i-1])*
			(float64(p.N[i]-p.N[i-1]+d)*(p.Q[i+1]-p.Q[i])/float64(p.N[i+1]-p.N[i])+
				float64(p.N[i+1]-p.N[i]-d)*(p.Q[i]-p.Q[i-1])/float64(p.N[i]-p.N[i-1]))
	//	return q[ i ] + d / (double)(n[ i + 1 ] - n[ i - 1 ]) * ((n[ i ] - n[ i - 1 ] + d) * (q[ i + 1 ] - q[ i ] ) / (n[ i + 1] - n[ i ] ) + (n[ i + 1 ] - n[ i ] - d) * (q[ i ] - q[ i - 1 ]) / (n[ i ] - n[ i - 1 ]) );
}

func (p *P2) linear(i int, d int) float64 {
	return p.Q[i] + float64(d)*(p.Q[i+d]-p.Q[i])/float64(p.N[i+d]-p.N[i])
	//	return q[ i ] + d * (q[ i + d ] - q[ i ] ) / (n[ i + d ] - n[ i ] );

}

func (p *P2) add(v float64) {
	if p.Count >= len(p.Dn) {
		p.Count++
		var k int

		// B1
		if v < p.Q[0] {
			p.Q[0] = v
			k = 1
		} else if v >= p.Q[len(p.Q)-1] {
			k = len(p.Q) - 1
			p.Q[len(p.Q)-1] = v
		} else {
			for i := range p.Q {
				if v < p.Q[i] {
					k = i
					break
				}
			}
		}

		// B2
		for i := range p.Np {
			p.Np[i] += p.Dn[i]
			if i >= k {
				p.N[i]++
			}
		}

		// B3
		for i := range p.Dn {
			d := p.Np[i] - float64(p.N[i])
			if (d >= 1.0 && p.N[i+1]-p.N[i] > 1) ||
				(d <= -1.0 && p.N[i-1]-p.N[i] < -1.0) {
				newQ := p.parabolic(i, sign(d))
				if p.Q[i-1] < newQ && newQ < p.Q[i+1] {
					p.Q[i] = newQ
				} else {
					p.Q[i] = p.linear(i, sign(d))
				}
				p.N[i] += sign(d)
			}
		}
	} else {
		p.Q[p.Count] = v
		p.Count++

		if p.Count == len(p.Dn) {
			sort.Float64s(p.Q)
			for i := range p.N {
				p.N[i] = i + 1
			}
		}
	}
}

func (p P2) getResult(quantile float64) float64 {
	closest := 1

	if p.Count < len(p.Dn) {
		sort.Float64s(p.Q)
		for i := 2; i < p.Count; i++ {
			if math.Abs(float64(i)/float64(p.Count)-quantile) < math.Abs(float64(closest)/float64(p.Count)-quantile) {
				closest = i
			}
		}
	} else {
		for i := 2; i < len(p.Dn); i++ {
			if math.Abs(p.Dn[i]-quantile) < math.Abs(p.Dn[closest]-quantile) {
				closest = i
			}
		}
	}

	return p.Q[closest]
}
