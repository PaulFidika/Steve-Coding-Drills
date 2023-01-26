type MapSchemaTypes = {
  ascii: string;
  'Option<ascii>': string | undefined | null;
  utf8: string;
  u8: number;
  u64: bigint;
  boolean: boolean;
};

type MapSchema<T extends Record<string, keyof MapSchemaTypes>> = {
  -readonly [K in keyof T]: MapSchemaTypes[T[K]];
};

const OutlawSchema = {
  name: 'Option<ascii>',
  image: 'ascii',
  power_level: 'u64'
} as const;

type Outlaw = MapSchema<typeof OutlawSchema>;

const person2: Outlaw = {
  name: undefined,
  image: 'https://www.wikipedia.org/',
  power_level: 199n
};

console.log(person2);
console.log(typeof person2);
