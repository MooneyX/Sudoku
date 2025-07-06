# 定义变量
CXX = g++
CXXFLAGS = -std=c++11 -fPIC
LDFLAGS = -shared
CGO_LDFLAGS = -L. -lsudoku
CGO_CXXFLAGS = -I.

# 定义目标
all: libsudoku.so sudoku_solve

# 编译 C++ 文件为共享库
libsudoku.so: sudoku_dancing_links.o sudoku_basic.o
	$(CXX) $(LDFLAGS) -o $@ $^

# 编译 C++ 源文件为对象文件
%.o: %.cc
	$(CXX) $(CXXFLAGS) -c $< -o $@

# 编译 Go 程序
sudoku_solve: sudoku_solve.go
	CGO_LDFLAGS="$(CGO_LDFLAGS)" CGO_CXXFLAGS="$(CGO_CXXFLAGS)" go build -o $@ $<

# 清理生成的文件
clean:
	rm -f *.o libsudoku.so sudoku_solve