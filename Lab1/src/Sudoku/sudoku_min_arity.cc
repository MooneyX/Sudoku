#include <assert.h>

#include <algorithm>

#include "sudoku.h"

extern "C" {

static int arity(int cell, int tid)
{
  int *board = Board[tid];
  bool occupied[10] = {false};
  for (int i = 0; i < NEIGHBOR; ++i) {
    int neighbor = neighbors[cell][i];
    occupied[board[neighbor]] = true;
  }
  return std::count(occupied+1, occupied+10, false);
}

static void find_min_arity(int space, int tid)
{
  int cell = spaces[space];
  int min_space = space;
  int min_arity = arity(cell, tid);

  for (int sp = space+1; sp < nspaces && min_arity > 1; ++sp) {
    int cur_arity = arity(spaces[sp], tid);
    if (cur_arity < min_arity) {
      min_arity = cur_arity;
      min_space = sp;
    }
  }

  if (space != min_space) {
    std::swap(spaces[min_space], spaces[space]);
  }
}

bool solve_sudoku_min_arity(int which_space, int tid)
{
  int *board = Board[tid];
  if (which_space >= nspaces) {
    return true;
  }

  find_min_arity(which_space, tid);
  int cell = spaces[which_space];

  for (int guess = 1; guess <= NUM; ++guess) {
    if (available(guess, cell, tid)) {
      // hold
      assert(board[cell] == 0);
      board[cell] = guess;

      // try
      if (solve_sudoku_min_arity(which_space+1, tid)) {
        return true;
      }

      // unhold
      assert(board[cell] == guess);
      board[cell] = 0;
    }
  }
  return false;
}

}