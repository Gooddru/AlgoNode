from fastapi import FastAPI
from pydantic import BaseModel
import networkx as nx

app = FastAPI(title="OptiGraph Python Worker")

class GraphData(BaseModel):
    edges: list[list[str]] # Ожидаем список ребер: [["A", "B"], ["B", "C"]]

@app.post("/api/v1/cycle")
def detect_cycle(data: GraphData):
    """Поиск циклов в ориентированном графе"""
    G = nx.DiGraph(data.edges)
    try:
        cycle = nx.find_cycle(G, orientation="original")
        return {"has_cycle": True, "cycle": cycle}
    except nx.NetworkXNoCycle:
        return {"has_cycle": False, "cycle": []}

@app.post("/api/v1/tikz")
def generate_tikz(data: GraphData):
    """Генерация LaTeX/TikZ кода для визуализации графа"""
    G = nx.DiGraph(data.edges)
    
    tikz_lines = [
        "\\begin{tikzpicture}[->, >=stealth, auto, node distance=2cm, thick, main node/.style={circle,draw,font=\\sffamily\\Large\\bfseries}]"
    ]
    
    # Упрощенная генерация узлов по кругу для демонстрации
    nodes = list(G.nodes())
    for i, node in enumerate(nodes):
        angle = 360 / len(nodes) * i
        tikz_lines.append(f"  \\node[main node] ({node}) at ({angle}:2) {{{node}}};")
        
    for edge in G.edges():
        tikz_lines.append(f"  \\path[every node/.style={{font=\\sffamily\\small}}] ({edge[0]}) edge node {{}} ({edge[1]});")
        
    tikz_lines.append("\\end{tikzpicture}")
    
    return {"latex_code": "\n".join(tikz_lines)}