/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "../templates/*.html",
    "../templates/**/*.html", // Scan semua file HTML kamu
    "./**/*.js", // Jika pakai JS di static folder
  ],
  darkMode: "class", // Aktifkan dark mode manual via class
  theme: {
    extend: {},
  },
  plugins: [],
};
