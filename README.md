# crossword-solver
------------------
Writtn to solve crosswords as given in the book "Codeword puzzles" that can be bought at Bletchley Park.

These codeword puzzles are laid out in a 15x15 type of crossword, where each cell contains a number 1-26.
The crossword is to be filled with English words. To help in solving the puzzles a "reference box" is printed below the crossword puzzle, consisting of two rows (first 1-26, second has spaces to note down the correct letter. 
Per puzzle three entries for the refernce box are given.

As an example codeword puzzle 7 is provided. The provided reference box entries are 2-T , 10-E and 12-N

This program does an exhausitve search based on a wordlist words.txt from https://github.com/dwyl/english-words
