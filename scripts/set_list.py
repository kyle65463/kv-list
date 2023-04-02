import requests
import sys
import base64

def encode(str):
    base64_bytes = base64.b64encode(str.encode("utf-8"))
    base64_string = base64_bytes.decode("ascii")
    return base64_string

if len(sys.argv) != 2:
    print("Usage: python get_list.py <list_key> [base_endpoint_url]")
    exit()

list_key = sys.argv[1]
base_url = sys.argv[2] if len(sys.argv) > 2 != None else "http://localhost:8080"

# Define the pages of the list
pages = [
    # TODO: Change the data
    'PAGE 1',
    'PAGE 2',
    'PAGE 3',
    'PAGE 4',
    'PAGE 5',
]
pages = [encode(page) for page in pages]

# Set the list
res = requests.post(f"{base_url}/api/v1/lists/{list_key}", json={"data": pages})
if res.status_code == 200:
    print(f"Set the list '{list_key}' successfully")
else:
    print(res.text)
    print(f"Set the list '{list_key}' failed")