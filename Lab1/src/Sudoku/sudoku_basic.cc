#include <assert.h>
#include <stdio.h>

#include <algorithm>

#include "sudoku.h"

int Board[SOLVEPTHNUM][N];
int spaces[N];
int nspaces;
//int (*chess)[COL] = (int (*)[COL])board;

static void find_spaces(int tid)
{
	int *board = Board[tid];
	nspaces = 0;
	for (int cell = 0; cell < N; ++cell) {
		if (board[cell] == 0)
			spaces[nspaces++] = cell;
	}
}

void input(const char in[N], int tid)
{
	int *board = Board[tid];
	for (int cell = 0; cell < N; ++cell) {
		board[cell] = in[cell] - '0';
		assert(0 <= board[cell] && board[cell] <= NUM);
	}
	find_spaces(tid);
}