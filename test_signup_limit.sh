#!/bin/bash
# filepath: /home/ayiko/Desktop/server/test_signup_limit.sh

BASE_URL="http://localhost:8080"
ENDPOINT="/users"
TOTAL_REQUESTS=100
SUCCESS=0
FAILURES=0

echo "Sending $TOTAL_REQUESTS signup requests to $BASE_URL$ENDPOINT"

for i in $(seq 1 $TOTAL_REQUESTS); do
    USER_DATA="{\"firstName\":\"Test\",\"lastName\":\"User\",\"phoneNumber\":\"1234567890\",\"email\":\"test${i}@example.com\",\"password\":\"password123\"}"
    
    response=$(curl -s -o /dev/null -w "%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "$USER_DATA" \
        $BASE_URL$ENDPOINT)
    
    if [ $response -eq 201 ] || [ $response -eq 200 ]; then
        SUCCESS=$((SUCCESS + 1))
        echo -ne "Request $i: OK ($response) | Success: $SUCCESS, Failures: $FAILURES\r"
    else
        FAILURES=$((FAILURES + 1))
        echo -e "Request $i: FAILED ($response) | Success: $SUCCESS, Failures: $FAILURES"
    fi
    
    sleep 0.05
done

echo -e "\nTest complete. Success: $SUCCESS, Failures: $FAILURES"
