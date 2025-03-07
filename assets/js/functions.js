export async function registerUser() {
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
