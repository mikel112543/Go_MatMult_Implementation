package main

import (
	"fmt"
	"time"
)

type Matrix [][]int

func printMat(inM Matrix) {

	for _, i := range inM {
		for _, j := range i {
			fmt.Print(" ", j)
		}
		fmt.Println()
	}
}

func rowCount(inM Matrix) int {
	return len(inM)
}

func colCount(inM Matrix) int {
	return len(inM[0])
}

func newMatrix(r, c int) [][]int {
	a := make([]int, c*r)
	m := make([][]int, r)
	lo, hi := 0, c
	for i := range m {
		m[i] = a[lo:hi:hi]
		lo, hi = hi, hi+c
	}
	return m
}

func splitMatrix(inA Matrix, rowStart int, rowEnd int, colStart int, colEnd int) [][]int {
	var i, j, p, q int
	rA := rowCount(inA)
	cA := rowCount(inA)
	p = 0
	q = 0
	halfrA := rA / 2
	halfcA := cA / 2

	c := newMatrix(halfrA, halfcA)
	//Currently does one separation of 4x4 matrix
	for i, p = 0, 0; i < halfrA || p < halfrA; i, p = i+1, p+1 {
		for j, q = 0, 0; j < halfcA; j, q = j+1, q+1 {
			fmt.Println(inA[i][j])
			c[p][q] = inA[i][j]
		}

	}
	return c
}

func doCalc(inA Matrix, inB Matrix) [][]int {
	var i, j int
	m := rowCount(inA) // number of rows the first matrix
	//   n := colCount(inA)     // number of columns the first matrix
	p := rowCount(inB) // number of rows the second matrix
	q := colCount(inB) // number of columns the second matrix
	k := 0
	total := 0

	nM := newMatrix(m, q) // create new matrix

	start := time.Now()
	time.Sleep(1111 * time.Millisecond) // just to max sure timer works delete later

	for i = 0; i < m; i++ {
		for j = 0; j < q; j++ {
			for k = 0; k < p; k++ {
				total = total + inA[i][k]*inB[k][j]
				//      fmt.Print("(", inA[i][k], " * ", inB[k][j], ") + ")
			}
			//          fmt.Println("giving", total)
			nM[i][j] = total
			total = 0
		}
		fmt.Println()
	}
	elapsed := time.Since(start)
	fmt.Printf("Time taken to calculate %s ", elapsed)
	return nM
}

func main() {
	// Use slices
	// Unlike arrays they are passed by reference,not by value
	a := Matrix{{2, 3}, {5, 6}, {9, 6}, {4, 7}}
	b := Matrix{{8, 18, 28, 14}, {38, 48, 58, 12}, {24, 56, 78, 34}, {12, 54, 76, 43}}

	fmt.Println("Matrix A")
	fmt.Println(" Number of cols in A ", colCount(a))
	printMat(a)

	fmt.Println("Matrix B")
	fmt.Println(" Number of rows in B ", rowCount(b))
	printMat(b)

	fmt.Println("Matrix SPlit")
	t := splitMatrix(b, 1, (rowCount(a)/2)-1, 1, (colCount(a)/2)-1)
	fmt.Println(t)
	printMat(t)
	/*
		fmt.Println("The Go Result of Matrix Multiplication = ")
		c := doCalc(a, b)
		printMat(c)*/
}
