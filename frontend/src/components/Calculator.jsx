import { useState, useEffect, useRef } from "react";

// ─── Funções matemáticas predefinidas ────────────────────────────────────────
const FUNCS = {
  x2:  x => x * x,
  x3:  x => x * x * x,
  sin: x => Math.sin(x),
  cos: x => Math.cos(x),
  exp: x => Math.exp(x),
  ln:  x => Math.log(x),
  sqrt:x => Math.sqrt(x),
  "1/x": x => 1 / x,
};

// ─── Utilitários ─────────────────────────────────────────────────────────────
function toRad(deg) { return deg * Math.PI / 180; }

function fmt(n) {
  if (isNaN(n))       return "NaN";
  if (!isFinite(n))   return n > 0 ? "+∞" : "-∞";
  if (Math.abs(n) < 1e-12 && n !== 0) return n.toExponential(4);
  if (Math.abs(n) > 1e10)              return n.toExponential(4);
  return parseFloat(n.toPrecision(10)).toString();
}

// Verifica divisão por zero na expressão antes de avaliar
function checkDivisionByZero(s) {
  // Detecta padrões como /0 ou / 0
  return /\/\s*0(?![.\d])/.test(s);
}

function parseExpr(s, mode, xval = 0) {
  // Verifica divisão por zero explícita
  if (checkDivisionByZero(s)) {
    throw new Error("DIV_ZERO");
  }

  s = s.replace(/π/g, "(Math.PI)").replace(/ℯ/g, "(Math.E)");
  s = s.replace(/\^/g, "**");
  s = s.replace(/x\*\*2/g, "(xv*xv)");
  s = s.replace(/\bx\b/g, "xv");

  if (mode === "DEG") {
    s = s.replace(/Math\.sin\(/g, "Math.sin(Math.PI/180*");
    s = s.replace(/Math\.cos\(/g, "Math.cos(Math.PI/180*");
    s = s.replace(/Math\.tan\(/g, "Math.tan(Math.PI/180*");
    s = s.replace(/Math\.asin\(/g, "(180/Math.PI)*Math.asin(");
    s = s.replace(/Math\.acos\(/g, "(180/Math.PI)*Math.acos(");
    s = s.replace(/Math\.atan\(/g, "(180/Math.PI)*Math.atan(");
  }

  s = s
    .replace(/sin\(/g,  "Math.sin(")  .replace(/cos\(/g,  "Math.cos(")
    .replace(/tan\(/g,  "Math.tan(")  .replace(/log\(/g,  "Math.log10(")
    .replace(/ln\(/g,   "Math.log(")  .replace(/sqrt\(/g, "Math.sqrt(")
    .replace(/cbrt\(/g, "Math.cbrt(") .replace(/abs\(/g,  "Math.abs(")
    .replace(/asin\(/g, "Math.asin(") .replace(/acos\(/g, "Math.acos(")
    .replace(/atan\(/g, "Math.atan(") .replace(/%/g, "/100");

  // eslint-disable-next-line no-new-func
  const result = Function('"use strict"; var xv=' + xval + '; return (' + s + ')')();

  // Verifica resultado infinito causado por divisão por zero implícita (ex: 1/x com x=0)
  if (!isFinite(result) && result !== Infinity && result !== -Infinity) {
    throw new Error("DIV_ZERO");
  }
  // Se o resultado for Infinity, também é divisão por zero
  if (result === Infinity || result === -Infinity) {
    throw new Error("DIV_ZERO");
  }

  return result;
}

// ─── Algoritmos numéricos ────────────────────────────────────────────────────
function simpson(f, a, b, n) {
  if (n % 2 !== 0) n++;
  const h = (b - a) / n;
  let s = f(a) + f(b);
  for (let i = 1; i < n; i++) {
    const x = a + i * h;
    s += i % 2 === 0 ? 2 * f(x) : 4 * f(x);
  }
  return (h / 3) * s;
}

function trapezoid(f, a, b, n) {
  const h = (b - a) / n;
  let s = (f(a) + f(b)) / 2;
  for (let i = 1; i < n; i++) s += f(a + i * h);
  return h * s;
}

function centralDiff(f, x, h = 1e-5) {
  return (f(x + h) - f(x - h)) / (2 * h);
}

function rk4(f, y0, t0, te, h = 0.05) {
  let t = t0, y = y0;
  const pts = [[t, y]];
  while (t < te - 1e-10) {
    const k1 = f(t, y);
    const k2 = f(t + h / 2, y + h * k1 / 2);
    const k3 = f(t + h / 2, y + h * k2 / 2);
    const k4 = f(t + h, y + h * k3);
    y += (h / 6) * (k1 + 2 * k2 + 2 * k3 + k4);
    t += h;
    pts.push([t, y]);
  }
  return pts;
}

// ─── Gerador de passos para cada operação ───────────────────────────────────
function buildSteps(raw, mode) {
  const s = raw.trim();
  const steps = [];

  const add = (label, val) => steps.push({ label, val });

  if (/^[\d.]+\s*[\+\-\*\/]\s*[\d.]+$/.test(s)) {
    const m = s.match(/^([\d.]+)\s*([\+\-\*\/])\s*([\d.]+)$/);
    if (m) {
      const a = +m[1], op = m[2], b = +m[3];
      if (op === "/" && b === 0) {
        add("Erro:", "Impossível — divisão por zero");
        return steps;
      }
      const names = { "+": "Adição", "-": "Subtração", "*": "Multiplicação", "/": "Divisão" };
      add("Operação:", names[op] || op);
      add("Valor A:", String(a));
      add("Valor B:", String(b));
      const r = op === "+" ? a + b : op === "-" ? a - b : op === "*" ? a * b : a / b;
      add("Resultado:", `${a} ${op} ${b} = ${fmt(r)}`);
    }
  } else if (/^(sin|cos|tan)\(/.test(s)) {
    const fn = s.match(/^(\w+)/)[1];
    const inner = s.slice(fn.length + 1, -1);
    const v = parseFloat(inner);
    add("Função:", fn + "(x)");
    add("Entrada:", `x = ${v}° ${mode === "DEG" ? "→ " + fmt(v * Math.PI / 180) + " rad" : ""}`);
    add("Série:", "Calculado via série de Taylor");
    const r = fn === "sin" ? Math.sin(toRad(v)) : fn === "cos" ? Math.cos(toRad(v)) : Math.tan(toRad(v));
    add("Resultado:", `${fn}(${v}°) = ${fmt(r)}`);
  } else if (/^sqrt\(/.test(s)) {
    const v = parseFloat(s.slice(5, -1));
    add("Operação:", "Raiz quadrada √x");
    add("Entrada:", `x = ${v}`);
    add("Identidade:", "√x = x^(1/2)");
    add("Cálculo:", `${v}^0.5 = ${fmt(Math.sqrt(v))}`);
  } else if (/^cbrt\(/.test(s)) {
    const v = parseFloat(s.slice(5, -1));
    add("Operação:", "Raiz cúbica ∛x");
    add("Entrada:", `x = ${v}`);
    add("Identidade:", "∛x = x^(1/3)");
    add("Cálculo:", `${v}^(1/3) = ${fmt(Math.cbrt(v))}`);
  } else if (/^log\(/.test(s)) {
    const v = parseFloat(s.slice(4, -1));
    add("Função:", "Logaritmo base 10");
    add("Entrada:", `x = ${v}`);
    add("Definição:", "log(x) = y onde 10ʸ = x");
    add("Resultado:", `log(${v}) = ${fmt(Math.log10(v))}`);
  } else if (/^ln\(/.test(s)) {
    const v = parseFloat(s.slice(3, -1));
    add("Função:", "Logaritmo natural");
    add("Entrada:", `x = ${v}`);
    add("Definição:", "ln(x) = y onde eʸ = x");
    add("Resultado:", `ln(${v}) = ${fmt(Math.log(v))}`);
  } else if (s.includes("^")) {
    const pts = s.split("^");
    if (pts.length === 2) {
      const base = parseFloat(pts[0]), exp = parseFloat(pts[1]);
      add("Operação:", "Potenciação bⁿ");
      add("Base:", `b = ${base}`);
      add("Expoente:", `n = ${exp}`);
      add("Definição:", "b × b × ... × b (n vezes)");
      add("Resultado:", `${base}^${exp} = ${fmt(Math.pow(base, exp))}`);
    }
  } else {
    add("Expressão:", s);
  }

  return steps;
}

// ─── Estilos inline (sem CSS externo) ───────────────────────────────────────
const S = {
  wrap:      { display: "flex", gap: 12, maxWidth: 780, margin: "0 auto", fontFamily: "'Courier New', monospace" },
  calc:      { background: "#111520", borderRadius: 14, padding: 14, width: 400, flexShrink: 0, border: "1px solid #1e2535" },
  screen:    { background: "#080c14", borderRadius: 10, padding: "10px 14px", marginBottom: 10, border: "1px solid #1a2030", minHeight: 110 },
  scTop:     { display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 4 },
  scMode:    { fontSize: 9, color: "#2a4560", letterSpacing: 2 },
  scMem:     { fontSize: 9, color: "#3a6a4a", letterSpacing: 1, minWidth: 30, textAlign: "right" },
  scExpr:    { fontSize: 12, color: "#3a6888", textAlign: "right", minHeight: 18, wordBreak: "break-all" },
  scResult:  (err) => ({ fontSize: 26, color: err ? "#e07070" : "#6dcfa0", textAlign: "right", fontWeight: 700, minHeight: 36, transition: "color .2s" }),
  scErr:     { fontSize: 12, color: "#e07070", textAlign: "right", marginTop: 4, minHeight: 18 },
  stepsWrap: { marginTop: 8, borderTop: "1px solid #131a26", paddingTop: 6 },
  stepRow:   { fontSize: 11, lineHeight: 1.9, display: "flex", gap: 6 },
  stepLabel: { color: "#2d6b4a", minWidth: 80, flexShrink: 0 },
  stepVal:   { color: "#6dcfa0" },
  modeRow:   { display: "flex", gap: 5, marginBottom: 10 },
  modeBtn:   (active) => ({ flex: 1, padding: "5px 2px", border: `1px solid ${active ? "#2a4a60" : "#1e2535"}`, borderRadius: 6, background: active ? "#162030" : "#0d1120", color: active ? "#5aaad0" : "#3a5a6a", fontSize: 9, fontFamily: "'Courier New', monospace", cursor: "pointer", letterSpacing: 1 }),
  keys:      { display: "grid", gridTemplateColumns: "repeat(5, 1fr)", gap: 5 },
  key:       (type) => {
    const base = { border: "none", borderRadius: 7, padding: "9px 3px", fontFamily: "'Courier New', monospace", fontSize: 11, fontWeight: 700, cursor: "pointer", letterSpacing: 0.3, lineHeight: 1.3, textAlign: "center", transition: "transform .1s" };
    const types = {
      fn:  { background: "#162030", color: "#5aaad0", border: "1px solid #243a50" },
      op:  { background: "#122820", color: "#6dcfa0", border: "1px solid #1e4030" },
      num: { background: "#0e1420", color: "#b0c8dc", border: "1px solid #1a2030" },
      eq:  { background: "#163528", color: "#50c080", border: "1px solid #2a5040", fontSize: 16 },
      cl:  { background: "#281414", color: "#e07070", border: "1px solid #402020" },
      mem: { background: "#1a1830", color: "#9880d0", border: "1px solid #2a2848", fontSize: 10 },
      sp:  { background: "#181828", color: "#8870c0", border: "1px solid #282840", fontSize: 10 },
    };
    return { ...base, ...types[type] };
  },
  sidebar:   { flex: 1, display: "flex", flexDirection: "column", gap: 10 },
  tabRow:    { display: "flex", gap: 4 },
  tab:       (active) => ({ flex: 1, padding: 5, border: `1px solid ${active ? "#2a4a60" : "#1e2535"}`, borderRadius: 5, background: active ? "#162030" : "#0d1120", color: active ? "#5aaad0" : "#3a5a6a", fontSize: 9, fontFamily: "'Courier New', monospace", cursor: "pointer", letterSpacing: 1 }),
  panel:     { background: "#080c14", borderRadius: 10, padding: 10, border: "1px solid #1a2030", flex: 1, overflowY: "auto", maxHeight: 260 },
  panelH3:   { fontSize: 9, color: "#2a4560", letterSpacing: 2, marginBottom: 8 },
  histItem:  { fontSize: 11, color: "#4a7a90", padding: "4px 6px", borderRadius: 5, cursor: "pointer", border: "1px solid transparent", marginBottom: 3 },
  histExpr:  { color: "#2a4a60", fontSize: 10 },
  calcPanel: { background: "#080c14", borderRadius: 10, padding: 12, border: "1px solid #1a2030" },
  ciLabel:   { fontSize: 10, color: "#3a5a6a", display: "block", marginBottom: 3 },
  ciSel:     { width: "100%", background: "#0d1220", color: "#6dcfa0", border: "1px solid #1e2e40", borderRadius: 6, padding: "5px 8px", fontSize: 11, fontFamily: "'Courier New', monospace", marginBottom: 8 },
  ciInput:   { width: "100%", background: "#0d1220", color: "#b0c8dc", border: "1px solid #1e2e40", borderRadius: 6, padding: "5px 8px", fontSize: 11, fontFamily: "'Courier New', monospace", marginBottom: 8 },
  runBtn:    { width: "100%", padding: 7, background: "#163528", color: "#50c080", border: "1px solid #2a5040", borderRadius: 7, fontFamily: "'Courier New', monospace", fontSize: 11, fontWeight: 700, cursor: "pointer", marginTop: 4 },
};

// ─── Componente principal ────────────────────────────────────────────────────
export default function Calculator() {
  const [expr, setExpr]       = useState("");
  const [result, setResult]   = useState("0");
  const [errMsg, setErrMsg]   = useState("");
  const [isErr, setIsErr]     = useState(false);
  const [mode, setMode]       = useState("DEG");
  const [mem, setMem]         = useState(null);
  const [history, setHistory] = useState([]);
  const [steps, setSteps]     = useState([]);
  const [visSteps, setVis]    = useState([]);
  const [tab, setTab]         = useState("hist");

  // Cálculo avançado
  const [calcOp, setCalcOp]   = useState("integral");
  const [calcFn, setCalcFn]   = useState("x2");
  const [iA, setIA]           = useState("0");
  const [iB, setIB]           = useState("1");
  const [iN, setIN]           = useState("1000");
  const [dX, setDX]           = useState("2");
  const [lX, setLX]           = useState("0");
  const [oY0, setOY0]         = useState("1");
  const [oT0, setOT0]         = useState("0");
  const [oTe, setOTe]         = useState("3");

  const timerRef = useRef([]);

  function p(v) { setExpr(e => e + v); }

  function clearAll() {
    setExpr(""); setResult("0"); setErrMsg(""); setIsErr(false);
    setSteps([]); setVis([]);
  }

  function delLast() { setExpr(e => e.slice(0, -1)); }

  function memStore() {
    const v = parseFloat(result);
    if (!isNaN(v)) setMem(v);
  }

  function memRecall() {
    if (mem !== null) setExpr(e => e + String(mem));
  }

  function animateSteps(newSteps) {
    timerRef.current.forEach(clearTimeout);
    timerRef.current = [];
    setSteps(newSteps);
    setVis([]);
    newSteps.forEach((_, i) => {
      const tid = setTimeout(() => setVis(v => [...v, i]), i * 120 + 50);
      timerRef.current.push(tid);
    });
  }

  function addHistory(e, r) {
    setHistory(h => [{ expr: e, result: r }, ...h].slice(0, 20));
  }

  function calculate() {
    if (!expr) return;
    setErrMsg("");
    try {
      const r = parseExpr(expr, mode);
      const formatted = fmt(r);
      setResult(formatted);
      setIsErr(false);
      setErrMsg("");
      if (mode === "STEP") animateSteps(buildSteps(expr, mode));
      else { setSteps([]); setVis([]); }
      addHistory(expr, formatted);
      setExpr(formatted);
    } catch (e) {
      if (e.message === "DIV_ZERO") {
        setResult("ERRO");
        setErrMsg("Impossível — divisão por zero");
        setIsErr(true);
        if (mode === "STEP") {
          animateSteps([{ label: "Erro:", val: "Impossível — divisão por zero" }]);
        }
      } else {
        setResult("ERRO");
        setErrMsg("Expressão inválida");
        setIsErr(true);
      }
    }
  }

  function runCalc() {
    const f = FUNCS[calcFn];
    let newResult = "", newSteps = [];
    const add = (l, v) => newSteps.push({ label: l, val: v });
    setErrMsg(""); setIsErr(false);

    try {
      if (calcOp === "integral") {
        const a = +iA, b = +iB, n = +iN;
        const s = simpson(f, a, b, n), tr = trapezoid(f, a, b, n);
        newResult = `Simpson: ${fmt(s)}`;
        add("Operação:", `∫ ${calcFn} dx de ${a} a ${b}`);
        add("Método 1:", "Simpson 1/3 (erro O(h⁴))");
        add("h =", `(${b}-${a})/${n} = ${fmt((b - a) / n)}`);
        add("Fórmula:", "h/3 × [f(a) + 4f(x₁) + 2f(x₂) + ... + f(b)]");
        add("Simpson:", fmt(s));
        add("Método 2:", "Trapézio composto (erro O(h²))");
        add("Trapézio:", fmt(tr));
        add("Diferença:", fmt(Math.abs(s - tr)) + " (erro numérico)");
        addHistory(`∫${calcFn}[${a},${b}]`, fmt(s));

      } else if (calcOp === "derivative") {
        const x = +dX, h = 1e-5;
        if (calcFn === "1/x" && x === 0) throw new Error("DIV_ZERO");
        const d = centralDiff(f, x, h);
        newResult = `f'(${x}) = ${fmt(d)}`;
        add("Operação:", `d/dx [${calcFn}] em x=${x}`);
        add("Método:", "Diferenças centrais");
        add("Fórmula:", "f'(x) ≈ [f(x+h) - f(x-h)] / 2h");
        add("h =", String(h));
        add("f(x+h) =", fmt(f(x + h)));
        add("f(x-h) =", fmt(f(x - h)));
        add("Numerador:", fmt(f(x + h) - f(x - h)));
        add(`f'(${x}) =`, fmt(d));
        addHistory(`d/dx ${calcFn} @ x=${x}`, fmt(d));

      } else if (calcOp === "limit") {
        const x = +lX, h = 1e-7;
        if (calcFn === "1/x" && x === 0) {
          newResult = "Impossível — divisão por zero";
          add("Operação:", `lim x→${x} [${calcFn}]`);
          add("Erro:", "Impossível — divisão por zero");
          add("Conclusão:", "Limite não existe em x=0 para 1/x");
          addHistory(`lim→${x} ${calcFn}`, "∄");
        } else {
          const left = f(x - h), right = f(x + h);
          const conv = Math.abs(left - right) < 1e-6;
          newResult = conv ? `lim = ${fmt((left + right) / 2)}` : "Não converge";
          add("Operação:", `lim x→${x} [${calcFn}]`);
          add("Método:", "Aproximação bilateral");
          add("h =", String(h));
          add("Limite esq.:", `f(${x}-h) = ${fmt(left)}`);
          add("Limite dir.:", `f(${x}+h) = ${fmt(right)}`);
          add("Diferença:", fmt(Math.abs(left - right)));
          add("Converge?", conv ? "Sim ✓" : "Não ✗");
          add("Limite:", conv ? fmt((left + right) / 2) : "indefinido");
          addHistory(`lim→${x} ${calcFn}`, conv ? fmt((left + right) / 2) : "∄");
        }

      } else if (calcOp === "ode") {
        const y0 = +oY0, t0 = +oT0, te = +oTe;
        const odeF = (t, y) => f(y);
        const pts = rk4(odeF, y0, t0, te);
        const last = pts[pts.length - 1];
        newResult = `y(${fmt(last[0])}) = ${fmt(last[1])}`;
        add("Operação:", `dy/dt = ${calcFn}(y)`);
        add("Método:", "Runge-Kutta 4 (erro O(h⁵))");
        add("Condição:", `y(${t0}) = ${y0}`);
        add("Intervalo:", `t ∈ [${t0}, ${te}]`);
        add("Passo h:", "0.05");
        add("Passos:", String(pts.length - 1));
        add("Pesos RK4:", "k1, 2k2, 2k3, k4 → média ponderada");
        add(`y(${fmt(last[0])}) =`, fmt(last[1]));
        addHistory(`RK4 ${calcFn} y₀=${y0}`, fmt(last[1]));
      }

      setResult(newResult);
      setIsErr(false);
      animateSteps(newSteps);
      setTab("hist");

    } catch (e) {
      if (e.message === "DIV_ZERO") {
        setResult("ERRO");
        setErrMsg("Impossível — divisão por zero");
        setIsErr(true);
        animateSteps([{ label: "Erro:", val: "Impossível — divisão por zero" }]);
      } else {
        setResult("ERRO");
        setErrMsg("Erro no cálculo");
        setIsErr(true);
      }
    }
  }

  // Teclado físico
  useEffect(() => {
    const handler = (e) => {
      if (e.key === "Enter")     calculate();
      else if (e.key === "Backspace") delLast();
      else if (e.key === "Escape")    clearAll();
      else if ("0123456789.+-*/()".includes(e.key)) p(e.key);
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [expr, mode]);

  // ─── Render ───────────────────────────────────────────────────────────────
  const Btn = ({ t, children, onClick, span }) => (
    <button style={{ ...S.key(t), ...(span ? { gridColumn: `span ${span}` } : {}) }} onClick={onClick}>
      {children}
    </button>
  );

  return (
    <div style={S.wrap}>
      {/* ── Calculadora ── */}
      <div style={S.calc}>
        {/* Modos */}
        <div style={S.modeRow}>
          {["DEG", "RAD", "STEP"].map(m => (
            <button key={m} style={S.modeBtn(mode === m)} onClick={() => setMode(m)}>{m}</button>
          ))}
        </div>

        {/* Tela */}
        <div style={S.screen}>
          <div style={S.scTop}>
            <span style={S.scMode}>{mode}</span>
            <span style={S.scMem}>{mem !== null ? `M=${fmt(mem)}` : ""}</span>
          </div>
          <div style={S.scExpr}>{expr || " "}</div>
          <div style={S.scResult(isErr)}>{result}</div>
          {errMsg && <div style={S.scErr}>{errMsg}</div>}

          {/* Passo a passo */}
          {steps.length > 0 && (
            <div style={S.stepsWrap}>
              {steps.map((s, i) => (
                <div key={i} style={{ ...S.stepRow, opacity: visSteps.includes(i) ? 1 : 0, transform: visSteps.includes(i) ? "none" : "translateY(5px)", transition: "opacity .28s, transform .28s" }}>
                  <span style={S.stepLabel}>{s.label}</span>
                  <span style={S.stepVal}>{s.val}</span>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Teclado */}
        <div style={S.keys}>
          <Btn t="fn" onClick={() => p("sin(")}>sin</Btn>
          <Btn t="fn" onClick={() => p("cos(")}>cos</Btn>
          <Btn t="fn" onClick={() => p("tan(")}>tan</Btn>
          <Btn t="fn" onClick={() => p("log(")}>log</Btn>
          <Btn t="fn" onClick={() => p("ln(")}>ln</Btn>

          <Btn t="sp" onClick={() => p("sqrt(")}>√x</Btn>
          <Btn t="sp" onClick={() => p("cbrt(")}>∛x</Btn>
          <Btn t="sp" onClick={() => p("^")}>xʸ</Btn>
          <Btn t="sp" onClick={() => p("π")}>π</Btn>
          <Btn t="sp" onClick={() => p("ℯ")}>e</Btn>

          <Btn t="mem" onClick={memStore}>M+</Btn>
          <Btn t="mem" onClick={memRecall}>MR</Btn>
          <Btn t="mem" onClick={() => setMem(null)}>MC</Btn>
          <Btn t="fn" onClick={() => p("asin(")}>sin⁻¹</Btn>
          <Btn t="fn" onClick={() => p("acos(")}>cos⁻¹</Btn>

          <Btn t="cl" onClick={clearAll}>AC</Btn>
          <Btn t="cl" onClick={delLast}>DEL</Btn>
          <Btn t="op" onClick={() => p("(")}>(</Btn>
          <Btn t="op" onClick={() => p(")")}>)</Btn>
          <Btn t="fn" onClick={() => p("atan(")}>tan⁻¹</Btn>

          <Btn t="num" onClick={() => p("7")}>7</Btn>
          <Btn t="num" onClick={() => p("8")}>8</Btn>
          <Btn t="num" onClick={() => p("9")}>9</Btn>
          <Btn t="op" onClick={() => p("/")}>÷</Btn>
          <Btn t="sp" onClick={() => p("x")}>x</Btn>

          <Btn t="num" onClick={() => p("4")}>4</Btn>
          <Btn t="num" onClick={() => p("5")}>5</Btn>
          <Btn t="num" onClick={() => p("6")}>6</Btn>
          <Btn t="op" onClick={() => p("*")}>×</Btn>
          <Btn t="sp" onClick={() => p("x^2")}>x²</Btn>

          <Btn t="num" onClick={() => p("1")}>1</Btn>
          <Btn t="num" onClick={() => p("2")}>2</Btn>
          <Btn t="num" onClick={() => p("3")}>3</Btn>
          <Btn t="op" onClick={() => p("-")}>−</Btn>
          <Btn t="sp" onClick={() => p("abs(")}>|x|</Btn>

          <Btn t="num" onClick={() => p("0")} span={2}>0</Btn>
          <Btn t="num" onClick={() => p(".")}>.</Btn>
          <Btn t="op" onClick={() => p("+")}>+</Btn>
          <Btn t="eq" onClick={calculate}>=</Btn>
        </div>
      </div>

      {/* ── Sidebar ── */}
      <div style={S.sidebar}>
        <div style={S.tabRow}>
          <button style={S.tab(tab === "hist")} onClick={() => setTab("hist")}>HISTÓRICO</button>
          <button style={S.tab(tab === "calc")} onClick={() => setTab("calc")}>CÁLCULO</button>
        </div>

        {/* Histórico */}
        {tab === "hist" && (
          <div style={S.panel}>
            <h3 style={S.panelH3}>ÚLTIMAS OPERAÇÕES</h3>
            {history.length === 0
              ? <span style={{ fontSize: 10, color: "#1e3040" }}>Nenhum cálculo ainda</span>
              : history.map((h, i) => (
                <div key={i} style={S.histItem} onClick={() => setExpr(String(h.result))}>
                  <div style={S.histExpr}>{h.expr}</div>
                  <div>{h.result}</div>
                </div>
              ))
            }
          </div>
        )}

        {/* Painel de Cálculo Avançado */}
        {tab === "calc" && (
          <div style={S.calcPanel}>
            <h3 style={S.panelH3}>CÁLCULO AVANÇADO</h3>

            <label style={S.ciLabel}>Operação</label>
            <select style={S.ciSel} value={calcOp} onChange={e => setCalcOp(e.target.value)}>
              <option value="integral">∫ Integral definida</option>
              <option value="derivative">d/dx Derivada num ponto</option>
              <option value="limit">lim Limite num ponto</option>
              <option value="ode">dy/dt EDO (Runge-Kutta 4)</option>
            </select>

            <label style={S.ciLabel}>Função f(x)</label>
            <select style={S.ciSel} value={calcFn} onChange={e => setCalcFn(e.target.value)}>
              {Object.keys(FUNCS).map(k => <option key={k} value={k}>{k}</option>)}
            </select>

            {calcOp === "integral" && <>
              <label style={S.ciLabel}>Limite inferior (a)</label>
              <input style={S.ciInput} type="number" value={iA} onChange={e => setIA(e.target.value)} />
              <label style={S.ciLabel}>Limite superior (b)</label>
              <input style={S.ciInput} type="number" value={iB} onChange={e => setIB(e.target.value)} />
              <label style={S.ciLabel}>Subintervalos (n)</label>
              <input style={S.ciInput} type="number" value={iN} onChange={e => setIN(e.target.value)} />
            </>}

            {calcOp === "derivative" && <>
              <label style={S.ciLabel}>Ponto x</label>
              <input style={S.ciInput} type="number" value={dX} onChange={e => setDX(e.target.value)} />
            </>}

            {calcOp === "limit" && <>
              <label style={S.ciLabel}>Ponto (x →)</label>
              <input style={S.ciInput} type="number" value={lX} onChange={e => setLX(e.target.value)} />
            </>}

            {calcOp === "ode" && <>
              <label style={S.ciLabel}>y₀ (condição inicial)</label>
              <input style={S.ciInput} type="number" value={oY0} onChange={e => setOY0(e.target.value)} />
              <label style={S.ciLabel}>t₀</label>
              <input style={S.ciInput} type="number" value={oT0} onChange={e => setOT0(e.target.value)} />
              <label style={S.ciLabel}>t final</label>
              <input style={S.ciInput} type="number" value={oTe} onChange={e => setOTe(e.target.value)} />
            </>}

            <button style={S.runBtn} onClick={runCalc}>CALCULAR</button>
          </div>
        )}
      </div>
    </div>
  );
}
