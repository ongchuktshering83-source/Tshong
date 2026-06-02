(() => {
  const productsGrid = document.getElementById("productsGrid");
  const searchInput = document.getElementById("searchInput");
  const categoryFilter = document.getElementById("categoryFilter");

  let products = [];
  let currentUser = null;

  try {
    currentUser = JSON.parse(localStorage.getItem("currentUser"));
  } catch {
    currentUser = null;
  }

  function getIcon(category) {
    if (category === "bike") return "🏍️";
    if (category === "parts") return "⚙️";
    if (category === "pickup") return "🚙";
    return "🚗";
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

  function setupSaveButtons() {
    const saveButtons = document.querySelectorAll(".save-btn:not(:disabled)");

    saveButtons.forEach((button) => {
      button.addEventListener("click", () => {
        if (!currentUser) {
          alert("Please login first to save products.");
          window.location.href = "login.html";
          return;
        }

        const productId = Number(button.dataset.id);

        const productToSave = products.find(
          (product) => Number(product.id) === productId
        );

        if (!productToSave) {
          alert("Product not found.");
          return;
        }

        const savedKey = `savedProducts_${currentUser.id}`;
        const savedProducts = JSON.parse(localStorage.getItem(savedKey)) || [];

        const alreadySaved = savedProducts.some(
          (product) => Number(product.id) === Number(productToSave.id)
        );

        if (alreadySaved) {
          alert("This product is already saved in your dashboard.");
          return;
        }

        savedProducts.push(productToSave);
        localStorage.setItem(savedKey, JSON.stringify(savedProducts));

        alert("Product saved to your dashboard.");
      });
    });
  }

  function displayProducts(productList) {
    if (!productsGrid) return;

    productsGrid.innerHTML = "";

    if (productList.length === 0) {
      productsGrid.innerHTML = `
        <div class="empty-state">
          <h3>No products found</h3>
          <p>Try searching with another keyword or category.</p>
        </div>
      `;
      return;
    }

    productList.forEach((product) => {
      const isOwnProduct =
        currentUser &&
        product.userId &&
        Number(product.userId) === Number(currentUser.id);

      const productCard = document.createElement("article");
      productCard.className = "product-card";

      productCard.innerHTML = `
        <div class="product-image">
          ${
            product.imagePath
              ? `<img src="${product.imagePath}" alt="${product.title}" />`
              : `<span>${getIcon(product.category)}</span>`
          }
        </div>

        <div class="product-content">
          <div class="product-badge">Live Listing</div>

          ${
            isOwnProduct
              ? `<div class="product-badge own-badge">Your Listing</div>`
              : ""
          }

          <div class="product-category">${product.category}</div>

          <h3>${product.title}</h3>
          <p>${product.description}</p>

          <div class="product-meta">
            <span>${product.price}</span>
            <span>${product.location}</span>
          </div>

          <div class="seller-info">
            <small>Seller: ${product.seller}</small>
          </div>

          <div class="product-actions">
            ${
              isOwnProduct
                ? `<button class="save-btn" disabled>Your Product</button>`
                : `
                  <a href="${makeContactLink(product.contact)}" target="_blank" class="btn btn-small">
                    Contact Seller
                  </a>
                  <button class="save-btn" data-id="${product.id}">Save</button>
                `
            }
          </div>
        </div>
      `;

      productsGrid.appendChild(productCard);
    });

    setupSaveButtons();
  }

  function filterProducts() {
    const searchValue = searchInput.value.toLowerCase();
    const selectedCategory = categoryFilter.value;

    const filteredProducts = products.filter((product) => {
      const matchesSearch =
        product.title.toLowerCase().includes(searchValue) ||
        product.location.toLowerCase().includes(searchValue) ||
        product.seller.toLowerCase().includes(searchValue) ||
        product.category.toLowerCase().includes(searchValue);

      const matchesCategory =
        selectedCategory === "all" || product.category === selectedCategory;

      return matchesSearch && matchesCategory;
    });

    displayProducts(filteredProducts);
  }

  async function loadProducts() {
    try {
      const response = await fetch("/api/products");
      products = await response.json();

      displayProducts(products);
    } catch (error) {
      console.error("Failed to load products:", error);

      productsGrid.innerHTML = `
        <div class="empty-state">
          <h3>Could not load products</h3>
          <p>Please refresh the page or try again later.</p>
        </div>
      `;
    }
  }

  if (searchInput && categoryFilter && productsGrid) {
    searchInput.addEventListener("input", filterProducts);
    categoryFilter.addEventListener("change", filterProducts);

    loadProducts();
  }
})();