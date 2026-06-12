function init() {
  document.querySelectorAll("obfuscated").forEach((el) => {
    const text = el.textContent || "";
    const length = text.length;
  });
}

init();
