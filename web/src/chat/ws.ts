let ws: WebSocket | null = null;
let reconnectTimeout: ReturnType<typeof setTimeout> | null = null;

const WS_URL = import.meta.env.VITE_BACKEND_URL || "ws://localhost:8080/ws";
const RECONNECT_DELAY = 2000;

function connect() {
  ws = new WebSocket(WS_URL);

  ws.addEventListener("open", () => {
    console.log("WebSocket connection established");

    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
      reconnectTimeout = null;
    }
  });

  ws.addEventListener("close", () => {
    console.log("WebSocket connection closed");

    reconnectTimeout = setTimeout(() => {
      connect();
    }, RECONNECT_DELAY);
  });

  ws.addEventListener("error", (e) => {
    console.error("WebSocket error: ", e);

    ws?.close();
  });
}

export const getWs = () => {
  if (!ws || ws.readyState === WebSocket.CLOSED) {
    connect();
  }
  return ws;
};

export const closeWs = () => {
  if (reconnectTimeout) {
    clearTimeout(reconnectTimeout);
    reconnectTimeout = null;
  }

  ws?.close();
  ws = null;
};
