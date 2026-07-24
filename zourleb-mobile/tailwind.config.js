/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./app/**/*.{ts,tsx}", "./src/**/*.{ts,tsx}"],
  presets: [require("nativewind/preset")],
  theme: {
    extend: {
      colors: {
        brand: {
          50: "#eef7f3",
          100: "#d6ebe1",
          500: "#1f8a5b",
          600: "#176b47",
          700: "#125539",
        },
      },
    },
  },
  plugins: [],
};
