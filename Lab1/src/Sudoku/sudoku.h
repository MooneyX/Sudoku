#ifndef SUDOKU_H
#define SUDOKU_H

#ifdef __cplusplus
extern "C" {
#endif

#define SOLVEPTHNUM 20

const bool DEBUG_MODE = false;
enum { ROW=9, COL=9, N = 81, NEIGHBOR = 20 };
const int NUM = 9;

extern int Board[SOLVEPTHNUM][N];
extern int spaces[N];
extern int nspaces;
extern int (*chess)[COL];

void input(const char in[N], int tid);


bool solve_sudoku_dancing_links(int unused, int tid);
bool solved();

#ifdef __cplusplus
}
#endif

#endif
