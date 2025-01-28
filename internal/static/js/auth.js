document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');

    if (loginForm) {
        loginForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const formData = new FormData(loginForm);
            const data = {
                email: formData.get('email'),
                password: formData.get('password')
            };

            console.log('Sending login request with data:', data);

            try {
                const response = await fetch('/api/auth/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(data)
                });

                const result = await response.json();
                console.log('Login response:', result);
                
                if (response.ok) {
                    localStorage.setItem('token', result.token);
                    window.location.href = '/dashboard';
                } else {
                    alert(result.error || 'Login failed');
                }
            } catch (error) {
                console.error('Login error:', error);
                alert('Error: ' + error.message);
            }
        });
    }

    if (registerForm) {
        registerForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const formData = new FormData(registerForm);
            const data = {
                name: formData.get('name'),
                email: formData.get('email'),
                password: formData.get('password'),
                phone: formData.get('phone')
            };

            try {
                const response = await fetch('/api/auth/register', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(data)
                });

                const result = await response.json();
                
                if (response.ok) {
                    alert('Registration successful! Please check your email for verification.');
                    window.location.href = '/login';
                } else {
                    alert(result.error || 'Registration failed');
                }
            } catch (error) {
                alert('Error: ' + error.message);
            }
        });
    }
});

async function handleLogin(event) {
    event.preventDefault();
    const form = event.target;
    const formData = new FormData(form);
    
    try {
        const response = await fetch('/api/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(Object.fromEntries(formData)),
        });

        const data = await response.json();
        
        if (response.ok) {
            // Store both session token and user ID in cookies
            setCookie('session_token', data.token);
            setCookie('user_id', data.user_id); // Make sure your backend sends user_id in response
            
            showMessage('Login successful!', 'success');
            setTimeout(() => {
                window.location.href = '/dashboard';
            }, 1000);
        } else {
            showMessage(data.error || 'Login failed', 'error');
        }
    } catch (error) {
        console.error('Login error:', error);
        showMessage('Error during login', 'error');
    }
}

async function logout() {
    const token = localStorage.getItem('token');
    const userId = localStorage.getItem('user_id');
    try {
        const response = await fetch('/api/auth/logout', {
            method: 'POST',
            headers: {
                'Authorization': token,
                'X-User-Id': userId
            }
        });

        if (response.ok) {
            localStorage.removeItem('token');
            window.location.href = '/login';
        } else {
            const result = await response.json();
            alert(result.error || 'Logout failed');
        }
    } catch (error) {
        console.error('Logout error:', error);
        alert('Error during logout');
    }
}

// Update the logout button handler
const logoutBtn = document.getElementById('logoutBtn');
if (logoutBtn) {
    logoutBtn.addEventListener('click', function(e) {
        e.preventDefault();
        logout();
    });
}