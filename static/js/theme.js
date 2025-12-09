// THEME TOGGLE SCRIPT
const themeToggleBtnDesktop = document.getElementById('theme-toggle-desktop');
const themeToggleDarkIconDesktop = document.getElementById('theme-toggle-dark-icon-desktop');
const themeToggleLightIconDesktop = document.getElementById('theme-toggle-light-icon-desktop');

const themeToggleBtnMobile = document.getElementById('theme-toggle-mobile');
const themeToggleDarkIconMobile = document.getElementById('theme-toggle-dark-icon-mobile');
const themeToggleLightIconMobile = document.getElementById('theme-toggle-light-icon-mobile');

function applyTheme(theme) {
    if (theme === 'dark') {
        document.documentElement.classList.add('dark');
        // Show light icon, hide dark icon for both buttons
        if (themeToggleLightIconDesktop) themeToggleLightIconDesktop.classList.remove('hidden');
        if (themeToggleDarkIconDesktop) themeToggleDarkIconDesktop.classList.add('hidden');
        if (themeToggleLightIconMobile) themeToggleLightIconMobile.classList.remove('hidden');
        if (themeToggleDarkIconMobile) themeToggleDarkIconMobile.classList.add('hidden');
    } else {
        document.documentElement.classList.remove('dark');
        // Show dark icon, hide light icon for both buttons
        if (themeToggleDarkIconDesktop) themeToggleDarkIconDesktop.classList.remove('hidden');
        if (themeToggleLightIconDesktop) themeToggleLightIconDesktop.classList.add('hidden');
        if (themeToggleDarkIconMobile) themeToggleDarkIconMobile.classList.remove('hidden');
        if (themeToggleLightIconMobile) themeToggleLightIconMobile.classList.add('hidden');
    }
}

// Initial theme application
const storedTheme = localStorage.getItem('color-theme');
if (storedTheme) {
    applyTheme(storedTheme);
} else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
    applyTheme('dark');
} else {
    applyTheme('light');
}

// Add event listeners to both buttons
if (themeToggleBtnDesktop) {
    themeToggleBtnDesktop.addEventListener('click', function() {
        let currentTheme = localStorage.getItem('color-theme');
        if (!currentTheme) {
            currentTheme = document.documentElement.classList.contains('dark') ? 'dark' : 'light';
        }
        const newTheme = currentTheme === 'light' ? 'dark' : 'light';
        localStorage.setItem('color-theme', newTheme);
        applyTheme(newTheme);
    });
}
if (themeToggleBtnMobile) {
    themeToggleBtnMobile.addEventListener('click', function() {
        let currentTheme = localStorage.getItem('color-theme');
        if (!currentTheme) {
            currentTheme = document.documentElement.classList.contains('dark') ? 'dark' : 'light';
        }
        const newTheme = currentTheme === 'light' ? 'dark' : 'light';
        localStorage.setItem('color-theme', newTheme);
        applyTheme(newTheme);
    });
}

// MOBILE MENU TOGGLE SCRIPT (assuming this is also needed globally)
const mobileMenuButton = document.getElementById('mobile-menu-button');
const mobileMenu = document.getElementById('mobile-menu');

if (mobileMenuButton && mobileMenu) {
    mobileMenuButton.addEventListener('click', () => {
        const isExpanded = mobileMenuButton.getAttribute('aria-expanded') === 'true';
        mobileMenuButton.setAttribute('aria-expanded', !isExpanded);
        mobileMenu.classList.toggle('hidden');
        // Toggle hamburger icons
        mobileMenuButton.querySelectorAll('svg').forEach(icon => icon.classList.toggle('hidden'));
    });
}
