export function getThemeFromLocalStorage() {
  const storedTheme = localStorage.getItem("theme");
  return storedTheme === "dark";
}

export function setThemeToLocalStorage(theme: boolean) {
  const newTheme = theme ? "dark" : "light";
  localStorage.setItem("theme", newTheme);
}
