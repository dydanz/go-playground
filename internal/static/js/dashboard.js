document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM fully loaded and parsed');
    fetchUserData();
    setupLogout();
});

async function fetchUserData() {
    console.log('Fetching user data...');
    try {
        const userId = getCookie('user_id');
        console.log('User ID from cookie:', userId);
        if (!userId) {
            showMessage('User ID not found', 'error');
            setTimeout(() => {
                window.location.href = '/login';
            }, 10000);
            return;
        }

        const response = await fetchWithCSRF(`/api/users/${userId}`, {
            headers: {
                'Authorization': `Bearer ${getCookie('session_token')}`
            }
        });

        console.log('User API response:', response);
        const userData = await response.json();
        if (!userData) {
            showMessage('User Data not found', 'error');
            setTimeout(() => {
                window.location.href = '/login';
            }, 5000);
            return;
        }

        if (response.ok) {
            console.log('User Data JSON:', userData);
            const userNameElement = document.getElementById('userName');
            if (userNameElement) {
                userNameElement.textContent = userData.name || 'Guest';
            }
        } else {
            if (response.status === 401) {
                showMessage('Session expired. Please login again.', 'error');
                setTimeout(() => {
                    window.location.href = '/login';
                }, 2000);
            }
        }
    } catch (error) {
        console.error('Error fetching user data:', error);
        showMessage('Error loading user data', 'error');
    }
}

function setupLogout() {
    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', async function(e) {
            e.preventDefault();
            try {
                const response = await fetchWithCSRF('/api/auth/logout', {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${getCSRFToken()}`
                    }
                });

                if (response.ok) {
                    const result = await response.json();
                    showMessage(result.message, 'success');
                    setTimeout(() => {
                        window.location.href = '/login';
                    }, 2000);
                } else {
                    const error = await response.json();
                    showMessage(error.error || 'Logout failed', 'error');
                }
            } catch (error) {
                console.error('Logout error:', error);
                showMessage('Error during logout', 'error');
            }
        });
    }
}

function showMessage(message, type = 'info') {
    // Create message element
    const messageDiv = document.createElement('div');
    messageDiv.className = `message-popup ${type}`;
    messageDiv.textContent = message;

    // Add to document
    document.body.appendChild(messageDiv);

    // Animate in
    setTimeout(() => {
        messageDiv.classList.add('show');
    }, 100);

    // Remove after delay
    setTimeout(() => {
        messageDiv.classList.remove('show');
        setTimeout(() => {
            messageDiv.remove();
        }, 300);
    }, 3000);
}

// Add this function to utils.js
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
} 