# 定义变量
CXX = g++
CXXFLAGS = -std=c++11 -fPIC
LDFLAGS = -shared
CGO_LDFLAGS = -L. -lsudoku
CGO_CXXFLAGS = -I.

# 定义目标
all: sudoku_solve

# 编译 Go 程序
sudoku_solve: sudoku_solve.go sudoku_dancing_links.go
	go build -o $@ *.go

# 清理生成的文件
clean:
	rm -f sudoku_solve