console.log("Hello from Vite + Go! (main.js)");

// Example: Add some dynamic content to the page
document.addEventListener('DOMContentLoaded', () => {
    const appDiv = document.getElementById('app');
    if (appDiv) {
        const heading = document.createElement('h2');
        heading.textContent = 'Content added by main.js';
        appDiv.appendChild(heading);
    } else {
        console.error('#app element not found!');
    }
});