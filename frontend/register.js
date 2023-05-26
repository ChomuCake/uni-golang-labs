// Set content type header to JSON for registration form
document
  .querySelector('form[action="/register"]')
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
          alert("Registration successful");
          window.location.href = "login.html"; // Перехід на login.html
        } else {
          alert("Registration failed");
        }
      })
      .catch((error) => {
        console.error("Error:", error);
      });
  });
