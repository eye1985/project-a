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

function testWs() {
  if (ws) {
    console.log("connection already established");
    return;
  }
  ws = new WebSocket("ws://localhost:8080/ws");

  ws.onopen = () => {
    console.log("Connected to WebSocket server");
    ws.send("Hello, server!"); // Send a message to the server
  };

  ws.onmessage = (event) => {
    const messages = document.getElementById("messages");
    messages.innerText += event.data + "\n\n";
    console.log("Received from server:", event.data);
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

document.getElementById("sendWs").addEventListener("click", async (e) => {
  ws.send("Hello");
});
