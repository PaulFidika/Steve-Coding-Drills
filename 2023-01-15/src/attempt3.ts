// // type PrimitiveType = 'ascii' | 'utf8' | 'u64';

// type JSInferred<Type> = {
//   [Property in keyof Type]: Type[Property] extends 'ascii' | 'utf8'
//     ? string
//     : Type[Property] extends 'Option<ascii>' | 'Option<utf8>'
//     ? string | null | undefined
//     : Type[Property] extends 'u64' | undefined
//     ? bigint
//     : never;
// };

// // TO DO:
// // - We could have 'name?:' to make a paramter optional, instead of 'Option<ascii'>
// //
// // - 'name: undefined' still needs to be written, even if the field is allowed to be undefined. It cannot just be excluded. There
// // are some craft typescript ways to fix this I think
// //
// // - Ideally we'd want to restrict types to just our PrimitiveType list; idk how to do that. If you define a field that has a
// // type that maps to 'never' in JSInferred<Type>, like say you invent some new type like 'monkey', then it will become impossible
// // to instantiate the object, because the object you try to create will always be missing a field, but the field's value is 'never'
// // which can never be implemented lol.

// // Create our schema in JS

// type OutlawSchema = {
//   name: 'Option<ascii>';
//   image: 'ascii';
//   power_level: 'u64';
// };

// // instantiate our schema
// let person: JSInferred<OutlawSchema> = {
//   name: undefined,
//   image: 'https://',
//   power_level: 99n
// };

// function update<T, Key extends keyof T>(object: T, key: Key, value: T[Key]) {
//   object[key] = value;
// }

// update(person, 'image', 'http://new');
// person.name = 'Kyrie';

// function getTypes<T extends Object>(object: T) {
//   let keys = Object.keys(object);
//   let z = typeof keys;
//   keys.map(key => {
//     console.log(typeof key);
//   });
// }

// getTypes(person);

// // I don't think it's possible to read the OutlawSchema values and convert it into something more interesting?
