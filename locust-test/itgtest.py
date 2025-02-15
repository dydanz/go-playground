#!/usr/bin/env python3

import requests 
import time
import uuid
from typing import Dict, Optional
from datetime import datetime, timezone, timedelta
import redis
import json
import psycopg2
import datamock
import random
import secrets
import string
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

def generate_random_email():
    # Generate a random string for the email prefix
    prefix = ''.join(secrets.choice(string.ascii_lowercase) for _ in range(8))

    return f"{prefix}@example.com"

def generate_random_name():
    return ''.join(secrets.choice(string.ascii_lowercase) for _ in range(8))

def generate_random_phone():
    return ''.join(secrets.choice(string.digits) for _ in range(10))  

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

def run_sequence(name_param: str, email_param: str):
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
        print("\n1. Registering user "+name_param+" email: "+email_param+"...")
        register_data = {
            "email": email_param,
            "name": name_param,
            "password": "Password123!",
            "phone": "+1234567890"
        }
        register_response = make_request("POST", "/api/auth/register", data=register_data, headers=auth_headers)
        
        # Check for 409 Conflict
        if register_response.status_code == 409:
            print("User already exists, skipping to re-login...")
            time.sleep(0.005)
            registration_success = False
        else:
            print("Registration successful")
            time.sleep(0.005)
            registration_success = True
        # 2. Get verification code
        if registration_success:
            print("\n2. Getting verification code...")
            verify_response = make_request("GET", f"/api/auth/test/get-verification/code?email={register_data['email']}", headers=auth_headers)
            print("Got verification code")
            otp = verify_response.json().get('otp')  # Assuming the OTP is returned in the response
            time.sleep(0.005)

            # 3. Verify code
            print("\n3. Verifying code...")
            verify_data = {
                "email": register_data['email'],
                "otp": otp
            }
            verify_response = make_request("POST", "/api/auth/verify", data=verify_data, headers=auth_headers)
            print("Verification successful")
            time.sleep(0.005)

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
        time.sleep(0.005)

        print("\n4.1 Checking if user exists in PostgreSQL...")
        execute_query("users", "email", register_data['email'])
        time.sleep(0.005)

        # Check Redis connection and data
        print("\n4.2 Checking Redis connection and data...")
        check_redis_data(user_id, auth_token)
        time.sleep(0.005)

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
        time.sleep(0.005)

        # 6. Logout
        print("\n6. Logging out...")
        logout_response = make_request("POST", "/api/auth/logout", headers=auth_headers)
        print("Logged out successfully")
        time.sleep(0.005)

        # 7. Re-login
        print("\n7. Re-logging in...")
        login_response = make_request("POST", "/api/auth/login", data=login_data)
        auth_data = login_response.json()
        auth_token = auth_data.get('token')
        user_id = auth_data.get('user_id')
        print("Re-login successful")
        time.sleep(0.005)
        print(f"Auth token: {auth_token}")
        print(f"User ID: {user_id}")

        # Continue with the rest of the sequence...
        # 8. Create merchant
        if registration_success:
            print("\n8. Creating merchant...")
            merchant_data = {
                "merchant_name": random.choice(["PT. ", "CV ", "Startup "]) + name_param,
                "merchant_type": random.choice(["repair_shop", "bank", "e-commerce"]),
                "user_id": user_id
            }
            merchant_response = make_request("POST", "/api/merchants", headers=auth_headers, data=merchant_data)
            merchant_id = merchant_response.json().get('id')
            print(f"Created merchant with ID: {merchant_id}")
            time.sleep(0.005)
        else:
            print("\n8. Get existing merchant...")
            merchant_response_list = make_request("GET", "/api/merchants", headers=auth_headers)
            for merchant in merchant_response_list.json():
                merchant_id = merchant.get('id')
            print(f"Got existing merchant with ID: {merchant_id}")
            time.sleep(0.005)   
        
        # 8.1 Create merchant customer user
        print("\n8.1 Creating merchant customer user...")
        merchant_customer_data = {
            "email": generate_random_email(),
            "merchant_id": merchant_id,
            "name": generate_random_name(),
            "password": "CustomerPassword123!",
            "phone": generate_random_phone()
        }
        
        merchant_customer_response = make_request("POST", "/api/merchant-customers", headers=auth_headers, 
                                                  data=merchant_customer_data)
        # Check for 409 Conflict
        if merchant_customer_response.status_code == 409 or merchant_customer_response.status_code == 401:
            print("Merchant customer user already exists, skipping to get existing merchant customer user...")
            time.sleep(0.005)
            registration_failed = True
        else:
            print(f"Created merchant customer user with ID: {merchant_customers_user_id}")
            time.sleep(0.005)
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
        time.sleep(0.005)

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
        time.sleep(0.005)

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
        time.sleep(0.005)

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
    
