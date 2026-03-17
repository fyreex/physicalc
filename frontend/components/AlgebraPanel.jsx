import { useState } from "react";
import ResultBox from "./ResultBox.jsx";

const card = {
  background: "var(--surface)", border: "1px solid var(--border)",
  borderRadius: "var(--radius)", padding: "1.5rem", marginBottom: "1.5rem",
};
const label = { display: "block", fontSize: 12, color: "var(--muted)", marginBottom: 6, marginTop: 12 };
const btn = {
  padding: "0.55rem 1.4rem", border: "none", borderRadius: 8,
  background: "#00d4aa", color: "#0f1117", fontWeight: 600,
  fontSize: 14, marginTop: 16,
};
const hint = { fontSize: 11, color: "var(--muted)", marginTop: 4 };

async function callApi(path, body) {
  const r = await fetch(`/api${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return r.json();
}

function parseVec(s) {
  return s.split(",").map(Number);
}

function parseMat(s) {
  return s.split(";").map(row => row.split(",").map(Number));
}

export default function AlgebraPanel() {
  // Vetores
  const [vecOp, setVecOp]   = useState("dot");
  const [vecA, setVecA]     = useState("1,2,3");
  const [vecB, setVecB]     = useState("4,5,6");
  const [vecRes, setVecRes] = useState(null);
  const [vecErr, setVecErr] = useState(null);
  const [vecLoad, setVecL]  = useState(false);

  // Matrizes
  const [matOp, setMatOp]   = useState("det");
  const [matA, setMatA]     = useState("1,2;3,4");
  const [matB, setMatB]     = useState("5,6;7,8");
  const [matRes, setMatRes] = useState(null);
  const [matErr, setMatErr] = useState(null);
  const [matLoad, setMatL]  = useState(false);

  async function calcVector() {
    setVecL(true); setVecErr(null);
    try {
      const data = await callApi("/algebra/vector", {
        operation: vecOp,
        a: parseVec(vecA),
        b: parseVec(vecB),
      });
      if (data.error) setVecErr(data.error); else setVecRes(data);
    } catch { setVecErr("Erro de conexão"); }
    finally { setVecL(false); }
  }

  async function calcMatrix() {
    setMatL(true); setMatErr(null);
    try {
      const data = await callApi("/algebra/matrix", {
        operation: matOp,
        a: parseMat(matA),
        b: parseMat(matB),
      });
      if (data.error) setMatErr(data.error); else setMatRes(data);
    } catch { setMatErr("Erro de conexão"); }
    finally { setMatL(false); }
  }

  return (
    <div>
      {/* Vetores */}
      <div style={card}>
        <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 4 }}>Operações com Vetores</h2>

        <label style={label}>Operação</label>
        <select value={vecOp} onChange={e => setVecOp(e.target.value)}>
          <option value="add">Adição (a + b)</option>
          <option value="sub">Subtração (a - b)</option>
          <option value="dot">Produto Escalar (a · b)</option>
          <option value="cross">Produto Vetorial (a × b) — 3D</option>
          <option value="norm">Norma (||a||)</option>
          <option value="normalize">Normalizar (â)</option>
          <option value="angle">Ângulo entre vetores (graus)</option>
        </select>

        <label style={label}>Vetor A</label>
        <input value={vecA} onChange={e => setVecA(e.target.value)} placeholder="1,2,3" />
        <p style={hint}>Componentes separados por vírgula</p>

        <label style={label}>Vetor B</label>
        <input value={vecB} onChange={e => setVecB(e.target.value)} placeholder="4,5,6" />

        <button style={btn} onClick={calcVector}>Calcular</button>
        <ResultBox result={vecRes} error={vecErr} loading={vecLoad} />
      </div>

      {/* Matrizes */}
      <div style={card}>
        <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 4 }}>Operações com Matrizes</h2>

        <label style={label}>Operação</label>
        <select value={matOp} onChange={e => setMatOp(e.target.value)}>
          <option value="det">Determinante</option>
          <option value="inverse">Inversa</option>
          <option value="transpose">Transposta</option>
          <option value="mul">Multiplicação (A × B)</option>
          <option value="trace">Traço</option>
        </select>

        <label style={label}>Matriz A</label>
        <input value={matA} onChange={e => setMatA(e.target.value)} placeholder="1,2;3,4" />
        <p style={hint}>Linhas separadas por <code>;</code>, colunas por <code>,</code> — ex: <code>1,2;3,4</code></p>

        <label style={label}>Matriz B (para multiplicação)</label>
        <input value={matB} onChange={e => setMatB(e.target.value)} placeholder="5,6;7,8" />

        <button style={btn} onClick={calcMatrix}>Calcular</button>
        <ResultBox result={matRes} error={matErr} loading={matLoad} />
      </div>
    </div>
  );
}