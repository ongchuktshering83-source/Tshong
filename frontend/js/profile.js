(() => {
  const loggedInStatus = localStorage.getItem("isLoggedIn");

  if (loggedInStatus !== "true") {
    window.location.href = "login.html";
    return;
  }

  const profileForm = document.getElementById("profileForm");
  const profileFullName = document.getElementById("profileFullName");
  const profileEmailInput = document.getElementById("profileEmailInput");
  const profileContactInput = document.getElementById("profileContactInput");
  const currentPassword = document.getElementById("currentPassword");
  const newPassword = document.getElementById("newPassword");
  const profileMessage = document.getElementById("profileMessage");
  const logoutBtn = document.getElementById("logoutBtn");

  const savedUser =
    JSON.parse(localStorage.getItem("currentUser")) ||
    JSON.parse(localStorage.getItem("tshongmartUser"));

  if (!savedUser || !savedUser.id) {
    localStorage.removeItem("isLoggedIn");
    window.location.href = "login.html";
    return;
  }

  profileFullName.value = savedUser.fullName || "";
  profileEmailInput.value = savedUser.email || "";
  profileContactInput.value = savedUser.contactInfo || "";

  if (logoutBtn) {
    logoutBtn.addEventListener("click", () => {
      localStorage.removeItem("isLoggedIn");
      localStorage.removeItem("currentUser");
      localStorage.removeItem("tshongmartUser");
      window.location.href = "index.html";
    });
  }

  profileForm.addEventListener("submit", async function (event) {
    event.preventDefault();

    const payload = {
      id: savedUser.id,
      fullName: profileFullName.value.trim(),
      email: profileEmailInput.value.trim(),
      contactInfo: profileContactInput.value.trim(),
      currentPassword: currentPassword.value,
      newPassword: newPassword.value.trim()
    };

    try {
      const response = await fetch("/api/profile", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
      });

      const result = await response.json();

      if (!result.success) {
        profileMessage.textContent = result.message;
        profileMessage.className = "form-message error";
        return;
      }

      localStorage.setItem("currentUser", JSON.stringify(result.user));
      localStorage.setItem("tshongmartUser", JSON.stringify(result.user));

      profileMessage.textContent = "Profile updated successfully. Redirecting to dashboard...";
      profileMessage.className = "form-message success";

      profileForm.reset();

      setTimeout(() => {
        window.location.href = "dashboard.html";
      }, 1200);
    } catch (error) {
      profileMessage.textContent = "Something went wrong. Please try again.";
      profileMessage.className = "form-message error";
    }
  });
})();