import React, { useEffect, useRef, useState } from 'react';

function emptyBoard() {
  return Array.from({ length: 6 }).map(() => Array.from({ length: 7 }).map(() => 0));
}

export default function Game({ username, wsUrl, setGameId, gameId, onExit }) {
  const [board, setBoard] = useState(emptyBoard());
  const [status, setStatus] = useState("waiting");
  const [you, setYou] = useState(1);
  const [opponent, setOpponent] = useState(null);
  const [msg, setMsg] = useState("");
  const [wins, setWins] = useState({});
  const [rematchRequested, setRematchRequested] = useState(false);
  const wsRef = useRef(null);

  useEffect(() => {
    const ws = new WebSocket(wsUrl);
    wsRef.current = ws;

    ws.onopen = () => {
      ws.send(JSON.stringify({ type: "join", username }));
      setMsg("Joined queue...");
    };

    ws.onmessage = (ev) => {
      const m = JSON.parse(ev.data);
      switch (m.type) {
        case "matched":
          setGameId(m.gameId);
          setOpponent(m.opponent);
          setYou(m.you);
          setStatus("ongoing");
          setMsg("Matched with " + m.opponent);
          break;

        case "state":
          setBoard(m.board);
          setMsg(m.nextTurn === username ? "Your turn" : `${m.nextTurn}'s turn`);
          break;

        case "result":
          if (m.result === "win") {
            setMsg("Winner: " + m.winner);
            setStatus("finished");
            updateWins(m.winner);
          } else {
            setMsg("Draw");
            setStatus("finished");
          }
          break;

        case "rematch_start":
          setBoard(emptyBoard());
          setStatus("ongoing");
          setMsg("Rematch started!");
          setRematchRequested(false);
          break;

        case "info":
          setMsg(m.message);
          break;

        case "error":
          setMsg("Error: " + m.message);
          break;

        default:
          console.warn("Unknown message type:", m);
      }
    };

    ws.onerror = () => setMsg("WebSocket error");
    return () => ws.close();
  }, [username, wsUrl, setGameId]);

  const updateWins = (winner) => {
    setWins((prev) => ({
      ...prev,
      [winner]: (prev[winner] || 0) + 1
    }));
  };

  const drop = (col) => {
    if (!wsRef.current || status !== "ongoing") return;
    wsRef.current.send(JSON.stringify({ type: "drop", column: col, gameId, username }));
  };

  const handleRematch = () => {
    if (!wsRef.current) return;
    setRematchRequested(true);
    wsRef.current.send(JSON.stringify({ type: "rematch_request", gameId }));
    setMsg("Rematch requested! Waiting for opponent...");
  };

  const cellColor = (cell) => {
    if (cell === 0) return "#eee";
    if (cell === 1) return "red";
    if (cell === 2) return "yellow";
    return "#eee";
  };

  return (
    <div style={{
      display: "flex",
      justifyContent: "center",
      alignItems: "flex-start",
      minHeight: "100vh",
      backgroundColor: "#f0f4f8",
      fontFamily: "Arial, sans-serif",
      padding: 20,
      gap: 40
    }}>
      <div style={{ textAlign: "center" }}>
        <h1 style={{ color: "#333" }}>4 in a Row (Connect Four)</h1>
        <div style={{ margin: "10px 0", fontWeight: "bold", color: "#555" }}>{msg}</div>

        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(7, 60px)', gap: 4 }}>
          {board[0].map((_, c) => (
            <button key={"col" + c} onClick={() => drop(c)}
              style={{ width: 60, height: 30, borderRadius: 5, backgroundColor: "#007BFF", color: "white", fontWeight: "bold", cursor: "pointer", border: "none" }}>
              Drop
            </button>
          ))}
        </div>

        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(7, 60px)',
          gap: 4,
          marginTop: 10,
          backgroundColor: "#0056b3",
          padding: 5,
          borderRadius: 10
        }}>
          {board.flatMap((row, r) => row.map((cell, c) => (
            <div key={r + "_" + c} style={{
              width: 60,
              height: 60,
              backgroundColor: cellColor(cell),
              border: '2px solid #333',
              borderRadius: '50%',
              transition: "background 0.3s"
            }} />
          )))}
        </div>

        <div style={{ marginTop: 12, fontSize: 16 }}>
          <strong>Opponent:</strong> {opponent || "-"}
        </div>

        <div style={{ marginTop: 20 }}>
          {status === "finished" && (
            <button
              onClick={handleRematch}
              disabled={rematchRequested}
              style={{
                padding: "8px 16px",
                backgroundColor: rematchRequested ? "#6c757d" : "#28a745",
                color: "white",
                border: "none",
                borderRadius: 5,
                cursor: rematchRequested ? "default" : "pointer",
                fontWeight: "bold",
                marginRight: 10
              }}
            >
              {rematchRequested ? "Waiting for opponent..." : "Rematch"}
            </button>
          )}
          <button onClick={onExit} style={{
            padding: "8px 16px",
            backgroundColor: "#dc3545",
            color: "white",
            border: "none",
            borderRadius: 5,
            cursor: "pointer",
            fontWeight: "bold"
          }}>
            Exit
          </button>
        </div>
      </div>

      <div style={{
        backgroundColor: "white",
        border: "2px solid #ccc",
        borderRadius: 10,
        padding: 20,
        minWidth: 200,
        textAlign: "center",
        boxShadow: "0px 0px 8px rgba(0,0,0,0.1)"
      }}>
        <h3 style={{ marginBottom: 10 }}>Leaderboard 🏆</h3>
        <div style={{ fontSize: 16, fontWeight: "bold", color: "#333" }}>
          <div>{username}: {wins[username] || 0}</div>
          <div>{opponent || "Opponent"}: {wins[opponent] || 0}</div>
        </div>
      </div>
    </div>
  );
}
