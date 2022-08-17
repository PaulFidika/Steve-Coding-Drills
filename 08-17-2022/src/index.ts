import data from './data';
import { Item } from './types';

function sortArray(data: Item[], largestFirst: boolean, key: keyof Item) {
  const sortedArray: Item[] = [];

  for (let [index, item] of data.entries()) {
    for (let i = 0; i < sortedArray.length; i++) {
      if (item[key] > sortedArray[i][key]) {
        sortedArray.splice(i, 0, item);
        break;
      }
    }

    if (index === 0) sortedArray.push(item);
  }

  return sortedArray;
}

console.log(sortArray(data, true, 'a'));
