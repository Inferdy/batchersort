package batchersort

import (
	"sync"

	"github.com/Inferdy/inferdymath"
)

type task struct {
	first  int
	second int
}

func getStartStopPair(first int, second int, total int, delta int) (start int, stop int, pair int) {
	if first > delta {
		start = first - delta
		return start, start + second - 1, start + delta
	}

	if second > delta {
		stop = first
	} else {
		stop = total - delta
	}
	stop--

	return 0, stop, delta
}

func solver(activeWG *sync.WaitGroup, x []int, tasks chan task) {
	var tmp int

	for task := range tasks {
		if x[task.first] > x[task.second] {
			tmp = x[task.first]
			x[task.first] = x[task.second]
			x[task.second] = tmp
		}

		activeWG.Done()
	}
}

// Sort1 sorts a slice containing two sorted parts
//
// first is the size of the first part
//
// second is the size of the second part
//
// total is the total size
//
// delta is the power of 2 closest and less or equal to total-1
//
// solversCount is the number of threads solving parallel tasks simultaneously
//
// extraCount defines how many tasks are prepared for the freed solvers
func Sort1(x []int, first int, second int, total int, delta int, solversCount int, extraCount int) {
	if (first == 0) || (second == 0) {
		return
	}

	var tasks chan task = make(chan task, solversCount+extraCount)

	var activeWG sync.WaitGroup

	//#region spawning solver threads
	for i := 1; ; i++ {
		go solver(&activeWG, x, tasks)

		if i == solversCount {
			break
		}
	}
	//#endergion spawning solver threads

	//#region solving first parallel tasks group with optimisations
	{
		var ptr, stop, pair int = getStartStopPair(first, second, total, delta)

		for {
			activeWG.Add(1)
			tasks <- task{ptr, pair}

			if ptr == stop {
				break
			}

			ptr++
			pair++
		}
	}
	activeWG.Wait()
	//#endregion solving first parallel tasks group with optimisations

	//#region solving all the other parallel task groups
	{
		var (
			//used for inexplicable calculations
			firstMirage int = first - 1
			totalMirage int = total - 1

			//stores max value of i shift for cycle, equals to delta-1
			stop int

			//tells if column is starting with comparators shift
			startShift bool
			//shows how soon startShift will change (on extraFirstLines == 0)
			extraFirstLines int

			//shows max last first index possible with this i shift
			maxJ int
			//shows how soon maxJ will decrease by 1 (on extraTotalLines == 0)
			extraTotalLines int

			//i shift (or ptr), changes from 0 to delta - 1, defines all <delta> spaces comparators of <delta> size may be in
			i int
			//j shift, 0 or 1, affected by startShift bool, will be multiplied by delta
			j int

			//ptr to the first comparator (task) index
			ptr int
			//stores max value of ptr for cycle
			maxptr int
			//ptr to the last comparator (task) index
			pair int

			//stores delta*2
			doubleDelta int
		)

		for {
			doubleDelta = delta
			delta >>= 1

			if delta == 0 {
				//break could be here

				//defer?
				close(tasks)

				//wait for activeWG?
				return
			}

			//stores max value of i shift for cycle
			stop = delta - 1

			{
				//inexplicable calculations
				var tmpFirstPointsMinus1 int = firstMirage / delta
				//shows how soon start shift will change (on extraFirstLines == 0)
				extraFirstLines = firstMirage % delta

				//tells if column is starting with comparators shift
				startShift = ((tmpFirstPointsMinus1 & 1) == 1)
			}

			//shows max last first index possible with this i shift
			maxJ = totalMirage/delta - 1
			//shows how soon maxJ will decrease by 1 (on extraTotalLines == 0)
			extraTotalLines = totalMirage % delta

			//i shift (or ptr), changes from 0 to delta - 1, defines all <delta> spaces comparators of <delta> size may be in
			//after this loop activeWG.Wait() is called
			for i = 0; ; i++ {
				if maxJ == 0 {
					//only one comparator location is possible, so comparator must be placed here
					activeWG.Add(1)
					tasks <- task{i, i + delta}

					//goto i_check?
				} else {
					//ptr to the first comparator (task) index
					ptr = i

					if startShift {
						ptr += delta
						j = 1
					} else {
						j = 0
					}

					//ptr to the first comparator (task) index
					maxptr = ((maxJ-j)>>1)*doubleDelta + ptr

					//ptr to the last comparator (task) index
					pair = ptr + delta

					for {
						activeWG.Add(1)
						tasks <- task{ptr, pair}

						if ptr == maxptr {
							break
						}

						ptr += doubleDelta
						pair += doubleDelta
					}
				}

				if i == stop {
					break
				}

				if extraTotalLines == 0 {
					extraTotalLines = totalMirage
					maxJ--
				} else {
					extraTotalLines--
				}

				if extraFirstLines == 0 {
					extraFirstLines = firstMirage
					startShift = !startShift
				} else {
					extraFirstLines--
				}
			}
			activeWG.Wait()
		}
	}
	//#endregion solving all the other parallel task groups
}

// Sort2 sorts a slice containing two sorted parts
//
// first is the size of the first part
//
// second is the size of the second part
//
// total is the total size
//
// solversCount is the number of threads solving parallel tasks simultaneously
//
// extraCount defines how many tasks are prepared for the freed solvers
func Sort2(x []int, first int, second int, total int, solversCount int, extraCount int) {
	if (first == 0) || (second == 0) {
		return
	}

	powerOf2, err := inferdymath.PowerOf2LessOrEqualTo(total - 1)

	if err == nil {
		Sort1(x, first, second, total, powerOf2, solversCount, extraCount)
	} /* else {
		//TODO: what to do? it is not possible!
	}
	*/
}
