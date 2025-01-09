function checkCaptchaOnSubmit(event) {
  const captcha = document.getElementById("captcha");
  if (!captcha.checked) {
    alert("Você precisa fazer o CAPTCHA!")
    return false;
  }
}
