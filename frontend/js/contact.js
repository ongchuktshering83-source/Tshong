(() => {
  const contactForm = document.getElementById("contactForm");
  const contactMessageStatus = document.getElementById("contactMessageStatus");

  if (!contactForm) return;

  contactForm.addEventListener("submit", async function (event) {
    event.preventDefault();

    const fullName = document.getElementById("contactName").value.trim();
    const email = document.getElementById("contactEmail").value.trim();
    const subject = document.getElementById("contactSubject").value;
    const message = document.getElementById("contactMessage").value.trim();

    try {
      const response = await fetch("/api/contact-messages", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          fullName,
          email,
          subject,
          message
        })
      });

      const result = await response.json();

      if (!result.success) {
        contactMessageStatus.textContent = result.message;
        contactMessageStatus.className = "form-message error";
        return;
      }

      contactMessageStatus.textContent =
        "Message sent successfully. Our team will get back to you soon.";
      contactMessageStatus.className = "form-message success";

      contactForm.reset();
    } catch (error) {
      contactMessageStatus.textContent = "Something went wrong. Please try again.";
      contactMessageStatus.className = "form-message error";
    }
  });
})();