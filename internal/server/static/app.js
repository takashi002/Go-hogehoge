document.addEventListener("DOMContentLoaded", () => {
  const waves = document.querySelector(".ocean__waves");
  if (!waves) return;

  let ticking = false;
  window.addEventListener("scroll", () => {
    if (ticking) return;
    ticking = true;
    requestAnimationFrame(() => {
      const y = Math.min(window.scrollY, 240);
      waves.style.transform = `translateY(${y * 0.12}px)`;
      ticking = false;
    });
  });
});
