import {
  Ed25519Keypair,
  JsonRpcProvider,
  RawSigner,
  GetObjectDataResponse,
  RawMoveCall,
  SuiJsonValue,
  MoveCallTransaction,
  SignableTransaction,
  UnserializedSignableTransaction,
  DevInspectResults
} from '@mysten/sui.js';
import { BCS, BcsConfig } from '@mysten/bcs';
import { string } from 'superstruct';

// public key (hex): 0xed2c39b73e055240323cf806a7d8fe46ced1cabb
// private key (base64): hDZ6+qWBigkbi40a+4Rpxd4NY9Y6+ZEiv0XO6OjQfzy9iW+TkgOZx2RKQIORP4bbY1XrG8Egc+Yo2Q74TNRYUw==
const privateKeyBytes = new Uint8Array([
  132, 54, 122, 250, 165, 129, 138, 9, 27, 139, 141, 26, 251, 132, 105, 197, 222, 13, 99, 214, 58,
  249, 145, 34, 191, 69, 206, 232, 232, 208, 127, 60, 189, 137, 111, 147, 146, 3, 153, 199, 100, 74,
  64, 131, 145, 63, 134, 219, 99, 85, 235, 27, 193, 32, 115, 230, 40, 217, 14, 248, 76, 212, 88, 83
]);

// Build a class to connect to Sui RPC servers
const provider = new JsonRpcProvider();

// Import the above keypair
let keypair = Ed25519Keypair.fromSecretKey(privateKeyBytes);
const signer = new RawSigner(keypair, provider);

// ===== First we instantiate bcs, and define Option<u64> =====

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

// ===== Next we register ut8f and ascii as custom primitive types =====

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

// ===== Do A Full CRUD (Create, Update, Read, Delete) Loop =====

fullLoop();

async function fullLoop() {
  const objectQuery = await provider.getObject('0xa37996c78eed10373b5f71cf6b5cd2b623494ba5');
  let schema = parseSchema(objectQuery);
  if (!schema) return;
  bcs.registerStructType('Outlaw', schema); // This assumes the parsing succeeded

  // Instantiate JS object from schema
  // let outlaw = instantiateSchema(schema, ['My Name', 'https:///somwhere', 65]);
  // let outlaw = {
  //   name: string,
  //   image: string,
  //   power_level: bigint
  // }

  // outlaw.power_level += 3;

  // console.log('-------------------');
  // console.log(schema);

  // TO DO: it would be really cool if we could instantiate the schema as a Javascript type, class
  // or object here. Like just be able to define a Typescript Type and then instantiate it easily
  // Instead we do it manually for now.
  // Superstruct library might be able to help? https://www.npmjs.com/package/superstruct

  // ===== We Define A Corresponding Schema Type, and Instantiate 2 Objects =====

  type Outlaw = {
    name: string;
    image: string;
    power_level: number;
  };

  let kyrie: Outlaw = {
    name: 'Kyrie',
    image: 'https://wikipedia.org/',
    power_level: 199990
  };

  let jin: Outlaw = {
    name: 'Jin',
    image: 'https://google.com/something/',
    power_level: 100000000
  };

  // ===== Read metatada from objects on-chain =====

  // let moveCall: RawMoveCall = {
  //   packageObjectId: '0x2f1c9c3610d58f793e936821b797b9b63d9e602a',
  //   module: 'outlaw_sky',
  //   function: 'view',
  //   typeArguments: [],
  //   arguments: ['outlaw id', 'schema id']
  // };

  // let viewQuery = await provider.devInspectMoveCall('@0x99', moveCall);

  // ===== Write metadata to objects on-chain =====

  let something = Array.from(bcs.ser('Outlaw', kyrie).toBytes());

  // Serializer that does vector<vector<u8>>

  let kyrieBytes = Array.from(bcs.ser('Outlaw', kyrie).toBytes());
  let jinBytes = bcs.ser('Outlaw', jin);

  // signer.signAndExecuteTransaction('', 'WaitForEffectsCert');
  // signer.executeMoveCall('', 'WaitForEffectsCert');

  const moveCallTxn = await signer.executeMoveCall({
    packageObjectId: '0x06958bb1c0368b0bc3af0c48d51bbf9b158d365a',
    module: 'outlaw_sky',
    function: 'create',
    typeArguments: [],
    arguments: ['0xa37996c78eed10373b5f71cf6b5cd2b623494ba5', kyrieBytes],
    gasBudget: 15000
  });

  console.log(kyrieBytes);

  let monkey = [3, 19, false, 'something'];

  let utf8Bytes_ser = bcs.ser('vector<utf8>', [
    'empty',
    'zero word',
    'first word',
    'hope this works',
    'large character set',
    'another one',
    'more more'
  ]);
  let utf8Bytes = Array.from(utf8Bytes_ser.toBytes());
  console.log(utf8Bytes);

  // let dataBuffer = await signer.serializer.serializeToBytes('', '', 'Commit');

  // let something2 = await signer.signAndExecuteTransaction(dataBuffer, 'WaitForEffectsCert');

  // let txnResults = await provider.executeTransaction(
  //   '',
  //   'ED25519',
  //   '',
  //   'pubkey',
  //   'WaitForEffectsCert'
  // );

  // const moveCallTxn = await signer.executeMoveCall({
  //   packageObjectId: '0xfe593d3458c272f73e2595ee7f6cfb00b9597315',
  //   module: 'bcs_maybe',
  //   function: 'give_bcs',
  //   typeArguments: [],
  //   arguments: [utf8Bytes2],
  //   gasBudget: 200000
  // });

  // RPC Error: Could not serialize argument of type Vector(U8) at 0 into vector<u8>

  console.log(moveCallTxn);
}

