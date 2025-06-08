console.log('hello there');

type Item = {
  a: number;
  b: number;
  c: number;
};

let data = [
  { a: 5, b: 10, c: 99 } as Item,
  { a: 9, b: 8, c: 15 } as Item,
  { a: -10, b: 11, c: 17 } as Item,
  { a: 100, b: 10, c: 15 } as Item,
  { a: 11, b: 10, c: 15 } as Item,
  { a: 18, b: 10, c: 15 } as Item,
  { a: 17, b: 10, c: 15 } as Item
];

const sorter = (data: Item[], largestFirst: boolean, key: keyof Item): Item[] => {
  let sortedData: Item[] = [];

  for (let item of data) {
    if (sortedData[0]) {
      if (sortedData[0][key] && item[key] > sortedData[0][key]) {
        sortedData.unshift(item);
      } else sortedData.push(item);
    } else sortedData.push(item);
  }

  return sortedData;
};

console.log(sorter(data, true, 'a'));
