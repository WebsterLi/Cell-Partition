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
	moved, leftside bool
	NetList []*Net
	prevcell, nextcell, endcell *Cell
}
var (
	cellcount, maxgain, mingain int
	degree float64
	netslice []*Net
	leftpart map[int]*Cell
	rightpart map[int]*Cell
	cellmap map[int]*Cell
	gainmap map[int]*Cell//bucketlist
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
			if degree > 0.5 { degree = 1 - degree }
		}
		for _, word := range netinfo {
			switch word[0] {
			case 'N':
				var clist []*Cell
				netptr = &Net{name:netid, leftnum:0, rightnum:0, CellList:clist}
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
					cellptr = &Cell{name:cellid, NetList:nlist, moved:false, gain:0}
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
			if _, ok := leftpart[cell.name]; !ok {
				leftpart[cell.name] = cell
				cell.leftside = true //update cell position
				//update net info 
				for _, net := range cell.NetList{
					net.leftnum ++
				}
			}
		} else {
			if _, ok := rightpart[cell.name]; !ok {
				rightpart[cell.name] = cell
				cell.leftside = false //update cell position
				//update net info 
				for _, net := range cell.NetList{
					net.rightnum ++
				}
			}
		}
		cell.moved = false //set to none moved.
	}
}

func InitialGain(){
	//Calculate gain of each cell
	for _, cell := range cellmap {
		var cellgain int
		for _, net := range cell.NetList {
			//no cell on the other side -> gain += -1
			//one self on this side -> gain += 1
			if cell.leftside {
				if net.rightnum == 0 { cellgain-- }
				if net.leftnum == 1 { cellgain++ }
			} else {
				if net.leftnum == 0 { cellgain-- }
				if net.rightnum == 1 { cellgain++ }
			}
		}
		cell.gain = cellgain
		//update bound
		if cellgain > maxgain {
			maxgain = cellgain
		}
		if cellgain < mingain {
			mingain = cellgain
		}
	}
}

func InitialBucket(){
	for _, cell := range cellmap {
		cellgain := cell.gain
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
	}
	/*
	//print gain map member
	for i := mingain; i <= maxgain; i++ {
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
	*/
}

func RemoveFromBucket(target *Cell) {
	index := target.gain
	//Target cell is the root of bucket
	if target.prevcell == nil {
		//Only member in this gain bucket.
		if target.nextcell == nil {
			target.endcell = nil
			delete (gainmap,target.gain)
			return
		}
		gainmap[index] = target.nextcell
		//next cell get endcell & delete prevcell
		target.nextcell.endcell = target.endcell
		target.nextcell.prevcell = nil
		//delete endcell & nextcell
		target.endcell = nil
		target.nextcell = nil
		return
	}
	//link nextcell and prevcell
	target.prevcell.nextcell = target.nextcell
	if target.nextcell != nil {
		target.nextcell.prevcell = target.prevcell
	}
	//delete self prevcell nextcell link
	target.nextcell = nil
	target.prevcell = nil
	return
}

func AppendToBucket(target *Cell) {
	index := target.gain
	if root, ok := gainmap[index]; ok {
		root.endcell.nextcell = target
		target.prevcell = root.endcell
		root.endcell = target
	} else {
		gainmap[index] = target
		target.endcell = target
	}
}

func UpdateGain(target *Cell) {
	var cellgain int
	if target.moved {
		cellgain = 0
		for _, net := range target.NetList {
			if !target.leftside {
				net.rightnum++
				net.leftnum--
				if net.leftnum == 0 { cellgain-- }
				if net.rightnum == 1 { cellgain++ }

			} else {
				net.leftnum++
				net.rightnum--
				if net.rightnum == 0 { cellgain-- }
				if net.leftnum == 1 { cellgain++ }
			}
		}
		target.gain = cellgain
		//Update gain of other realated cell.
		for _, net := range target.NetList {
			for _, cell := range net.CellList {
				if !cell.moved { UpdateGain(cell) }
			}
		}
	} else {
		cellgain = 0
		for _, net := range target.NetList {
			//no cell on the other side -> gain += -1
			//one self on this side -> gain += 1
			if target.leftside {
				if net.rightnum == 0 { cellgain-- }
				if net.leftnum == 1 { cellgain++ }
			} else {
				if net.leftnum == 0 { cellgain-- }
				if net.rightnum == 1 { cellgain++ }
			}
		}
		if target.gain != cellgain {
			RemoveFromBucket(target)//Need to be done before update gain!
			target.gain = cellgain
			AppendToBucket(target)
		}
	}
	//update bound
	if cellgain > maxgain {
		maxgain = cellgain
	}
	if cellgain < mingain {
		mingain = cellgain
	}
}

func MoveCell(target *Cell) {
	var move bool
	if target.leftside {
		move = len(leftpart) - 1 > int(float64(cellcount) * degree)
	} else {
		move = len(rightpart) - 1 > int(float64(cellcount) * degree)
	}
	if move {
		//Remove operation need to be done before update gain!
		RemoveFromBucket(target)
		//move cell to other side.
		if target.leftside {
			delete (leftpart, target.name)
			rightpart[target.name] = target
		} else {
			delete (rightpart, target.name)
			leftpart[target.name] = target
		}
		target.leftside = !target.leftside
		target.moved = true
		//calculate gain.
		UpdateGain(target)
		//Target cell don't need to append back to bucket? TODO
	}
}
func FMLoop() {
	if len(gainmap) == 0 {
		InitialPartition()
		InitialGain()
		InitialBucket()
	}
	for i := maxgain; i > 0; i-- {
		if gcell, ok := gainmap[i]; ok {
			for gcell.nextcell != nil {
				queuecell := gcell.nextcell
				MoveCell(gcell)
				gcell = queuecell
			}
			MoveCell(gcell)
		}
	}
	fmt.Println(len(leftpart), len(rightpart))
	if len(leftpart)+len(rightpart) == cellcount { fmt.Println("FM iteration success!") }
}
func main() {
	//Initial map
	cellmap = make(map[int]*Cell)
	leftpart = make(map[int]*Cell)
	rightpart = make(map[int]*Cell)
	gainmap = make(map[int]*Cell)//Initial map
	// Loop over lines in file.
	lines := LinesInFile(`input_data/input_0.dat`)
	LinesToGraph(lines)
	FMLoop()
}
