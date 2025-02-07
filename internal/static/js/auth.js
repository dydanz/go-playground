// Constants
const API_BASE_URL = 'http://localhost:8080/api';

// Utility function to handle API errors
const handleApiError = (error) => {
    console.error('API Error:', error);
    alert(error.message || 'An error occurred. Please try again.');
};

// Function to set secure cookies
const setCookie = (name, value, days = 7) => {
    const expires = new Date(Date.now() + days * 864e5).toUTCString();
    document.cookie = `${name}=${encodeURIComponent(value)}; expires=${expires}; path=/; secure; samesite=strict`;
};

// Function to get cookie value
const getCookie = (name) => {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return decodeURIComponent(parts.pop().split(';').shift());
};

// Function to handle registration
const handleRegister = async (event) => {
    event.preventDefault();
    const form = event.target;
    const formData = new FormData(form);

    try {
        const response = await fetch(`${API_BASE_URL}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({
                email: formData.get('email'),
                name: formData.get('name'),
                password: formData.get('password'),
                phone: formData.get('phone')
            })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Registration failed');
        }

        alert('Registration successful! Please login.');
        window.location.href = '/login';
    } catch (error) {
        handleApiError(error);
    }
};

// Function to handle login
const handleLogin = async (event) => {
    event.preventDefault();
    const form = event.target;
    const formData = new FormData(form);

    try {
        const response = await fetch(`${API_BASE_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({
                email: formData.get('email'),
                password: formData.get('password')
            })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Login failed');
        }

        const data = await response.json();
        
        // Store auth token and user ID in secure cookies
        setCookie('auth_token', data.token);
        setCookie('user_id', data.user_id);

        window.location.href = '/dashboard';
    } catch (error) {
        handleApiError(error);
    }
};

// Function to load user profile
const loadUserProfile = async () => {
    const authToken = getCookie('auth_token');
    if (!authToken) {
        window.location.href = '/login';
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/users/me`, {
            headers: {
                'Accept': 'application/json',
                'Authorization': authToken
            }
        });

        if (!response.ok) {
            throw new Error('Failed to load user profile');
        }

        const userData = await response.json();
        
        // Update UI with user data
        const userNameElement = document.getElementById('userName');
        if (userNameElement) {
            userNameElement.textContent = userData.name;
        }

        return userData;
    } catch (error) {
        handleApiError(error);
        window.location.href = '/login';
    }
};

// Function to handle logout
const handleLogout = () => {
    // Clear auth cookies
    document.cookie = 'auth_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    document.cookie = 'user_id=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    window.location.href = '/login';
};

// Add event listeners when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    // Register form handler
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', handleRegister);
    }

    // Login form handler
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }

    // Logout button handler
    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', (e) => {
            e.preventDefault();
            handleLogout();
        });
    }

    // Load user profile if on dashboard
    if (window.location.pathname === '/dashboard') {
        loadUserProfile();
    }

    // Load user profile if on profile page
    if (window.location.pathname === '/profile') {
        displayUserProfile();
    }
});

// Password visibility toggle function
const togglePassword = (inputId) => {
    const passwordInput = document.getElementById(inputId);
    const icon = passwordInput.parentElement.querySelector('.password-toggle i');
    
    if (passwordInput.type === 'password') {
        passwordInput.type = 'text';
        icon.classList.remove('fa-eye');
        icon.classList.add('fa-eye-slash');
    } else {
        passwordInput.type = 'password';
        icon.classList.remove('fa-eye-slash');
        icon.classList.add('fa-eye');
    }
}; 