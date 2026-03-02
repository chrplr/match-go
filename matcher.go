package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"
)

type Matcher struct {
	Config      *Config
	DataSets    map[string]*DataSet
	SetIDs      []string
	BestSol     *Solution
	BestDist    float64
	Running     bool
	
	// State for DFS
	UsedItems map[string]map[int]bool // FileID -> ItemID -> bool
	Tuples    []*Tuple
	
	StartTime time.Time
	OnUpdate  func(float64)
}

func NewMatcher(config *Config, dataSets map[string]*DataSet) *Matcher {
	setIDs := make([]string, 0, len(config.InputFiles))
	for id := range config.InputFiles {
		setIDs = append(setIDs, id)
	}
	sort.Strings(setIDs)

	return &Matcher{
		Config:    config,
		DataSets:  dataSets,
		SetIDs:    setIDs,
		UsedItems: make(map[string]map[int]bool),
		Tuples:    make([]*Tuple, config.OutputSize),
		BestDist:  math.MaxFloat64,
	}
}

func (m *Matcher) calcTupleDist(items []*Item) float64 {
	distSq := 0.0
	for dIdx, dim := range m.Config.Dimensions {
		var item1, item2 *Item
		idx1, idx2 := -1, -1
		for i, id := range m.SetIDs {
			if id == dim.File1ID {
				idx1 = i
			}
			if id == dim.File2ID {
				idx2 = i
			}
		}
		if idx1 != -1 && idx2 != -1 {
			item1 = items[idx1]
			item2 = items[idx2]
			if item1 != nil && item2 != nil {
				diff := item1.Values[dIdx] - item2.Values[dIdx]
				distSq += diff * diff
			}
		}
	}
	return distSq
}

func (m *Matcher) SeedInitialSolution() {
	fmt.Println("Seeding initial solution...")
	rand.Seed(time.Now().UnixNano())
	
	m.UsedItems = make(map[string]map[int]bool)
	for _, id := range m.SetIDs {
		m.UsedItems[id] = make(map[int]bool)
	}

	for tIdx := 0; tIdx < m.Config.OutputSize; tIdx++ {
		m.Tuples[tIdx] = &Tuple{Items: make([]*Item, len(m.SetIDs))}
		for sIdx, sID := range m.SetIDs {
			ds := m.DataSets[sID]
			var available []*Item
			for _, item := range ds.Items {
				if !m.UsedItems[sID][item.ID] {
					available = append(available, item)
				}
			}
			if len(available) == 0 {
				break
			}
			picked := available[rand.Intn(len(available))]
			m.Tuples[tIdx].Items[sIdx] = picked
			m.UsedItems[sID][picked.ID] = true
		}
	}

	for iter := 0; iter < 5; iter++ {
		perm := rand.Perm(len(m.SetIDs))
		for _, sIdx := range perm {
			sID := m.SetIDs[sIdx]
			ds := m.DataSets[sID]
			
			for tIdx := 0; tIdx < m.Config.OutputSize; tIdx++ {
				currentPicked := m.Tuples[tIdx].Items[sIdx]
				delete(m.UsedItems[sID], currentPicked.ID)
				
				bestItem := currentPicked
				minDist := m.calcTupleDist(m.Tuples[tIdx].Items)
				
				for _, item := range ds.Items {
					if m.UsedItems[sID][item.ID] {
						continue
					}
					m.Tuples[tIdx].Items[sIdx] = item
					d := m.calcTupleDist(m.Tuples[tIdx].Items)
					if d < minDist {
						minDist = d
						bestItem = item
					}
				}
				m.Tuples[tIdx].Items[sIdx] = bestItem
				m.UsedItems[sID][bestItem.ID] = true
			}
		}
	}

	totalDist := 0.0
	for _, t := range m.Tuples {
		totalDist += m.calcTupleDist(t.Items)
	}
	m.BestSol = &Solution{
		Tuples:         m.cloneTuples(m.Tuples),
		TotalDistance: totalDist,
	}
	m.BestDist = totalDist
	fmt.Printf("Initial solution distance: %.4f\n", totalDist)
}

func (m *Matcher) cloneTuples(src []*Tuple) []*Tuple {
	dst := make([]*Tuple, len(src))
	for i, t := range src {
		dst[i] = &Tuple{Items: make([]*Item, len(t.Items))}
		copy(dst[i].Items, t.Items)
	}
	return dst
}

func (m *Matcher) Run() {
	m.Running = true
	m.StartTime = time.Now()
	
	// Handle termination signal (CLI only, GUI will use Stop())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-sigChan:
			fmt.Println("\nReceived termination signal. Saving best solution...")
			m.Stop()
		}
	}()

	m.SeedInitialSolution()
	
	// Clear state for DFS
	m.UsedItems = make(map[string]map[int]bool)
	for _, id := range m.SetIDs {
		m.UsedItems[id] = make(map[int]bool)
	}
	m.Tuples = make([]*Tuple, m.Config.OutputSize)
	for i := range m.Tuples {
		m.Tuples[i] = &Tuple{Items: make([]*Item, len(m.SetIDs))}
	}

	fmt.Println("Starting DFS...")
	m.dfs(0, 0)
	
	m.Running = false
	m.SaveOutput()
}

