import React, { useState } from "react";
import Game from "./Game";

export default function App() {
  const [username, setUsername] = useState("");
  const [page, setPage] = useState("start"); // "start" or "game"
  const [wsUrl, setWsUrl] = useState("ws://localhost:8080/ws");
  const [gameId, setGameId] = useState(null);

  return (
    <div style={{ padding: 20, textAlign: "center" }}>
      <h2>4 in a Row (Connect Four)</h2>

      {page === "start" && (
        <div>
          <input
            placeholder="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <br />
          <input
            placeholder="ws url"
            value={wsUrl}
            onChange={(e) => setWsUrl(e.target.value)}
            style={{ width: 400, marginTop: 8 }}
          />
          <br />
          <button
            onClick={() => setPage("game")}
            disabled={!username}
            style={{ marginTop: 10 }}
          >
            Join
          </button>
        </div>
      )}

      {page === "game" && (
        <Game
          username={username}
          wsUrl={wsUrl}
          setGameId={setGameId}
          gameId={gameId}
          onExit={() => setPage("start")} // exit button will go back to start
        />
      )}
    </div>
  );
}
