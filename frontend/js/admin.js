(() => {
  const loggedInStatus = localStorage.getItem("isLoggedIn");
  const adminAccess = localStorage.getItem("adminAccess");

  if (loggedInStatus !== "true") {
    window.location.href = "login.html";
    return;
  }

  if (adminAccess !== "true") {
    window.location.href = "admin-login.html";
    return;
  }

  const adminUserCount = document.getElementById("adminUserCount");
  const adminProductCount = document.getElementById("adminProductCount");
  const adminProductTable = document.getElementById("adminProductTable");
  const adminUsersTable = document.getElementById("adminUsersTable");
  const logoutBtn = document.getElementById("logoutBtn");
  const adminMessagesList = document.getElementById("adminMessagesList");

  function formatDate(dateText) {
    if (!dateText) return "N/A";

    const date = new Date(dateText);

    if (Number.isNaN(date.getTime())) {
      return dateText.split(".")[0];
    }

    return date.toLocaleDateString();
  }

  async function loadAdminUsers() {
    if (!adminUsersTable) return;

    try {
      const response = await fetch("/api/users");
      const users = await response.json();

      if (adminUserCount) {
        adminUserCount.textContent = users.length;
      }

      if (!users || users.length === 0) {
        adminUsersTable.innerHTML = `
          <tr>
            <td colspan="7" class="admin-empty-cell">
              No registered users found.
            </td>
          </tr>
        `;
        return;
      }

      adminUsersTable.innerHTML = "";

      users.forEach((user) => {
        const row = document.createElement("tr");
        const isBanned = user.status === "banned";

        row.innerHTML = `
          <td>${user.fullName || "Unnamed User"}</td>
          <td>${user.email || "N/A"}</td>
          <td>${user.contactInfo || "N/A"}</td>
          <td>${user.role || "user"}</td>
          <td>
            <span class="user-status-badge ${isBanned ? "status-banned" : "status-active"}">
              ${isBanned ? "Banned" : "Active"}
            </span>
          </td>
          <td>${formatDate(user.createdAt)}</td>
          <td>
            <button
              class="${isBanned ? "unban-user-btn" : "ban-user-btn"}"
              data-id="${user.id}"
              data-status="${isBanned ? "active" : "banned"}"
            >
              ${isBanned ? "Unban" : "Ban"}
            </button>
          </td>
        `;

        adminUsersTable.appendChild(row);
      });

      const statusButtons = document.querySelectorAll(".ban-user-btn, .unban-user-btn");

      statusButtons.forEach((button) => {
        button.addEventListener("click", async () => {
          const userId = button.dataset.id;
          const newStatus = button.dataset.status;

          const confirmMessage =
            newStatus === "banned"
              ? "Are you sure you want to ban this user?"
              : "Are you sure you want to unban this user?";

          const confirmAction = confirm(confirmMessage);

          if (!confirmAction) return;

          try {
            const response = await fetch(`/api/users/status?id=${userId}&status=${newStatus}`, {
              method: "PUT"
            });

            const result = await response.json();

            if (result.success) {
              loadAdminUsers();
              loadAdminProducts();
            } else {
              alert(result.message);
            }
          } catch (error) {
            alert("Could not update user status. Please try again.");
          }
        });
      });
    } catch (error) {
      adminUsersTable.innerHTML = `
        <tr>
          <td colspan="7" class="admin-empty-cell">
            Could not load users. Please refresh the page.
          </td>
        </tr>
      `;
    }
  }

  async function loadAdminProducts() {
    if (!adminProductTable) return;

    try {
      const response = await fetch("/api/products");
      const products = await response.json();

      if (adminProductCount) {
        adminProductCount.textContent = products.length;
      }

      renderProducts(products);
    } catch (error) {
      adminProductTable.innerHTML = `
        <tr>
          <td colspan="7" class="admin-empty-cell">
            Could not load products. Please refresh the page.
          </td>
        </tr>
      `;
    }
  }

  function renderProducts(products) {
    adminProductTable.innerHTML = "";

    if (!products || products.length === 0) {
      adminProductTable.innerHTML = `
        <tr>
          <td colspan="7" class="admin-empty-cell">
            No active product listings yet.
          </td>
        </tr>
      `;
      return;
    }

    products.forEach((product) => {
      const row = document.createElement("tr");

      row.innerHTML = `
        <td>
          <div class="admin-product-thumb">
            ${
              product.imagePath
                ? `<img src="${product.imagePath}" alt="${product.title}" />`
                : `<span>🚗</span>`
            }
          </div>
        </td>

        <td>${product.title || "Untitled Product"}</td>
        <td>${product.category || "N/A"}</td>
        <td>${product.seller || "Unknown Seller"}</td>
        <td>${product.location || "N/A"}</td>
        <td>${product.price || "N/A"}</td>
        <td>
          <button class="delete-listing-btn admin-delete-btn" data-id="${product.id}">
            Delete
          </button>
        </td>
      `;

      adminProductTable.appendChild(row);
    });

    const deleteButtons = document.querySelectorAll(".admin-delete-btn");

    deleteButtons.forEach((button) => {
      button.addEventListener("click", async () => {
        const productId = button.dataset.id;

        const confirmDelete = confirm("Are you sure you want to delete this listing?");

        if (!confirmDelete) return;

        try {
          const response = await fetch(`/api/products?id=${productId}`, {
            method: "DELETE"
          });

          const result = await response.json();

          if (result.success) {
            loadAdminProducts();
          } else {
            alert(result.message);
          }
        } catch (error) {
          alert("Could not delete product. Please try again.");
        }
      });
    });
  }

  async function loadContactMessages() {
    if (!adminMessagesList) return;

    try {
      const response = await fetch("/api/contact-messages");
      const messages = await response.json();

      if (!messages || messages.length === 0) {
        adminMessagesList.innerHTML = `
          <div class="empty-dashboard-box">
            <h3>No contact messages yet</h3>
            <p>Messages from the contact form will appear here.</p>
          </div>
        `;
        return;
      }

      adminMessagesList.innerHTML = "";

      messages.forEach((message) => {
        const messageCard = document.createElement("article");
        messageCard.className = "admin-message-card";

        messageCard.innerHTML = `
          <div class="admin-message-top">
            <div>
              <h3>${message.fullName}</h3>
              <p>${message.email}</p>
            </div>
            <span>${message.subject}</span>
          </div>

          <p class="admin-message-text">${message.message}</p>

          <small>${formatDate(message.createdAt)}</small>
        `;

        adminMessagesList.appendChild(messageCard);
      });
    } catch (error) {
      adminMessagesList.innerHTML = `
        <div class="empty-dashboard-box">
          <h3>Could not load contact messages</h3>
          <p>Please refresh the page or try again later.</p>
        </div>
      `;
    }
  }

  if (logoutBtn) {
    logoutBtn.addEventListener("click", () => {
      localStorage.removeItem("isLoggedIn");
      localStorage.removeItem("currentUser");
      localStorage.removeItem("adminAccess");
      window.location.href = "index.html";
    });
  }

  loadAdminUsers();
  loadAdminProducts();
  loadContactMessages();
})();