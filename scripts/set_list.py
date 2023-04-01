import requests
import sys

if len(sys.argv) != 2:
    print("Usage: python get_list.py <list_key> [base_endpoint_url]")
    exit()

list_key = sys.argv[1]
base_url = sys.argv[2] if len(sys.argv) > 2 != None else "http://localhost:8080"

# Define the pages of the list
pages = [
    # TODO: Change the data
    "ffaa", # The first page of the list
    "ccbb", 
    "aacc",
]

# Set the list
res = requests.post(f"{base_url}/api/v1/lists/{list_key}", json={"data": pages})
if res.status_code == 200:
    print(f"Set the list '{list_key}' successfully")
else:
    print(res.text)
    print(f"Set the list '{list_key}' failed")