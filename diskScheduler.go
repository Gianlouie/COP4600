/* I Gianlouie Molinary (gi713278) affirm that this program is entirely my own work and
that I have neither developed my code together with any other person, nor copied any code
from any other person, nor permitted my code to be copied or otherwise used by any other 
person, nor have I copied, modified, or otherwise used programs created by others.
I acknowledge that any violation of the above terms will be treated as academic dishonesty.
*/

package main

import (
	"os"
	"fmt"
	"log"
	"bufio"
	"strconv"
	"sort"
)

var algo Algorithm; 

func main () {

	if len(os.Args) < 2 { // the program has to be built this way with at least 2 arguments
		fmt.Println("Invalid command, should be <program name> <input file>");
	}

	input, err := os.Open(os.Args[1]);

	if err != nil {
		log.Fatal(err);
	}
	
	i := 0;

	scanner := bufio.NewScanner(input);
	scanner.Split(bufio.ScanWords);

	for scanner.Scan() {
		if scanner.Text() == "use" {
			scanner.Scan();
			algo.style = scanner.Text();
		}

		if scanner.Text() == "lowerCYL" {
			scanner.Scan();
			algo.lowerCYL, _ = strconv.Atoi(scanner.Text());
		}

		if scanner.Text() == "upperCYL" {
			scanner.Scan();
			algo.upperCYL, _ = strconv.Atoi(scanner.Text());
		}

		if scanner.Text() == "initCYL" {
			scanner.Scan();
			algo.initCYL, _ = strconv.Atoi(scanner.Text());

			// should be the one directly before the cylinder requests
			scanner.Scan();

			for scanner.Text() != "end" { // grabs the amount of cylinder requests
				scanner.Scan();
				algo.cylreq[i], _ = strconv.Atoi(scanner.Text());
				scanner.Scan();
				i++;
			}

			algo.numOfCylReq = i;
		}
	}

	switch algo.style {
		case "fcfs": FirstComeFirstServed();
		case "sstf": ShortestSeekTimeFirst();
		case "scan": SCAN();
		case "c-scan": C_SCAN();
		case "look": LOOK();
		case "c-look": C_LOOK(); 
	}	
}

// Had to make my own abs function since go lang doesnt inherently support absolute value for ints in the math library
func Abs(x int) int { // Yeah it's pretty wack I know
	if x < 0 {
		return -x;
	}

	return x;
}

func traversalCount(arr []int) int { // Of course computes the traversal count of a given array

	traversal := 0;

	current := algo.initCYL;

	for i := 1; i <= len(arr); i++ {
		traversal += Abs((current - arr[i-1]));
		current = arr[i-1];
	}

	return traversal;
}


func FirstComeFirstServed() {

	fmt.Println("Seek algorithm: FCFS");
	fmt.Printf("\tLower cylinder: %d\n", algo.lowerCYL);
	fmt.Printf("\tUpper cylinder: %d\n", algo.upperCYL);
	fmt.Printf("\tInit cylinder: %d\n", algo.initCYL);

	fmt.Println("\tCylinder requests:");
	for i := 0; i < algo.numOfCylReq; i++ {
		fmt.Printf("\t\tCylinder %d\n", algo.cylreq[i]);
	}

	// We service them by which ever arrives first so just leave the array as is
	for i := 0; i < algo.numOfCylReq; i++ {
		fmt.Printf("Servicing %5d\n", algo.cylreq[i]);
	}

	newcylreq := algo.cylreq[0:algo.numOfCylReq];

	fmt.Printf("FCFS traversal count = %d\n", traversalCount(newcylreq));
}

