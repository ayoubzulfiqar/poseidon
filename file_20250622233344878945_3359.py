import math

PLAYER_X = 'X'
PLAYER_O = 'O'
EMPTY = ' '

def create_board():
    return [[EMPTY, EMPTY, EMPTY],
            [EMPTY, EMPTY, EMPTY],
            [EMPTY, EMPTY, EMPTY]]

def print_board(board):
    for row in board:
        print("|".join(row))
        print("-" * 5)

def check_win(board, player):
    # Check rows
    for row in board:
        if all(s == player for s in row):
            return True
    # Check columns
    for col in range(3):
        if all(board[row][col] == player for row in range(3)):
            return True
    # Check diagonals
    if all(board[i][i] == player for i in range(3)):
        return True
    if all(board[i][2 - i] == player for i in range(3)):
        return True
    return False

def check_draw(board):
    for row in board:
        for cell in row:
            if cell == EMPTY:
                return False
    return True

def is_game_over(board):
    return check_win(board, PLAYER_X) or check_win(board, PLAYER_O) or check_draw(board)

def get_empty_cells(board):
    cells = []
    for r in range(3):
        for c in range(3):
            if board[r][c] == EMPTY:
                cells.append((r, c))
    return cells

def evaluate(board):
    if check_win(board, PLAYER_O):
        return 10  # AI wins
    elif check_win(board, PLAYER_X):
        return -10 # Human wins
    else:
        return 0   # Draw or game in progress

def minimax(board, depth, is_maximizing):
    score = evaluate(board)

    if score == 10:
        return score - depth # Prioritize faster wins
    if score == -10:
        return score + depth # Penalize slower losses
    if check_draw(board):
        return 0

    if is_maximizing:
        best_score = -math.inf
        for r, c in get_empty_cells(board):
            board[r][c] = PLAYER_O
            best_score = max(best_score, minimax(board, depth + 1, False))
            board[r][c] = EMPTY  # Undo move
        return best_score
    else:
        best_score = math.inf
        for r, c in get_empty_cells(board):
            board[r][c] = PLAYER_X
            best_score = min(best_score, minimax(board, depth + 1, True))
            board[r][c] = EMPTY  # Undo move
        return best_score

def find_best_move(board):
    best_score = -math.inf
    best_move = (-1, -1)

    for r, c in get_empty_cells(board):
        board[r][c] = PLAYER_O
        score = minimax(board, 0, False)
        board[r][c] = EMPTY  # Undo move

        if score > best_score:
            best_score = score
            best_move = (r, c)
    return best_move

def play_game():
    board = create_board()
    current_player = PLAYER_X # Human starts

    while not is_game_over(board):
        print_board(board)
        if current_player == PLAYER_X:
            while True:
                try:
                    row = int(input("Enter row (0-2): "))
                    col = int(input("Enter column (0-2): "))
                    if 0 <= row <= 2 and 0 <= col <= 2 and board[row][col] == EMPTY:
                        board[row][col] = PLAYER_X
                        break
                    else:
                        print("Invalid move. Try again.")
                except ValueError:
                    print("Invalid input. Please enter numbers.")
            current_player = PLAYER_O
        else: # AI's turn
            print("AI is making a move...")
            row, col = find_best_move(board)
            board[row][col] = PLAYER_O
            current_player = PLAYER_X
    
    print_board(board)
    if check_win(board, PLAYER_X):
        print("Congratulations! You won!")
    elif check_win(board, PLAYER_O):
        print("AI wins! You lost!")
    else:
        print("It's a draw!")

play_game()

# Additional implementation at 2025-06-22 23:34:06
import math

