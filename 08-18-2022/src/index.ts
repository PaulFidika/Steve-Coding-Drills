import someFunction from './second';

/// <reference path = "ambient.d.ts" />

someFunction();

let baboon: Monkey = { name: 'Harambe' };

declare global {
  type SampleType = {
    a: String;
  };

  type Item = {
    a: number;
    b: number;
    c: number;
  };
}
