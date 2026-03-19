// ResultBox: exibe resultados da API de forma legível e organizada

const S = {
  box: (error) => ({
    marginTop: "1.5rem",
    background: "var(--surface)",
    border: `1px solid ${error ? "var(--error)" : "var(--border)"}`,
    borderRadius: "var(--radius)",
    padding: "1.2rem",
    minHeight: 60,
  }),
  label: {
    fontSize: 11,
    color: "var(--muted)",
    textTransform: "uppercase",
    letterSpacing: 1,
    marginBottom: 12,
  },
  mainValue: {
    fontSize: 28,
    fontWeight: 700,
    color: "#00d4aa",
    marginBottom: 16,
    wordBreak: "break-all",
  },
  grid: {
    display: "grid",
    gridTemplateColumns: "repeat(auto-fill, minmax(200px, 1fr))",
    gap: 10,
  },
  card: {
    background: "#0d1117",
    borderRadius: 8,
    padding: "10px 14px",
    border: "1px solid var(--border)",
  },
  cardLabel: {
    fontSize: 10,
    color: "var(--muted)",
    marginBottom: 4,
    textTransform: "uppercase",
    letterSpacing: 0.8,
  },
  cardValue: {
    fontSize: 15,
    fontWeight: 600,
    color: "#e8eaf0",
    wordBreak: "break-all",
  },
  arrayWrap: {
    display: "flex",
    gap: 6,
    flexWrap: "wrap",
    marginTop: 4,
  },
  arrayItem: {
    background: "#1a1d27",
    border: "1px solid var(--border)",
    borderRadius: 6,
    padding: "3px 10px",
    fontSize: 13,
    color: "#00d4aa",
    fontFamily: "monospace",
  },
  matrixWrap: {
    fontFamily: "monospace",
    fontSize: 13,
    color: "#e8eaf0",
    lineHeight: 1.8,
  },
};

function formatValue(val) {
  if (val === null || val === undefined) return "—";
  if (typeof val === "number") {
    if (!isFinite(val)) return val > 0 ? "+∞" : "-∞";
    return parseFloat(val.toPrecision(8)).toString();
  }
  if (typeof val === "boolean") return val ? "Sim" : "Não";
  return String(val);
}

const KEY_LABELS = {
  simpson:            "Simpson 1/3",
  trapezoid:          "Trapézio",
  derivative:         "Derivada f'(x)",
  limit_left:         "Limite pela esquerda",
  limit_right:        "Limite pela direita",
  converges:          "Converge?",
  position_m:         "Posição (m)",
  velocity_ms:        "Velocidade (m/s)",
  displacement_m:     "Deslocamento (m)",
  distance_m:         "Distância percorrida (m)",
  velocity_sq_ms2:    "v² (m²/s²)",
  kinetic_energy_j:   "Energia cinética (J)",
  momentum_kgms:      "Momento (kg·m/s)",
  movement_type:      "Tipo de movimento",
  force_n:            "Força (N)",
  acceleration_ms2:   "Aceleração (m/s²)",
  weight_n:           "Peso (N)",
  work_j:             "Trabalho (J)",
  potential_energy_j: "Energia potencial (J)",
  total_energy_j:     "Energia total (J)",
  power_w:            "Potência (W)",
  result:             "Resultado",
};

const SKIP_KEYS = new Set(["a", "b", "n", "x", "h", "func", "y0", "t0",
  "t_end", "x0_m", "v0_ms", "a_ms2", "t_s", "steps", "method", "function", "operation"]);

function VectorDisplay({ data }) {
  const items = data?.Components || data;
  if (!Array.isArray(items)) return <span style={S.cardValue}>{formatValue(data)}</span>;
  return (
    <div style={S.arrayWrap}>
      {items.map((v, i) => <span key={i} style={S.arrayItem}>{formatValue(v)}</span>)}
    </div>
  );
}

function MatrixDisplay({ data }) {
  const rows = data?.Data || data;
  if (!Array.isArray(rows)) return <span style={S.cardValue}>{formatValue(data)}</span>;
  return (
    <div style={S.matrixWrap}>
      {rows.map((row, i) => (
        <div key={i} style={{ display: "flex", gap: 12 }}>
          <span style={{ color: "var(--muted)" }}>│</span>
          {(Array.isArray(row) ? row : [row]).map((v, j) => (
            <span key={j} style={{ minWidth: 70, color: "#00d4aa" }}>{formatValue(v)}</span>
          ))}
          <span style={{ color: "var(--muted)" }}>│</span>
        </div>
      ))}
    </div>
  );
}

function isVector(val) {
  return val && typeof val === "object" && (Array.isArray(val?.Components) || Array.isArray(val));
}
function isMatrix(val) {
  return val && typeof val === "object" && Array.isArray(val?.Data);
}

function FieldCard({ label, value }) {
  return (
    <div style={S.card}>
      <div style={S.cardLabel}>{label}</div>
      {isMatrix(value)
        ? <MatrixDisplay data={value} />
        : isVector(value)
        ? <VectorDisplay data={value} />
        : <div style={S.cardValue}>{formatValue(value)}</div>
      }
    </div>
  );
}

export default function ResultBox({ result, error, loading }) {
  if (loading) return (
    <div style={S.box(false)}>
      <span style={{ color: "var(--muted)", fontSize: 13 }}>⏳ Calculando...</span>
    </div>
  );

  if (error) return (
    <div style={S.box(true)}>
      <span style={{ color: "var(--error)", fontSize: 13 }}>❌ {error}</span>
    </div>
  );

  if (!result) return null;

  const mainVal = result.result ?? result.simpson ?? result.derivative ??
                  result.position_m ?? result.force_n ?? result.y;

  const entries = Object.entries(result).filter(([k, v]) => {
    if (SKIP_KEYS.has(k)) return false;
    if (k === "t" || k === "y") return false;
    if (Array.isArray(v) && v.length > 20) return false;
    return true;
  });

  return (
    <div style={S.box(false)}>
      <div style={S.label}>Resultado</div>

      {mainVal !== undefined && (
        <div style={S.mainValue}>
          {isMatrix(mainVal)
            ? <MatrixDisplay data={mainVal} />
            : isVector(mainVal)
            ? <VectorDisplay data={mainVal} />
            : formatValue(mainVal)
          }
        </div>
      )}

      {entries.length > 0 && (
        <div style={S.grid}>
          {entries.map(([key, val]) => {
            if (key === "result" && mainVal !== undefined) return null;
            if (key === "simpson" && mainVal !== undefined) return null;
            const label = KEY_LABELS[key] || key.replace(/_/g, " ");
            return <FieldCard key={key} label={label} value={val} />;
          })}
        </div>
      )}
    </div>
  );
}