func (m *Matcher) Stop() {
	m.Running = false
}

type candidate struct {
	item *Item
	dist float64
}

func (m *Matcher) dfs(filledSlots int, runningDist float64) {
	if !m.Running {
		return
	}
	if filledSlots == m.Config.OutputSize*len(m.SetIDs) {
		if runningDist < m.BestDist {
			m.BestDist = runningDist
			m.BestSol = &Solution{
				Tuples:         m.cloneTuples(m.Tuples),
				TotalDistance: runningDist,
			}
			fmt.Printf("New best solution: %.4f (after %v)\n", runningDist, time.Since(m.StartTime))
			if m.OnUpdate != nil {
				m.OnUpdate(runningDist)
			}
		}
		return
	}

	if runningDist > m.BestDist*1.10 {
		return
	}

	type slotInfo struct {
		tIdx, sIdx int
		diff       float64
		candidates []candidate
	}

	var bestSlot *slotInfo

	for tIdx := 0; tIdx < m.Config.OutputSize; tIdx++ {
		for sIdx := 0; sIdx < len(m.SetIDs); sIdx++ {
			if m.Tuples[tIdx].Items[sIdx] != nil {
				continue
			}

			sID := m.SetIDs[sIdx]
			var cands []candidate
			for _, item := range m.DataSets[sID].Items {
				if m.UsedItems[sID][item.ID] {
					continue
				}
				m.Tuples[tIdx].Items[sIdx] = item
				d := m.calcTupleDist(m.Tuples[tIdx].Items)
				cands = append(cands, candidate{item, d})
			}
			m.Tuples[tIdx].Items[sIdx] = nil

			if len(cands) == 0 {
				return
			}

			sort.Slice(cands, func(i, j int) bool {
				return cands[i].dist < cands[j].dist
			})

			diff := 0.0
			if len(cands) > 1 {
				diff = cands[1].dist - cands[0].dist
			} else {
				diff = 1e9
			}

			if bestSlot == nil || diff > bestSlot.diff {
				bestSlot = &slotInfo{
					tIdx: tIdx, sIdx: sIdx,
					diff:     diff,
					candidates: cands,
				}
			}
		}
	}

	if bestSlot == nil {
		return
	}

	tIdx, sIdx := bestSlot.tIdx, bestSlot.sIdx
	sID := m.SetIDs[sIdx]

	for _, cand := range bestSlot.candidates {
		if !m.Running {
			return
		}
		m.Tuples[tIdx].Items[sIdx] = cand.item
		m.UsedItems[sID][cand.item.ID] = true
		
		newRunningDist := 0.0
		for _, t := range m.Tuples {
			newRunningDist += m.calcTupleDist(t.Items)
		}

		if newRunningDist <= m.BestDist*1.10 {
			m.dfs(filledSlots+1, newRunningDist)
		}
		
		m.UsedItems[sID][cand.item.ID] = false
		m.Tuples[tIdx].Items[sIdx] = nil
	}
}

func (m *Matcher) SaveOutput() {
	if m.BestSol == nil {
		fmt.Println("No solution found.")
		return
	}

	fmt.Printf("Saving best solution (dist: %.4f) to output files...\n", m.BestDist)
	
	for sIdx, sID := range m.SetIDs {
		fc := m.Config.InputFiles[sID]
		f, err := os.Create(fc.OutputPath)
		if err != nil {
			fmt.Printf("Error creating output file %s: %v\n", fc.OutputPath, err)
			continue
		}
		writer := bufio.NewWriter(f)
		for _, tuple := range m.BestSol.Tuples {
			item := tuple.Items[sIdx]
			fmt.Fprintln(writer, strings.Join(item.OriginalValues, "\t"))
		}
		writer.Flush()
		f.Close()
	}

	if m.Config.SummaryFile != "" {
		f, err := os.Create(m.Config.SummaryFile)
		if err != nil {
			fmt.Printf("Error creating summary file %s: %v\n", m.Config.SummaryFile, err)
			return
		}
		defer f.Close()
		fmt.Fprintf(f, "Best solution distance: %.4f\n", m.BestDist)
		fmt.Fprintf(f, "Matching Dimensions:\n")
		for _, dim := range m.Config.Dimensions {
			fmt.Fprintf(f, "  %s: %s (col %d) <-> %s (col %d)\n", dim.Name, dim.File1ID, dim.Col1, dim.File2ID, dim.Col2)
		}
	}
}
