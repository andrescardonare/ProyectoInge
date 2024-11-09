/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./cmd/main.go",
        "./templates/**/*.{html,js,templ,go}",
        "./templates/index.templ",
        "./templates/login.templ",
        "./templates/*.{html,js,templ,go}"
    ],
    theme: {
        extend: {},
        fontFamily: {
            sans: ["Quicksand"],
        },
    },
    plugins: [
        require("@tailwindcss/forms"),
        require("@tailwindcss/typography"),
        require('@tailwindcss/aspect-ratio'),
        require('@tailwindcss/container-queries'),
        require("autoprefixer")
    ],
};