(() => {
  const loggedInStatus = localStorage.getItem("isLoggedIn");

  if (loggedInStatus !== "true") {
    window.location.href = "login.html";
    return;
  }

  const addProductForm = document.getElementById("addProductForm");
  const productMessage = document.getElementById("productMessage");
  const logoutBtn = document.getElementById("logoutBtn");

  const currentUser = JSON.parse(localStorage.getItem("currentUser"));

  if (!currentUser || !currentUser.id) {
    localStorage.removeItem("isLoggedIn");
    window.location.href = "login.html";
    return;
  }

  if (logoutBtn) {
    logoutBtn.addEventListener("click", () => {
      localStorage.removeItem("isLoggedIn");
      localStorage.removeItem("currentUser");
      localStorage.removeItem("adminAccess");
      window.location.href = "index.html";
    });
  }

  if (!addProductForm) return;

  addProductForm.addEventListener("submit", async function (event) {
    event.preventDefault();

    const imageInput = document.getElementById("productImage");

    const formData = new FormData();
    formData.append("userId", currentUser.id);
    formData.append("title", document.getElementById("productTitle").value.trim());
    formData.append("category", document.getElementById("productCategory").value);
    formData.append("price", document.getElementById("productPrice").value.trim());
    formData.append("location", document.getElementById("productLocation").value.trim());
    formData.append("contact", document.getElementById("productContact").value.trim());
    formData.append("description", document.getElementById("productDescription").value.trim());

    if (imageInput.files && imageInput.files[0]) {
      formData.append("image", imageInput.files[0]);
    }

    try {
      const response = await fetch("/api/products", {
        method: "POST",
        body: formData
      });

      const result = await response.json();

      if (!result.success) {
        productMessage.textContent = result.message;
        productMessage.className = "form-message error";
        return;
      }

      productMessage.textContent = "Product added successfully. Redirecting to dashboard...";
      productMessage.className = "form-message success";

      addProductForm.reset();

      setTimeout(() => {
        window.location.href = "dashboard.html";
      }, 1200);
    } catch (error) {
      productMessage.textContent = "Something went wrong. Please try again.";
      productMessage.className = "form-message error";
    }
  });
})();