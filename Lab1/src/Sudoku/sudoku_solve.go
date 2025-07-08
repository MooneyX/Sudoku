package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L.
#include <stdio.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>
#include "sudoku.h"
*/
import "C"
import (
	"container/heap"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
)

var threadNum = flag.Int("thread_num", 12, "num of threads")
var (
	inNUm  int = 0
	outNum int = 0
)

type Puzzle struct {
	puzzle string
	id     int
}
type PuzzleHeap []Puzzle

func (h PuzzleHeap) Len() int           { return len(h) }
func (h PuzzleHeap) Less(i, j int) bool { return h[i].id < h[j].id }
func (h PuzzleHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push 将元素添加到堆中
func (h *PuzzleHeap) Push(x interface{}) {
	*h = append(*h, x.(Puzzle))
}

// Pop 从堆中移除并返回最小元素
func (h *PuzzleHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Peek 查看堆顶元素而不取出
func (h *PuzzleHeap) Peek() (interface{}, bool) {
	if h.Len() == 0 {
		return nil, false
	}
	return (*h)[0], true
}

func input_thread(waitGroup *sync.WaitGroup, inCh chan Puzzle) {
	defer waitGroup.Done()
	for {
		var filename string
		fmt.Scan(&filename)
		if filename == "end" {
			break
		}
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("file", filename, "read failed:", err)
			continue
		}
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if len(line) > 0 {
				// slog.Info("input_thread", "line", line)
				inCh <- Puzzle{line, inNUm}
				inNUm++
				// slog.Info("input_thread", "inNUm", inNUm)
			}
		}
	}
}

func worker_thread(waitGroup *sync.WaitGroup, inCh chan Puzzle, outCh chan Puzzle, threadId int) {
	defer waitGroup.Done()
	solveCnt := 0
	for {
		puzzleWithId, ok := <-inCh
		if ok {
			puzzle := puzzleWithId.puzzle
			id := puzzleWithId.id
			for i := 0; i < 81; i++ {
				c := puzzle[i]
				// slog.Debug("worker_thread", "", threadId, "", c)
				C.Board[threadId][i] = C.int(c - '0')
			}
			ok := C.solve_sudoku_dancing_links(0, C.int(threadId))
			if !ok {
				slog.Error("worker_thread", "", id, "solve sudoku", "failed")
			}
			// for i := 0; i < 81; i++ {
			// 	fmt.Printf("%d", C.board[i])
			// }
			// fmt.Println()
			var result strings.Builder
			for i := 0; i < 81; i++ {
				result.WriteString(fmt.Sprintf("%d", C.Board[threadId][i]))
			}
			// slog.Info("worker_thread", "_", threadId, "puzzle_id", id, "ok", ok, "result", result.String())
			solveCnt++
			slog.Info("worker_thread", "_", threadId, "puzzle_id", id, "solveCnt", solveCnt)
			outCh <- Puzzle{result.String(), id}
		}
	}
}

func output_thread(waitGroup *sync.WaitGroup, outCh chan Puzzle) {
	// 打开或创建一个文件用于写入
	// file, err := os.OpenFile("./test_group_Folder/Basic_Result", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()
	defer waitGroup.Done()
	resultHeap := &PuzzleHeap{}
	// loopCnt := 0
	for {
		// loopCnt++
		// slog.Info("output_thread", "loopCnt", loopCnt)
		select {
		case resultWithId := <-outCh:
			heap.Push(resultHeap, resultWithId)
			// slog.Info("output_thread", "resultWithId", resultWithId)
		default:
			// slog.Info("output_thread", "heap", "empty")
		}
		peek, ok1 := resultHeap.Peek()
		// slog.Info("output_thread", "peek", peek, "ok1", ok1, "outNum", outNum, "heapLen", resultHeap.Len())
		if ok1 && peek.(Puzzle).id == outNum {
			pop := heap.Pop(resultHeap).(Puzzle)
			// fmt.Fprintln(file, pop.puzzle)
			fmt.Println(pop.puzzle)
			outNum++
		}
	}
}

func main() {

	// fmt.Println("Start")
	// defer fmt.Println("End")
	flag.Parse()
	if *threadNum < 3 {
		slog.Error("main", "args errer", "threadNum must > 2")
		return
	}
	var waitGroup sync.WaitGroup
	inCh := make(chan Puzzle)
	outCh := make(chan Puzzle)
	waitGroup.Add(1)
	go input_thread(&waitGroup, inCh)
	workerNum := *threadNum - 3
	for i := 0; i < workerNum; i++ {
		waitGroup.Add(1)
		go worker_thread(&waitGroup, inCh, outCh, i)
	}
	waitGroup.Add(1)
	go output_thread(&waitGroup, outCh)
	waitGroup.Add(1)
	go output_thread(&waitGroup, outCh)
	defer waitGroup.Wait()
}
