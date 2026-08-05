/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,jsx}",
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          gold: '#f59e0b',
          'gold-hover': '#d97706',
          red: '#dc2626',
          navy: '#0f172a',
          surface: '#1e293b',
          border: '#334155',
        },
        seat: {
          available: '#10b981',
          premium: '#3b82f6',
          vip: '#8b5cf6',
          selected: '#f59e0b',
          booked: '#475569',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        display: ['Playfair Display', 'serif'],
      },
    },
  },
  plugins: [],
}
