#!/bin/bash

# CONFIGURATION
REDIS_CLI="redis-cli -u redis://localhost:6379"
TYPE="$1" # "customer" or "event"
CSV_FILE="${TYPE}_sample.csv"
FIELDS_CUSTOMER=("email" "phone" "visitor_id")
FIELDS_EVENT=("visitor_id" "call_id" "chat_id" "external_id" "form2lead_id" "tickets_id")

if [[ "$TYPE" != "customer" && "$TYPE" != "event" ]]; then
  echo "Usage: $0 [customer|event]"
  exit 1
fi

echo "[1/6] Fetching all keys for pattern ${TYPE}:* ..."
# 1. Get all keys
KEYS=($($REDIS_CLI --raw keys "${TYPE}:*"))
TOTAL=${#KEYS[@]}
if (( TOTAL == 0 )); then
  echo "No keys found for pattern ${TYPE}:*"
  exit 1
fi

echo "[2/6] Calculating sample size ..."
# 2. Calculate 5% sample size (at least 1)
SAMPLE_SIZE=$(( (TOTAL + 19) / 20 ))
(( SAMPLE_SIZE == 0 )) && SAMPLE_SIZE=1

echo "[3/6] Randomly sampling keys ..."
# 3. Randomly sample keys
SAMPLED_KEYS=($(printf "%s\n" "${KEYS[@]}" | shuf -n $SAMPLE_SIZE))

echo "[4/6] Writing CSV header ..."
# 4. Write CSV header
if [ "$TYPE" == "customer" ]; then
  echo "key,email,phone,visitor_id" > "$CSV_FILE"
  FIELDS=("${FIELDS_CUSTOMER[@]}")
else
  echo "key,visitor_id,call_id,chat_id,external_id,form2lead_id,tickets_id" > "$CSV_FILE"
  FIELDS=("${FIELDS_EVENT[@]}")
fi

echo "[5/6] Extracting fields and writing to CSV ..."
# 5. Extract fields and write to CSV using the Go API
API_URL="http://localhost:8080/document_by_key"
for KEY in "${SAMPLED_KEYS[@]}"; do
  RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "${API_URL}?key=${KEY}")
  # Extract body and status
  BODY=$(echo "$RESPONSE" | sed -e 's/HTTPSTATUS:.*//g')
  STATUS=$(echo "$RESPONSE" | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
  if [ "$STATUS" != "200" ]; then
    echo "[WARN] Skipping $KEY: API returned status $STATUS" >&2
    continue
  fi
  DOCUMENT=$(echo "$BODY" | jq -c '.document')
  if [ "$DOCUMENT" == "null" ] || [ -z "$DOCUMENT" ]; then
    echo "[WARN] Skipping $KEY: No document found" >&2
    continue
  fi
  RECORD=$(echo "$DOCUMENT" | jq 'if type=="array" then .[0] else . end')
  if [ "$TYPE" == "customer" ]; then
    EMAIL=$(echo "$RECORD" | jq -r '.primaryIdentifiers.email // empty')
    PHONE=$(echo "$RECORD" | jq -r '.primaryIdentifiers.phone // empty')
    VISITOR_ID=$(echo "$RECORD" | jq -r '.primaryIdentifiers.cmec_visitor_id // empty')
    echo "\"$KEY\",\"$EMAIL\",\"$PHONE\",\"$VISITOR_ID\"" >> "$CSV_FILE"
  else
    VISITOR_ID=$(echo "$RECORD" | jq -r '.identifiers.cmec_visitor_id // empty')
    CALL_ID=$(echo "$RECORD" | jq -r '.identifiers.cmec_contact_call_id // empty')
    CHAT_ID=$(echo "$RECORD" | jq -r '.identifiers.cmec_contact_chat_id // empty')
    EXTERNAL_ID=$(echo "$RECORD" | jq -r '.identifiers.cmec_contact_external_id // empty')
    FORM2LEAD_ID=$(echo "$RECORD" | jq -r '.identifiers.cmec_contact_form2lead_id // empty')
    TICKETS_ID=$(echo "$RECORD" | jq -r '.identifiers.cmec_contact_tickets_id // empty')
    echo "\"$KEY\",\"$VISITOR_ID\",\"$CALL_ID\",\"$CHAT_ID\",\"$EXTERNAL_ID\",\"$FORM2LEAD_ID\",\"$TICKETS_ID\"" >> "$CSV_FILE"
  fi
done

echo "[6/6] Finished! Sampled $SAMPLE_SIZE records out of $TOTAL. Output written to $CSV_FILE"

