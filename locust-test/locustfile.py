from locust import HttpUser, task
import json
import random
import string

class HelloWorldUser(HttpUser):
    host = "http://localhost:8080"

    @task(1)
    def ping_connection_status(self):
        self.client.get("/ping")

class RegisterUser(HttpUser):
    host = "http://localhost:8080"

    # TODO: add a method to generate random email, name, phone, password 
    # using https://faker.readthedocs.io/en/master/index.html
    def generate_random_email(self):
        # Generate a random string for the email prefix
        prefix = ''.join(random.choices(string.ascii_lowercase, k=8))
        return f"{prefix}@example.com"

    def generate_random_name(self):
        return ''.join(random.choices(string.ascii_letters, k=8))

    def generate_random_phone(self):
        return ''.join(random.choices(string.digits, k=10))   
     
    @task(1)
    def registration_flow(self):
        # Generate random user data
        email = self.generate_random_email()
        name = self.generate_random_name()
        phone = self.generate_random_phone()
        password = "12345678"

        # Step 1: Register new user
        register_data = {
            "email": email,
            "name": name,
            "password": password,
            "phone": phone
        }
        register_response = self.client.post(
            "/api/auth/register",
            json=register_data,
            headers={"Content-Type": "application/json"}
        )

        # Check if registration was successful
        if register_response.status_code != 201:
            return

        # Step 2: Get verification code
        verify_code_response = self.client.get(
            f"/api/auth/test/get-verification/code?email={email}",
            headers={"accept": "application/json"}
        )

        # Check if we got the verification code
        if verify_code_response.status_code != 200:
            return

        # Extract OTP from response
        otp = verify_code_response.json().get("otp")
        if not otp:
            return

        # Step 3: Verify registration
        verify_data = {
            "email": email,
            "otp": otp
        }
        verify_response = self.client.post(
            "/api/auth/verify",
            json=verify_data,
            headers={"Content-Type": "application/json"}
        )

        # Check if verification was successful
        if verify_response.status_code!= 200:
            return

class MultipleLoginUser(HttpUser):
    host = "http://localhost:8080"

    # TODO: add a method to generate random email, name, phone, password 
    # using https://faker.readthedocs.io/en/master/index.html
    def generate_random_email(self):
        # Generate a random string for the email prefix
        prefix = ''.join(random.choices(string.ascii_lowercase, k=8))
        return f"{prefix}@example.com"

    def generate_random_name(self):
        return ''.join(random.choices(string.ascii_letters, k=8))

    def generate_random_phone(self):
        return ''.join(random.choices(string.digits, k=10))
    
    @task(1)
    def login_and_profile_flow(self):
        # Use fixed test credentials for login
        login_data = {
            "email": "john@doe.com",
            "password": "12345678"
        }

        # Step 1: Login
        login_response = self.client.post(
            "/api/auth/login",
            json=login_data,
            headers={"Content-Type": "application/json"}
        )

        # Check if login was successful
        if login_response.status_code != 200:
            return

        # Extract token and user_id from response
        login_result = login_response.json()
        token = login_result.get("token")
        user_id = login_result.get("user_id")

        if not token or not user_id:
            return

        # Step 2: Get user profile
        profile_response = self.client.get(
            f"/api/users/{user_id}",
            headers={
                "accept": "application/json",
                "Authorization": f"Bearer {token}",
                "X-User-Id": user_id
            }
        )

        # Check if profile fetch was successful
        if profile_response.status_code != 200:
            return

class LoginAndUpdatePointsUser(HttpUser):
    host = "http://localhost:8080"

    # TODO: add a method to generate random email, name, phone, password 
    # using https://faker.readthedocs.io/en/master/index.html
    def generate_random_email(self):
        # Generate a random string for the email prefix
        prefix = ''.join(random.choices(string.ascii_lowercase, k=8))
        return f"{prefix}@example.com"

    def generate_random_name(self):
        return ''.join(random.choices(string.ascii_letters, k=8))

    def generate_random_phone(self):
        return ''.join(random.choices(string.digits, k=10))
    
    @task(1)
    def random_user_login_and_points_flow(self):
        # Step 1: Get random active user
        random_user_response = self.client.get(
            "/api/auth/test/random-user",
            headers={"accept": "application/json"}
        )

        # Check if we got a random user
        if random_user_response.status_code != 200:
            return

        # Extract credentials
        user_creds = random_user_response.json()
        email = user_creds.get("email")
        password = user_creds.get("password")

        if not email or not password:
            return

        # Step 2: Login with random user
        login_data = {
            "email": email,
            "password": '12345678'
        }
        login_response = self.client.post(
            "/api/auth/login",
            json=login_data,
            headers={"Content-Type": "application/json"}
        )
        
        # Check if login was successful
        if login_response.status_code != 200:
            return

        # Extract token and user_id from response
        login_result = login_response.json()
        token = login_result.get("token")
        user_id = login_result.get("user_id")

        if not token or not user_id:
            return

        # Step 3: Add one point to user's balance
        points_data = 1  # Add one point
        points_response = self.client.put(
            f"/api/points/{user_id}",
            json=points_data,
            headers={
                "Content-Type": "application/json",
                "Authorization": f"Bearer {token}",
                "X-User-Id": user_id
            }
        )

        # Check if points update was successful
        if points_response.status_code != 200:
            return
