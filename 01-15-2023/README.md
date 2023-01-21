export SUI_RPC_HOST='https://fullnode.devnet.sui.io:443'

signer : <SuiAddress> - the transaction signer's Sui address
coin_object_id : <ObjectID> - the coin object to be spilt
split_amounts : <[]> - the amounts to split out from the coin
gas : <ObjectID> - gas object to be used in this transaction, node will pick one from the signer's possession if not provided
gas_budget : <uint64> - the gas budget, the transaction will fail if the gas cost exceed the budget

curl --location --request POST $SUI_RPC_HOST \
--header 'Content-Type: application/json' \
--data-raw '{
"jsonrpc": "2.0",
"id": 1,
"method": "sui_batchTransaction",
"params": [
"0xed2c39b73e055240323cf806a7d8fe46ced1cabb",
[

        "splitCoinRequestParams": {
        "signer": "0xed2c39b73e055240323cf806a7d8fe46ced1cabb",
          "coin_object_id": "0x5f72138198f1e5706938a9e90588bf4264de2f73",
          "split_amounts": ["99999000"],
          "gas_budget": "3000"
        },

{
"moveCallRequestParams": {
"packageObjectId": "0x9c3cc9999b66afad398eeaf0f2da4923508e0ad1",
"module": "exact_balance",
"function": "pay_me",
"typeArguments": [],
"arguments": ["0x0591b9e2121f9bbb01ad876a9fabc7c06c58ea39"]
}
}
],
"0x5f72138198f1e5706938a9e90588bf4264de2f73",
2000
]
}'
