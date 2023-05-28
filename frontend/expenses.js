// Function to save token in local storage
function saveToken(token) {
  localStorage.setItem("token", token);
}

// Function to retrieve token from local storage
function getToken() {
  return localStorage.getItem("token");
}

function deleteExpense(expenseID) {
  const options = {
    method: "DELETE",
    headers: {
      Authorization: getToken(),
    },
  };
  fetch("/expenses/" + expenseID, options)
    .then((response) => {
      if (response.ok) {
        alert("Expense deleted successfully");
        fetchExpenses(); // Refresh the expenses table
      } else {
        alert("Failed to delete expense");
      }
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}

// Fetch expenses data and display them in the table
function fetchExpenses(sortBy) {
  let url = "/expenses";
  if (sortBy) {
    url += `?sort=${sortBy}`;
  }

  const options = {
    headers: {
      Authorization: getToken(),
    },
  };

  fetch(url, options)
    .then((response) => response.json())
    .then((expenses) => {
      const expensesList = document.getElementById("expenses-list");
      expensesList.innerHTML = "";

      let totalAmount = 0;

      expenses.forEach((expense) => {
        const row = document.createElement("tr");
        const categoryCell = document.createElement("td");
        const amountCell = document.createElement("td");
        const actionCell = document.createElement("td");
        const deleteButton = document.createElement("button");
        const updateButton = document.createElement("button");
      
        categoryCell.innerText = expense.category;
        amountCell.innerText = expense.amount;
        deleteButton.innerText = "Delete";
        updateButton.innerText = "Update";
      
        deleteButton.addEventListener("click", function () {
          deleteExpense(expense.id);
        });
      
        updateButton.addEventListener("click", function () {
            openUpdateExpensePage(expense.id);
        });
      
        totalAmount += expense.amount;
        actionCell.appendChild(deleteButton);
        actionCell.appendChild(updateButton);
        row.appendChild(categoryCell);
        row.appendChild(amountCell);
        row.appendChild(actionCell);
        expensesList.appendChild(row);
      });
      

      const totalExpenses = document.getElementById("total-expenses");
      totalExpenses.innerText = `Total: $${totalAmount}`;
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}

document
  .querySelector('form[action="/expenses"]')
  .addEventListener("submit", function (e) {
    e.preventDefault();
    const form = e.target;
    const formData = new FormData(form);
    const data = {
      category: formData.get("category"),
      amount: parseInt(formData.get("amount")),
    };
    const options = {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: getToken(),
      },
      body: JSON.stringify(data),
    };

    fetch(form.action, options)
      .then((response) => {
        if (response.ok) {
          alert("Expenses add successful");
        } else {
          alert("Expenses add failure");
        }
      })
      .catch((error) => {
        console.error("Error:", error);
      });
  });

// Get Expenses Button Event Listener
document.getElementById("get-expenses").addEventListener("click", function () {
  const sortBy = document.getElementById("sort-by").value;
  fetchExpenses(sortBy);
});

function openUpdateExpensePage(expenseID) {
    window.location.href = "expensesupdate.html?expenseID=" + expenseID;
  }
  
