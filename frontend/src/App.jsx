import { useState } from "react";
import CalculusPanel from "./components/CalculusPanel.jsx";
import AlgebraPanel from "./components/AlgebraPanel.jsx";
import PhysicsPanel from "./components/PhysicsPanel.jsx";
import ODEPanel from "./components/ODEPanel.jsx";

const TABS = [
  { id: "calculus", label: "∫ Cálculo" },
  { id: "algebra",  label: "⊕ Álgebra Linear" },
  { id: "physics",  label: "⚡ Física" },
  { id: "ode",      label: "∂ EDO" },
];

const styles = {
  app: { maxWidth: 900, margin: "0 auto", padding: "2rem 1.5rem" },
  header: { marginBottom: "2rem", textAlign: "center" },
  title: { fontSize: "2rem", fontWeight: 700, letterSpacing: "-0.5px",
    background: "linear-gradient(135deg, #6c63ff, #00d4aa)",
    WebkitBackgroundClip: "text", WebkitTextFillColor: "transparent" },
  subtitle: { color: "var(--muted)", marginTop: "0.4rem", fontSize: "14px" },
  tabs: { display: "flex", gap: "0.5rem", marginBottom: "2rem",
    borderBottom: "1px solid var(--border)", paddingBottom: "0" },
  tab: (active) => ({
    padding: "0.6rem 1.2rem", border: "none", borderRadius: "8px 8px 0 0",
    background: active ? "var(--accent)" : "transparent",
    color: active ? "#fff" : "var(--muted)",
    fontWeight: active ? 600 : 400, fontSize: "14px",
    transition: "all 0.15s", borderBottom: active ? "2px solid var(--accent)" : "2px solid transparent",
  }),
};

export default function App() {
  const [tab, setTab] = useState("calculus");

  return (
    <div style={styles.app}>
      <header style={styles.header}>
        <h1 style={styles.title}>Physicalc</h1>
        <p style={styles.subtitle}>
          Calculadora de Cálculo, Álgebra Linear, Física e EDOs — feita em Go
        </p>
      </header>

      <nav style={styles.tabs}>
        {TABS.map(t => (
          <button key={t.id} style={styles.tab(tab === t.id)} onClick={() => setTab(t.id)}>
            {t.label}
          </button>
        ))}
      </nav>

      {tab === "calculus" && <CalculusPanel />}
      {tab === "algebra"  && <AlgebraPanel />}
      {tab === "physics"  && <PhysicsPanel />}
      {tab === "ode"      && <ODEPanel />}
    </div>
  );
}