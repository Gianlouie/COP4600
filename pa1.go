/* I Gianlouie Molinary (gi713278) affirm that this program is entirely my own work and
that I have neither developed my code together with any other person, nor copied any code
from any other person, nor permitted my code to be copied or otherwise used by any other 
person, nor have I copied, modified, or otherwise used programs created by others.
I acknowledge that any violation of the above terms will be treated as academic dishonesty.
*/

package main

import (
	"log"
	"fmt"
	"os"
	"strings"
	"bufio"
	"strconv"
)

var algo Algorithm;

func main() {

	if len(os.Args) < 3 { // the program has to be built this way with at least 3 arguments
		fmt.Println("Invalid command, should be <program name> <input file> <output file>");
	}

	input, err := os.Open(os.Args[1]);

	if err != nil {
		log.Fatal(err);
	}

	i := 0;

	scanner := bufio.NewScanner(input);
	scanner.Split(bufio.ScanWords); // splits it word by word seperated by spaces thus allowing us to check each word
	for scanner.Scan() { // pretty much does the job of ignoring everything expect the desired information

		if strings.Compare(scanner.Text(), "processcount") == 0 {
			scanner.Scan();
			algo.processcount, _ = strconv.Atoi(scanner.Text());
		}

		if strings.Compare(scanner.Text(), "runfor") == 0 {
			scanner.Scan();
			algo.runfor, _ = strconv.Atoi(scanner.Text());
		}

		if strings.Compare(scanner.Text(), "use") == 0 {
			scanner.Scan();
			algo.style = scanner.Text();
		}
		
		if strings.Compare(scanner.Text(), "name") == 0 {
			scanner.Scan();
			algo.processes[i].name = scanner.Text();
		}

		if strings.Compare(scanner.Text(), "quantum") == 0 {
			scanner.Scan();
			if (scanner.Text() != "-") {
				algo.quantum, _ = strconv.Atoi(scanner.Text());
			}
		}

		if strings.Compare(scanner.Text(), "arrival") == 0 {
			scanner.Scan();
			algo.processes[i].arrival, _ = strconv.Atoi(scanner.Text());
		}

		if strings.Compare(scanner.Text(), "burst") == 0 {
			scanner.Scan();
			algo.processes[i].burst, _ = strconv.Atoi(scanner.Text());
			i++; // only really want to increment the currentProcess when it hits the burst stat for the next process
		}
	}

	// shrinks the array down to the needed size if needed
	var processes []Process = algo.processes[0:algo.processcount];

	switch algo.style {
		case "rr": roundRobin(processes);
		case "fcfs": firstComeFirstServed(processes);
		case "sjf": shortestJobFirst(processes);
	}
}

func sortByArrival(processes []Process) []Process {

	// Sort the processes by arrival
	for i := 0; i < len(processes); i++ {
		min := i;

		for j := i + 1; j < len(processes); j++ {
			if processes[j].arrival < processes[min].arrival {
				min = j;
			}
		}

		temp := processes[min];
		processes[min] = processes[i];
		processes[i] = temp;
	}

	return processes;
}

func sortByName(processes []Process) []Process {

	// Sort the processes by name
	for i := 0; i < len(processes); i++ {
		min := i;

		for j := i + 1; j < len(processes); j++ {
			if processes[j].name < processes[min].name {
				min = j;
			}
		}

		temp := processes[min];
		processes[min] = processes[i];
		processes[i] = temp;
	}

	return processes;
}

func sortByBurst(processes []Process) []Process {

	// Sort the processes by burst
	for i := 0; i < len(processes); i++ {
		min := i;

		for j := i + 1; j < len(processes); j++ {
			if processes[j].burst < processes[min].burst {
				min = j;
			}
		}

		temp := processes[min];
		processes[min] = processes[i];
		processes[i] = temp;
	}

	return processes;
}

