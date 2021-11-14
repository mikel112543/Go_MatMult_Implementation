package main

import (
	"fmt"
	"sync"
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

func addMatrix(A Matrix, B Matrix) [][]int {
	var row, col int
	row = rowCount(A)
	col = colCount(A)
	C := newMatrix(row, col)

	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			C[i][j] = A[i][j] + B[i][j]
		}
	}
	return C
}

func subtractMatrix(A Matrix, B Matrix) [][]int {
	var row, col int
	row = rowCount(A)
	col = colCount(A)
	C := newMatrix(row, col)

	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			C[i][j] = A[i][j] - B[i][j]
		}
	}
	return C
}

func splitMatrix(matrix Matrix, rowStart int, rowEnd int, colStart int, colEnd int, ret chan [][]int, splitWg *sync.WaitGroup) {
	var i, j, p, q int
	rA := rowCount(matrix)
	cA := rowCount(matrix)
	halfrA := rA / 2
	halfcA := cA / 2

	c := newMatrix(halfrA, halfcA)
	//Partitions matrix into 4 sub-matrices
	for i, p = rowStart, 0; i < rowEnd || p < halfrA; i, p = i+1, p+1 {
		for j, q = colStart, 0; j < colEnd; j, q = j+1, q+1 {
			c[p][q] = matrix[i][j]
		}
	}
	splitWg.Done()
	ret <- c //Pass resulting Matrix into the channel
}

func combineMatrices(subMat [][]int, res [][]int, rowStart int, colStart int) [][]int {
	var i, j, p, q int
	size := rowCount(subMat)
	for i, p = 0, rowStart; i < size; i, p = i+1, p+1 {
		for j, q = 0, colStart; j < size; j, q = j+1, q+1 {
			res[p][q] = subMat[i][j]
		}
	}
	return res
}

func doCalc(inA Matrix, inB Matrix) [][]int {
	var i, j int
	m := rowCount(inA) // number of rows the first matrix
	//   n := colCount(inA)     // number of columns the first matrix
	p := rowCount(inB) // number of rows the second matrix
	q := colCount(inB) // number of columns the second matrix
	k := 0
	total := 0
	res := newMatrix(m, q)

	if p <= 2 {
		C := newMatrix(m, q) // create new matrix

		//time.Sleep(1111 * time.Millisecond) // just to max sure timer works delete later

		for i = 0; i < m; i++ {
			for j = 0; j < q; j++ {
				for k = 0; k < p; k++ {
					total = total + inA[i][k]*inB[k][j]
					//      fmt.Print("(", inA[i][k], " * ", inB[k][j], ") + ")
				}
				//          fmt.Println("giving", total)
				C[i][j] = total
				total = 0
			}
		}
		//calWg.Done()
		return C
	} else {
		//Set up wait group
		splitWg := new(sync.WaitGroup)
		splitWg.Add(8)

		//Create separate channels as order of execution matters
		returnA11 := make(chan [][]int)
		returnA12 := make(chan [][]int)
		returnA21 := make(chan [][]int)
		returnA22 := make(chan [][]int)

		returnB11 := make(chan [][]int)
		returnB12 := make(chan [][]int)
		returnB21 := make(chan [][]int)
		returnB22 := make(chan [][]int)

		//7 channels for Strassen equation outputs
		/*calReturnp1 := make(chan [][]int)
		calReturnp2 := make(chan [][]int)
		calReturnp3 := make(chan [][]int)
		calReturnp4 := make(chan [][]int)
		calReturnp5 := make(chan [][]int)
		calReturnp6 := make(chan [][]int)
		calReturnp7 := make(chan [][]int)*/

		//Split Matrix A into 4 sub-matrices
		go splitMatrix(inA, 0, rowCount(inA)/2, 0, colCount(inA)/2, returnA11, splitWg)
		go splitMatrix(inA, 0, rowCount(inA)/2, rowCount(inA)/2, rowCount(inA), returnA12, splitWg)
		go splitMatrix(inA, rowCount(inA)/2, rowCount(inA), 0, colCount(inA)/2, returnA21, splitWg)
		go splitMatrix(inA, rowCount(inA)/2, rowCount(inA), rowCount(inA)/2, colCount(inA), returnA22, splitWg)

		//splitMatrixA := StrassenMatrix{{<-returnA11, <-returnA12},{<-returnA21, <-returnA22}}

		//Split Matrix B into 4 sub-matrices
		go splitMatrix(inB, 0, rowCount(inB)/2, 0, colCount(inB)/2, returnB11, splitWg)
		go splitMatrix(inB, 0, rowCount(inB)/2, rowCount(inB)/2, rowCount(inB), returnB12, splitWg)
		go splitMatrix(inB, rowCount(inB)/2, rowCount(inB), 0, colCount(inB)/2, returnB21, splitWg)
		go splitMatrix(inB, rowCount(inB)/2, rowCount(inB), rowCount(inB)/2, colCount(inB), returnB22, splitWg)

		//Wait for split routines to complete
		splitWg.Wait()
		//Recursively call doCall 8 times to multiply each of the sub-matrices

		//Let channel values equal variables for multiple use
		a := <-returnA11
		b := <-returnA12
		c := <-returnA21
		d := <-returnA22

		e := <-returnB11
		f := <-returnB12
		g := <-returnB21
		h := <-returnB22

		//Strassen's 7 equations
		/**
		  p1 = (a + d)(e + h)
		  p2 = (c + d)e
		  p3 = a(f - h)
		  p4 = d(g - e)
		  p5 = (a + b)h
		  p6 = (c - a) (e + f)
		  p7 = (b - d) (g + h)
		**/
		/*var calWg sync.WaitGroup
		calWg.Add(7)*/

		p1 := doCalc(addMatrix(a, d), addMatrix(e, h)) //p1
		p2 := doCalc(addMatrix(c, d), e)               //p2
		p3 := doCalc(a, subtractMatrix(f, h))          //p3
		p4 := doCalc(d, subtractMatrix(g, e))          //p4 doCalc(addMatrix(a, b), h, &calWg, calReturnp5)                    //p5
		p5 := doCalc(addMatrix(a, b), h)
		p6 := doCalc(subtractMatrix(c, a), addMatrix(e, f)) //p6
		p7 := doCalc(subtractMatrix(b, d), addMatrix(g, h)) //p7

		//time.Sleep(5* time.Second)
		/*
			calWg.Wait()

			p1 := <-calReturnp1
			p2 := <-calReturnp2
			p3 := <-calReturnp3
			p4 := <-calReturnp4
			p5 := <-calReturnp5
			p6 := <-calReturnp6
			p7 := <-calReturnp7*/

		/**
		  C11 = p1 + p4 - p5 + p7
		  C12 = p3 + p5
		  C21 = p2 + p4
		  C22 = p1 - p2 + p3 + p6
		**/

		C11 := addMatrix(subtractMatrix(addMatrix(p1, p4), p5), p7)
		C12 := addMatrix(p3, p5)
		C21 := addMatrix(p2, p4)
		C22 := addMatrix(subtractMatrix(addMatrix(p1, p3), p2), p6)

		//Combine Submatrices into 1
		combineMatrices(C11, res, 0, 0)
		combineMatrices(C12, res, 0, m/2)
		combineMatrices(C21, res, m/2, 0)
		combineMatrices(C22, res, m/2, m/2)

	}
	return res
}

