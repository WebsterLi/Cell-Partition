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
	CellList []int
}
type Cell struct{
	name int
	NetList []int
}
var (
	degree float64
	cellmap map[int]Cell
	netslice []Net
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
			fmt.Println(degree)
		}
		for _, word := range netinfo {
			switch word[0] {
			case 'n':
				netid, err = strconv.Atoi(strings.Trim(word,"n"))
				if err != nil {fmt.Println(word)}
				var clist []int
				netslice = append(netslice, Net{netid,clist})
			case 'c':
				cellid, err = strconv.Atoi(strings.Trim(word,"c"))
				if err != nil {fmt.Println(word)}
				if curcell, ok := cellmap[cellid]; ok {
					curcell.NetList = append(curcell.NetList, netid)
					cellmap[cellid] = curcell
				} else {
					//Initial a cell
					nlist := []int{netid}
					cellmap[cellid] = Cell{cellid, nlist}
				}
				netslice[len(netslice)-1].CellList = append(netslice[len(netslice)-1].CellList, cellid)
			default :
			}
		}
	}
}

func main() {
	cellmap = make(map[int]Cell)//Initial map
	// Loop over lines in file.
	lines := LinesInFile(`input_data/input_0.dat`)
	LinesToCell(lines)
}
