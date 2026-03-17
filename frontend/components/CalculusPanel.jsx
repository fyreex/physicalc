import { useState } from "react";
import ResultBox from "./ResultBox.jsx";

const FUNCTIONS = ["x2","x3","sin","cos","tan","exp","ln","sqrt","1/x","x2+1","sinx2"];

const card = {
  background: "var(--surface)", border: "1px solid var(--border)",
  borderRadius: "var(--radius)", padding: "1.5rem", marginBottom: "1.5rem",
};
const label = { display: "block", fontSize: 12, color: "var(--muted)", marginBottom: 6, marginTop: 12 };
const row = { display: "flex", gap: "1rem", alignItems: "flex-end" };
const btn = {
  padding: "0.55rem 1.4rem", border: "none", borderRadius: 8,
  background: "var(--accent)", color: "#fff", fontWeight: 600,
  fontSize: 14, marginTop: 16, transition: "opacity 0.15s",
};

async function callApi(path, body) {
  const r = await fetch(`/api${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return r.json();
}

export default function CalculusPanel() {
  // Integral state
  const [intFn, setIntFn]     = useState("x2");
  const [intA, setIntA]       = useState("0");
  const [intB, setIntB]       = useState("1");
  const [intN, setIntN]       = useState("1000");
  const [intResult, setIR]    = useState(null);
  const [intErr, setIE]       = useState(null);
  const [intLoad, setIL]      = useState(false);

  // Derivative state
  const [derFn, setDerFn]     = useState("x2");
  const [derX, setDerX]       = useState("3");
  const [derResult, setDR]    = useState(null);
  const [derErr, setDE]       = useState(null);
  const [derLoad, setDL]      = useState(false);

  async function calcIntegral() {
    setIL(true); setIE(null);
    try {
      const data = await callApi("/calculus/integrate", {
        function: intFn, a: +intA, b: +intB, n: +intN,
      });
      if (data.error) setIE(data.error); else setIR(data);
    } catch { setIE("Erro de conexão com a API"); }
    finally { setIL(false); }
  }

  async function calcDerivative() {
    setDL(true); setDE(null);
    try {
      const data = await callApi("/calculus/derivative", {
        function: derFn, x: +derX, h: 1e-5,
      });
      if (data.error) setDE(data.error); else setDR(data);
    } catch { setDE("Erro de conexão com a API"); }
    finally { setDL(false); }
  }

  return (
    <div>
      {/* Integral */}
      <div style={card}>
        <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 4 }}>Integração Numérica</h2>
        <p style={{ fontSize: 13, color: "var(--muted)" }}>
          Métodos de Simpson 1/3 e Trapézio
        </p>

        <label style={label}>Função</label>
        <select value={intFn} onChange={e => setIntFn(e.target.value)}>
          {FUNCTIONS.map(f => <option key={f} value={f}>{f}</option>)}
        </select>

        <div style={row}>
          <div style={{ flex: 1 }}>
            <label style={label}>Limite inferior (a)</label>
            <input type="number" value={intA} onChange={e => setIntA(e.target.value)} />
          </div>
          <div style={{ flex: 1 }}>
            <label style={label}>Limite superior (b)</label>
            <input type="number" value={intB} onChange={e => setIntB(e.target.value)} />
          </div>
          <div style={{ flex: 1 }}>
            <label style={label}>Subintervalos (n)</label>
            <input type="number" value={intN} onChange={e => setIntN(e.target.value)} />
          </div>
        </div>

        <button style={btn} onClick={calcIntegral}>Calcular Integral</button>
        <ResultBox result={intResult} error={intErr} loading={intLoad} />
      </div>

      {/* Derivada */}
      <div style={card}>
        <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 4 }}>Derivada Numérica</h2>
        <p style={{ fontSize: 13, color: "var(--muted)" }}>
          Diferenças centrais — f'(x) ≈ [f(x+h) - f(x-h)] / 2h
        </p>

        <label style={label}>Função</label>
        <select value={derFn} onChange={e => setDerFn(e.target.value)}>
          {FUNCTIONS.map(f => <option key={f} value={f}>{f}</option>)}
        </select>

        <label style={label}>Ponto x</label>
        <input type="number" value={derX} onChange={e => setDerX(e.target.value)} />

        <button style={btn} onClick={calcDerivative}>Calcular Derivada</button>
        <ResultBox result={derResult} error={derErr} loading={derLoad} />
      </div>
    </div>
  );
}