func ShortestSeekTimeFirst() {

	newcylreq := algo.cylreq[0:algo.numOfCylReq];
	current := algo.initCYL;
	var shortest int;
	var maybe int;
	var shortestDistance int;

	fmt.Println("Seek algorithm: SSTF");
	fmt.Printf("\tLower cylinder: %d\n", algo.lowerCYL);
	fmt.Printf("\tUpper cylinder: %d\n", algo.upperCYL);
	fmt.Printf("\tInit cylinder: %d\n", algo.initCYL);

	fmt.Println("\tCylinder requests:");
	for i := 0; i < algo.numOfCylReq; i++ {
		fmt.Printf("\t\tCylinder %d\n", algo.cylreq[i]);
	}

	sort.Ints(newcylreq);

	a := make([]int, 0);

	for i := 0; i < algo.numOfCylReq; i++ {
		shortest = algo.upperCYL;
		
		for j := 0; j < algo.numOfCylReq; j++ {
			maybe = Abs(current - newcylreq[j]);
	
			if maybe < shortest && maybe != 0 {
				shortest = maybe;
				shortestDistance = j;
			}
		}

		a = append(a, newcylreq[shortestDistance]);
		current = newcylreq[shortestDistance];
		newcylreq[shortestDistance] = -algo.upperCYL; // this is to avoid it going back forward between the same values
	}

	for i := 0; i < algo.numOfCylReq; i++ {
		fmt.Printf("Servicing %d\n", a[i]);
	}
	
	fmt.Printf("SSTF traversal count = %d\n", traversalCount(a));
}

func SCAN() {

	newcylreq := algo.cylreq[0:algo.numOfCylReq];
	shortest := algo.upperCYL;
	var shortestDistance int;
	var maybe int; 

	fmt.Println("Seek algorithm: SCAN");
	fmt.Printf("\tLower cylinder: %d\n", algo.lowerCYL);
	fmt.Printf("\tUpper cylinder: %d\n", algo.upperCYL);
	fmt.Printf("\tInit cylinder: %d\n", algo.initCYL);

	fmt.Println("\tCylinder requests:");
	for i := 0; i < algo.numOfCylReq; i++ {
		fmt.Printf("\t\tCylinder %d\n", algo.cylreq[i]);
	}

	sort.Ints(newcylreq);
	
	for j := 0; j < algo.numOfCylReq; j++ {
		maybe = Abs(algo.initCYL - newcylreq[j]);

		if maybe < shortest && maybe != 0 {
			shortest = maybe;
			shortestDistance = j;
			if newcylreq[shortestDistance] < algo.initCYL {
				shortestDistance = j+1;
			}
		}
	}

	a := make([]int, 0);

	for i := shortestDistance; i < algo.numOfCylReq; i++ {
		fmt.Printf("Servicing %d\n", newcylreq[i]);
		a = append(a, newcylreq[i]);
	}

	if newcylreq[0] < algo.initCYL {
		a = append(a, algo.upperCYL);
	}

	for i := shortestDistance-1; i >= 0; i-- { // go back to service those left behind if any
		fmt.Printf("Servicing %d\n", newcylreq[i]);
		a = append(a, newcylreq[i]);
	}

	fmt.Printf("SCAN traversal count = %d\n", traversalCount(a));
	
}

func C_SCAN() {

	newcylreq := algo.cylreq[0:algo.numOfCylReq];
	shortest := algo.upperCYL;
	var shortestDistance int;
	var maybe int; 

	fmt.Println("Seek algorithm: C-SCAN");
	fmt.Printf("\tLower cylinder: %d\n", algo.lowerCYL);
	fmt.Printf("\tUpper cylinder: %d\n", algo.upperCYL);
	fmt.Printf("\tInit cylinder: %d\n", algo.initCYL);

	fmt.Println("\tCylinder requests:");
	for i := 0; i < algo.numOfCylReq; i++ {
		fmt.Printf("\t\tCylinder %d\n", algo.cylreq[i]);
	}

	sort.Ints(newcylreq);
	
	for j := 0; j < algo.numOfCylReq; j++ {
		maybe = Abs(algo.initCYL - newcylreq[j]);

		if maybe < shortest && maybe != 0 {
			shortest = maybe;
			shortestDistance = j;
			if newcylreq[shortestDistance] < algo.initCYL {
				shortestDistance = j+1;
			}
		}
	}

	a := make([]int, 0);

	for i := shortestDistance; i < algo.numOfCylReq; i++ {
		fmt.Printf("Servicing %d\n", newcylreq[i]);
		a = append(a, newcylreq[i]);
	}

	if newcylreq[0] < algo.initCYL {
		a = append(a, algo.upperCYL);
		a = append(a, 0);
	}

	for i := 0; i < shortestDistance; i++ {
		fmt.Printf("Servicing %d\n", newcylreq[i]);
		a = append(a, newcylreq[i]);
	}

	fmt.Printf("C-SCAN Traversal Count = %d\n", traversalCount(a));
}

