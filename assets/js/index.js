async function registerUser() {
  const username = document.querySelector("[name='username']");
  const email = document.querySelector("[name='email']");

  if (username.value.trim().length === 0 && email.value.trim().length === 0) {
    alert("Please enter email or username");
  }

  try {
    await fetch("/users", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username: username.value + "",
        email: email.value + "",
      }),
    });
  } catch (err) {
    console.error(err);
  }
}

let ws;
let messages = document.getElementById("messages");

function testWs() {
  if (ws) {
    console.log("connection already established");
    return;
  }
  ws = new WebSocket("ws://localhost:8080/ws");

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

document.getElementById("wsButton").addEventListener("click", async (e) => {
  testWs();
  e.currentTarget.setAttribute("disabled", "disabled");
});

document.getElementById("closeWs").addEventListener("click", async (e) => {
  if (!ws) {
    console.log("No ws");
    return;
  }

  document.getElementById("wsButton").removeAttribute("disabled");

  ws.close();
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