func firstComeFirstServed(processes []Process) {

	output, err := os.Create(os.Args[2]);

	if err != nil {
		log.Fatal(err);
	}
	
	fmt.Fprintf(output, "%3d processes\n", algo.processcount);

	fmt.Fprintln(output,"Using First-Come First-Served");

	processes = sortByArrival(processes);

	var Queue []Process;
	Queue = make([]Process, algo.processcount);
	QCap := 0;

	time := 0;
	currentProcess := 0;

	for time < algo.runfor {
		for i := 0; i < algo.processcount; i++ {
			if processes[i].arrival == time {
				fmt.Fprintf(output, "Time %3d : %s arrived\n", time, processes[i].name);
				Queue[i] = processes[i];
				QCap++;
			}
		}

		if QCap == 0 {
			fmt.Fprintf(output, "Time %3d : Idle\n", time);
		}

		if QCap > 0 {
			if Queue[currentProcess].chosen && ((Queue[currentProcess].selection + Queue[currentProcess].burst) == time) {
				Queue[currentProcess].completed = true;
				Queue[currentProcess].completion = time;
				Queue[currentProcess].chosen = false; 
				QCap--;

				fmt.Fprintf(output, "Time %3d : %s finished\n", time, processes[currentProcess].name);

				if (currentProcess < (algo.processcount - 1)) {
					currentProcess++;
				}
			}

			if !Queue[currentProcess].chosen && !Queue[currentProcess].completed && QCap > 0 {
				Queue[currentProcess].chosen = true;
				Queue[currentProcess].selection = time;

				fmt.Fprintf(output, "Time %3d : %s selected (burst %3d)\n", time, processes[currentProcess].name, processes[currentProcess].burst);
			} else if QCap == 0 {
				fmt.Fprintf(output, "Time %3d : Idle\n", time);
			}
		}
		time++;
	}

	processes = Queue;

	for i:= 0; i < algo.processcount; i++ {
		processes[i].turnaround = processes[i].completion - processes[i].arrival;
		processes[i].wait = processes[i].turnaround - processes[i].burst;
	}

	fmt.Fprintf(output, "Finished at time %3d\n\n", algo.runfor);

	processes = sortByName(processes);

	for i := 0; i < algo.processcount; i++ {
		fmt.Fprintf(output, "%s wait %3d turnaround %3d\n", processes[i].name, processes[i].wait,
		processes[i].turnaround);
	}
}

func roundRobin(processes []Process) { 

	output, err := os.Create(os.Args[2]);

	if err != nil {
		log.Fatal(err);
	}

	fmt.Fprintf(output, "%3d processes\n", algo.processcount);
	fmt.Fprintln(output, "Using Round-Robin");
	fmt.Fprintf(output, "Quantum %3d\n\n", algo.quantum);

	processes = sortByArrival(processes);

	var currentProcess Process;

	var currentQuantum int;

	var time int;

	next := 0;

	running := 1;

	finished := 0;

	selected := -1;

	for time := 0; running == 1; time++ {
		for j := 0; j < algo.processcount; j++ {
			if processes[j].arrival == time {
				fmt.Fprintf(output, "Time %3d: %s arrived\n", time, processes[j].name);
				processes[j].arrived = true;
				if selected < 0 {
					currentQuantum = algo.quantum;
					currentProcess = processes[j];
					fmt.Fprintf(output, "Time %3d: %s selected (burst %3d)\n", time, currentProcess.name, currentProcess.burst);
					selected = j;
				}
			}
		}

		if next == 1 {
			searching := 1;
			i := 0;
			peek := selected+1;

			for searching == 1 {
				if peek > (algo.processcount - 1) {
					peek = 0;
				}

				if processes[peek].burst > 0 && processes[peek].arrived == true {
					currentQuantum = algo.quantum;
					currentProcess = processes[peek];
					fmt.Fprintf(output, "Time %3d: %s selected (burst %3d)\n", time, currentProcess.name, currentProcess.burst);
					selected = peek;
					searching = 0;
					next = 0;
				}

				if i > (algo.processcount+1) {
					selected = -1;
					searching = 0;
					next = 1;
				}

				peek++;
				i++;
			}
		}

		if selected >= 0 {
			if currentProcess.burst == 0 {
				fmt.Fprintf(output, "Time %3d: %s finished\n", time, currentProcess.name);
				currentProcess.turnaround = time - currentProcess.arrival;
				finished++;

				searching := 1;
				i := 0;
				peek := selected + 1;

				for searching == 1 {
					if peek > (algo.processcount-1) {
						peek = 0;
					}

					if processes[peek].burst > 0 && processes[peek].arrived == true {
						currentQuantum = algo.quantum;
						currentProcess = processes[peek];
						fmt.Fprintf(output, "Time %3d: %s selected (burst %3d)\n", time, currentProcess.name, currentProcess.burst);
						selected = peek;
						searching = 0;
					}

					if i > (algo.processcount+1) {
						selected = -1;
						searching = 0;
					}

					peek++;
					i++;
				}
			}

			currentProcess.burst--;
			currentQuantum--;

			if currentQuantum == 0 {
				next = 1;
			}
		}

		if selected == -1 {
			fmt.Fprintf(output, "Time %3d: IDLE\n", time);
		}

		if finished == algo.processcount {
			running = 0;
		}
	}

	fmt.Fprintf(output, "Finished at time %3d\n\n", time);

	processes = sortByName(processes);

	for i := 0; i < algo.processcount; i++ {
		fmt.Fprintf(output, "%s wait %3d turnaround %3d\n", processes[i].name, processes[i].wait,
		processes[i].turnaround);
	}
}

