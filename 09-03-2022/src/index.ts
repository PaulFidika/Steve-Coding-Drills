// Finished this in only 26.5 mins; it includes the ability to handle undefined
// values. LargestFirst toggle is not used.
// This is the same test they gave me at Phantom Labs

type Item = {
  a?: number;
  b?: number;
  c?: number;
};

const sortFunction = (data: Item[], key: keyof Item, largestFirst: boolean): Item[] => {
  let sortedData: Item[] = [];

  for (let [index, item] of data.entries()) {
    if (!item[key]) {
      sortedData.splice(sortedData.length, 0, item);
      continue;
    }
    for (let x = 0; x <= sortedData.length; x++) {
      let comparisonItem = sortedData[x]?.[key];
      if (!comparisonItem) {
        sortedData.splice(x, 0, item);
        break;
      } else if (item[key]! > comparisonItem) {
        sortedData.splice(x, 0, item);
        break;
      }
    }
  }
  return sortedData;
};

let data = [
  { a: 5, b: 69, c: 0 } as Item,
  { a: 9, b: 5, c: 99999 } as Item,
  { a: 99, b: 69 } as Item,
  { a: 0, b: 17, c: 7 } as Item,
  { a: -100, b: 11, c: 420 } as Item,
  { a: -100, b: 11 } as Item
];

let sortedData = sortFunction(data, 'c', true);
console.log(sortedData);
