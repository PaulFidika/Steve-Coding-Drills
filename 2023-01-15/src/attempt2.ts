interface Circle {
  type: 'circle';
  radius: number;
}

interface Square {
  type: 'square';
  length: number;
}

type TypeName = 'circle' | 'square';

type ObjectType<T> = T extends 'circle' ? Circle : T extends 'square' ? Square : never;

const shapes: (Circle | Square)[] = [
  { type: 'circle', radius: 1 },
  { type: 'circle', radius: 2 },
  { type: 'square', length: 10 }
];

function getItems<T extends TypeName>(type: T): ObjectType<T>[] {
  return shapes.filter(s => s.type == type) as ObjectType<T>[];
}

const circles = getItems('circle');
for (const circle of circles) {
  console.log(circle.radius);
}

type PrimitiveType = 'ascii' | 'utf8' | 'u64';
type PrimitiveJSType = string | number;

type InferredJSType<T> = T extends 'ascii' | 'utf8' ? string : T extends 'u64' ? number : never;

function instantiate2<T extends PrimitiveType>(
  _fieldType: T,
  fieldValue: InferredJSType<T>
): InferredJSType<T> {
  return fieldValue as InferredJSType<T>;
}

let whatever2 = instantiate2('ascii', 'Kyrie');
let whatever3 = instantiate2('u64', 97);

function tryAgain<T extends PrimitiveType>(
  types: T[],
  values: InferredJSType<T>[]
): InferredJSType<T>[] {
  let response: InferredJSType<T>[] = [];
  let i = 0;
  for (const fieldType of types) {
    response[i] = instantiate2(fieldType, values[i]);
    i += 1;
  }
  return response;
}

let d = tryAgain(['ascii', 'u64'], [13, 14, 'thingy']);

function instantiate<T extends PrimitiveType>(
  types: { [key: string]: T },
  fields: { [key: string]: InferredJSType<T> }
): { [key: string]: InferredJSType<T> } {
  let response = {};

  let k: keyof { [key: string]: T };
  for (k in types) {
    let type = types[k];
    response[k] = instantiate2(type, fields[k]);
  }

  //   let keys = Object.keys(types);

  //   keys.map((key, i) => {
  //     let type = types[key];
  //     response[key] = fields[i];
  //   });

  return response;
}

let someone = instantiate(
  { name: 'ascii', power_level: 'u64' },
  { name: 'Kyrie', power_level: '16' }
);

someone.name = 15;

// ============= try to implement =============

let recipe: { [key: string]: PrimitiveType } = {
  name: 'utf8',
  image: 'ascii',
  power_level: 'u64'
};

let fields = ['Kyrie', 'https://something', 197];

let whatever = instantiate(recipe, fields);

console.log(whatever);
