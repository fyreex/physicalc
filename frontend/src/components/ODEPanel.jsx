import { useState } from "react";
import ResultBox from "./ResultBox.jsx";

const card = {
  background: "var(--surface)", border: "1px solid var(--border)",
  borderRadius: "var(--radius)", padding: "1.5rem", marginBottom: "1.5rem",
};
const label = { display: "block", fontSize: 12, color: "var(--muted)", marginBottom: 6, marginTop: 12 };
const row = { display: "flex", gap: "1rem" };
const btn = {
  padding: "0.55rem 1.4rem", border: "none", borderRadius: 8,
  background: "#f59e0b", color: "#0f1117", fontWeight: 600,
  fontSize: 14, marginTop: 16,
};

async function callApi(path, body) {
  const r = await fetch(`/api${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return r.json();
}

function Field({ label: l, value, onChange, placeholder = "0" }) {
  return (
    <div style={{ flex: 1 }}>
      <label style={label}>{l}</label>
      <input type="number" value={value} onChange={e => onChange(e.target.value)} placeholder={placeholder} />
    </div>
  );
}

export default function PhysicsPanel() {
  // Cinemática
  const [x0, setX0] = useState("0");
  const [v0, setV0] = useState("10");
  const [a, setA]   = useState("2");
  const [t, setT]   = useState("5");
  const [kinRes, setKR] = useState(null);
  const [kinErr, setKE] = useState(null);
  const [kinLoad, setKL] = useState(false);

  // Dinâmica
  const [mass, setMass]   = useState("5");
  const [force, setForce] = useState("20");
  const [disp, setDisp]   = useState("10");
  const [vel, setVel]     = useState("4");
  const [height, setH]    = useState("2");
  const [dynRes, setDR]   = useState(null);
  const [dynErr, setDE]   = useState(null);
  const [dynLoad, setDL]  = useState(false);

  async function calcKinematics() {
    setKL(true); setKE(null);
    try {
      const data = await callApi("/physics/kinematics", {
        x0: +x0, v0: +v0, a: +a, t: +t,
      });
      if (data.error) setKE(data.error); else setKR(data);
    } catch { setKE("Erro de conexão"); }
    finally { setKL(false); }
  }

  async function calcDynamics() {
    setDL(true); setDE(null);
    try {
      const data = await callApi("/physics/dynamics", {
        mass: +mass, force: +force, displacement: +disp,
        velocity: +vel, height: +height, gravity: 9.81,
      });
      if (data.error) setDE(data.error); else setDR(data);
    } catch { setDE("Erro de conexão"); }
    finally { setDL(false); }
  }

  return (
    <div>
      {/* Cinemática */}
      <div style={card}>
        <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 4 }}>Cinemática — MRUA</h2>
        <p style={{ fontSize: 13, color: "var(--muted)" }}>
          x = x₀ + v₀t + ½at² &nbsp;|&nbsp; v = v₀ + at
        </p>

        <div style={row}>
          <Field label="Posição inicial x₀ (m)" value={x0} onChange={setX0} />
          <Field label="Velocidade inicial v₀ (m/s)" value={v0} onChange={setV0} />
        </div>
        <div style={row}>
          <Field label="Aceleração a (m/s²)" value={a} onChange={setA} />
          <Field label="Tempo t (s)" value={t} onChange={setT} />
        </div>

        <button style={btn} onClick={calcKinematics}>Calcular</button>
        <ResultBox result={kinRes} error={kinErr} loading={kinLoad} />
      </div>

      {/* Dinâmica */}
      <div style={card}>
        <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 4 }}>Dinâmica — Leis de Newton</h2>
        <p style={{ fontSize: 13, color: "var(--muted)" }}>
          F = ma &nbsp;|&nbsp; W = Fd &nbsp;|&nbsp; Ec = ½mv² &nbsp;|&nbsp; Ep = mgh
        </p>

        <div style={row}>
          <Field label="Massa m (kg)" value={mass} onChange={setMass} placeholder="5" />
          <Field label="Força F (N)" value={force} onChange={setForce} placeholder="20" />
        </div>
        <div style={row}>
          <Field label="Deslocamento d (m)" value={disp} onChange={setDisp} placeholder="10" />
          <Field label="Velocidade v (m/s)" value={vel} onChange={setVel} placeholder="4" />
        </div>
        <div style={row}>
          <Field label="Altura h (m)" value={height} onChange={setH} placeholder="2" />
        </div>

        <button style={btn} onClick={calcDynamics}>Calcular</button>
        <ResultBox result={dynRes} error={dynErr} loading={dynLoad} />
      </div>
    </div>
  );
}