// ResultBox: componente reutilizável para exibir resultados da API
export default function ResultBox({ result, error, loading }) {
  const box = {
    marginTop: "1.5rem",
    background: "var(--surface)",
    border: `1px solid ${error ? "var(--error)" : "var(--border)"}`,
    borderRadius: "var(--radius)",
    padding: "1.2rem",
    minHeight: 60,
  };

  if (loading) return (
    <div style={box}>
      <span style={{ color: "var(--muted)", fontSize: 13 }}>⏳ Calculando...</span>
    </div>
  );

  if (error) return (
    <div style={box}>
      <span style={{ color: "var(--error)", fontSize: 13 }}>❌ {error}</span>
    </div>
  );

  if (!result) return null;

  return (
    <div style={box}>
      <div style={{ fontSize: 12, color: "var(--muted)", marginBottom: 8 }}>Resultado</div>
      <pre>{JSON.stringify(result, null, 2)}</pre>
    </div>
  );
}