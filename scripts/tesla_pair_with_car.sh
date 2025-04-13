if [ -z "$1" ]; then
  echo "Error: Missing VIN of the car you would like to pair with."
  echo "Usage: $0 <vin>"
  exit 1
fi

tesla-control \
    -debug \
    -vin "$1" \
    -ble \
    add-key-request public_key.pem owner cloud_key
