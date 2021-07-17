package main

import (
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"os"
)
type Net struct{
	name int
	leftnum int
	rightnum int
	CellList []*Cell
}
type Cell struct{
	name, gain int
	moved, leftpart bool
	NetList []*Net
}
var (
	cellcount int
	degree float64
	cellmap map[int]*Cell
	netslice []*Net
	bucketlist [][]int
	leftpart []*Cell
	rightpart []*Cell
)

func LinesInFile(fileName string) []string {
	f, _ := os.Open(fileName)
	// Create new Scanner.
	scanner := bufio.NewScanner(f)
	result := []string{}
	// Use Scan.
	for scanner.Scan() {
		line := scanner.Text()
		// Append line to result.
		result = append(result, line)
	}
	return result
}

func LinesToGraph(lines []string){
	var (
		netid, cellid int
		err error
		cellptr *Cell
		netptr *Net
	)
	for iter, line := range lines {
		netinfo := strings.Fields(line)
		if iter == 0 {
			degree, err = strconv.ParseFloat(netinfo[0], 64)
			if err != nil {fmt.Println(netinfo)}
		}
		for _, word := range netinfo {
			switch word[0] {
			case 'N':
				var clist []*Cell
				netptr = &Net{name:netid, CellList:clist}
				netslice = append(netslice, netptr)
				netid++
			case 'c':
				cellid, err = strconv.Atoi(strings.Trim(word,"c"))
				if err != nil {fmt.Println(word)}
				if curcell, ok := cellmap[cellid]; ok {
					curcell.NetList = append(curcell.NetList, netptr)
					cellmap[cellid] = curcell
				} else {
					//Initial a cell
					nlist := []*Net{netptr}
					cellptr = &Cell{name:cellid, NetList:nlist}
					cellmap[cellid] = cellptr
					cellcount++
				}
				netslice[len(netslice)-1].CellList = append(netslice[len(netslice)-1].CellList, cellptr)
			default :
			}
		}
	}
}

func InitialPartition(){
	var counter int
	for _, cell := range cellmap{
		if float64(len(leftpart)+1) <= float64 (cellcount) * (0.5) {
			leftpart = append(leftpart, cell)
			cell.leftpart = true //update cell position
			//update net info 
			for _, net := range cell.NetList{
				net.leftnum ++
			}
		} else {
			rightpart = append(rightpart, cell)
			cell.leftpart = false //update cell position
			//update net info 
			for _, net := range cell.NetList{
				net.rightnum ++
			}
		}
		counter++
	}
}

func InitialBucket(){
	//Calculate gain of each cell
	maxgain := 0
	for _, cell := range cellmap {
		var cellgain int
		for _, net := range cell.NetList {
			if cell.leftpart {
				cellgain += net.rightnum - net.leftnum
			} else {
				cellgain += net.leftnum - net.rightnum
			}
		}
		cell.gain = cellgain
		if cellgain > maxgain {
			maxgain = cellgain
		}
	}
	fmt.Println(maxgain)
	fmt.Printf("len=%d cap=%d %v\n", len(bucketlist), cap(bucketlist), bucketlist)
}

func main() {
	cellmap = make(map[int]*Cell)//Initial map
	// Loop over lines in file.
	lines := LinesInFile(`input_data/input_0.dat`)
	LinesToGraph(lines)
	InitialPartition()
	InitialBucket()
}
