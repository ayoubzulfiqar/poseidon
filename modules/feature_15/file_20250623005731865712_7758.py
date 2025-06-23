def solve_maze(maze, start, end):
    rows = len(maze)
    cols = len(maze[0])
    
    directions = [(-1, 0), (1, 0), (0, -1), (0, 1)]

    visited = set()

    def dfs_solver(r, c, current_path):
        if not (0 <= r < rows and 0 <= c < cols):
            return None
        if maze[r][c] == '#':
            return None
        if (r, c) in visited:
            return None

        current_path.append((r, c))
        visited.add((r, c))

        if (r, c) == end:
            return list(current_path)

        for dr, dc in directions:
            nr, nc = r + dr, c + dc
            
            solution_path = dfs_solver(nr, nc, current_path)
            if solution_path:
                return solution_path

        current_path.pop()
        visited.remove((r, c))

        return None

    start_row, start_col = start
    return dfs_solver(start_row, start_col, [])

# Additional implementation at 2025-06-23 00:58:15
