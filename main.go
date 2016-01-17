package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	var fp *os.File
	var err error

	if len(os.Args) < 2 {
		fp = os.Stdin
	} else {
		fp, err = os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer fp.Close()
	}

	data := make([]float64, 0, 1000)

	p210 := createP2()
	p250 := createP2()
	p290 := createP2()
	p2M := createP2()
	p2E := createP2()
	p2DE := createP2()
	p2D2E := createP2()
	//fmt.Printf("%d %d\n", cap((*p2).Dn), (*p2).Dn[0])
	p210.addQuantile(0.1)
	p250.addQuantile(0.5)
	p290.addQuantile(0.9)
	p2M.addQuantile(0.1)
	p2M.addQuantile(0.5)
	p2M.addQuantile(0.9)
	p2E.addEqualSpacing(10)
	p2DE.addEqualSpacing(20)
	p2D2E.addEqualSpacing(100)

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		l := scanner.Text()
		d, _ := strconv.ParseFloat(l, 64)
		//fmt.Println(l, d)
		p210.add(d)
		p250.add(d)
		p290.add(d)
		p2M.add(d)
		p2E.add(d)
		p2DE.add(d)
		p2D2E.add(d)
		data = append(data, d)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	//fmt.Println(p250.Dn)
	//fmt.Println(p2E.Q)
	//fmt.Printf("%f\n", p210.getResult(0.1))
	//fmt.Printf("%f %f %f\n", p210.getResult(0.1), p250.getResult(0.5), p290.getResult(0.9))
	fmt.Printf("%f %f %f\n", p210.getResult(0.1), p250.getResult(0.5), p290.getResult(0.9))
	fmt.Printf("%f %f %f\n", p2M.getResult(0.1), p2M.getResult(0.5), p2M.getResult(0.9))
	fmt.Printf("%f %f %f\n", p2E.getResult(0.1), p2E.getResult(0.5), p2E.getResult(0.9))
	fmt.Printf("%f %f %f\n", p2DE.getResult(0.1), p2DE.getResult(0.5), p2DE.getResult(0.9))
	fmt.Printf("%f %f %f\n", p2D2E.getResult(0.1), p2D2E.getResult(0.5), p2D2E.getResult(0.9))
	sort.Float64s(data)
	fmt.Printf("%f %f %f\n", data[len(data)/10], data[len(data)/2], data[len(data)-len(data)/10])
}
