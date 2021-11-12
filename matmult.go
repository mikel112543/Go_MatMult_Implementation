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

func doCalc(inA Matrix, inB Matrix, calWg *sync.WaitGroup, ret chan [][]int) {
	defer calWg.Done()
	var i, j int
	m := rowCount(inA) // number of rows the first matrix
	//   n := colCount(inA)     // number of columns the first matrix
	p := rowCount(inB) // number of rows the second matrix
	q := colCount(inB) // number of columns the second matrix
	k := 0
	total := 0

	if p <= 2 {
		C := newMatrix(m, q) // create new matrix

		start := time.Now()
		time.Sleep(1111 * time.Millisecond) // just to max sure timer works delete later

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
			fmt.Println()
		}
		elapsed := time.Since(start)
		fmt.Printf("Time taken to calculate %s ", elapsed)
		ret <- C
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

		calReturnC1 := make(chan [][]int)
		calReturnC2 := make(chan [][]int)

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

		//Let channel values equal variables for multiple use
		A11 := <-returnA11
		A12 := <-returnA12
		/*A21 := <-returnA21
		A22 := <-returnA22*/

		B11 := <-returnB11
		/*B12 := <-returnB12*/
		B21 := <-returnB21
		/*B22 := <-returnB22*/

		//Wait for split routines to complete
		splitWg.Wait()
		//Recursively call doCall 8 times to multiply each of the sub-matrices
		go doCalc(A11, B11, calWg, calReturnC1)
		go doCalc(A12, B21, calWg, calReturnC2)
		/*		go doCalc(A11, B12, calWg, calReturnC2)
				go doCalc(A12, B22, calWg, calReturnC2)
				go doCalc(A21, B11, calWg, calReturnC2)
				go doCalc(A22, B21, calWg, calReturnC2)
				go doCalc(A21, B12, calWg, calReturnC2)
				go doCalc(A22, B22, calWg, calReturnC2)*/

		_A := <-calReturnC1
		_B := <-calReturnC2

		C11 := addMatrix(_A, _B)
		printMat(C11)
		/*
			A11B12 := doCalc(A11, B12)
			A12B22 := doCalc(A12, B22)

			A21B11 := doCalc(A21, B11)
			A22B21 := doCalc(A22, B21)

			A21B12 := doCalc(A21, B12)
			A22B22 := doCalc(A22, B22)

			C11 := addMatrix(A11B11, A12B21)
			C12 := addMatrix(A11B12, A12B22)
			C21 := addMatrix(A21B11, A22B21)
			C22 := addMatrix(A21B12, A22B22)

			fmt.Println("C11 =")
			printMat(C11)
			fmt.Println("C12 =")
			printMat(C12)
			fmt.Println("C21 =")
			printMat(C21)
			fmt.Println("C22 =")
			printMat(C22)*/
	}
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

	calReturnMatrix := make(chan [][]int)
	calWg := new(sync.WaitGroup)
	calWg.Add(3)
	start := time.Now()
	//splitWg := new(sync.WaitGroup)
	//Create Wait group for split go routines
	//
	// Use slices
	// Unlike arrays they are passed by reference,not by value
	a := Matrix{{2, 3, 6, 4}, {5, 6, 4, 23}, {9, 6, 12, 23}, {4, 7, 12, 43}}
	b := Matrix{{8, 18, 28, 14}, {38, 48, 58, 12}, {24, 56, 78, 34}, {12, 54, 76, 43}}

	fmt.Println("Matrix A")
	fmt.Println(" Number of cols in A ", colCount(a))
	printMat(a)

	fmt.Println("Matrix B")
	fmt.Println(" Number of rows in B ", rowCount(b))
	printMat(b)

	fmt.Println("Matrix Split")

	doCalc(a, b, calWg, calReturnMatrix)
	calWg.Wait()

	//Separate Channels for the Go routines as the order of execution matters
	/*returnA11 := make(chan [][]int)
	returnA12 := make(chan [][]int)
	returnA21 := make(chan [][]int)
	returnA22 := make(chan [][]int)
	returnB11 := make(chan [][]int)
	returnB12 := make(chan [][]int)
	returnB21 := make(chan [][]int)
	returnB22 := make(chan [][]int)*/

	//Go routine for each sub matrix
	//Split Matrix A

	stop := time.Since(start)
	//Wait for Go routines to finish
	//splitWg.Wait()
	//Close Channels once routines are finished

	fmt.Println("Time taken to split matrices sequentially", stop)
	/*mat1 := <-returnA11
	printMat(mat1)*/
	/*
		x := 1
		fmt.Println("----------------Matrix A Split-------------------")
		for mat := range returnA {
			fmt.Println("Matrix", x)
			printMat(mat)
			x++
		}*/
	/*	var tempA [4]Matrix
		x := 0
		for mat := range returnA {
			tempA[x] = mat
			x++
		}
		for i:=0; i < len(tempA); i++ {
			printMat(tempA[i])
		}
		fmt.Println("MATRIX 1")
		fmt.Println(tempA[1])*/
	/*	for i :=0; i < 4; i++ {
			mat := <-returnA
			tempA[i] = mat
		}
		for i := 0; i < len(tempA); i++ {
			printMat(tempA[i])
		}*/
	/*
		fmt.Println("---------------Matrix B Split--------------------")
		for mat := range returnB {
			fmt.Println("Matrix")
			printMat(mat)
		}*/

	/*	printMat(<-returnA11)
		fmt.Println("A12")
		printMat(<-returnA12)
		fmt.Println("A21")
		printMat(<-returnA21)
		fmt.Println("A22")
		printMat(<-returnA22)*/

	/*
		fmt.Println("The Go Result of Matrix Multiplication = ")
		c := doCalc(a, b)
		printMat(c)*/
}