/*func doCalc(inA Matrix, inB Matrix) [][]int {
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
}*/

func main() {

	//calReturnMatrix := make(chan [][]int)
	start := time.Now()
	calgWReturn := new(sync.WaitGroup)
	calgWReturn.Add(8)
	//splitWg := new(sync.WaitGroup)
	//Create Wait group for split go routines
	//
	// Use slices
	// Unlike arrays they are passed by reference,not by value
	/*	a := Matrix{{2, 3, 6, 4, 2, 3, 6, 4}, {5, 6, 4, 23, 5, 6, 4, 23},
						{9, 6, 12, 23, 5, 6, 4, 23},{4, 7, 12, 43, 5, 6, 4, 23},
							{4, 7, 12, 43, 5, 6, 4, 23}, {9, 6, 12, 23, 5, 6, 4, 23},
								{9, 6, 12, 23, 5, 6, 4, 23}, {5, 6, 4, 23, 5, 6, 4, 23}}

		b:= Matrix{{2, 3, 6, 4, 2, 3, 6, 4}, {5, 6, 4, 23, 5, 6, 4, 23},
						{9, 6, 12, 23, 5, 6, 4, 23},{4, 7, 12, 43, 5, 6, 4, 23},
							{4, 7, 12, 43, 5, 6, 4, 23}, {9, 6, 12, 23, 5, 6, 4, 23},
								{9, 6, 12, 23, 5, 6, 4, 23}, {5, 6, 4, 23, 5, 6, 4, 23}}*/

	a := Matrix{{8, 18, 28, 14}, {38, 48, 58, 12},
		{24, 56, 78, 34}, {12, 54, 76, 43}}

	b := Matrix{{8, 18, 28, 14}, {38, 48, 58, 12},
		{24, 56, 78, 34}, {12, 54, 76, 43}}

	fmt.Println("Matrix A")
	fmt.Println(" Number of cols in A ", colCount(a))
	printMat(a)

	fmt.Println("Matrix B")
	fmt.Println(" Number of rows in B ", rowCount(b))
	printMat(b)

	fmt.Println("Matrix Split")

	mat := doCalc(a, b)
	printMat(mat)

	elapsed := time.Since(start)
	fmt.Printf("Time taken to calculate %s ", elapsed)
}
