import os

ROWS = 6
COLS = 7
EMPTY = ' '
PLAYER_ONE_PIECE = 'R'
PLAYER_TWO_PIECE = 'Y'

def create_board():
    return [[EMPTY for _ in range(COLS)] for _ in range(ROWS)]

def print_board(board):
    os.system('cls' if os.name == 'nt' else 'clear')
    
    col_numbers = " "
    for c in range(COLS):
        col_numbers += f" {c + 1}  "
    print(col_numbers)

    print("+" + "---+" * COLS)

    for r in range(ROWS - 1, -1, -1):
        row_str = "|"
        for c in range(COLS):
            row_str += f" {board[r][c]} |"
        print(row_str)
        print("+" + "---+" * COLS)

def is_valid_location(board, col):
    if not (0 <= col < COLS):
        return False
    return board[ROWS - 1][col] == EMPTY

def get_next_open_row(board, col):
    for r in range(ROWS):
        if board[r][col] == EMPTY:
            return r
    return -1

def drop_piece(board, row, col, piece):
    board[row][col] = piece

def check_win(board, piece):
    for c in range(COLS - 3):
        for r in range(ROWS):
            if board[r][c] == piece and board[r][c+1] == piece and board[r][c+2] == piece and board[r][c+3] == piece:
                return True

    for c in range(COLS):
        for r in range(ROWS - 3):
            if board[r][c] == piece and board[r+1][c] == piece and board[r+2][c] == piece and board[r+3][c] == piece:
                return True

    for c in range(COLS - 3):
        for r in range(ROWS - 3):
            if board[r][c] == piece and board[r+1][c+1] == piece and board[r+2][c+2] == piece and board[r+3][c+3] == piece:
                return True

    for c in range(COLS - 3):
        for r in range(3, ROWS):
            if board[r][c] == piece and board[r-1][c+1] == piece and board[r-2][c+2] == piece and board[r-3][c+3] == piece:
                return True
    return False

def is_board_full(board):
    for c in range(COLS):
        if board[ROWS - 1][c] == EMPTY:
            return False
    return True

def play_game():
    board = create_board()
    game_over = False
    turn = 0

    while not game_over:
        print_board(board)

        current_player_piece = PLAYER_ONE_PIECE if turn == 0 else PLAYER_TWO_PIECE
        player_num = 1 if turn == 0 else 2

        col = -1
        valid_input = False
        while not valid_input:
            try:
                col_input = input(f"Player {player_num} ({current_player_piece}), choose a column (1-{COLS}): ")
                col = int(col_input) - 1
                if is_valid_location(board, col):
                    valid_input = True
                else:
                    print("Invalid column. Column is full or out of range. Try again.")
            except ValueError:
                print("Invalid input. Please enter a number.")

        row = get_next_open_row(board, col)
        drop_piece(board, row, col, current_player_piece)

        if check_win(board, current_player_piece):
            print_board(board)
            print(f"Player {player_num} ({current_player_piece}) wins!")
            game_over = True
        elif is_board_full(board):
            print_board(board)
            print("It's a draw!")
            game_over = True
        
        turn = (turn + 1) % 2

if __name__ == "__main__":
    play_game()

# Additional implementation at 2025-06-21 02:23:27
import os