func shortestJobFirst(processes []Process) {

	output, err := os.Create(os.Args[2]);

	if err != nil {
		log.Fatal(err);
	}

	fmt.Fprintf(output, "%3d processes\n", algo.processcount);

	fmt.Fprintln(output, "Using preemptive Shortest Job First");

	processes = sortByArrival(processes);

	time := 0;

	var Queue []Process = make([]Process, 0, algo.processcount);

	QCap := 0;

	for time < algo.runfor {
		for i := 0; i < algo.processcount; i++ {
			if processes[i].arrival == time {
				fmt.Fprintf(output, "Time %3d : %s arrived\n", time, processes[i].name)
				Queue = append(Queue, processes[i]);
				QCap++;
			}
		}

		if QCap == 0 {
			fmt.Fprintf(output, "Time %3d : Idle\n", time)
		}

		currentProcess := 0;
		var prevProcess string;

		if QCap > 0 {
			for Queue[currentProcess].completed {
				currentProcess++;
			}

			if (currentProcess < algo.processcount) && Queue[currentProcess].chosen {
				prevProcess = Queue[currentProcess].name;
			}

			for i := 0; i < len(Queue); i++ {
				Queue[i].chosen = false;
			}

			Queue = sortByBurst(Queue);

			for (currentProcess < algo.processcount) && Queue[currentProcess].completed {
				currentProcess++;
			}

			if (currentProcess < algo.processcount) && (Queue[currentProcess].name == prevProcess) {
				Queue[currentProcess].chosen = true;
			}

			if (currentProcess < algo.processcount) && Queue[currentProcess].chosen && (Queue[currentProcess].burst == 0) && !Queue[currentProcess].completed {
				Queue[currentProcess].chosen = false;
				Queue[currentProcess].completed = true;
				Queue[currentProcess].completion = time;
				QCap--;

				fmt.Fprintf(output, "Time %3d : %s finished\n", time, Queue[currentProcess].name);
			}

			if (currentProcess < algo.processcount) && Queue[currentProcess].completed && QCap > 0 {
				currentProcess++;
			}

			if (currentProcess < algo.processcount) && !Queue[currentProcess].chosen && !Queue[currentProcess].completed && QCap > 0 {
				Queue[currentProcess].chosen = true;
				Queue[currentProcess].selection = time;

				fmt.Fprintf(output, "Time %3d : %s selected (burst %3d)\n", time, Queue[currentProcess].name, Queue[currentProcess].burst);
			}

			if QCap == 0 {
				fmt.Fprintf(output, "Time %3d : Idle\n", time)
			}
		}

		time++;

		if (currentProcess < algo.processcount) && Queue[currentProcess].burst > 0 {
			Queue[currentProcess].burst--;
		}
	}

	for i := 0; i < algo.processcount; i++ {
		Queue[i].turnaround = Queue[i].completion - Queue[i].arrival;
	}

	Queue = sortByName(Queue);
	processes = sortByName(processes);

	for i := 0; i < algo.processcount; i++ {
		Queue[i].wait = Queue[i].turnaround - processes[i].burst;
	}

	fmt.Fprintf(output, "Finished at time %3d\n\n", algo.runfor);

	for i := 0; i < algo.processcount; i++ {
		fmt.Fprintf(output, "%s wait %3d turnaround %3d\n", Queue[i].name, Queue[i].wait, Queue[i].turnaround);
	}
}

type Process struct {
	name string; // P1, P2, etc.
	arrival int; // process arrival time
	burst int; // process burst time
	wait int; // how long a process waits
	turnaround int; // time between wait and completion for a process
	completion int; // time to completion for a process
	selection int; // time to be selected for a process
	chosen bool; // when a process is selected, this becomes true
	completed bool; // when a process is completed, this becomes true
	arrived bool; // when a process has arrived yet or not
}

type Algorithm struct {
	style string; // either fcfs, rr, or psjf
	processcount int; // how many process will be read
	runfor int; // time units
	processes[10] Process; // per the rubric, max # of processes will be 10
	quantum int; // only if using rr
}