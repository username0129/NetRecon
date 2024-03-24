/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./src/**/*.{html,vue,js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            backgroundImage: {
                'login-pattern': "url('/assets/login.svg')",
            }
        },
    },
    plugins: [],
}