class ConnectFour:
    def __init__(self, rows=6, cols=7):
        self.rows = rows
        self.cols = cols
        self.board = [[' ' for _ in range(cols)] for _ in range(rows)]
        self.players = ['R', 'Y'] # Red and Yellow
        self.current_player_idx = 0
        self.last_move = None # (row, col) of the last dropped piece

    def _clear_screen(self):
        os.system('cls' if os.name == 'nt' else 'clear')

    def print_board(self):
        self._clear_screen()
        # Print column numbers
        print(" " + " ".join(str(i + 1) for i in range(self.cols)))
        print("+" + "---" * self.cols + "+")
        for r in range(self.rows):
            row_str = "|"
            for c in range(self.cols):
                piece = self.board[r][c]
                row_str += f" {piece} |"
            print(row_str)
            print("+" + "---" * self.cols + "+")
        print("\n")

    def drop_piece(self, col):
        col_idx = col - 1 # Convert 1-based input to 0-based index

        if not (0 <= col_idx < self.cols):
            print(f"Invalid column number. Please choose a column between 1 and {self.cols}.")
            return False

        for r in range(self.rows - 1, -1, -1): # Start from bottom row
            if self.board[r][col_idx] == ' ':
                self.board[r][col_idx] = self.players[self.current_player_idx]
                self.last_move = (r, col_idx)
                return True
        print(f"Column {col} is full. Please choose another column.")
        return False

    def check_win(self):
        if self.last_move is None:
            return False

        r, c = self.last_move
        piece = self.board[r][c]

        # Directions to check: (dr, dc)
        # Horizontal: (0, 1)
        # Vertical: (1, 0)
        # Diagonal / : (1, -1)
        # Diagonal \ : (1, 1)
        directions = [(0, 1), (1, 0), (1, -1), (1, 1)]

        for dr, dc in directions:
            count = 1 # Count the piece itself
            # Check in one direction
            for i in range(1, 4):
                nr, nc = r + dr * i, c + dc * i
                if 0 <= nr < self.rows and 0 <= nc < self.cols and self.board[nr][nc] == piece:
                    count += 1
                else:
                    break
            # Check in opposite direction
            for i in range(1, 4):
                nr, nc = r - dr * i, c - dc * i
                if 0 <= nr < self.rows and 0 <= nc < self.cols and self.board[nr][nc] == piece:
                    count += 1
                else:
                    break
            if count >= 4:
                return True
        return False

    def check_tie(self):
        for r in range(self.rows):
            for c in range(self.cols):
                if self.board[r][c] == ' ':
                    return False # Found an empty spot, not a tie
        return True # Board is full

    def switch_player(self):
        self.current_player_idx = (self.current_player_idx + 1) % len(self.players)

    def reset_game(self):
        self.board = [[' ' for _ in range(self.cols)] for _ in range(self.rows)]
        self.current_player_idx = 0
        self.last_move = None

    def play_game(self):
        while True: # Loop for playing multiple games
            self.reset_game()
            game_over = False
            while not game_over:
                self.print_board()
                current_player_piece = self.players[self.current_player_idx]
                print(f"Player {current_player_piece}'s turn.")

                valid_move = False
                while not valid_move:
                    try:
                        col_choice = int(input(f"Enter column number (1-{self.cols}): "))
                        valid_move = self.drop_piece(col_choice)
                    except ValueError:
                        print("Invalid input. Please enter a number.")

                if self.check_win():
                    self.print_board()
                    print(f"Player {current_player_piece} wins!")
                    game_over = True
                elif self.check_tie():
                    self.print_board()
                    print("It's a tie!")
                    game_over = True
                else:
                    self.switch_player()
            
            play_again = input("Play again? (yes/no): ").lower()
            if play_again != 'yes':
                break

if __name__ == "__main__":
    game = ConnectFour()
    game.play_game()

# Additional implementation at 2025-06-21 02:24:03
import os
import sys

BOARD_ROWS = 6
BOARD_COLS = 7
EMPTY_SLOT = ' '
PLAYER1_PIECE = 'R'
PLAYER2_PIECE = 'Y'

if sys.stdout.isatty():
    PLAYER1_COLOR = '\033[91m'
    PLAYER2_COLOR = '\033[93m'
    RESET_COLOR = '\033[0m'
    BOARD_COLOR = '\033[94m'
else:
    PLAYER1_COLOR = ''
    PLAYER2_COLOR = ''
    RESET_COLOR = ''
    BOARD_COLOR = ''

def clear_screen():
    if os.name == 'nt':
        _ = os.system('cls')
    else:
        _ = os.system('clear')

def create_board():
    return [[EMPTY_SLOT for _ in range(BOARD_COLS)] for _ in range(BOARD_ROWS)]

