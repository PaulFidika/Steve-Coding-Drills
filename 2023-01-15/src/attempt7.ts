import { BCS } from '@mysten/bcs';
import {
  assert,
  object,
  integer,
  bigint,
  string,
  boolean,
  number,
  array,
  is,
  define,
  optional,
  never,
  union,
  any,
  Struct
} from 'superstruct';

type JSTypes<T extends Record<string, keyof MoveToJSTypes>> = {
  -readonly [K in keyof T]: MoveToJSTypes[T[K]];
};

type MoveToJSTypes = {
  address: Uint8Array;
  bool: boolean;
  id: Uint8Array;
  u8: number;
  u16: number;
  u32: number;
  u64: bigint;
  u128: bigint;
  u256: bigint;
  ascii: string;
  utf8: string;
  'vector<u8>': number[];
  'Option<ascii>': { some: string } | { none: null };
};

const MoveToStruct: Record<string, any> = {
  address: array(number()),
  bool: boolean(),
  id: array(number()),
  u8: integer(),
  u16: integer(),
  u32: integer(),
  u64: bigint(),
  u128: bigint(),
  u256: bigint(),
  ascii: string(),
  utf8: string(),
  'vector<u8>': array(integer()),
  'Option<ascii>': union([object({ none: any() }), object({ some: string() })])
};

// ====== Helper Functions ======

function moveStructValidator(
  schema: Record<string, string>
): Struct<{ [x: string]: any }, Record<string, any>> {
  const dynamicStruct: Record<string, any> = {};

  Object.keys(schema).map(key => {
    dynamicStruct[key] = MoveToStruct[schema[key]];
  });

  return object(dynamicStruct);
}

function serialize<T>(bcs: BCS, structName: string, data: T): number[] {
  return Array.from(bcs.ser(structName, data).toBytes());
}

// ====== Define Move Schema, JSType, and StructValidator ======

const outlawSchema = {
  name: 'ascii',
  description: 'Option<ascii>',
  image: 'ascii',
  power_level: 'u64'
} as const;

let outlawValidator = moveStructValidator(outlawSchema);

type Outlaw = JSTypes<typeof outlawSchema>;

// ====== Instantiate and Validate ======

const kyrie: Outlaw = {
  name: 'Kyrie',
  description: { none: null },
  image: 'https://www.wikipedia.org/',
  power_level: 199n
};

assert(kyrie, outlawValidator);

// ====== Interact with the Sui Network ======

import { bcs, signer } from './configured_bcs';

const OUTLAW_SKY_PACKAGE_ID = '0x8af0cbf8380f7738907ef9aee9aab4e34c3d0716';
const SCHEMA_ID = '0x37cef7c69de4b1cea22f1ef445940432d6968ac6';

// How do we know that our JS Schema here == the Move Schema on-chain?
// We can enter an ObjectID, and assume that its contents == our JS Schema here

bcs.registerStructType('Outlaw', outlawSchema);
let kyrieBytes = serialize(bcs, 'Outlaw', kyrie);

console.log(kyrieBytes);

// This will post data as vector<u8>
async function create() {
  const moveCallTxn = await signer.executeMoveCall({
    packageObjectId: OUTLAW_SKY_PACKAGE_ID,
    module: 'outlaw_sky',
    function: 'create',
    typeArguments: [],
    arguments: [SCHEMA_ID, kyrieBytes],
    gasBudget: 15000
  });
}
// create();

async function overwriteSubset() {
  // How do we serialize just a subset of keys, rather than the entire object?
  let KeysToUpdate = ['name', 'description'];
  let bytes = Array.from(bcs.ser('Outlaw', outlaw).toBytes());

  let objectID = '0xc6c3028a0df2eb49af8cf766971c9b2cf5a8d0c2';
  let schemaID = '0x37cef7c69de4b1cea22f1ef445940432d6968ac6';

  const moveCallTxn = await signer.executeMoveCall({
    packageObjectId: OUTLAW_SKY_PACKAGE_ID,
    module: 'outlaw_sky',
    function: 'overwrite',
    typeArguments: [],
    // outlaw id, keys, data bytes, schema
    arguments: [objectID, KeysToUpdate, bytes, SCHEMA_ID],
    gasBudget: 15000
  });
}
// overwriteSubset()
