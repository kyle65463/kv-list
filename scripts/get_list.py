import requests
import sys

if len(sys.argv) != 2:
    print("Usage: python get_list.py <list_key> [base_endpoint_url]")
    exit(1)

list_key = sys.argv[1]
base_url = sys.argv[2] if len(sys.argv) > 2 != None else "http://localhost:8080"

# Get the head of the list
res = requests.get(f"{base_url}/api/v1/lists/{list_key}")
if res.status_code != 200:
    print(f"Get the list '{list_key}' failed")
    exit(1)
body = res.json()
next_page_key = body["nextPageKey"]

# Get the page recursively
pages_data = []
while next_page_key != None:
    res = requests.get(f"{base_url}/api/v1/pages/{next_page_key}")
    if res.status_code != 200:
        print(f"Get the page '{next_page_key}' failed")
        exit(1)
    body = res.json()
    next_page_key = body["nextPageKey"]
    page_data = body["data"]
    pages_data.append(page_data)

# Print the result
print({
    "key": list_key,
    "pages": pages_data
})