def print_board(board):
    print(BOARD_COLOR + " " + "   ".join(str(i + 1) for i in range(BOARD_COLS)) + RESET_COLOR)
    print(BOARD_COLOR + "+---" * BOARD_COLS + "+" + RESET_COLOR)
    for r in range(BOARD_ROWS):
        row_str = []
        for c in range(BOARD_COLS):
            piece = board[r][c]
            if piece == PLAYER1_PIECE:
                row_str.append(PLAYER1_COLOR + piece + RESET_COLOR)
            elif piece == PLAYER2_PIECE:
                row_str.append(PLAYER2_COLOR + piece + RESET_COLOR)
            else:
                row_str.append(EMPTY_SLOT)
        print(BOARD_COLOR + "| " + " | ".join(row_str) + " |" + RESET_COLOR)
        print(BOARD_COLOR + "+---" * BOARD_COLS + "+" + RESET_COLOR)

def is_valid_location(board, col):
    return 0 <= col < BOARD_COLS and board[0][col] == EMPTY_SLOT

def get_next_open_row(board, col):
    for r in range(BOARD_ROWS - 1, -1, -1):
        if board[r][col] == EMPTY_SLOT:
            return r
    return -1

def drop_piece(board, row, col, piece):
    board[row][col] = piece

def winning_move(board, piece):
    for c in range(BOARD_COLS - 3):
        for r in range(BOARD_ROWS):
            if board[r][c] == piece and board[r][c+1] == piece and board[r][c+2] == piece and board[r][c+3] == piece:
                return True

    for c in range(BOARD_COLS):
        for r in range(BOARD_ROWS - 3):
            if board[r][c] == piece and board[r+1][c] == piece and board[r+2][c] == piece and board[r+3][c] == piece:
                return True

    for c in range(BOARD_COLS - 3):
        for r in range(BOARD_ROWS - 3):
            if board[r][c] == piece and board[r+1][c+1] == piece and board[r+2][c+2] == piece and board[r+3][c+3] == piece:
                return True

    for c in range(BOARD_COLS - 3):
        for r in range(3, BOARD_ROWS):
            if board[r][c] == piece and board[r-1][c+1] == piece and board[r-2][c+2] == piece and board[r-3][c+3] == piece:
                return True
    return False

def check_draw(board):
    for r in range(BOARD_ROWS):
        for c in range(BOARD_COLS):
            if board[r][c] == EMPTY_SLOT:
                return False
    return True

def get_player_input(player_name):
    while True:
        try:
            col_choice = int(input(f"{player_name}, choose a column (1-{BOARD_COLS}): ")) - 1
            return col_choice
        except ValueError:
            print("Invalid input. Please enter a number.")

def main():
    player_names = ["", ""]
    player_scores = [0, 0]

    print("Welcome to Connect Four!")
    player_names[0] = input("Enter name for Player 1 (Red): ")
    player_names[1] = input("Enter name for Player 2 (Yellow): ")

    game_on = True
    while game_on:
        board = create_board()
        turn = 0
        game_over = False

        while not game_over:
            clear_screen()
            print_board(board)
            print(f"\nScores: {player_names[0]} ({PLAYER1_PIECE}): {player_scores[0]} | {player_names[1]} ({PLAYER2_PIECE}): {player_scores[1]}\n")

            current_player_name = player_names[turn]
            current_player_piece = PLAYER1_PIECE if turn == 0 else PLAYER2_PIECE

            col = get_player_input(current_player_name)

            if is_valid_location(board, col):
                row = get_next_open_row(board, col)
                drop_piece(board, row, col, current_player_piece)

                if winning_move(board, current_player_piece):
                    clear_screen()
                    print_board(board)
                    print(f"\n{current_player_name} wins! Congratulations!")
                    player_scores[turn] += 1
                    game_over = True
                elif check_draw(board):
                    clear_screen()
                    print_board(board)
                    print("\nIt's a draw!")
                    game_over = True
                else:
                    turn = (turn + 1) % 2
            else:
                print("Column is full or invalid. Try again.")

        play_again = input("Play again? (yes/no): ").lower()
        if play_again != 'yes':
            game_on = False
    
    clear_screen()
    print("Final Scores:")
    print(f"{player_names[0]} ({PLAYER1_PIECE}): {player_scores[0]}")
    print(f"{player_names[1]} ({PLAYER2_PIECE}): {player_scores[1]}")
    print("Thanks for playing!")

if __name__ == "__main__":
    main()