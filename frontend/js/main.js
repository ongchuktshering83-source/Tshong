const navLinks = document.getElementById("navLinks");
const navActions = document.querySelector(".nav-actions");

const currentPage = window.location.pathname.split("/").pop() || "index.html";
const isLoggedIn = localStorage.getItem("isLoggedIn") === "true";
const hasAdminAccess = localStorage.getItem("adminAccess") === "true";

function setActive(pageName) {
  return currentPage === pageName ? "active" : "";
}

function renderNavbar() {
  if (navLinks) {
    navLinks.innerHTML = `
      <li><a href="index.html" class="${setActive("index.html")}">Home</a></li>
      <li><a href="about.html" class="${setActive("about.html")}">About</a></li>
      <li><a href="products.html" class="${setActive("products.html")}">Products</a></li>
      <li><a href="contact.html" class="${setActive("contact.html")}">Contact</a></li>
     ${
       isLoggedIn
         ? `
      <li><a href="dashboard.html" class="${setActive("dashboard.html")}">Dashboard</a></li>
      ${
        hasAdminAccess
          ? `<li><a href="admin.html" class="${setActive("admin.html")}">Admin</a></li>`
          : `<li><a href="admin-login.html" class="${setActive("admin-login.html")}">Admin Login</a></li>`
      }
    `
         : ""
     }
    `;
  }
  if (navActions) {
    if (isLoggedIn) {
      navActions.innerHTML = `
        <button class="logout-btn" id="globalLogoutBtn">Logout</button>
        <button class="menu-btn" id="menuBtn">☰</button>
      `;
    } else {
      navActions.innerHTML = `
        <a href="login.html" class="login-link">Login</a>
        <a href="signup.html" class="btn btn-small">Sign Up</a>
        <button class="menu-btn" id="menuBtn">☰</button>
      `;
    }
  }

  const menuBtn = document.getElementById("menuBtn");

  if (menuBtn && navLinks) {
    menuBtn.addEventListener("click", () => {
      navLinks.classList.toggle("show");
    });
  }

  const globalLogoutBtn = document.getElementById("globalLogoutBtn");

  if (globalLogoutBtn) {
    globalLogoutBtn.addEventListener("click", () => {
      localStorage.removeItem("isLoggedIn");
      localStorage.removeItem("currentUser");
      localStorage.removeItem("adminAccess");
      window.location.href = "index.html";
    });
  }
}

renderNavbar();
