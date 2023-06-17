import { useState, useEffect, useCallback } from "react";
import { Moon, Sun } from "lucide-react";
import { getThemeFromLocalStorage, setThemeToLocalStorage } from "@/utils/theme";

const ThemeToggle = () => {
  const [isDarkMode, setIsDarkMode] = useState(getThemeFromLocalStorage() === "dark");

  useEffect(() => {
    const systemTheme = window.matchMedia("(prefers-color-scheme: dark)").matches;
    const storedTheme = localStorage.getItem("theme");
    const initialTheme = storedTheme !== null ? storedTheme === "dark" : systemTheme;
    setIsDarkMode(initialTheme);
    if (initialTheme) {
      document.documentElement.classList.add("dark");
    }
  }, []);

  useEffect(() => {
    if (isDarkMode) {
      document.documentElement.classList.add("dark");
    } else {
      document.documentElement.classList.remove("dark");
    }
  }, [isDarkMode]);

  const toggleTheme = useCallback(() => {
    setThemeToLocalStorage(isDarkMode ? "dark" : "light");
    setIsDarkMode((prevState) => !prevState);
  }, [isDarkMode]);

  return (
    <div className="flex items-center">
      <button className={`relative w-12 h-6 bg-secondary rounded-full p-1 transition-transform flex items-center`} onClick={toggleTheme}>
        <div
          className={`absolute ${isDarkMode ? "right-0" : "left-0"} bg-white dark:bg-black w-5 h-5 rounded-full transition-transform ${
            isDarkMode ? "-translate-x-0.5" : "translate-x-0.5"
          }`}
        />
        {isDarkMode ? <Moon size={16} color="white" /> : <Sun size={16} color="black" className="ml-auto" />}
      </button>
    </div>
  );
};

export default ThemeToggle;
