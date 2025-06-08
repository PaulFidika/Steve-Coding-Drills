type Item = {
  a: number;
  b: number;
  c: number;
};

const data = [
  { a: 5, b: 19, c: 0.0001 } as Item,
  { a: 7, b: -15, c: 0 } as Item,
  { a: 6, b: 19, c: 25 } as Item,
  { a: -1, b: 0, c: -10 } as Item
];

const sortingFunction = (data: Item[], largestFirst: boolean, key: keyof Item) => {
  data.map((item, index) => {
    const nextItem = data[index + 1];
    const nextValue = nextItem ? nextItem[key] : Infinity;
    console.log(nextValue);

    if (item[key] > item[key]) {
    }
  });
};

sortingFunction(data, true, 'a');
