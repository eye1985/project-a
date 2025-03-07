import { registerUser } from "./functions.js";

let ws;
let messages = document.getElementById("messages");
const username = document.getElementById("username");

function testWs(username) {
  if (ws) {
    console.log("connection already established");
    return;
  }

  ws = new WebSocket(`ws://localhost:8080/ws?username=${username}`);

  ws.onopen = () => {
    console.log("Connected to WebSocket server");
  };

  ws.onmessage = (event) => {
    messages = document.getElementById("messages");
    const newMessage = document.createElement("p");
    newMessage.innerText = event.data;

    messages.appendChild(newMessage);
    messages.scrollTo(0, messages.scrollHeight);
  };

  ws.onclose = (event) => {
    console.log("WebSocket closed:", event.code, event.reason);
  };

  ws.onerror = (error) => {
    console.error("WebSocket error:", error);
  };
}

document.getElementById("registerBtn").addEventListener("click", async (e) => {
  e.preventDefault();
  await registerUser();
});

document
  .getElementById("connectToWsBtn")
  .addEventListener("click", async (e) => {
    if (username.value === "") {
      console.log("username is required");
      return;
    }

    testWs(username.value);
    e.currentTarget.setAttribute("disabled", "disabled");
  });

document.getElementById("closeWs").addEventListener("click", async (e) => {
  if (!ws) {
    console.log("No ws");
    return;
  }

  document.getElementById("connectToWsBtn").removeAttribute("disabled");

  ws.close();
  ws = null;
});

const messageInput = document.getElementById("messageInput");
messageInput.addEventListener("keypress", async (e) => {
  if (ws && e.key === "Enter") {
    ws.send(e.currentTarget.value);
    e.currentTarget.value = "";
  }
});
document.getElementById("sendWs").addEventListener("click", async (e) => {
  const value = document.getElementById("messageInput").value;
  ws.send(value);
});
