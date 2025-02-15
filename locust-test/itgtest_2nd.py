import random
import time
import uuid
from datetime import datetime, timezone, timedelta
from .itgtest import make_request, execute_query, check_redis_data, generate_random_email, generate_random_name, generate_random_phone

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
            time.sleep(0.45)
            registration_success = False
        else:
            print("Registration successful")
            time.sleep(0.45)
            registration_success = True
            # 2. Get verification code
            print("\n2. Getting verification code...")
            verify_response = make_request("GET", f"/api/auth/test/get-verification/code?email={register_data['email']}", headers=auth_headers)
            print("Got verification code")
            otp = verify_response.json().get('otp')
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
                time.sleep(0.45)

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
                    time.sleep(0.45)

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
                    time.sleep(0.45)

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
                        time.sleep(0.45)

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
                    time.sleep(0.45)

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
                    
                    # Create transaction
                    transaction_data = {
                        "merchant_id": merchant_id,
                        "merchant_customers_id": customer_id,
                        "program_id": program_id,
                        "transaction_type": "purchase",
                        "transaction_amount": random.uniform(50.0, 500.0),
                        "transaction_date": random_date.strftime("%Y-%m-%dT%H:%M:%SZ"),
                        "status": "completed"
                    }
                    transaction_response = make_request("POST", "/api/transactions", headers=auth_headers, 
                                                      data=transaction_data)
                    print(f"Created transaction for customer {customer_id}")
                    time.sleep(0.45)

                    # Randomly create redemption (30% chance)
                    if random.random() < 0.3:
                        # Get rewards for the program
                        rewards_response = make_request("GET", f"/api/rewards/program/{program_id}", headers=auth_headers)
                        rewards = rewards_response.json()
                        if rewards:
                            reward = random.choice(rewards)
                            redemption_data = {
                                "merchant_customers_id": customer_id,
                                "reward_id": reward.get('id'),
                                "points_used": reward.get('points_required'),
                                "point_required": reward.get('points_required'),
                                "status": "completed"
                            }
                            redemption_response = make_request("POST", "/api/redemptions", headers=auth_headers, 
                                                             data=redemption_data)
                            print(f"Created redemption for customer {customer_id}")
                            time.sleep(0.45)

        print("\nSequence completed successfully!")

    except Exception as e:
        print(f"\nError in sequence: {str(e)}")
        raise


if __name__ == "__main__":
    try:
        name, email = datamock.get_random_name()
        run_sequence(name, email)
        print("\nSequence completed successfully!")
    except Exception as e:
        print(f"\nSequence failed: {str(e)}")