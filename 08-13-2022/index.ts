class arrayItem {
  a: number;
  b: number;
  c: number;

  constructor(a: number, b: number, c: number) {
    this.a = a;
    this.b = b;
    this.c = c;
  }
}

const array = [new arrayItem(5, 9, 0), new arrayItem(0, 18, -1), new arrayItem(-5, 5, 1)];

function sortArray(array: arrayItem[]) {
  for (let item of array) {
    if (item.a === 0) console.log('found it!');
  }
}

const whatever = { x: 9, y: 10, z: 11 };

for (let key in whatever) {
  console.log(key);
}

// continue here:
// https://medium.com/swlh/making-objects-iterable-in-javascript-252d9e270be6#:~:text=entries.,iterator%5D()%20method.

// also why is ts-node types required?
