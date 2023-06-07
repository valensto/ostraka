type Theme = "light" | "dark"

export function getThemeFromLocalStorage(): Theme {
  const storedTheme = localStorage.getItem("theme");
  return storedTheme === "dark" ? "dark" : "light";
}

export function setThemeToLocalStorage(theme: Theme) {
  localStorage.setItem("theme", theme);
}
