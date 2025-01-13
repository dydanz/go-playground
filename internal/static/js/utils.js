// Function to get CSRF token from cookie
function getCSRFToken() {
    return getCookie('csrf_token');
}

// Update fetch calls to include CSRF token
async function fetchWithCSRF(url, options = {}) {
    const csrfToken = getCSRFToken();
    const headers = {
        ...options.headers,
        'X-CSRF-Token': csrfToken
    };

    return fetch(url, { ...options, headers });
}

// Add getCookie function
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
} 