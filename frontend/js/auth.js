(() => {
  const signupForm = document.getElementById("signupForm");
  const loginForm = document.getElementById("loginForm");

  if (signupForm) {
    signupForm.addEventListener("submit", async function (event) {
      event.preventDefault();

      const fullName = document.getElementById("fullName").value.trim();
      const email = document.getElementById("email").value.trim();
      const contactInfo = document.getElementById("contactInfo").value.trim();
      const password = document.getElementById("password").value;
      const confirmPassword = document.getElementById("confirmPassword").value;
      const signupMessage = document.getElementById("signupMessage");

      if (password !== confirmPassword) {
        signupMessage.textContent = "Passwords do not match.";
        signupMessage.className = "form-message error";
        return;
      }

      try {
        const response = await fetch("/api/signup", {
          method: "POST",
          headers: {
            "Content-Type": "application/json"
          },
          body: JSON.stringify({
            fullName,
            email,
            contactInfo,
            password
          })
        });

        const result = await response.json();

        if (!result.success) {
          signupMessage.textContent = result.message;
          signupMessage.className = "form-message error";
          return;
        }

        signupMessage.textContent = "Account created successfully. Redirecting to login...";
        signupMessage.className = "form-message success";

        signupForm.reset();

        setTimeout(() => {
          window.location.href = "login.html";
        }, 1200);
      } catch (error) {
        signupMessage.textContent = "Something went wrong. Please try again.";
        signupMessage.className = "form-message error";
      }
    });
  }

  if (loginForm) {
    loginForm.addEventListener("submit", async function (event) {
      event.preventDefault();

      const email = document.getElementById("loginEmail").value.trim();
      const password = document.getElementById("loginPassword").value;
      const loginMessage = document.getElementById("loginMessage");

      try {
        const response = await fetch("/api/login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json"
          },
          body: JSON.stringify({
            email,
            password
          })
        });

        const result = await response.json();

        if (!result.success) {
          loginMessage.textContent = result.message;
          loginMessage.className = "form-message error";
          return;
        }

        localStorage.setItem("isLoggedIn", "true");
        localStorage.setItem("currentUser", JSON.stringify(result.user));
        localStorage.setItem("tshongmartUser", JSON.stringify(result.user));

        loginMessage.textContent = "Login successful. Redirecting to dashboard...";
        loginMessage.className = "form-message success";

        setTimeout(() => {
          window.location.href = "dashboard.html";
        }, 1200);
      } catch (error) {
        loginMessage.textContent = "Something went wrong. Please try again.";
        loginMessage.className = "form-message error";
      }
    });
  }
})();