"""
curl --request POST \
--url https://onefootball.atlassian.net/gateway/api/compass/v1/metrics \
--user "$USER_EMAIL:$USER_API_KEY" \
--header "Accept: application/json" \
--header "Content-Type: application/json" \
--data "{
  \"metricSourceId\": \"ari:cloud:compass:fca6a80f-888b-4079-82e6-3c2f61c788e2:metric-source/4d010f50-96c4-48c0-bab5-a3dd5112b464/0f593a5d-b19f-4896-aa04-119a63861009\",
  \"value\": $METRIC_VALUE,
  \"timestamp\": \"$(date -u +'%Y-%m-%dT%H:%M:%SZ')\"
}"
"""

"""

curl --request POST \
--url https://onefootball.atlassian.net/gateway/api/compass/v1/metrics \
--user "$USER_EMAIL:$USER_API_KEY" \
--header "Accept: application/json" \
--header "Content-Type: application/json" \
--data "{
  "componendId": "ari:cloud:compass:fca6a80f-888b-4079-82e6-3c2f61c788e2:metric-source/4d010f50-96c4-48c0-bab5-a3dd5112b464/0f593a5d-b19f-4896-aa04-119a63861009",
  "metricSourceId": "ari:cloud:compass:fca6a80f-888b-4079-82e6-3c2f61c788e2:metric-source/4d010f50-96c4-48c0-bab5-a3dd5112b464/0f593a5d-b19f-4896-aa04-119a63861009",
  \"value\": $METRIC_VALUE,
  \"timestamp\": \"$(date -u +'%Y-%m-%dT%H:%M:%SZ')\"
}"

"""