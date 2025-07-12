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

# 1. Get all keys
KEYS=($($REDIS_CLI --raw keys "${TYPE}:*"))
TOTAL=${#KEYS[@]}
if (( TOTAL == 0 )); then
  echo "No keys found for pattern ${TYPE}:*"
  exit 1
fi

# 2. Calculate 5% sample size (at least 1)
SAMPLE_SIZE=$(( (TOTAL + 19) / 20 ))
(( SAMPLE_SIZE == 0 )) && SAMPLE_SIZE=1

# 3. Randomly sample keys
SAMPLED_KEYS=($(printf "%s\n" "${KEYS[@]}" | shuf -n $SAMPLE_SIZE))

# 4. Write CSV header
if [ "$TYPE" == "customer" ]; then
  echo "key,email,phone,visitor_id" > "$CSV_FILE"
  FIELDS=("${FIELDS_CUSTOMER[@]}")
else
  echo "key,visitor_id,call_id,chat_id,external_id,form2lead_id,tickets_id" > "$CSV_FILE"
  FIELDS=("${FIELDS_EVENT[@]}")
fi

# 5. Extract fields and write to CSV
for KEY in "${SAMPLED_KEYS[@]}"; do
  JSON=$($REDIS_CLI --raw json.get "$KEY" '$')
  # Remove outer array if present (from JSON.GET $)
  RECORD=$(echo "$JSON" | jq 'if type=="array" then .[0] else . end')
  if [ "$TYPE" == "customer" ]; then
    EMAIL=$(echo "$RECORD" | jq -r '.primaryIdentifiers.email // empty')
    PHONE=$(echo "$RECORD" | jq -r '.primaryIdentifiers.phone // empty')
    VISITOR_ID=$(echo "$RECORD" | jq -r '.identifiers.visitor_ids[0] // empty')
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

echo "Sampled $SAMPLE_SIZE records out of $TOTAL. Output written to $CSV_FILE"
