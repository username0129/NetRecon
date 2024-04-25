/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      backgroundImage: {
        'login-pattern': 'url(\'/assets/login.svg\')'
      },
      backgroundColor: {
        'main': '#F5F5F5'
      }
    }
  },
  plugins: []
}
