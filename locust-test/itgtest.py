#!/usr/bin/env python3

import requests 
import time
import uuid
from typing import Dict, Optional
from datetime import datetime, timezone, timedelta
import redis
import json
import psycopg2

"""
python itgtest.py 

to run the test, run the command above

$ python -m venv venv
$ source venv/bin/activate
$ pip install requests
$ pip install redis
$ pip install psycopg2
$ python itgtest.py

"""

BASE_URL = "http://localhost:8080"

def make_request(method: str, endpoint: str, headers: Optional[Dict] = None, data: Optional[Dict] = None) -> requests.Response:
    url = f"{BASE_URL}{endpoint}"
    try:
        if method == "GET":
            response = requests.get(url, headers=headers, timeout=5)
        elif method == "POST":
            response = requests.post(url, headers=headers, json=data, timeout=5)
        else:
            raise ValueError(f"Unsupported HTTP method: {method}")
        
        # Handle HTTP 409 Conflict without raising an exception
        if response.status_code == 409:
            print(f"Conflict occurred: {response.json()}")
            return response  # Return the response for further handling

        response.raise_for_status()  # Raise for other HTTP errors
        return response
    except requests.exceptions.RequestException as e:
        print(f"Error making {method} request to {endpoint}: {str(e)}")
        raise

pg_conn = None
def get_pgsql_client_connection():
    pg_conn = psycopg2.connect(
        host='host.docker.internal',
        port=5432,
        user='postgres',
        password='postgres',
        dbname='go_cursor'
    )
    return pg_conn  # Return the connection object

def execute_query(table_name: str, param_key: str, param_value: str):
    # Establish a connection to PostgreSQL
    pg_conn = get_pgsql_client_connection()
    cursor = pg_conn.cursor()

    query = "SELECT "+param_key+" FROM "+table_name+" WHERE "+param_key+" = '"+param_value+"';"
    
    try:
        # Execute the provided query
        cursor.execute(query)
        
        # Fetch results if the query is a SELECT statement
        if query.strip().lower().startswith("select"):
            results = cursor.fetchall()
            for row in results:
                assert row[0] == param_value 
                print("DB query result : "+row[0])
                print("Expected value  : "+param_value)
        else:
            pg_conn.commit()  # Commit changes for INSERT, UPDATE, DELETE statements
            print("Query executed successfully.")
    except Exception as e:
        print(f"Error executing query: {str(e)}")
    finally:
        # Close the cursor and connection
        cursor.close()
        pg_conn.close()

def check_redis_data(user_id: str, auth_token: str):
     # Redis configuration from .env
    redis_client = redis.Redis(
        host='host.docker.internal',  # from REDIS_HOST
        port=6379,                    # from REDIS_PORT
        password='redis123',          # from REDIS_PASSWORD
        decode_responses=True         # automatically decode responses to strings
    )

    try:
        # Test Redis connection
        redis_client.ping()
        print("Successfully connected to Redis")

        # Check if the auth token is stored in Redis
        token_key = f"session:userid:{user_id}"
        stored_token = redis_client.get(token_key)
        if stored_token:
            print(f"Found auth token in Redis for user {user_id}")
            print(f"Token from Redis: {json.loads(stored_token).get('token_hash')}")
            print(f"Token from Login: {auth_token}")
            assert json.loads(stored_token).get('token_hash') == auth_token
        else:
            print("No auth token found in Redis")

    except redis.ConnectionError as e:
        print(f"Failed to connect to Redis: {str(e)}")
    except Exception as e:
        print(f"Error checking Redis data: {str(e)}")
    finally:
        redis_client.close()

