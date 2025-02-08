// Auth handling functions
const auth = {
    // Register new user
    async register(userData) {
        try {
            const response = await fetch('http://localhost:8080/api/auth/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'accept': 'application/json'
                },
                body: JSON.stringify(userData)
            });

            const data = await response.json();
            
            if (response.ok) {
                // Show OTP modal on successful registration
                this.showOTPModal(userData.email);
                return data;
            } else {
                throw new Error(data.message || 'Registration failed');
            }
        } catch (error) {
            console.error('Registration error:', error);
            throw error;
        }
    },

    // Verify OTP
    async verifyOTP(email, otp) {
        try {
            const response = await fetch('http://localhost:8080/api/auth/verify', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'accept': 'application/json'
                },
                body: JSON.stringify({
                    email: email,
                    otp: otp
                })
            });

            const data = await response.json();
            
            if (response.ok) {
                // Redirect to sign-in page on successful verification
                window.location.href = '/sign-in';
                return data;
            } else {
                throw new Error(data.message || 'OTP verification failed');
            }
        } catch (error) {
            console.error('OTP verification error:', error);
            throw error;
        }
    },

    // Show OTP Modal
    showOTPModal(email) {
        // Create modal HTML
        const modalHTML = `
        <div class="modal fade" id="otpModal" tabindex="-1" role="dialog" aria-labelledby="modal-form" aria-hidden="true">
            <div class="modal-dialog modal-dialog-centered modal-md" role="document">
                <div class="modal-content">
                    <div class="modal-body p-0">
                        <div class="card card-plain">
                            <div class="card-header pb-0 text-left">
                                <h3 class="font-weight-bolder text-info text-gradient">Verify Your Email</h3>
                                <p class="mb-0">Enter the OTP sent to your email</p>
                            </div>
                            <div class="card-body">
                                <form id="otpForm" role="form text-left">
                                    <label>OTP Code</label>
                                    <div class="input-group mb-3">
                                        <input type="text" id="otpInput" class="form-control" placeholder="Enter OTP" aria-label="OTP" aria-describedby="otp-addon">
                                    </div>
                                    <div class="text-center">
                                        <button type="submit" class="btn btn-round bg-gradient-info btn-lg w-100 mt-4 mb-0">Verify OTP</button>
                                    </div>
                                </form>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>`;

        // Add modal to body
        document.body.insertAdjacentHTML('beforeend', modalHTML);

        // Get modal element
        const modalElement = document.getElementById('otpModal');
        const modal = new bootstrap.Modal(modalElement);

        // Show modal
        modal.show();

        // Handle form submission
        const form = document.getElementById('otpForm');
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            const otp = document.getElementById('otpInput').value;
            
            try {
                await this.verifyOTP(email, otp);
                modal.hide();
                // Success notification
                $.notify({
                    icon: "check",
                    message: "OTP verified successfully! Redirecting to login..."
                }, {
                    type: 'success',
                    timer: 3000,
                    placement: {
                        from: 'top',
                        align: 'right'
                    }
                });
            } catch (error) {
                // Error notification
                $.notify({
                    icon: "error",
                    message: error.message
                }, {
                    type: 'danger',
                    timer: 3000,
                    placement: {
                        from: 'top',
                        align: 'right'
                    }
                });
            }
        });
    },

    // Login function
    async login(credentials) {
        try {
            const response = await fetch('http://localhost:8080/api/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'accept': 'application/json'
                },
                body: JSON.stringify(credentials)
            });

            const data = await response.json();
            
            if (response.ok) {
                // Store session data in cookies (these will be handled by the server)
                // The server will set secure cookies via Set-Cookie header
                
                // Redirect to dashboard
                window.location.href = '/dashboard';
                return data;
            } else {
                throw new Error(data.message || 'Login failed');
            }
        } catch (error) {
            console.error('Login error:', error);
            throw error;
        }
    },

    // Add new logout method
    async logout() {
        try {
            // Get tokens from cookies
            const userID = getCookie('user_id');
            const sessionToken = getCookie('session_token');

            const response = await fetch('http://localhost:8080/api/auth/logout', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'accept': 'application/json',
                    'User-ID': userID,
                    'Authorization': `Bearer ${sessionToken}`
                }
            });

            if (response.ok) {
                // Clear cookies
                document.cookie = 'session_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
                document.cookie = 'user_id=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
                document.cookie = 'user_name=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
                
                // Redirect to login page
                window.location.href = '/sign-in';
            } else {
                throw new Error('Logout failed');
            }
        } catch (error) {
            console.error('Logout error:', error);
            throw error;
        }
    },
};

// Add helper function to get cookie value
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

// Handle registration form submission
document.addEventListener('DOMContentLoaded', () => {
    const registrationForm = document.querySelector('form[role="form"]');
    if (registrationForm) {
        registrationForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const userData = {
                name: registrationForm.querySelector('input[placeholder="Name"]').value,
                email: registrationForm.querySelector('input[placeholder="Email"]').value,
                password: registrationForm.querySelector('input[placeholder="Password"]').value,
                phone: registrationForm.querySelector('input[placeholder="Phone"]').value,
            };

            try {
                await auth.register(userData);
                // Success notification
                $.notify({
                    icon: "check",
                    message: "Registration successful! Please check your email for OTP."
                }, {
                    type: 'success',
                    timer: 3000,
                    placement: {
                        from: 'top',
                        align: 'right'
                    }
                });
            } catch (error) {
                // Error notification
                $.notify({
                    icon: "error",
                    message: error.message
                }, {
                    type: 'danger',
                    timer: 3000,
                    placement: {
                        from: 'top',
                        align: 'right'
                    }
                });
            }
        });
    }
      
}); 