// Function to get CSRF token from cookie
function getCSRFToken() {
    return getCookie('csrf_token');
}

// Update fetch calls to include CSRF token
async function fetchWithCSRF(url, options = {}) {
    const csrfToken = getCSRFToken();
    const userId = getCookie('user_id');
    
    const headers = {
        ...options.headers,
        'X-CSRF-Token': csrfToken,
        'X-User-Id': userId || '' // Ensure userId is always sent even if null
    };

    return fetch(url, { ...options, headers });
}

// Add this function to check authentication
function checkAuth() {
    const userId = getCookie('user_id');
    const sessionToken = getCookie('session_token');
    
    if (!userId || !sessionToken) {
        window.location.href = '/login';
        return false;
    }
    return true;
}

// Add getCookie function
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
}

// Add these utility functions
function setCookie(name, value, days = 1) {
    const date = new Date();
    date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
    const expires = `expires=${date.toUTCString()}`;
    document.cookie = `${name}=${value};${expires};path=/;secure`;
}

function clearAuthCookies() {
    setCookie('session_token', '', -1);
    setCookie('user_id', '', -1);
    setCookie('csrf_token', '', -1);
}