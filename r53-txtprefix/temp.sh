  aws route53 list-resource-record-sets --hosted-zone-id ZNCSEA5FAW4K2 --query "ResourceRecordSets[?Type == 'TXT']" > txtrecords
  jq '.[] | select(.ResourceRecords[].Value | test("heritage")) | select(.Name | test("^prefix")| not)' txtrecords

