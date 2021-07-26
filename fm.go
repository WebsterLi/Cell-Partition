package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Net struct {
	name     int
	leftnum  int
	rightnum int
	CellList map[int]*Cell
}

type Cell struct {
	name, gain                  int
	moved, leftside             bool
	NetList                     map[int]*Net
	prevcell, nextcell, endcell *Cell
}

type Partitioner struct {
	iter, cellcount, maxgain, mingain, currcut, prevcut int
	degree                                              float64
	netslice                                            []*Net
	leftpart                                            map[int]*Cell
	rightpart                                           map[int]*Cell
	cellmap                                             map[int]*Cell
	gainmap                                             map[int]*Cell //bucketlist
}

func (c *Cell) Reset() {
	c.prevcell = nil
	c.nextcell = nil
	c.endcell = nil
}

func NewPartitioner() *Partitioner {
	p := new(Partitioner)
	//Initial maps
	p.cellmap = make(map[int]*Cell)
	p.leftpart = make(map[int]*Cell)
	p.rightpart = make(map[int]*Cell)
	p.gainmap = make(map[int]*Cell)
	return p
}

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

func OutputFile(name string, pter *Partitioner) {
	content := fmt.Sprintf("Cutsize = %d", pter.currcut)
	g1 := fmt.Sprintf("G1 %d\n", len(pter.leftpart))
	for _, cell := range pter.leftpart {
		g1 = fmt.Sprintf("%sc%d ", g1, cell.name)
	}
	g1 = fmt.Sprintf("%s;", g1)
	g2 := fmt.Sprintf("G2 %d\n", len(pter.rightpart))
	for _, cell := range pter.rightpart {
		g2 = fmt.Sprintf("%sc%d ", g2, cell.name)
	}
	g2 = fmt.Sprintf("%s;", g2)
	content = fmt.Sprintf("%s\n%s\n%s\n", content, g1, g2)
	word := []byte(content)
	// write the whole body at once
	err := ioutil.WriteFile(name, word, 0777)
	if err != nil {
		panic(err)
	}
}

func LinesToGraph(lines []string, pter *Partitioner) {
	var (
		netid, cellid int
		err           error
		cellptr       *Cell
		netptr        *Net
	)
	for iter, line := range lines {
		netinfo := strings.Fields(line)
		if iter == 0 {
			pter.degree, err = strconv.ParseFloat(netinfo[0], 64)
			if err != nil {
				fmt.Println(netinfo)
			}
			if pter.degree > 0.5 {
				pter.degree = 1 - pter.degree
			}
		}
		for _, word := range netinfo {
			switch word[0] {
			case 'N':
				clist := make(map[int]*Cell)
				netptr = &Net{name: netid, leftnum: 0, rightnum: 0, CellList: clist}
				pter.netslice = append(pter.netslice, netptr)
				netid++
			case 'c':
				cellid, err = strconv.Atoi(strings.Trim(word, "c"))
				if err != nil {
					fmt.Println(word)
				}
				if curcell, ok := pter.cellmap[cellid]; ok {
					curcell.NetList[netptr.name] = netptr
					pter.cellmap[cellid] = curcell
				} else {
					//Initial a cell
					nlist := make(map[int]*Net)
					nlist[netptr.name] = netptr
					cellptr = &Cell{name: cellid, NetList: nlist, moved: false, gain: 0}
					pter.cellmap[cellid] = cellptr
					pter.cellcount++
				}
				pter.netslice[len(pter.netslice)-1].CellList[cellptr.name] = cellptr
			default:
			}
		}
	}
}

func PrintInfo(pter *Partitioner) {
	if pter.prevcut == 0 {
		fmt.Println("--------------Initial info---------------")
	} else {
		fmt.Println("------------FM partition info------------")
	}
	fmt.Println("	total iteration:", pter.iter)
	fmt.Println("	gain range:", pter.maxgain, pter.mingain)
	fmt.Println("	total remain gain:", pter.currcut)
	fmt.Println("	partition status:", len(pter.leftpart), len(pter.rightpart))
}

func (pter *Partitioner) InitialPartition() {
	var cell_by_netnum []*Cell
	for _, cell := range pter.cellmap {
		cell_by_netnum = append(cell_by_netnum, cell)
	}
	sort.Slice(cell_by_netnum, func(i, j int) bool {
		return len(cell_by_netnum[i].NetList) < len(cell_by_netnum[j].NetList)
	})
	for _, cell := range cell_by_netnum {
		_, isleft := pter.leftpart[cell.name]
		_, isright := pter.rightpart[cell.name]
		if !(isleft || isright) {
			if len(pter.leftpart)+1 <= pter.cellcount/2 {
				pter.leftpart[cell.name] = cell
				cell.leftside = true //update cell position
				//update net info
				for _, net := range cell.NetList {
					net.leftnum++
				}
			} else {
				pter.rightpart[cell.name] = cell
				cell.leftside = false //update cell position
				//update net info
				for _, net := range cell.NetList {
					net.rightnum++
				}
			}
			cell.moved = false //set to none moved.
		}
	}
}

func (pter *Partitioner) GetGain() {
	pter.maxgain = 0
	pter.mingain = 0
	//Calculate gain of each cell
	for _, cell := range pter.cellmap {
		var cellgain int
		for _, net := range cell.NetList {
			//no cell on the other side -> gain += -1
			//one self on this side -> gain += 1
			if cell.leftside {
				if net.rightnum == 0 {
					cellgain--
				}
				if net.leftnum == 1 {
					cellgain++
				}
			} else {
				if net.leftnum == 0 {
					cellgain--
				}
				if net.rightnum == 1 {
					cellgain++
				}
			}
		}
		cell.gain = cellgain
		//update bound
		if cellgain > pter.maxgain {
			pter.maxgain = cellgain
		}
		if cellgain < pter.mingain {
			pter.mingain = cellgain
		}
	}
}

