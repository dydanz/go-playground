import json
import yaml

swagger_file = "./../internal/docs/swagger.json"

def load_swagger(file_path):
    """Load and parse Swagger JSON or YAML file"""
    if file_path.endswith(".json"):
        with open(file_path, "r") as file:
            return json.load(file)
    elif file_path.endswith(".yaml") or file_path.endswith(".yml"):
        with open(file_path, "r") as file:
            return yaml.safe_load(file)
    else:
        raise ValueError("Unsupported file format. Use JSON or YAML.")

def generate_curl_from_swagger(swagger_file, base_url="http://localhost:8080/swagger/index.htm"):
    """Generate cURL commands from Swagger (OpenAPI) file"""
    swagger_data = load_swagger(swagger_file)
    curl_commands = []

    for path, methods in swagger_data.get("paths", {}).items():
        for method, details in methods.items():
            headers = '-H "Content-Type: application/json"'
            body = ""

            # Add Authorization header if defined
            if "security" in details:
                headers += ' -H "Authorization: Bearer <TOKEN>"'

            # Check if the request has a body (for POST, PUT, PATCH)
            if method.lower() in ["post", "put", "patch"] and "requestBody" in details:
                example_body = details["requestBody"].get("content", {}).get("application/json", {}).get("example")
                if not example_body:
                    example_body = "{}"  # Default empty JSON
                body = f"-d '{json.dumps(example_body, indent=2)}'"

            # Construct cURL command
            curl_cmd = f'curl -X {method.upper()} "{base_url}{path}" {headers} {body}'
            curl_commands.append(curl_cmd.strip())

    return curl_commands

# Example Usage
curl_commands = generate_curl_from_swagger("./../internal/docs/swagger.json")
for cmd in curl_commands:
    print(cmd)
