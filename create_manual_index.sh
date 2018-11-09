 curl --request PUT \
   --url http://localhost:9200/transactions \
   --header 'Content-Type: application/json' \
   --data '{
    "mappings":{
       "_doc":{
          "properties":{
             "sequence":{
                "type":"integer"
             },
             "date":{
                "type":"date"
             },
             "price":{
                "type":"double"
             },
             "amount":{
                "type":"double"
             },
             "value":{
                "type":"double"
             },
             "side1_account_id":{
                "type":"keyword"
             },
             "side2_account_id":{
                "type":"keyword"
             },
             "currency_origin":{
                "type":"keyword"
             },
             "currency_target":{
                "type":"keyword"
             }

    }
 }}}'
