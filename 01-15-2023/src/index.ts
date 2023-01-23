import {
  Ed25519Keypair,
  JsonRpcProvider,
  RawSigner,
  GetObjectDataResponse,
  RawMoveCall,
  SuiJsonValue,
  MoveCallTransaction,
  SignableTransaction,
  UnserializedSignableTransaction
} from '@mysten/sui.js';
import { BCS, BcsConfig } from '@mysten/bcs';

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

// ===== We Fetch A Schema =====

createOutlaw();

async function createOutlaw() {
  const objectQuery = await provider.getObject('0xd4f5b8cd2fc7cead99e70cbd9d0667ee94c862c8');

  let schema = parseSchema(objectQuery);

  if (!schema) return; // We got a bad response, nothing more we can do for now
  bcs.registerStructType('Outlaw', schema);

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

  let kyrieBytes = Array.from(bcs.ser('Outlaw', kyrie).toBytes());
  let jinBytes = bcs.ser('Outlaw', jin);

  // signer.signAndExecuteTransaction('', 'WaitForEffectsCert');
  // signer.executeMoveCall('', 'WaitForEffectsCert');

  const moveCallTxn = await signer.executeMoveCall({
    packageObjectId: '0x6d4d5cd3a4ab962539513d011c03c8a10b5a31f6',
    module: 'outlaw_sky',
    function: 'create',
    typeArguments: [],
    arguments: ['0x5169e6aa4136ce330960559ec8416f32b2d01645', kyrieBytes],
    gasBudget: 15000
  });

  console.log(kyrieBytes);

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

async function readOutlaw() {
  // TO DO
  // Do a read and modify cycle; just keep going back and forth
}
