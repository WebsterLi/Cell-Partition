package main

import (
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"sort"
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
	prevcell, nextcell, endcell *Cell
}
var (
	cellcount, maxgain, mingain int
	degree float64
	cellmap map[int]*Cell
	netslice []*Net
	leftpart []*Cell
	rightpart []*Cell
	bucketroot *Cell
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
	var cell_by_netnum []*Cell
	for _, cell := range cellmap {
		cell_by_netnum =  append(cell_by_netnum, cell)
	}
	sort.Slice(cell_by_netnum, func(i, j int) bool {
		return len(cell_by_netnum[i].NetList) < len(cell_by_netnum[j].NetList)
	})
	for _, cell := range cell_by_netnum {
		if len(leftpart)+1 <= cellcount/2 {
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
	}
}

func InitialBucket(){
	//Calculate gain of each cell
	//gain calculate
	maxgain, mingain = 0, 0
	gainmap := make(map[int]*Cell)//Initial map

	for _, cell := range cellmap {
		var cellgain int
		for _, net := range cell.NetList {
			if cell.leftpart {
				if net.rightnum == 0 { cellgain-- }
				if net.leftnum == 1 { cellgain++ }
			} else {
				if net.leftnum == 0 { cellgain-- }
				if net.rightnum == 1 { cellgain++ }
			}
		}
		cell.gain = cellgain
		//Initial gain(bucket) list
		if root, ok := gainmap[cellgain]; ok {
			root.endcell.nextcell = cell
			cell.prevcell = root.endcell
			root.endcell = cell
		} else {
			//Initial a cell
			gainmap[cellgain] = cell
			cell.endcell = cell
		}

		if cellgain > maxgain {
			maxgain = cellgain
		}
		if cellgain < mingain {
			mingain = cellgain
		}
	}
	for i:= mingain; i<=maxgain; i++ {
		if gcell, ok := gainmap[i]; ok {
			count := 1
			fmt.Println("")
			fmt.Printf("Gain %d : ", i)
			for gcell.nextcell != nil {
				count++
				gcell = gcell.nextcell
			}
			fmt.Println(count)
		}
	}
}

func main() {
	cellmap = make(map[int]*Cell)//Initial map
	// Loop over lines in file.
	lines := LinesInFile(`input_data/input_0.dat`)
	LinesToGraph(lines)
	InitialPartition()
	InitialBucket()
}
