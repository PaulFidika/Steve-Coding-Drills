### BCS Tutorial

This is a tutorial for the '@mysten/bcs' library.

The following primitive types are available in BCS by default:

`address`, `vector<T>`, `u8`, `u16`, `u32`, `u64`, `u128`, `u256`, `bool`, `string`

In the bcs library, `string` means extended-ascii (a single byte from [0, 255]), NOT ASCII (single byte from [0, 127]), or UTF8 (1 - 4 bytes per character), which are the types used within Move. In `string`, special characters outside of the extended-ascii range get smooshed to somewhere within the ascii range randomly. These extended-ASCII bytes do not convert into UTF8 or ASCII due to conflicts. As such NO ONE should use the `string` type, because (de)serialization will fail when interacting with on-chain modules. Instead we define our own custom `utf8` and `ascii` types within bcs, which should always be used instead.

Vector types work as expected, like vector<utf8> or vector<u64>. If `T` is defined, then `vector<T>` is automatically defined.

`Option<T>` has to be registered manually in bcs as an enum for each `T`. For example, we can define `Option<u128>` as follows:

```
bcs.registerEnumType('Option<u128>', {
  none: null,
  some: 'u128'
});
```

It's important that 'none' come first; when bcs serializes this, it will prepend the u128 with an option byte. We want to ensure that option-byte `0 = none`, and `1 = some` at all times. We can then use our new `Option<u128>` like below:

```
type BalanceMaybe = {
  owner: string;
  balance: { none: null } | { some: bigint };
};

bcs.registerStructType('BalanceMaybe', {
  owner: 'utf8',
  balance: 'Option<u128>'
});

let balance1: BalanceMaybe = {
  owner: 'Paul',
  balance: { some: 999007 }
};

let balance2: BalanceMaybe = {
  owner: 'Tara',
  balance: { none: null }
};
```

Above, we first define a type in Typescript, then register a type in bcs.

Notice that Typescript's types are less precise than Move's types (i.e., string VS utf8, or bigint VS u128), as such it's up to the Typescript application to ensure that the values it's using can be seralized into the corresponding Move type. For example, if Move expects a u8 for a field, then you have to ensure that the corresponding value you're using in Typescript is always [0, 255], and not something like -15 or 312, which are both valid `number` types in Typescript, but invalid for a u8 in Move.

### Binary Canonical Serialization (BCS)

BCS prepends types of variable length with a ULEB128 number, which is always 1 - 4 bytes long. For example, if you have an array of strings, or an array of arrays, you'll get something like:

[length of array as ULEB128]
[length of first array / string as ULEB128]
[...bytes of first array]
...

If you are serializing a struct alone however, the fields will be listed in the order they were defined, like:

[length of first field (if it's variable-length, such as vector<T> or utf8)]
[...bytes of first field]

For fixed-length fields, such as `u8` or `address`, there will be no length prepended.

### Using BCS As An Argument

The proper way to serialize an object is like this:

`let kyrieBytes = Array.from(bcs.ser('Outlaw', kyrie).toBytes());`

If you use just `.toBytes()` it will fail, because `Uint8Array` is not a valid type argument for the Sui Typescript SDK (yet--support coming maybe?), hence why we turn the `Unit8Array` into a `number[]`. Also `.toString()` for both `base64` and `hex` fail as well when submitted as arguments.

### RPC Stuff

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
