(() => {
  const loggedInStatus = localStorage.getItem("isLoggedIn");

  if (loggedInStatus !== "true") {
    window.location.href = "login.html";
    return;
  }

  const savedUser =
    JSON.parse(localStorage.getItem("currentUser")) ||
    JSON.parse(localStorage.getItem("tshongmartUser"));

  if (!savedUser || !savedUser.id) {
    localStorage.removeItem("isLoggedIn");
    window.location.href = "login.html";
    return;
  }

  const dashboardName = document.getElementById("dashboardName");
  const profileName = document.getElementById("profileName");
  const profileEmail = document.getElementById("profileEmail");
  const profileContact = document.getElementById("profileContact");
  const profileAvatar = document.getElementById("profileAvatar");
  const logoutBtn = document.getElementById("logoutBtn");
  const myListings = document.getElementById("myListings");
  const listingCount = document.getElementById("listingCount");
  const savedProductsList = document.getElementById("savedProductsList");

  dashboardName.textContent = savedUser.fullName || "User";
  profileName.textContent = savedUser.fullName || "User Name";
  profileEmail.textContent = savedUser.email || "user@email.com";
  profileContact.textContent = savedUser.contactInfo || "Not added";

  const initials = (savedUser.fullName || "TM")
    .split(" ")
    .map((name) => name.charAt(0))
    .join("")
    .slice(0, 2)
    .toUpperCase();

  profileAvatar.textContent = initials;

  function displayUserListings(products) {
    listingCount.textContent = products.length;

    if (products.length === 0) {
      myListings.innerHTML = `
        <div class="empty-dashboard-box">
          <h3>You have not listed any product yet</h3>
          <p>Add your first automotive product and it will appear here.</p>
          <a href="add-product.html" class="btn btn-small">Add Product</a>
        </div>
      `;
      return;
    }

    myListings.innerHTML = "";

    products.forEach((product) => {
      const item = document.createElement("article");
      item.className = "listing-item";

      item.innerHTML = `
        <div class="listing-icon">
  ${
    product.imagePath
      ? `<img src="${product.imagePath}" alt="${product.title}" />`
      : `<span>🚗</span>`
  }
</div>

        <div class="listing-info">
          <h3>${product.title}</h3>
          <p>${product.category} • ${product.location}</p>
        </div>

        <strong>${product.price}</strong>

        <button class="delete-listing-btn" data-id="${product.id}">Delete</button>
      `;

      myListings.appendChild(item);
    });

    const deleteButtons = document.querySelectorAll(".delete-listing-btn");

    deleteButtons.forEach((button) => {
      button.addEventListener("click", async () => {
        const productId = button.dataset.id;
        const confirmDelete = confirm(
          "Are you sure you want to delete this product?",
        );

        if (!confirmDelete) return;

        try {
          const response = await fetch(`/api/products?id=${productId}`, {
            method: "DELETE",
          });

          const result = await response.json();

          if (result.success) {
            loadMyProducts();
          } else {
            alert(result.message);
          }
        } catch (error) {
          alert("Could not delete product. Please try again.");
        }
      });
    });
  }

  function makeContactLink(contact) {
  if (!contact) return "#";

  if (contact.startsWith("http")) {
    return contact;
  }

  const cleanNumber = contact.replace(/\D/g, "");

  if (cleanNumber.length >= 8) {
    return `https://wa.me/975${cleanNumber}`;
  }

  return "#";
}

function displaySavedProducts() {
  if (!savedProductsList) return;

  const savedKey = `savedProducts_${savedUser.id}`;
  const savedProducts = JSON.parse(localStorage.getItem(savedKey)) || [];

  if (savedProducts.length === 0) {
    savedProductsList.innerHTML = `
      <div class="empty-dashboard-box">
        <h3>No saved products yet</h3>
        <p>Products you save from the marketplace will appear here.</p>
        <a href="products.html" class="btn btn-small">Browse Products</a>
      </div>
    `;
    return;
  }

  savedProductsList.innerHTML = "";

  savedProducts.forEach((product) => {
    const item = document.createElement("article");
    item.className = "listing-item";

    item.innerHTML = `
      <div class="listing-icon">
        ${
          product.imagePath
            ? `<img src="${product.imagePath}" alt="${product.title}" />`
            : `<span>🚗</span>`
        }
      </div>

      <div class="listing-info">
        <h3>${product.title}</h3>
        <p>${product.category} • ${product.location} • Seller: ${product.seller}</p>
      </div>

      <strong>${product.price}</strong>

      <a href="${makeContactLink(product.contact)}" target="_blank" class="btn btn-small">
        Contact
      </a>
    `;

    savedProductsList.appendChild(item);
  });
}

  async function loadMyProducts() {
    try {
      const response = await fetch(`/api/my-products?userId=${savedUser.id}`);
      const products = await response.json();

      displayUserListings(products);
    } catch (error) {
      myListings.innerHTML = `
        <div class="empty-dashboard-box">
          <h3>Could not load your products</h3>
          <p>Please refresh the page or try again later.</p>
        </div>
      `;
    }
  }

  if (logoutBtn) {
    logoutBtn.addEventListener("click", () => {
      localStorage.removeItem("isLoggedIn");
      localStorage.removeItem("currentUser");
      window.location.href = "index.html";
    });
  }

  loadMyProducts();
  displaySavedProducts();
})();
