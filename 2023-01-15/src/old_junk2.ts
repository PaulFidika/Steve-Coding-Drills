import { Ed25519Keypair, JsonRpcProvider, RawSigner } from '@mysten/sui.js';
import { BCS, BcsConfig } from '@mysten/bcs';

// let monkey2 = Buffer.from('àϰ😂😂', 'utf8');
// console.log(monkey2.toJSON().data);

// let whatever = JSON.parse(monkey2.toJSON().data);
// console.log(whatever);

// let encoder = new TextEncoder();
// let bytes2 = encoder.encode('àϰ😂😂');
// console.log(bytes2);

// let decoder = new TextDecoder();
// console.log(decoder.decode(bytes2.buffer));

// private key: hDZ6+qWBigkbi40a+4Rpxd4NY9Y6+ZEiv0XO6OjQfzy9iW+TkgOZx2RKQIORP4bbY1XrG8Egc+Yo2Q74TNRYUw==
// public key: 0xed2c39b73e055240323cf806a7d8fe46ced1cabb
const privateKeyBytes = new Uint8Array([
  132, 54, 122, 250, 165, 129, 138, 9, 27, 139, 141, 26, 251, 132, 105, 197, 222, 13, 99, 214, 58,
  249, 145, 34, 191, 69, 206, 232, 232, 208, 127, 60, 189, 137, 111, 147, 146, 3, 153, 199, 100, 74,
  64, 131, 145, 63, 134, 219, 99, 85, 235, 27, 193, 32, 115, 230, 40, 217, 14, 248, 76, 212, 88, 83
]);

// The following primitive types are available in BCS by default:
// address, vector<T>, u8, u16, u32, u64, u128, u256, bool, string
//
// In this library 'string' means extended-ascii [0, 255], NOT ascii [0, 127] or UTF8, which are the
// types used within Move. Special characters outside of the extended-ascii set get smooshed to
// somewhere within the ascii range randomly
//
// vector works as expected, like vector<string> or vector<u64>
// Options have to be registered manually as enums

// Caveats:
// -

// For BCS, if you have an array of strings or an array of arrays, you'll get something like:
// [length of array] [length of first array / string] [...bytes of first array] ...
//
// If you are serializing a struct alone however, the items will be listed in logical order, like
// [length of first field (if it's not a fixed length, like u8 or address)] [...bytes of first field]
//
// If you are serializing an array of structs, then again the overall array-length (nubmer of items)
// will be serialized at first (prepended)

let bcsConfig: BcsConfig = {
  vectorType: 'vector',
  addressLength: 20,
  addressEncoding: 'hex',
  types: {
    enums: {
      'Option<u64>': {
        none: null,
        some: 'u64'
      }
    }
  },
  withPrimitives: true
};

let bcs = new BCS(bcsConfig);

bcs.registerType(
  'utf8',
  (writer, data: string) => {
    let bytes = new TextEncoder().encode(data);
    writer.writeVec(Array.from(bytes), (w, el) => w.write8(el));
    return writer;
  },
  reader => {
    let bytes = reader.readBytes(reader.readULEB());
    return new TextDecoder('utf8').decode(bytes);
  },
  value => typeof value == 'string'
);

bcs.registerType(
  'ascii',
  (writer, data: string) => {
    let bytes = new TextEncoder().encode(data);
    if (bytes.length > data.length) throw Error('Not ASCII string');

    writer.writeVec(Array.from(bytes), (w, el: number) => {
      if (el > 127) throw Error('Not ASCII string');
      return w.write8(el);
    });

    return writer;
  },
  reader => {
    let bytes = reader.readBytes(reader.readULEB());
    bytes.forEach(byte => {
      if (byte > 127) throw Error('Not ASCII string');
    });

    return new TextDecoder('ascii').decode(bytes);
  },
  value => typeof value == 'string'
);

let bytes3 = bcs.ser('utf8', 'àϰ😂😂').toBytes();
console.log(bytes3);
console.log(bcs.de('utf8', bytes3));

let stupid = bcs.ser('vector<utf8>', ['abc', '123', 'àϰ😂😂']).toBytes();
console.log(stupid);
console.log(bcs.de('vector<utf8>', stupid));

let bytes4 = bcs.ser('ascii', '1234 hello world!').toBytes();
console.log(bytes4);
console.log(bcs.de('ascii', bytes4));

// let bytes3 = bcs.ser('utf8', 'àϰ😂😂').toBytes();
// console.log(bytes3);
// console.log(bcs.de('utf8', bytes3));

let serialized = bcs.ser('u32', 55555).toBytes();
let deserialized = bcs.de('u32', serialized);

console.log(deserialized);
console.log(typeof deserialized);

bcs.registerStructType('Coin', { value: 'u64' });
bcs.registerEnumType('MyEnum', {
  single: 'Coin',
  multi: 'vector<Coin>',
  empty: null
});