def run_sequence():
    # Store session data
    auth_token = None
    user_id = None
    merchant_id = None
    program_id = None
    merchant_customers_user_id = None

    try:
        auth_headers = {
            'Content-Type': 'application/json'
        }
        # 1. Register
        print("\n1. Registering user...")
        register_data = {
            "email": f"user_test_001@example.com",
            "name": "User Test 001",
            "password": "Password123!",
            "phone": "+1234567890"
        }
        register_response = make_request("POST", "/api/auth/register", data=register_data, headers=auth_headers)
        
        # Check for 409 Conflict
        if register_response.status_code == 409:
            print("User already exists, skipping to re-login...")
            time.sleep(0.45)
            registration_success = False
        else:
            print("Registration successful")
            time.sleep(0.45)
            registration_success = True
        # 2. Get verification code
        if registration_success:
            print("\n2. Getting verification code...")
            verify_response = make_request("GET", f"/api/auth/test/get-verification/code?email={register_data['email']}", headers=auth_headers)
            print("Got verification code")
            otp = verify_response.json().get('otp')  # Assuming the OTP is returned in the response
            time.sleep(0.45)

            # 3. Verify code
            print("\n3. Verifying code...")
            verify_data = {
                "email": register_data['email'],
                "otp": otp
            }
            verify_response = make_request("POST", "/api/auth/verify", data=verify_data, headers=auth_headers)
            print("Verification successful")
            time.sleep(0.45)

        # 4. Login
        print("\n4. Logging in...")
        login_data = {
            "email": register_data['email'],
            "password": register_data['password']
        }
        login_response = make_request("POST", "/api/auth/login", data=login_data)
        auth_data = login_response.json()
        auth_token = auth_data.get('token')
        user_id = auth_data.get('user_id')
        print("Login successful")
        time.sleep(0.45)

        print("\n4.1 Checking if user exists in PostgreSQL...")
        execute_query("users", "email", register_data['email'])
        time.sleep(0.45)

        # Check Redis connection and data
        print("\n4.2 Checking Redis connection and data...")
        check_redis_data(user_id, auth_token)
        time.sleep(0.45)

        # Set auth headers for subsequent requests
        auth_headers = {
            'Authorization': f'Bearer {auth_token}',
            'X-User-Id': user_id,
            'Content-Type': 'application/json'
        }

        # 5. Get profile
        print("\n5. Getting profile...")
        profile_response = make_request("GET", "/api/users/me", headers=auth_headers)
        print("Got profile "+str(profile_response.json()))
        time.sleep(0.45)

        # 6. Logout
        print("\n6. Logging out...")
        logout_response = make_request("POST", "/api/auth/logout", headers=auth_headers)
        print("Logged out successfully")
        time.sleep(0.45)

        # 7. Re-login
        print("\n7. Re-logging in...")
        login_response = make_request("POST", "/api/auth/login", data=login_data)
        auth_data = login_response.json()
        auth_token = auth_data.get('token')
        user_id = auth_data.get('user_id')
        print("Re-login successful")
        time.sleep(0.45)
        print(f"Auth token: {auth_token}")
        print(f"User ID: {user_id}")

        # Continue with the rest of the sequence...
        # 8. Create merchant
        if registration_success:
            print("\n8. Creating merchant...")
            merchant_data = {
                "merchant_name": "Merchant Test 001",
                "merchant_type": "bank",
                "user_id": user_id
            }
            merchant_response = make_request("POST", "/api/merchants", headers=auth_headers, data=merchant_data)
            merchant_id = merchant_response.json().get('id')
            print(f"Created merchant with ID: {merchant_id}")
            time.sleep(0.45)
        else:
            print("\n8. Get existing merchant...")
            merchant_response_list = make_request("GET", "/api/merchants", headers=auth_headers)
            for merchant in merchant_response_list.json():
                merchant_id = merchant.get('id')
            print(f"Got existing merchant with ID: {merchant_id}")
            time.sleep(0.45)   
        
        # 8.1 Create merchant customer user
        print("\n8.1 Creating merchant customer user...")
        merchant_customer_data = {
            "email": "customer_test@example.com",
            "merchant_id": merchant_id,
            "name": "Customer Test",
            "password": "CustomerPassword123!",
            "phone": "+1234567890"
        }
        
        merchant_customer_response = make_request("POST", "/api/merchant-customers", headers=auth_headers, 
                                                  data=merchant_customer_data)
        # Check for 409 Conflict
        if merchant_customer_response.status_code == 409:
            print("Merchant customer user already exists, skipping to get existing merchant customer user...")
            time.sleep(0.45)
            registration_failed = True
        else:
            print(f"Created merchant customer user with ID: {merchant_customers_user_id}")
            time.sleep(0.45)
            registration_failed = False
            merchant_customers_user_id = merchant_customer_response.json().get('id')
            print(f"New merchant customer user with ID: {merchant_customers_user_id}")

        if registration_failed:
            print("\n8.2 Get existing merchant customer user...")
            merchant_response_list = make_request("GET", "/api/merchant-customers/merchant/"+merchant_id, 
                                                      headers=auth_headers)
            for merchant in merchant_response_list.json():
                merchant_customers_user_id = merchant.get('id')    
            print(f"Got existing merchant customer user with ID: {merchant_customers_user_id}")


        # 9. Create program
        print("\n9. Creating program...")
        program_data = {
            "merchant_id": merchant_id,
            "user_id": user_id,
            "program_name": f"Test Program {uuid.uuid4()}",
            "point_currency_name": "Points"
        }
        program_response = make_request("POST", "/api/programs", headers=auth_headers, data=program_data)
        program_id = program_response.json().get('program_id')
        print(f"Created program with ID: {program_id}")
        time.sleep(0.45)

        # 10. Create program rule
        print("\n10. Creating program rule...")
        rule_data = {
            "program_id": program_id,
            "rule_name": "Standard Rule",
            "condition_type": "amount",
            "condition_value": "100",
            "multiplier": 1.0,
            "points_awarded": 10,
            "effective_from": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
            "effective_to": (datetime.now(timezone.utc) + timedelta(days=1)).strftime("%Y-%m-%dT%H:%M:%SZ")
        }
        rule_response = make_request("POST", "/api/program-rules", headers=auth_headers, data=rule_data)
        rule_id = rule_response.json().get('id')
        print(f"Created program rule with ID: {rule_id}")
        time.sleep(0.45)

        # 11. Create transaction
        print("\n11. Creating transaction...")
        transaction_data = {
            "merchant_id": merchant_id,
            "merchant_customers_id": merchant_customers_user_id,
            "program_id": program_id,
            "transaction_type": "purchase",
            "transaction_amount": 100.0,
            "status": "completed"
        }
        transaction_response = make_request("POST", "/api/transactions", headers=auth_headers, data=transaction_data)
        transaction_id = transaction_response.json().get('transaction_id')
        print(f"Created transaction with ID: {transaction_id}")
        time.sleep(0.45)

        # 12. Create reward
        print("\n12. Creating reward...")
        reward_data = {
            "program_id": program_id,
            "name": "Test Reward",
            "description": "Test Reward Description",
            "points_required": 1,
            "quantity": 10,
            "is_active": True
        }

        reward_response = make_request("POST", "/api/rewards", headers=auth_headers, data=reward_data)
        reward_id = reward_response.json().get('id')
        print(f"Created reward with ID: {reward_id}")
        time.sleep(2)

        # 13. Create redemption
        print("\n13. Creating redemption...")
        redemption_data = {
            "merchant_customers_id": merchant_customers_user_id,
            "reward_id": reward_id,
            "points_used": 1,
            "point_required": 1,
            "status": "completed"
        }
        redemption_response = make_request("POST", "/api/redemptions", headers=auth_headers, data=redemption_data)
        redemption_id = redemption_response.json().get('id')
        print(f"Created redemption with ID: {redemption_id}")
        time.sleep(2)

    except Exception as e:
        print(f"\nError in sequence: {str(e)}")
        raise

if __name__ == "__main__":
    try:
        run_sequence()
        print("\nSequence completed successfully!")
    except Exception as e:
        print(f"\nSequence failed: {str(e)}")

