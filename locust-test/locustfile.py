from locust import HttpUser, task, events
import json
import secrets
import string


host = "http://localhost:8080"

# TODO: add a method to generate random email, name, phone, password 
# using https://faker.readthedocs.io/en/master/index.html
def generate_random_email():
    # Generate a random string for the email prefix
    prefix = ''.join(secrets.choice(string.ascii_lowercase) for _ in range(8))

    return f"{prefix}@example.com"

def generate_random_name():
    return ''.join(secrets.choice(string.ascii_lowercase) for _ in range(8))

def generate_random_phone():
    return ''.join(secrets.choice(string.digits) for _ in range(10))   

@events.test_start.add_listener
def on_test_start(environment, **kwargs):
    print("A new test is starting")

@events.test_stop.add_listener
def on_test_stop(environment, **kwargs):
    print("A new test is ending")

class HelloWorldUser(HttpUser):
    host = "http://localhost:8080"

    @task(1)
    def ping_connection_status(self):
        self.client.get("/ping")
     
class RegisterUser(HttpUser):

    @task(1)
    def registration_flow(self):
        # Generate random user data
        email = generate_random_email()
        name = generate_random_name()
        phone = generate_random_phone()
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

class MerchantOnboarding(HttpUser):
    @task(1)
    def merchant_onboarding_flow(self):
        # Step 0: Login to get authentication token
        login_data = {
            "email": "john@doe.com",  # Replace with admin credentials
            "password": "12345678"
        }
        login_response = self.client.post(
            "/api/auth/login",
            json=login_data,
            headers={"Content-Type": "application/json"}
        )

        if login_response.status_code != 200:
            return

        # Extract token and user_id
        login_result = login_response.json()
        token = 'Bearer '+login_result.get("token")
        user_id = login_result.get("user_id")

        if not token or not user_id:
            return

        auth_headers = {
            "Content-Type": "application/json",
            "Authorization": token,
            "X-User-Id": user_id
        }

        # Step 1: Create a merchant (bank)
        merchant_data = {
            "merchant_name": "Sample Bank",
            "merchant_type": "bank"
        }
        merchant_response = self.client.post(
            "/api/merchants",
            json=merchant_data,
            headers=auth_headers
        )

        if merchant_response.status_code != 201:
            return

        merchant_id = merchant_response.json().get("merchant_id")
        if not merchant_id:
            return

        # Step 2: Create loyalty program
        program_data = {
            "merchant_id": merchant_id,
            "program_name": "Bank Rewards",
            "point_currency_name": "Reward Points"
        }
        program_response = self.client.post(
            "/api/programs",
            json=program_data,
            headers=auth_headers
        )

        if program_response.status_code != 201:
            return

        program_id = program_response.json().get("id")
        if not program_id:
            return

        # Step 3: Create program rules
        rules = [
            {
                "program_id": program_id,
                "rule_name": "Base Points",
                "condition_type": "transaction_amount",
                "condition_value": ">0",
                "multiplier": 1.0,
                "points_awarded": 123,
                "effective_from": "2023-01-01T00:00:00Z",
                "effective_to": "2023-12-01T00:00:00Z"
            },
            {
                "program_id": program_id,
                "rule_name": "Dining Bonus",
                "condition_type": "item_category",
                "condition_value": "Dining",
                "multiplier": 5.0,
                "points_awarded": 123,
                "effective_from": "2023-01-01T00:00:00Z",
                "effective_to": "2023-12-01T00:00:00Z"
            },
            {
                "program_id": program_id,
                "rule_name": "Travel Bonus",
                "condition_type": "item_category",
                "condition_value": "Travel",
                "multiplier": 3.0,
                "points_awarded": 123,
                "effective_from": "2023-01-01T00:00:00Z",
                "effective_to": "2023-12-01T00:00:00Z"
            }
        ]

        for rule in rules:
            rule_response = self.client.post(
                "/api/program-rules",
                json=rule,
                headers=auth_headers
            )
            if rule_response.status_code != 201:
                return

        # Step 4: Process a dining transaction
        transaction_data = {
            "merchant_id": merchant_id,
            "customer_id": user_id,  # Using the logged-in user as customer
            "transaction_type": "purchase",
            "transaction_amount": 100.00,
            "status": "completed"
        }
        transaction_response = self.client.post(
            "/api/transactions",
            json=transaction_data,
            headers=auth_headers
        )

        if transaction_response.status_code != 201:
            return

        transaction_id = transaction_response.json().get("transaction_id")
        if not transaction_id:
            return

        # Step 5: Check points balance
        balance_response = self.client.get(
            f"/api/points/{user_id}/{program_id}/balance",
            headers=auth_headers
        )

        if balance_response.status_code != 200:
            return

        # Check points ledger
        ledger_response = self.client.get(
            f"/api/points/{user_id}/{program_id}/ledger",
            headers=auth_headers
        )

        if ledger_response.status_code != 200:
            return

        # Step 6: Redeem points
        redeem_data = {
            "points": 10000
        }
        redeem_response = self.client.post(
            f"/api/points/{user_id}/{program_id}/redeem",
            json=redeem_data,
            headers=auth_headers
        )

        if redeem_response.status_code != 200:
            return

        # Step 7: Verify final balance
        final_balance_response = self.client.get(
            f"/api/points/{user_id}/{program_id}/balance",
            headers=auth_headers
        )

        if final_balance_response.status_code != 200:
            return