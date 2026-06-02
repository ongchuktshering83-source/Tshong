(() => {
  const adminLoginForm = document.getElementById("adminLoginForm");
  const adminPassword = document.getElementById("adminPassword");
  const adminLoginMessage = document.getElementById("adminLoginMessage");

  if (!adminLoginForm) return;

  adminLoginForm.addEventListener("submit", async function (event) {
    event.preventDefault();

    try {
      const response = await fetch("/api/admin-login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          password: adminPassword.value
        })
      });

      const result = await response.json();

      if (!result.success) {
        adminLoginMessage.textContent = result.message;
        adminLoginMessage.className = "form-message error";
        return;
      }

      localStorage.setItem("adminAccess", "true");

      adminLoginMessage.textContent = "Admin access granted. Redirecting...";
      adminLoginMessage.className = "form-message success";

      setTimeout(() => {
        window.location.href = "admin.html";
      }, 1000);
    } catch (error) {
      adminLoginMessage.textContent = "Something went wrong. Please try again.";
      adminLoginMessage.className = "form-message error";
    }
  });
})();