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
	CellList []int
}
type Cell struct{
	name, gain int
	moved, leftpart bool
	NetList []int
}
var (
	cellcount int
	degree float64
	cellmap map[int]*Cell
	netslice []Net
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

func LinesToCell(lines []string){
	var (
		netid, cellid int
		err error
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
				var clist []int
				netslice = append(netslice, Net{name:netid,CellList:clist})
				netid++
			case 'c':
				cellid, err = strconv.Atoi(strings.Trim(word,"c"))
				if err != nil {fmt.Println(word)}
				if curcell, ok := cellmap[cellid]; ok {
					curcell.NetList = append(curcell.NetList, netid)
					cellmap[cellid] = curcell
				} else {
					//Initial a cell
					nlist := []int{netid}
					cellmap[cellid] = &Cell{name:cellid, NetList:nlist}
					cellcount++
				}
				netslice[len(netslice)-1].CellList = append(netslice[len(netslice)-1].CellList, cellid)
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
			for _, netid := range cell.NetList{
				netslice[netid-1].leftnum ++// netindex = netid - 1
			}
		} else {
			rightpart = append(rightpart, cell)
			cell.leftpart = false //update cell position
			//update net info 
			for _, netid := range cell.NetList{
				netslice[netid-1].rightnum ++// netindex = netid - 1
			}
		}
		counter++
	}
}

func InitialBucket(){
	//Calculate gain of each cell
	var maxgain int
	for _, cell := range cellmap {
		var cellgain int
		for _, netid := range cell.NetList {
			if cell.leftpart {
				cellgain += netslice[netid-1].rightnum - netslice[netid-1].leftnum
			} else {
				cellgain += netslice[netid-1].leftnum - netslice[netid-1].rightnum
			}
		}
		cell.gain = cellgain
		if cellgain > maxgain { maxgain = cellgain }
	}
	fmt.Println(maxgain)
}

func main() {
	cellmap = make(map[int]*Cell)//Initial map
	// Loop over lines in file.
	lines := LinesInFile(`input_data/input_0.dat`)
	LinesToCell(lines)
	InitialPartition()
	InitialBucket()
}
