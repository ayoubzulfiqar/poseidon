import os
import time
import sys
import random

try:
    import termios
    import tty
    import select

    def get_input_char_non_blocking():
        fd = sys.stdin.fileno()
        old_settings = termios.tcgetattr(fd)
        try:
            tty.setcbreak(fd)
            if select.select([sys.stdin], [], [], 0) == ([sys.stdin], [], []):
                return sys.stdin.read(1)
        finally:
            termios.tcsetattr(fd, termios.TCSADRAIN, old_settings)
        return None
except ImportError:
    def get_input_char_non_blocking():
        return None

BOARD_WIDTH = 40
BOARD_HEIGHT = 25
FPS = 10

EMPTY = ' '
WALL = '#'
BALL = 'O'
BUMPER = '*'
DRAIN = 'V'

board = []
ball_pos = [0, 0]
ball_vel = [0, 0]
score = 0
lives = 3
game_over = False

def clear_screen():
    os.system('cls' if os.name == 'nt' else 'clear')

def initialize_board():
    global board, ball_pos, ball_vel, score, lives, game_over

    board = [[EMPTY for _ in range(BOARD_WIDTH)] for _ in range(BOARD_HEIGHT)]

    for i in range(BOARD_HEIGHT):
        board[i][0] = WALL
        board[i][BOARD_WIDTH - 1] = WALL
    for j in range(BOARD_WIDTH):
        board[0][j] = WALL
        board[BOARD_HEIGHT - 1][j] = WALL

    for j in range(1, BOARD_WIDTH - 1):
        board[BOARD_HEIGHT - 1][j] = DRAIN

    board[5][BOARD_WIDTH // 2 - 5] = BUMPER
    board[5][BOARD_WIDTH // 2 + 5] = BUMPER
    board[10][BOARD_WIDTH // 2] = BUMPER
    board[15][BOARD_WIDTH // 2 - 8] = BUMPER
    board[15][BOARD_WIDTH // 2 + 8] = BUMPER

    reset_ball()
    score = 0
    lives = 3
    game_over = False

def reset_ball():
    global ball_pos, ball_vel
    ball_pos = [BOARD_WIDTH // 2, BOARD_HEIGHT - 3]
    ball_vel = [random.choice([-1, 1]), -1]

def draw_board():
    board_copy = [row[:] for row in board]
    if 0 <= ball_pos[1] < BOARD_HEIGHT and 0 <= ball_pos[0] < BOARD_WIDTH:
        board_copy[ball_pos[1]][ball_pos[0]] = BALL

    for row in board_copy:
        print("".join(row))

    print(f"Score: {score} | Lives: {lives}")
    if game_over:
        print("GAME OVER! Press 'r' to restart or 'q' to quit.")
    else:
        print("Press 'a' for left flipper, 'd' for right flipper (if supported).")

def update_game_state(input_char):
    global ball_pos, ball_vel, score, lives, game_over

    if game_over:
        if input_char == 'r':
            initialize_board()
        elif input_char == 'q':
            sys.exit()
        return

    ball_pos[0] += ball_vel[0]
    ball_pos[1] += ball_vel[1]

    if ball_pos[0] <= 0 or ball_pos[0] >= BOARD_WIDTH - 1:
        ball_vel[0] *= -1
        ball_pos[0] = max(1, min(ball_pos[0], BOARD_WIDTH - 2))
    if ball_pos[1] <= 0:
        ball_vel[1] *= -1
        ball_pos[1] = 1

    current_char = board[ball_pos[1]][ball_pos[0]]

    if current_char == BUMPER:
        score += 100
        ball_vel[0] *= -1
        ball_vel[1] *= -1
    elif current_char == DRAIN:
        lives -= 1
        if lives <= 0:
            game_over = True
        else:
            reset_ball()
        return

    left_flipper_zone_y = BOARD_HEIGHT - 4
    right_flipper_zone_y = BOARD_HEIGHT - 4
    left_flipper_zone_x_start = BOARD_WIDTH // 2 - 10
    left_flipper_zone_x_end = BOARD_WIDTH // 2 - 3
    right_flipper_zone_x_start = BOARD_WIDTH // 2 + 3
    right_flipper_zone_x_end = BOARD_WIDTH // 2 + 10

    if ball_pos[1] == left_flipper_zone_y and \
       left_flipper_zone_x_start <= ball_pos[0] <= left_flipper_zone_x_end:
        if input_char == 'a':
            ball_vel[0] = 1
            ball_vel[1] = -abs(ball_vel[1]) - 1
            score += 50
    elif ball_pos[1] == right_flipper_zone_y and \
         right_flipper_zone_x_start <= ball_pos[0] <= right_flipper_zone_x_end:
        if input_char == 'd':
            ball_vel[0] = -1
            ball_vel[1] = -abs(ball_vel[1]) - 1
            score += 50

    if ball_vel[1] < 1:
        ball_vel[1] += 0.1

    if abs(ball_vel[0]) < 0.5: ball_vel[0] = random.choice([-1, 1])
    if abs(ball_vel[1]) < 0.5: ball_vel[1] = -1
    if abs(ball_vel[0]) > 2: ball_vel[0] = 2 * (1 if ball_vel[0] > 0 else -1)
    if abs(ball_vel[1]) > 2: ball_vel[1] = 2 * (1 if ball_vel[1] > 0 else -1)

    ball_vel[0] = int(ball_vel[0])
    ball_vel[1] = int(ball_vel[1])

def game_loop():
    initialize_board()
    while True:
        input_char = get_input_char_non_blocking()
        clear_screen()
        update_game_state(input_char)
        draw_board()
        time.sleep(1 / FPS)

game_loop()

# Additional implementation at 2025-06-18 01:06:08