class TicTacToe:
    def __init__(self):
        self.board = [' ' for _ in range(9)]
        self.human_player = 'X'
        self.ai_player = 'O'
        self.current_player = 'X' # Default starting player

    def print_board(self):
        print("-------------")
        for i in range(3):
            print(f"| {self.board[i*3]} | {self.board[i*3+1]} | {self.board[i*3+2]} |")
            print("-------------")

    def is_board_full(self, board_state=None):
        if board_state is None:
            board_state = self.board
        return ' ' not in board_state

    def check_win(self, player, board_state=None):
        if board_state is None:
            board_state = self.board

        # Check rows
        for i in range(3):
            if all(board_state[i*3 + j] == player for j in range(3)):
                return True
        # Check columns
        for i in range(3):
            if all(board_state[i + j*3] == player for j in range(3)):
                return True
        # Check diagonals
        if (board_state[0] == player and board_state[4] == player and board_state[8] == player) or \
           (board_state[2] == player and board_state[4] == player and board_state[6] == player):
            return True
        return False

    def get_empty_cells(self, board_state=None):
        if board_state is None:
            board_state = self.board
        return [i for i, spot in enumerate(board_state) if spot == ' ']

    def minimax(self, board, depth, is_maximizing):
        if self.check_win(self.ai_player, board):
            return 10 - depth
        if self.check_win(self.human_player, board):
            return -10 + depth
        if self.is_board_full(board):
            return 0

        if is_maximizing:
            best_score = -math.inf
            for cell in self.get_empty_cells(board):
                board[cell] = self.ai_player
                score = self.minimax(board, depth + 1, False)
                board[cell] = ' ' # Undo the move
                best_score = max(best_score, score)
            return best_score
        else:
            best_score = math.inf
            for cell in self.get_empty_cells(board):
                board[cell] = self.human_player
                score = self.minimax(board, depth + 1, True)
                board[cell] = ' ' # Undo the move
                best_score = min(best_score, score)
            return best_score

    def find_best_move(self):
        best_score = -math.inf
        best_move = -1
        temp_board = list(self.board) # Create a copy of the current board for simulation

        for cell in self.get_empty_cells():
            temp_board[cell] = self.ai_player
            score = self.minimax(temp_board, 0, False) # AI makes a move, then it's human's turn (minimizing)
            temp_board[cell] = ' ' # Undo the move
            if score > best_score:
                best_score = score
                best_move = cell
        return best_move

    def play_game(self):
        print("Welcome to Tic-Tac-Toe!")
        print("Enter a number 1-9 to place your mark (top-left is 1, bottom-right is 9).")

        while True:
            choice = input("Do you want to go first? (y/n): ").lower()
            if choice == 'y':
                self.human_player = 'X'
                self.ai_player = 'O'
                self.current_player = self.human_player
                print(f"You are '{self.human_player}', AI is '{self.ai_player}'. You go first.")
                break
            elif choice == 'n':
                self.human_player = 'O'
                self.ai_player = 'X'
                self.current_player = self.ai_player
                print(f"You are '{self.human_player}', AI is '{self.ai_player}'. AI goes first.")
                break
            else:
                print("Invalid choice. Please enter 'y' or 'n'.")

        while True:
            self.print_board()

            if self.current_player == self.human_player:
                move = -1
                while True:
                    try:
                        move_input = input(f"Player {self.human_player}, enter your move (1-9): ")
                        move = int(move_input) - 1
                        if 0 <= move < 9 and self.board[move] == ' ':
                            self.board[move] = self.human_player
                            break
                        else:
                            print("Invalid move. That spot is taken or out of range. Try again.")
                    except ValueError:
                        print("Invalid input. Please enter a number.")
            else: # AI's turn
                print(f"AI ({self.ai_player}) is thinking...")
                move = self.find_best_move()
                self.board[move] = self.ai_player
                print(f"AI placed '{self.ai_player}' at position {move + 1}.")

            if self.check_win(self.current_player):
                self.print_board()
                print(f"Player {self.current_player} wins!")
                break
            elif self.is_board_full():
                self.print_board()
                print("It's a draw!")
                break
            else:
                self.current_player = self.ai_player if self.current_player == self.human_player else self.human_player

if __name__ == "__main__":
    game = TicTacToe()
    game.play_game()