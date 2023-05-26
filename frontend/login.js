// Function to save token in local storage
function saveToken(token) {
  localStorage.setItem("token", token);
}

// Function to retrieve token from local storage
function getToken() {
  return localStorage.getItem("token");
}

// JSON for log form
document
  .querySelector('form[action="/login"]')
  .addEventListener("submit", function (e) {
    e.preventDefault();
    const form = e.target;
    const formData = new FormData(form);
    const data = Object.fromEntries(formData.entries());
    const options = {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    };

    fetch(form.action, options)
      .then((response) => {
        if (response.ok) {
          alert("Login successful");
          const token = response.headers.get("Authorization");
          saveToken(token);
          window.location.href = "expenses.html"; // Перехід на expenses.html
        } else {
          alert("Login failed");
        }
      })
      .catch((error) => {
        console.error("Error:", error);
      });
  });