async function defineSchema() {
  // const moveCallTxn = await signer.signAndExecuteTransaction(txBytes, 'WaitForEffectsCert');

  // ===== Create a Schema Manually =====

  const moveCallTxn = await signer.executeMoveCall({
    packageObjectId: '0x4f2801f232f4cd689e7d1791b74e7fad1dfa068c',
    module: 'schema',
    function: 'define',
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

// defineSchema();

readAndWriteOutlaw();

async function readAndWriteOutlaw() {
  // Do a read and modify cycle; just keep going back and forth

  let result = await provider.devInspectMoveCall('0xed2c39b73e055240323cf806a7d8fe46ced1cabb', {
    packageObjectId: '0x06958bb1c0368b0bc3af0c48d51bbf9b158d365a',
    module: 'outlaw_sky',
    function: 'view',
    typeArguments: [],
    arguments: [
      '0x655b7e5c9a0011c67aebcffe05f4a07f99d7785d',
      '0xa37996c78eed10373b5f71cf6b5cd2b623494ba5'
    ]
  });

  console.log(result.results);
  // view returns:
  //    38,   3,   5,  75, 121, 114, 105, 101,  22,
  // 104, 116, 116, 112, 115,  58,  47,  47, 119,
  // 105, 107, 105, 112, 101, 100, 105,  97,  46,
  // 111, 114, 103,  47,  54,  13,   3,   0,   0,
  //   0,   0,   0
  // Total number of bytes: 38
  // Total number of items: 3
  // Length of first item: 5
  // Bytes for first item: 75, 121, 114, 105, 101
  // ..

  // view_alternate vector<vector<u8>> returns:
  // [ 3, 6, 5,  75, 121, 114, 105, 101, 23,  22, 104,
  //116, 116, 112, 115,  58,  47, 47, 119, 105,
  // 107, 105, 112, 101, 100, 105, 97,  46, 111,
  // 114, 103,  47,   8,  54,  13,  3,   0,   0,
  // 0,   0,   0 ]
  // total number of items: 3 (vectors)
  // total length of of vector: 6
  // total length of bytes for first item: 5
  // bytes for first item: 75, 121, 114, 105, 101
  // total length of vector: 23
  // total length of bytes for vector: 22
  // bytes for second item: 104, 116...
  // total length of vector: 8
  // bytes for u64: 54, 13, ...
  //
  // In other words, returning vector<vector<u8>> is stupid because it gets REALLY garbled
  // it simply takes the vectors, prepends each of them with their length, smooshes them together,
  // then adds ANOTHER byte that specifies length (total number of vectors)

  let data = parseViewResults(result);
  console.log(data);

  const objectQuery = await provider.getObject('0xd4f5b8cd2fc7cead99e70cbd9d0667ee94c862c8');
  let schema = parseSchema(objectQuery);

  if (!schema) return; // We got a bad response, nothing more we can do for now
  bcs.registerStructType('Outlaw', schema);

  // We have to manually define the TS type for now
  type Outlaw = {
    name: string;
    image: string;
    power_level: number;
  };

  let outlaw = bcs.de('Outlaw', new Uint8Array(data)) as Outlaw;
  console.log(outlaw);

  // Modify the metadata
  outlaw.name = 'Kizande';
  outlaw.power_level = 917;

  let kyrieBytes = Array.from(bcs.ser('Outlaw', outlaw).toBytes());

  const moveCallTxn = await signer.executeMoveCall({
    packageObjectId: '0x06958bb1c0368b0bc3af0c48d51bbf9b158d365a',
    module: 'outlaw_sky',
    function: 'overwrite',
    typeArguments: [],
    // outlaw id, keys, data bytes, schema
    arguments: [
      '0x655b7e5c9a0011c67aebcffe05f4a07f99d7785d',
      ['name', 'image', 'power_level'],
      kyrieBytes,
      '0xa37996c78eed10373b5f71cf6b5cd2b623494ba5'
    ],
    gasBudget: 15000
  });
}

// ====== Helper Parse Functions ======

// TO DO: there is probably a better way to do this, without the ts-ignore and the null response
function parseSchema(objectQuery: GetObjectDataResponse): { [key: string]: string } | null {
  try {
    let response: { [key: string]: string } = {};

    // @ts-ignore
    objectQuery.details.data.fields.schema.forEach(item => {
      if (item.fields.optional) response[item.fields.key] = `Option<${item.fields.type}>`;
      else response[item.fields.key] = item.fields.type;
    });
    return response;
  } catch (err) {
    return null;
  }
}

// TO DO: we can't just remove the 0th byte; that will fail if there are more than 128 items
// in the response. Instead we need to parse the ULEB128 and remove that
function parseViewResults(result: DevInspectResults): number[] {
  // @ts-ignore
  let data = result.results.Ok[0][1].returnValues[0][0] as number[];

  // Delete the first tunnecessary ULEB128 length auto-added by the sui bcs view-function response
  data.splice(0, 1);
  // data.splice(0, 1);

  return data;
}