func LOOK() {

	newcylreq := algo.cylreq[0:algo.numOfCylReq];
	shortest := algo.upperCYL;
	var shortestDistance int;
	var maybe int; 

	fmt.Println("Seek algorithm: LOOK");
	fmt.Printf("\tLower cylinder: %d\n", algo.lowerCYL);
	fmt.Printf("\tUpper cylinder: %d\n", algo.upperCYL);
	fmt.Printf("\tInit cylinder: %d\n", algo.initCYL);

	fmt.Println("\tCylinder requests:");
	for i := 0; i < algo.numOfCylReq; i++ {
		fmt.Printf("\t\tCylinder %d\n", algo.cylreq[i]);
	}

	sort.Ints(newcylreq);
	
	for j := 0; j < algo.numOfCylReq; j++ {
		maybe = Abs(algo.initCYL - newcylreq[j]);

		if maybe < shortest && maybe != 0 {
			shortest = maybe;
			shortestDistance = j;
			if newcylreq[shortestDistance] < algo.initCYL {
				shortestDistance = j+1;
			}
		}
	}

	a := make([]int, 0);

	for i := shortestDistance; i < algo.numOfCylReq; i++ {
		fmt.Printf("Servicing %d\n", newcylreq[i]);
		a = append(a, newcylreq[i]);
	}

	for i := shortestDistance-1; i >= 0; i-- { // go back to service those left behind if any
		fmt.Printf("Servicing %d\n", newcylreq[i]);
		a = append(a, newcylreq[i]);
	}

	fmt.Printf("LOCK traversal count = %d\n", traversalCount(a));
	
}

func C_LOOK() {

	newcylreq := algo.cylreq[0:algo.numOfCylReq];
	shortest := algo.upperCYL;
	var shortestDistance int;
	var maybe int; 

	fmt.Println("Seek algorithm: C-LOOK");
	fmt.Printf("\tLower cylinder: %d\n", algo.lowerCYL);
	fmt.Printf("\tUpper cylinder: %d\n", algo.upperCYL);
	fmt.Printf("\tInit cylinder: %d\n", algo.initCYL);

	fmt.Println("\tCylinder requests:");
	for i := 0; i < algo.numOfCylReq; i++ {
		fmt.Printf("\t\tCylinder %d\n", algo.cylreq[i]);
	}

	sort.Ints(newcylreq);
	
	for j := 0; j < algo.numOfCylReq; j++ {
		maybe = Abs(algo.initCYL - newcylreq[j]);

		if maybe < shortest && maybe != 0 {
			shortest = maybe;
			shortestDistance = j;
			if newcylreq[shortestDistance] < algo.initCYL {
				shortestDistance = j+1;
			}
		}
	}

	a := make([]int, 0);

	for i := shortestDistance; i < algo.numOfCylReq; i++ {
		fmt.Printf("Servicing %d\n", newcylreq[i]);
		a = append(a, newcylreq[i]);
	}

	for i := 0; i < shortestDistance; i++ {
		fmt.Printf("Servicing %d\n", newcylreq[i]);
		a = append(a, newcylreq[i]);
	}

	fmt.Printf("C-LOCK Traversal count = %d\n", traversalCount(a));
}

type Algorithm struct {
	style string;
	lowerCYL int;
	upperCYL int;
	initCYL int;
	cylreq[20] int;
	numOfCylReq int;
}