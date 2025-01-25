/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        "primary": "#0d6efd",
        "success": "#32a852",
        "error": "#ff453a",

        "btn": "#484a4d",
        "btn-hover": "#525357",

        "subtle": "#333338",
        "sidebar": "#1E293C",
        "background": "#161E2B",
      },
      fontFamily: {
        "avant": "ITC Avant Garde Std Md, sans-serif",
        "now": "Helvetica Now Text, Helvetica, sans-serif",
      },
      height: {
        "112": "28rem",
        "128": "32rem",
        "144": "36rem",
        "160": "40rem",
        "176": "44rem",
        "192": "48rem",
      },
      borderWidth: {
        "1": "1px"
      }
    },
  },
  plugins: [],
}

