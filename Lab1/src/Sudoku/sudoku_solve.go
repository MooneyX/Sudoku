package main

// /*
// #cgo CFLAGS: -I.
// #cgo LDFLAGS: -L.
// #include <stdio.h>
// #include <stdbool.h>
// #include <stdlib.h>
// #include <string.h>
// #include "sudoku.h"
// */
// import "C"
import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
)

type Puzzle struct {
	puzzle string
	id     int
}

var (
	threadNum     = flag.Int("thread_num", 12, "num of threads")
	inNUm     int = 0
	outNum    int = 0
	waitGroup sync.WaitGroup
	inCh      = make(chan Puzzle)
	outCh     = make(chan Puzzle)
)

func input_thread() {
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

func worker_thread(threadId int) {
	defer waitGroup.Done()
	// puzzleChar := C.malloc(C.size_t(82 * C.sizeof_char))
	// defer C.free(unsafe.Pointer(puzzleChar))
	solveCnt := 0
	for {
		puzzleWithId, ok := <-inCh
		if ok {
			puzzle := puzzleWithId.puzzle
			id := puzzleWithId.id
			for i := 0; i < 81; i++ {
				c := puzzle[i]
				// slog.Debug("worker_thread", "", threadId, "", c)
				Board[threadId][i] = int(c - '0')
			}
			// C.memcpy(puzzleChar, unsafe.Pointer(C.CString(puzzle)), C.size_t(82))
			// CStr := (*C.char)(puzzleChar)
			// C.input(CStr, C.int(threadId))
			// ok := C.solve_sudoku_dancing_links(0, C.int(threadId))
			// fmt.Println(Board[threadId])
			SolveSudokuDancingLinks(0, threadId)
			// fmt.Println(Board[threadId])
			if !ok {
				slog.Error("worker_thread", "", id, "solve sudoku", "failed")
			}
			// for i := 0; i < 81; i++ {
			// 	fmt.Printf("%d", C.board[i])
			// }
			// fmt.Println()
			var result strings.Builder
			for i := 0; i < 81; i++ {
				// result.WriteString(fmt.Sprintf("%d", C.Board[threadId][i]))
				result.WriteString(fmt.Sprintf("%d", Board[threadId][i]))
			}
			// slog.Info("worker_thread", "_", threadId, "puzzle_id", id, "ok", ok, "result", result.String())
			solveCnt++
			// slog.Info("worker_thread", "_", threadId, "puzzle_id", id, "solveCnt", solveCnt)
			outCh <- Puzzle{result.String(), id}
		}
	}
}

func output_thread() {
	// 打开或创建一个文件用于写入
	// file, err := os.OpenFile("./test_group_Folder/Basic_Result", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()
	defer waitGroup.Done()
	buffer := make(map[int]string)
	// loopCnt := 0
	for {
		// loopCnt++
		// slog.Info("output_thread", "loopCnt", loopCnt)
		select {
		case resultWithId := <-outCh:
			buffer[resultWithId.id] = resultWithId.puzzle
			// slog.Info("output_thread", "resultWithId", resultWithId)
		default:
			// slog.Info("output_thread", "heap", "empty")
		}
		for {
			result, ok := buffer[outNum]
			if ok {
				fmt.Println(result)
				outNum++
			} else {
				break
			}
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
	waitGroup.Add(1)
	go input_thread()
	workerNum := *threadNum - 2
	for i := 0; i < workerNum; i++ {
		waitGroup.Add(1)
		go worker_thread(i)
	}
	waitGroup.Add(1)
	go output_thread()
	defer waitGroup.Wait()
}
