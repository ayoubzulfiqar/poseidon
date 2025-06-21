def find_connected_components(graph):
    visited = set()
    components = []

    all_nodes = set()
    for node, neighbors in graph.items():
        all_nodes.add(node)
        for neighbor in neighbors:
            all_nodes.add(neighbor)

    for node in all_nodes:
        if node not in visited:
            current_component = []
            stack = [node]
            visited.add(node)

            while stack:
                u = stack.pop()
                current_component.append(u)

                for neighbor in graph.get(u, []):
                    if neighbor not in visited:
                        visited.add(neighbor)
                        stack.append(neighbor)
            components.append(current_component)
    return components

# Additional implementation at 2025-06-21 03:55:26
class Graph:
    def __init__(self):
        self.adj = {}

    def add_node(self, node):
        if node not in self.adj:
            self.adj[node] = set()

    def add_edge(self, u, v):
        self.add_node(u)
        self.add_node(v)
        self.adj[u].add(v)
        self.adj[v].add(u)

    def get_nodes(self):
        return list(self.adj.keys())

    def get_neighbors(self, node):
        return list(self.adj.get(node, set()))

class ConnectedComponentsFinder:
    def __init__(self, graph):
        self.graph = graph
        self._components = None

    def _dfs(self, start_node, visited, current_component):
        visited.add(start_node)
        current_component.append(start_node)
        for neighbor in self.graph.get_neighbors(start_node):
            if neighbor not in visited:
                self._dfs(neighbor, visited, current_component)

    def find_components(self):
        if self._components is not None:
            return self._components

        components = []
        visited = set()
        for node in self.graph.get_nodes():
            if node not in visited:
                current_component = []
                self._dfs(node, visited, current_component)
                components.append(sorted(current_component))
        self._components = components
        return components

    def get_number_of_components(self):
        return len(self.find_components())

    def is_graph_fully_connected(self):
        return self.get_number_of_components() == 1

    def is_connected(self, node1, node2):
        if node1 not in self.graph.get_nodes() or node2 not in self.graph.get_nodes():
            return False

        if node1 == node2:
            return True

        components = self.find_components()
        for component in components:
            if node1 in component and node2 in component:
                return True
        return False

    def get_component_of_node(self, node):
        if node not in self.graph.get_nodes():
            return None

        components = self.find_components()
        for component in components:
            if node in component:
                return sorted(component)
        return None

# Additional implementation at 2025-06-21 03:55:55


# Additional implementation at 2025-06-21 03:56:41
class Graph:
    def __init__(self):
        self.graph = {}

    def add_edge(self, u, v):
        if u not in self.graph:
            self.graph[u] = []
        if v not in self.graph:
            self.graph[v] = []
        self.graph[u].append(v)
        self.graph[v].append(u)

    def add_node(self, node):
        if node not in self.graph:
            self.graph[node] = []

    def get_nodes(self):
        return list(self.graph.keys())

    def find_connected_components(self):
        visited = set()
        components = []
        node_to_component_id = {}
        component_id_counter = 0

        for node in self.get_nodes():
            if node not in visited:
                current_component = []
                stack = [node]
                visited.add(node)

                while stack:
                    vertex = stack.pop()
                    current_component.append(vertex)
                    node_to_component_id[vertex] = component_id_counter

                    for neighbor in self.graph.get(vertex, []):
                        if neighbor not in visited:
                            visited.add(neighbor)
                            stack.append(neighbor)

                components.append(current_component)
                component_id_counter += 1

        return components, node_to_component_id

if __name__ == "__main__":
    g1 = Graph()
    g1.add_edge(0, 1)
    g1.add_edge(1, 2)
    g1.add_edge(3, 4)
    g1.add_node(5)
    g1.add_edge(6, 7)
    g1.add_edge(7, 8)
    g1.add_edge(6, 8)

    print("--- Graph 1 ---")
    components1, node_map1 = g1.find_connected_components()
    print("Connected Components:", components1)
    print("Node to Component ID Map:", node_map1)
    print(f"Number of Components: {len(components1)}\n")

    g2 = Graph()
    g2.add_edge('A', 'B')
    g2.add_edge('B', 'C')
    g2.add_edge('C', 'A')
    g2.add_edge('C', 'D')

    print("--- Graph 2 ---")
    components2, node_map2 = g2.find_connected_components()
    print("Connected Components:", components2)
    print("Node to Component ID Map:", node_map2)
    print(f"Number of Components: {len(components2)}\n")

    g3 = Graph()
    g3.add_node('X')
    g3.add_node('Y')
    g3.add_node('Z')

    print("--- Graph 3 ---")
    components3, node_map3 = g3.find_connected_components()
    print("Connected Components:", components3)
    print("Node to Component ID Map:", node_map3)
    print(f"Number of Components: {len(components3)}\n")

    g4 = Graph()

    print("--- Graph 4 ---")
    components4, node_map4 = g4.find_connected_components()
    print("Connected Components:", components4)
    print("Node to Component ID Map:", node_map4)
    print(f"Number of Components: {len(components4)}\n")