def run_sequence_2nd(name_param: str, email_param: str):
    # Store session data
    auth_token = None
    user_id = None
    merchant_ids = []
    program_ids = {}
    merchant_customer_ids = {}

    try:
        auth_headers = {
            'Content-Type': 'application/json'
        }
        # 1. Register
        print("\n1. Registering user "+name_param+" email: "+email_param+"...")
        register_data = {
            "email": email_param,
            "name": name_param,
            "password": "Password123!",
            "phone": "+1234567890"
        }
        register_response = make_request("POST", "/api/auth/register", data=register_data, headers=auth_headers)
        
        # Check for 409 Conflict
        if register_response.status_code == 409:
            print("User already exists, skipping to re-login...")
            time.sleep(0.005)
            registration_success = False
        else:
            print("Registration successful")
            time.sleep(0.005)
            registration_success = True
            # 2. Get verification code
            print("\n2. Getting verification code...")
            verify_response = make_request("GET", f"/api/auth/test/get-verification/code?email={register_data['email']}", headers=auth_headers)
            print("Got verification code")
            otp = verify_response.json().get('otp')
            time.sleep(0.005)

            # 3. Verify code
            print("\n3. Verifying code...")
            verify_data = {
                "email": register_data['email'],
                "otp": otp
            }
            verify_response = make_request("POST", "/api/auth/verify", data=verify_data, headers=auth_headers)
            print("Verification successful")
            time.sleep(0.005)

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
        time.sleep(0.005)

        # Set auth headers for subsequent requests
        auth_headers = {
            'Authorization': f'Bearer {auth_token}',
            'X-User-Id': user_id,
            'Content-Type': 'application/json'
        }

        if registration_success:
            # 5. Create 5 merchants
            print("\n5. Creating 5 merchants...")
            for i in range(5):
                merchant_data = {
                    "merchant_name": random.choice(["PT. ", "CV ", "Startup "]) + f"{name_param}_{i}",
                    "merchant_type": random.choice(["repair_shop", "bank", "e-commerce"]),
                    "user_id": user_id
                }
                merchant_response = make_request("POST", "/api/merchants", headers=auth_headers, data=merchant_data)
                merchant_id = merchant_response.json().get('id')
                merchant_ids.append(merchant_id)
                program_ids[merchant_id] = []
                merchant_customer_ids[merchant_id] = []
                print(f"Created merchant with ID: {merchant_id}")
                time.sleep(0.0005)

                # 6. Create 3-7 programs per merchant
                num_programs = random.randint(3, 7)
                print(f"\n6. Creating {num_programs} programs for merchant {merchant_id}...")
                for j in range(num_programs):
                    program_data = {
                        "merchant_id": merchant_id,
                        "user_id": user_id,
                        "program_name": f"Program {j} - {uuid.uuid4()}",
                        "point_currency_name": "Points"
                    }
                    program_response = make_request("POST", "/api/programs", headers=auth_headers, data=program_data)
                    program_id = program_response.json().get('program_id')
                    program_ids[merchant_id].append(program_id)
                    print(f"Created program with ID: {program_id}")
                    time.sleep(0.0005)

                    # 7. Create default program rule
                    rule_data = {
                        "program_id": program_id,
                        "rule_name": f"Standard Rule {j}",
                        "condition_type": "amount",
                        "condition_value": str(random.randint(50, 200)),
                        "multiplier": random.uniform(0.5, 2.0),
                        "points_awarded": random.randint(5, 20),
                        "effective_from": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                        "effective_to": (datetime.now(timezone.utc) + timedelta(days=365)).strftime("%Y-%m-%dT%H:%M:%SZ")
                    }
                    rule_response = make_request("POST", "/api/program-rules", headers=auth_headers, data=rule_data)
                    print(f"Created program rule for program {program_id}")
                    time.sleep(0.0005)

                    # 8. Create 2-5 rewards per program
                    num_rewards = random.randint(2, 5)
                    for k in range(num_rewards):
                        reward_data = {
                            "program_id": program_id,
                            "name": f"Reward {k} for Program {j}",
                            "description": f"Test Reward {k} Description",
                            "points_required": random.randint(5, 100),
                            "quantity": random.randint(10, 100),
                            "is_active": True
                        }
                        reward_response = make_request("POST", "/api/rewards", headers=auth_headers, data=reward_data)
                        print(f"Created reward for program {program_id}")
                        time.sleep(0.00005)

                # 9. Create 5-10 merchant clients
                num_clients = random.randint(5, 10)
                print(f"\n9. Creating {num_clients} merchant clients for merchant {merchant_id}...")
                for l in range(num_clients):
                    merchant_customer_data = {
                        "email": generate_random_email(),
                        "merchant_id": merchant_id,
                        "name": generate_random_name(),
                        "password": "CustomerPassword123!",
                        "phone": generate_random_phone()
                    }
                    merchant_customer_response = make_request("POST", "/api/merchant-customers", headers=auth_headers, 
                                                          data=merchant_customer_data)
                    merchant_customer_id = merchant_customer_response.json().get('id')
                    merchant_customer_ids[merchant_id].append(merchant_customer_id)
                    print(f"Created merchant customer with ID: {merchant_customer_id}")
                    time.sleep(0.00005)

        else:
            # Get existing merchants and their data
            print("\n5. Getting existing merchants...")
            merchant_response_list = make_request("GET", "/api/merchants", headers=auth_headers)
            for merchant in merchant_response_list.json():
                merchant_id = merchant.get('id')
                merchant_ids.append(merchant_id)
                program_ids[merchant_id] = []
                merchant_customer_ids[merchant_id] = []

                # Get programs for this merchant
                programs_response = make_request("GET", f"/api/programs/merchant/{merchant_id}", headers=auth_headers)
                for program in programs_response.json():
                    program_ids[merchant_id].append(program.get('program_id'))

                # Get merchant customers for this merchant
                customers_response = make_request("GET", f"/api/merchant-customers/merchant/{merchant_id}", 
                                                headers=auth_headers)
                for customer in customers_response.json():
                    merchant_customer_ids[merchant_id].append(customer.get('id'))

        # 10. Create transactions and redemptions
        print("\n10. Creating transactions and redemptions...")
        start_date = datetime(2024, 8, 1, tzinfo=timezone.utc)
        end_date = datetime.now(timezone.utc)

        for merchant_id in merchant_ids:
            for customer_id in merchant_customer_ids[merchant_id]:
                num_transactions = random.randint(30, 50)
                for _ in range(num_transactions):
                    # Generate random date
                    random_date = start_date + timedelta(seconds=random.randint(0, int((end_date - start_date).total_seconds())))
                    
                    # Select random program
                    program_id = random.choice(program_ids[merchant_id])
                    tx_date = random_date.strftime("%Y-%m-%dT%H:%M:%SZ")
                    # Create transaction
                    transaction_data = {
                        "merchant_id": merchant_id,
                        "merchant_customers_id": customer_id,
                        "program_id": program_id,
                        "transaction_type": "purchase",
                        "transaction_amount": random.uniform(50.0, 500.0),
                        "transaction_date": tx_date,
                        "status": random.choice(["pending", "completed", "failed", "cancelled"])
                    }
                    transaction_response = make_request("POST", "/api/transactions", headers=auth_headers, 
                                                      data=transaction_data)
                    print(f"Created transaction for customer {customer_id} trx-date: {tx_date}")
                    time.sleep(0.0005)

                    # Randomly create redemption (30% chance)
                    if random.random() < 0.3:
                        # Get rewards for the program
                        rewards_response = make_request("GET", f"/api/rewards/program/{program_id}", headers=auth_headers)
                        rewards = rewards_response.json()

                        # Parse the string into a datetime object
                        dt = datetime.strptime(tx_date, "%Y-%m-%dT%H:%M:%SZ")
                        # Add 1 minute
                        dt += timedelta(minutes=1)
                        new_time_str = dt.strftime("%Y-%m-%dT%H:%M:%SZ")

                        if rewards:
                            reward = random.choice(rewards)
                            redemption_data = {
                                "merchant_customers_id": customer_id,
                                "reward_id": reward.get('id'),
                                "points_used": reward.get('points_required'),
                                "point_required": reward.get('points_required'),
                                "redempton_date": new_time_str,
                                "status":  random.choice(["pending", "completed", "failed"]),
                            }
                            redemption_response = make_request("POST", "/api/redemptions", headers=auth_headers, 
                                                             data=redemption_data)
                            print(f"Created redemption for customer {customer_id}")
                            time.sleep(0.0005)

        print("\nSequence completed successfully!")

    except Exception as e:
        print(f"\nError in sequence: {str(e)}")
        raise

if __name__ == "__main__":
    try:
        name, email = datamock.get_random_name()
        run_sequence_2nd(name, email)
        print("\nSequence completed successfully!")
    except Exception as e:
        print(f"\nSequence failed: {str(e)}")

