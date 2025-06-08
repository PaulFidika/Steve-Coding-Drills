const { nextui } = require("@nextui-org/react");

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
    "./node_modules/@nextui-org/theme/dist/**/*.{js,ts,jsx,tsx}"
  ],
  theme: {
    extend: {},
  },
  darkMode: "class",
  plugins: [
    nextui({
      themes: {
        dark: {
          colors: {
            primary: {
              DEFAULT: "#4ade80",
              foreground: "#000000",
            },
            secondary: {
              DEFAULT: "#9333ea",
              foreground: "#ffffff",
            },
            background: "#1a1a1a",
            foreground: "#ffffff",
          },
          layout: {
            radius: {
              small: "4px",
              medium: "6px",
              large: "8px",
            },
            borderWidth: {
              small: "1px",
              medium: "2px",
              large: "3px",
            },
          },
          fonts: {
            sans: "'Inter', sans-serif",
          },
        },
      },
    })
  ],
}