console.log(
  bcs.de('MyEnum', 'AICWmAAAAAAA', 'base64'), // { single: { value: 10000000 } }
  bcs.de('MyEnum', 'AQIBAAAAAAAAAAIAAAAAAAAA', 'base64') // { multi: [ { value: 1 }, { value: 2 } ] }
);

// and serialization
bcs.ser('MyEnum', { single: { value: 10000000 } }).toBytes();
bcs.ser('MyEnum', { multi: [{ value: 1 }, { value: 2 }] });

bcs.registerStructType('Coin', {
  value: 'u64',
  owner: 'd',
  is_locked: 'd'
});

bcs.registerType(
  'number_string',
  (writer, data) => writer.writeVec(data, (w, el) => w.write8(el)),
  reader => reader.readVec(r => r.read8()).join(''), // read each value as u8
  value => /[0-9]+/.test(value) // test that it has at least one digit
);

// console.log(Array.from(bcs.ser('number_string', '12345').toBytes()) == [5, 1, 2, 3, 4, 5]);

// * bcs.registerVectorType('vector<u8>', 'u8');
//    *
//    * let serialized = BCS
//    *   .set('vector<u8>', [1,2,3,4,5,6])
//    *   .toBytes();

// bcs.registerEnumType('Option<u64>', {
//   none: null,
//   some: 'u64'
// });

type Outlaw = {
  name: string;
  image: string;
  power_level: string[];
  monkey: { none: null } | { some: number };
};

bcs.registerStructType('Outlaw', {
  name: 'utf8',
  image: 'utf8',
  power_level: 'vector<utf8>',
  monkey: 'Option<u64>'
});

let t_interface = bcs.getTypeInterface('Outlaw');

let metadataObject: Outlaw = {
  name: 'Kyrie😂',
  image: 'https://wikipedia.org/',
  power_level: ['hey', 'how', 'are you?'],
  monkey: { some: 19 }
};

let typeInterfaces = bcs.getTypeInterface('Outlaw');

console.log(metadataObject);

let bytes = bcs.ser('Outlaw', metadataObject).toBytes();
console.log(bytes);
let reserialized = bcs.de('Outlaw', bytes);
console.log(reserialized);

async function getObject() {
  let keypair = Ed25519Keypair.fromSecretKey(privateKeyBytes);
  const provider = new JsonRpcProvider();
  // const signer = new RawSigner(keypair, provider);

  let result = await provider.getTotalTransactionNumber();
  console.log(result);

  const object = await provider.getObject('0xbf4a8f90818aad78cc5a10bc1f4f6d0067e2cca7');

  console.log(object);

  // @ts-ignore
  console.log(object.details.data.fields);
  // @ts-ignore
  console.log(object.details.data.fields.schema[2].fields);

  type Schema = { [key: string]: string };
  let object2: Schema = {};

  // @ts-ignore
  // TO DO: handle cases where details.data does not exist, or this fails
  object.details.data.fields.schema.forEach(item => {
    object2[item.fields.key] = item.fields.type;
  });

  console.log(object2);

  bcs.registerStructType('Outlaw2', object2);

  let metadataObject2 = {
    name: 'Kyrie',
    image: 'https://wikipedia.org/',
    power_level: 199990
  };

  let metadataObject3 = {
    name: 'Jin',
    image: 'https://google.com/something/',
    power_level: 100000000
  };

  let bytes = bcs.ser('Outlaw2', metadataObject2).toBytes();
  console.log(bytes);
  let reserial = bcs.de('Outlaw2', bytes);
  console.log(reserial);

  let bytes5 = bcs.ser('vector<Outlaw2>', [metadataObject2, metadataObject3]).toBytes();
  console.log(bytes5);
  console.log(bcs.de('vector<Outlaw2>', bytes5));
}

getObject();

async function mintNFT() {
  // Import an existing keyapri
  let keypair = Ed25519Keypair.fromSecretKey(privateKeyBytes);
  const provider = new JsonRpcProvider();
  const signer = new RawSigner(keypair, provider);

  //   console.log(keypair.getPublicKey().toSuiAddress());

  // const moveCallTxn = await signer.signAndExecuteTransaction(txBytes, 'WaitForEffectsCert');

  const moveCallTxn = await signer.executeMoveCall({
    packageObjectId: '0x4f2801f232f4cd689e7d1791b74e7fad1dfa068c',
    module: 'schema',
    function: 'create',
    typeArguments: [],
    arguments: [
      ['name', 'image', 'power level'],
      ['ascii', 'ascii', 'u64'],
      [false, false, false]
    ],
    gasBudget: 10000
  });

  console.log('moveCallTxn', moveCallTxn);
}

// mintNFT();