func (pter *Partitioner) GetBucket() {
	pter.gainmap = make(map[int]*Cell) //Reset map
	for _, cell := range pter.cellmap {
		cell.Reset()
		cellgain := cell.gain
		//Initial gain(bucket) list
		if root, ok := pter.gainmap[cellgain]; ok {
			root.endcell.nextcell = cell
			cell.prevcell = root.endcell
			root.endcell = cell
		} else {
			//Initial a cell
			pter.gainmap[cellgain] = cell
			cell.endcell = cell
		}
	}
}

func (target *Cell) RemoveFromBucket(pter *Partitioner) {
	index := target.gain
	//Target cell is the root of bucket
	if target.prevcell == nil {
		//Only member in this gain bucket.
		if target.nextcell == nil {
			target.endcell = nil
			delete(pter.gainmap, target.gain)
			return
		}
		pter.gainmap[index] = target.nextcell
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
	} else {
		pter.gainmap[index].endcell = target.prevcell
	}
	//delete self pointer link
	target.nextcell = nil
	target.prevcell = nil
	target.endcell = nil
	return
}

func (target *Cell) AppendToBucket(pter *Partitioner) {
	index := target.gain
	if root, ok := pter.gainmap[index]; ok {
		root.endcell.nextcell = target
		target.prevcell = root.endcell
		root.endcell = target
	} else {
		pter.gainmap[index] = target
		target.endcell = target
	}
}

func (target *Cell) UpdateGain(pter *Partitioner) {
	var cellgain int
	if target.moved {
		cellgain = 0
		for _, net := range target.NetList {
			if !target.leftside {
				if net.leftnum == 0 {
					cellgain--
				}
				if net.rightnum == 1 {
					cellgain++
				}

			} else {
				if net.rightnum == 0 {
					cellgain--
				}
				if net.leftnum == 1 {
					cellgain++
				}
			}
		}
		target.gain = cellgain
		//Update gain of other realated cell.
		for _, net := range target.NetList {
			for _, cell := range net.CellList {
				if !cell.moved {
					cell.UpdateGain(pter)
				}
			}
		}
	} else {
		cellgain = 0
		for _, net := range target.NetList {
			//no cell on the other side -> gain += -1
			//one self on this side -> gain += 1
			if target.leftside {
				if net.rightnum == 0 {
					cellgain--
				}
				if net.leftnum == 1 {
					cellgain++
				}
			} else {
				if net.leftnum == 0 {
					cellgain--
				}
				if net.rightnum == 1 {
					cellgain++
				}
			}
		}
		if target.gain != cellgain {
			target.RemoveFromBucket(pter) //Need to be done before update gain!
			target.gain = cellgain
			target.AppendToBucket(pter)
		}
	}
	//update bound
	if cellgain > pter.maxgain {
		pter.maxgain = cellgain
	}
	if cellgain < pter.mingain {
		pter.mingain = cellgain
	}
}

func (target *Cell) MoveCell(pter *Partitioner) {
	var move bool
	if target.leftside {
		move = len(pter.leftpart)-1 > int(float64(pter.cellcount)*pter.degree)
	} else {
		move = len(pter.rightpart)-1 > int(float64(pter.cellcount)*pter.degree)
	}
	if move {
		//Remove operation need to be done before update gain!
		target.RemoveFromBucket(pter)
		//move cell to other side.
		if target.leftside {
			delete(pter.leftpart, target.name)
			pter.rightpart[target.name] = target
			/*
				for _, net := range target.NetList {
					net.leftnum--
					net.rightnum++
				}
			*/
		} else {
			delete(pter.rightpart, target.name)
			pter.leftpart[target.name] = target
			/*
				for _, net := range target.NetList {
					net.leftnum++
					net.rightnum--
				}
			*/
		}
		target.leftside = !target.leftside
		target.moved = true
		//calculate gain.
		target.UpdateGain(pter)
		//Update after moving cell(may increase time?)
		for _, net := range target.NetList {
			left := 0
			right := 0
			for _, cell := range net.CellList {
				if cell.leftside {
					left++
				} else {
					right++
				}
			}
			net.leftnum = left
			net.rightnum = right
		}
	}
}

func (pter *Partitioner) PartitionSet() {
	pter.iter++
	pter.GetGain()
	pter.GetBucket()
	pter.currcut = 0
	//Calculate cutsize
	for _, net := range pter.netslice {
		if net.leftnum != 0 && net.rightnum != 0 {
			pter.currcut++
		}
	}
	//PrintInfo(pter)
}
func (pter *Partitioner) FMLoop() {
	if len(pter.gainmap) == 0 {
		pter.InitialPartition()
		pter.PartitionSet()
		pter.prevcut = pter.currcut
	}
	for i := pter.maxgain; i > 0; i-- {
		if gcell, ok := pter.gainmap[i]; ok {
			for gcell.nextcell != nil {
				queuecell := gcell.nextcell
				gcell.MoveCell(pter)
				gcell = queuecell
			}
			gcell.MoveCell(pter)
		}
	}
	//prepare for next loop
	pter.PartitionSet()
	if pter.currcut < pter.prevcut {
		pter.prevcut = pter.currcut
		pter.FMLoop()
	} else {
		PrintInfo(pter)
	}
}
func main() {
	pter := NewPartitioner()
	// Loop over lines in file.
	lines := LinesInFile(`input_data/input_5.dat`)
	LinesToGraph(lines, pter)
	pter.FMLoop()
	OutputFile("result_5.dat", pter)
}
