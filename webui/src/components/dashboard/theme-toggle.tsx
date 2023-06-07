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
    setIsDarkMode(prevState => !prevState);
  }, [isDarkMode]);

  return (
    <div className="flex items-center">
      <button
        className={`relative w-12 h-6 bg-primary rounded-full p-1 transition-transform`}
        onClick={toggleTheme}
      >
        <div
          className={`absolute top-1 ${
            isDarkMode ? 'right-5' : 'left-1'
          } bg-white dark:bg-black w-4 h-4 rounded-full transition-transform ${
            isDarkMode ? 'translate-x-full' : 'translate-x-0'
          }`}
        />
        {isDarkMode ? (
          <Moon size={16} color="black" />
        ) : (
          <Sun size={16} color="white" className="ml-auto" />
        )}
      </button>
    </div>
  );
};

export default ThemeToggle;
