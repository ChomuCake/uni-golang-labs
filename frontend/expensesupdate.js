const urlParams = new URLSearchParams(window.location.search);
const expenseID = urlParams.get("expenseID");

fetch("/expenses/" + expenseID)
  .then((response) => response.json())
  .then((expense) => {
    const categoryInput = document.getElementById("update-category");
    const amountInput = document.getElementById("update-amount");
    const dateInput = document.getElementById("update-date");

    categoryInput.value = expense.category;
    amountInput.value = expense.amount;
    dateInput.value = expense.date; 
  })
  .catch((error) => {
    console.error("Error:", error);
  });

document.getElementById("update-expense-form").addEventListener("submit", function (e) {
  e.preventDefault();
  const form = e.target;
  const formData = new FormData(form);
  const data = {
    id: parseInt(expenseID),
    rawdate: formData.get("rawdate"), 
    category: formData.get("category"),
    amount: parseInt(formData.get("amount")),
  };
  const options = {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      Authorization: getToken(),
    },
    body: JSON.stringify(data),
  };

  fetch("/expenses/" + expenseID, options)
    .then((response) => {
      if (response.ok) {
        alert("Expense updated successfully");
        window.location.href = "expenses.html";
      } else {
        alert("Failed to update expense");
      }
    })
    .catch((error) => {
      console.error("Error:", error);
    });
});

// Function to retrieve token from local storage
function getToken() {
  return localStorage.getItem("token");
}